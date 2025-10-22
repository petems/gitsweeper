package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

const (
	// MaxCommitsToCheck limits how many commits to check for merged branches.
	MaxCommitsToCheck = 10000
	// ConcurrentWorkers defines how many goroutines to use for concurrent processing.
	ConcurrentWorkers = 4
	// BatchSize for processing commits in batches.
	BatchSize = 100
)

// BranchInfo holds branch information.
type BranchInfo struct {
	Name   string
	Hash   plumbing.Hash
	Remote string
	Short  string
}

// commitBatch represents a batch of commits to process.
type commitBatch struct {
	commits  []*object.Commit
	startIdx int
}

// RemoteBranches returns an iterator over all remote branch references in the repository.
// It filters the reference store to return only references that represent remote branches.
func RemoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}

// ParseBranchname splits a branch string of the form "remote/branch" into the remote and branch name.
// If the input contains no slash, the entire input is returned as the remote and the branch name is empty.
func ParseBranchname(branchString string) (remote, branchname string) {
	if idx := strings.IndexByte(branchString, '/'); idx > 0 {
		return branchString[:idx], branchString[idx+1:]
	}
	return branchString, ""
}

// DeleteBranch deletes the named branch from the given remote by invoking
// `git push <remote> --delete <branchShortName>`.
//
// We shell out to git instead of using go-git's push operations to avoid complex
// authentication handling. The go-git library has significant limitations with various
// authentication methods (SSH keys with passphrases, SSH agents, credential helpers,
// tokens, deploy keys, etc.). By using the system git command, we leverage the user's
// existing authentication configuration automatically.
// See: https://github.com/go-git/go-git/issues/28
//
// The command runs with a 30-second timeout and returns an error containing
// the combined output if the deletion fails.
func DeleteBranch(repo *git.Repository, remote, branchShortName string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}
	repoPath := worktree.Filesystem.Root()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "push", remote, "--delete", branchShortName)
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to delete branch %s on remote %s: %w\nOutput: %s",
			branchShortName, remote, err, string(output))
	}

	return nil
}

func RemoteBranchesToStrings(gitRemoteArray []*git.Remote) []string {
	stringArray := make([]string, len(gitRemoteArray))
	for i, v := range gitRemoteArray {
		stringArray[i] = v.Config().Name
	}
	return stringArray
}

func GetCurrentDirAsGitRepo() (*git.Repository, error) {
	LogInfo("Getting current working directory")

	dir, err := os.Getwd()
	if err != nil {
		LogFatalError("Error opening current directory:", err)
		return nil, err
	}

	LogInfof("Attempting to open Git directory at %s", dir)

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// GetMergedBranches finds branches that have been merged into the master branch.
func GetMergedBranches(remoteOrigin, masterBranchName, skipBranches string) ([]string, error) {
	repo, err := GetCurrentDirAsGitRepo()
	if err != nil {
		return nil, err
	}

	// Convert skip branches to a set for O(1) lookups
	var skipSet map[string]bool
	if skipBranches != "" {
		skipSet = StringSliceToSet(strings.Split(skipBranches, ","))
	} else {
		skipSet = make(map[string]bool)
	}

	LogInfo("Attempting to get master information from branches from repo")

	// Get branch heads
	branchHeads, err := getBranchHeads(repo)
	if err != nil {
		return nil, err
	}

	fmt.Println("Fetching from the remote...")

	// Validate remote exists
	listRemotes, err := repo.Remotes()
	if err != nil {
		LogFatalError("Error looking for remotes", err)
		return nil, err
	}

	remoteBranchesAsStrings := RemoteBranchesToStrings(listRemotes)
	if !IsStringInSet(remoteOrigin, StringSliceToSet(remoteBranchesAsStrings)) {
		return nil, errors.New("Could not find the remote named " + remoteOrigin)
	}

	// Get master commits with context and timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	masterHash, exists := branchHeads[masterBranchName]
	if !exists {
		return nil, fmt.Errorf("master branch %s not found", masterBranchName)
	}

	masterCommits, err := repo.Log(&git.LogOptions{From: masterHash})
	if err != nil {
		LogFatalError("get commits from master failed:", err)
		return nil, err
	}

	// Get remote branches
	remoteBranches, err := getRemoteBranches(repo, remoteOrigin, masterBranchName, skipSet)
	if err != nil {
		return nil, err
	}

	// Early exit if no branches to check
	if len(remoteBranches) == 0 {
		LogInfo("No branches found for the specified origin")
		return []string{}, nil
	}

	LogInfof("Origin has been set to '%s', checking %d branches", remoteOrigin, len(remoteBranches))

	// Use concurrent processing for large branch sets
	if len(remoteBranches) > 10 {
		return findMergedBranchesConcurrent(ctx, masterCommits, remoteBranches)
	}

	// Use sequential processing for smaller sets
	return findMergedBranchesSequential(ctx, masterCommits, remoteBranches)
}

// getBranchHeads gets all branch heads.
func getBranchHeads(repo *git.Repository) (map[string]plumbing.Hash, error) {
	branchRefs, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("list branches failed: %w", err)
	}

	branchHeads := make(map[string]plumbing.Hash)

	err = branchRefs.ForEach(func(reference *plumbing.Reference) error {
		branchName := strings.TrimPrefix(reference.Name().String(), "refs/heads/")
		branchHeads[branchName] = reference.Hash()
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("iterating branches failed: %w", err)
	}

	return branchHeads, nil
}

// getRemoteBranches gets remote branches with filtering.
func getRemoteBranches(
	repo *git.Repository,
	remoteOrigin string,
	masterBranchName string,
	skipSet map[string]bool,
) ([]BranchInfo, error) {
	remoteBranches, err := RemoteBranches(repo.Storer)
	if err != nil {
		return nil, fmt.Errorf("list remote branches failed: %w", err)
	}

	var branches []BranchInfo
	masterBranchRemote := fmt.Sprintf("%s/%s", remoteOrigin, masterBranchName)

	err = remoteBranches.ForEach(func(branch *plumbing.Reference) error {
		remoteBranchName := strings.TrimPrefix(branch.Name().String(), "refs/remotes/")

		// Skip master branch
		if remoteBranchName == masterBranchRemote {
			return nil
		}

		remote, shortBranchName := ParseBranchname(remoteBranchName)

		// Filter by origin and skip list
		if remote == remoteOrigin {
			// Check if this branch should be skipped
			if IsStringInSet(shortBranchName, skipSet) {
				LogInfof("Branch '%s' matches skip branch string '[%s]'", remoteBranchName, shortBranchName)
				return nil
			}

			branches = append(branches, BranchInfo{
				Name:   remoteBranchName,
				Hash:   branch.Hash(),
				Remote: remote,
				Short:  shortBranchName,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("iterating remote branches failed: %w", err)
	}

	return branches, nil
}

// findMergedBranchesSequential processes branches sequentially with optimizations.
func findMergedBranchesSequential(
	ctx context.Context,
	masterCommits object.CommitIter,
	branches []BranchInfo,
) ([]string, error) {
	// Create hash lookup map
	branchHashMap := make(map[string][]BranchInfo, len(branches))
	for _, branch := range branches {
		hashStr := branch.Hash.String()
		branchHashMap[hashStr] = append(branchHashMap[hashStr], branch)
	}

	var mergedBranches []string
	foundBranches := make(map[string]bool)
	commitCount := 0

	err := masterCommits.ForEach(func(commit *object.Commit) error {
		// Check context for cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Limit the number of commits to check
		commitCount++
		if commitCount > MaxCommitsToCheck {
			LogInfof("Reached maximum commit limit (%d), stopping search", MaxCommitsToCheck)
			return errors.New("max commits reached")
		}

		// Early termination when all branches found
		if len(foundBranches) == len(branches) {
			return errors.New("all branches found")
		}

		// Check if this commit hash matches any branch
		commitHash := commit.Hash.String()
		if branchInfos, exists := branchHashMap[commitHash]; exists {
			for _, branchInfo := range branchInfos {
				if !foundBranches[branchInfo.Name] {
					LogInfof(
						"Branch %s head (%s) was found in master, so has been merged!",
						branchInfo.Name,
						commitHash,
					)
					mergedBranches = append(mergedBranches, branchInfo.Name)
					foundBranches[branchInfo.Name] = true
				}
			}
		}

		return nil
	})

	// Handle expected early termination errors
	if err != nil && err.Error() != "all branches found" && err.Error() != "max commits reached" {
		return nil, fmt.Errorf("looking for merged commits failed: %w", err)
	}

	sort.Strings(mergedBranches)
	return mergedBranches, nil
}

// findMergedBranchesConcurrent processes branches using concurrent workers.
func findMergedBranchesConcurrent(
	ctx context.Context,
	masterCommits object.CommitIter,
	branches []BranchInfo,
) ([]string, error) {
	// Create hash lookup map
	branchHashMap := make(map[string][]BranchInfo, len(branches))
	for _, branch := range branches {
		hashStr := branch.Hash.String()
		branchHashMap[hashStr] = append(branchHashMap[hashStr], branch)
	}

	// Channel for commit batches
	commitBatches := make(chan commitBatch, ConcurrentWorkers*2)
	results := make(chan []string, ConcurrentWorkers)

	// Start worker goroutines
	var wg sync.WaitGroup
	numWorkers := minInt(ConcurrentWorkers, runtime.NumCPU())

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mergedInWorker := processCommitBatches(ctx, commitBatches, branchHashMap, len(branches))
			results <- mergedInWorker
		}()
	}

	// Producer goroutine to batch commits
	go func() {
		defer close(commitBatches)

		var batch []*object.Commit
		commitCount := 0
		batchStartIdx := 0

		err := masterCommits.ForEach(func(commit *object.Commit) error {
			// Check context for cancellation
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Limit the number of commits to check
			commitCount++
			if commitCount > MaxCommitsToCheck {
				return errors.New("max commits reached")
			}

			batch = append(batch, commit)

			// Send batch when it's full
			if len(batch) >= BatchSize {
				select {
				case commitBatches <- commitBatch{commits: batch, startIdx: batchStartIdx}:
					batch = make([]*object.Commit, 0, BatchSize)
					batchStartIdx = commitCount
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			return nil
		})

		// Send remaining commits
		if len(batch) > 0 && err == nil {
			select {
			case commitBatches <- commitBatch{commits: batch, startIdx: batchStartIdx}:
			case <-ctx.Done():
			}
		}
	}()

	// Wait for workers and collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and merge results from all workers
	var allMerged []string
	seenBranches := make(map[string]bool)

	for workerResults := range results {
		for _, branch := range workerResults {
			if !seenBranches[branch] {
				allMerged = append(allMerged, branch)
				seenBranches[branch] = true
			}
		}
	}

	sort.Strings(allMerged)

	// Check if context was cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return allMerged, nil
	}
}

// processCommitBatches processes batches of commits in a worker goroutine.
func processCommitBatches(
	ctx context.Context,
	batches <-chan commitBatch,
	branchHashMap map[string][]BranchInfo,
	totalBranches int,
) []string {
	var mergedBranches []string
	foundBranches := make(map[string]bool)

	for batch := range batches {
		// Check context for cancellation
		select {
		case <-ctx.Done():
			return mergedBranches
		default:
		}

		// Early termination if all branches found
		if len(foundBranches) >= totalBranches {
			return mergedBranches
		}

		// Process commits in this batch
		for _, commit := range batch.commits {
			commitHash := commit.Hash.String()
			if branchInfos, exists := branchHashMap[commitHash]; exists {
				for _, branchInfo := range branchInfos {
					if !foundBranches[branchInfo.Name] {
						LogInfof(
							"Branch %s head (%s) was found in master, so has been merged!",
							branchInfo.Name,
							commitHash,
						)
						mergedBranches = append(mergedBranches, branchInfo.Name)
						foundBranches[branchInfo.Name] = true
					}
				}
			}
		}
	}

	return mergedBranches
}

// min returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
