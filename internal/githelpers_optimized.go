//go:build optimized

package internal

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

func RemoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}

func ParseBranchname(branchString string) (remote, branchname string) {
	branchArray := strings.SplitN(branchString, "/", 2)
	remote = branchArray[0]
	branchname = branchArray[1]
	return
}

func DeleteBranch(repo *git.Repository, remote, branchShortName string) error {
	// Create refspec for deleting the remote branch
	deleteRefSpec := config.RefSpec(fmt.Sprintf(":%s", plumbing.NewBranchReferenceName(branchShortName)))

	// Push with delete refspec
	err := repo.Push(&git.PushOptions{
		RemoteName: remote,
		RefSpecs:   []config.RefSpec{deleteRefSpec},
	})

	if err != nil {
		return fmt.Errorf("failed to delete branch %s on remote %s: %w", branchShortName, remote, err)
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
	LogInfof("Current working directory is %s", dir)

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

func GetMergedBranches(remoteOrigin, masterBranchName, skipBranches string) ([]string, error) {
	repo, err := GetCurrentDirAsGitRepo()
	if err != nil {
		return nil, err
	}

	skipSlice := strings.Split(skipBranches, ",")

	LogInfo("Attempting to get master information from branches from repo")

	branchRefs, err := repo.Branches()
	if err != nil {
		LogFatalError("list branches failed:", err)
	}

	branchHeads := make(map[string]plumbing.Hash)

	fmt.Println("Fetching from the remote...")

	listRemotes, err := repo.Remotes()
	if err != nil {
		LogFatalError("Error looking for remotes", err)
		return nil, err
	}

	remoteBranchesAsStrings := RemoteBranchesToStrings(listRemotes)

	if !IsStringInSlice(remoteOrigin, remoteBranchesAsStrings) {
		return nil, errors.New("Could not find the remote named " + remoteOrigin)
	}

	err = branchRefs.ForEach(func(reference *plumbing.Reference) error {
		branchName := strings.TrimPrefix(reference.Name().String(), "refs/heads/")
		branchHead := reference.Hash()
		branchHeads[branchName] = branchHead
		return nil
	})

	if err != nil {
		LogFatalError("list branches failed:", err)
		return nil, err
	}

	masterCommits, err := repo.Log(&git.LogOptions{From: branchHeads[masterBranchName]})
	if err != nil {
		LogFatalError("get commits from master failed:", err)
		return nil, err
	}

	remoteBranches, err := RemoteBranches(repo.Storer)
	if err != nil {
		LogFatalError("list remote branches failed:", err)
		return nil, err
	}

	remoteBranchHeads := make(map[string]plumbing.Hash)

	err = remoteBranches.ForEach(func(branch *plumbing.Reference) error {
		remoteBranchName := strings.TrimPrefix(branch.Name().String(), "refs/remotes/")
		remoteBranchHead := branch.Hash()
		_, shortBranchName := ParseBranchname(remoteBranchName)
		if slices.Contains(skipSlice, shortBranchName) {
			LogInfof("Branch '%s' matches skip branch string '%s'", remoteBranchName, skipSlice)
		} else {
			remoteBranchHeads[remoteBranchName] = remoteBranchHead
		}
		return nil
	})

	if err != nil {
		LogFatalError("iterating remote branches failed:", err)
		return nil, err
	}

	mergedBranches := make([]string, 0)
	masterBranchRemote := fmt.Sprintf("%s/%s", remoteOrigin, masterBranchName)
	delete(remoteBranchHeads, masterBranchRemote)

	LogInfof("Origin has been set to '%s', restricting branches to preview to that origin", remoteOrigin)

	// Filter branches by origin early to reduce search space
	originBranches := make(map[string]plumbing.Hash)
	for branchName, branchHead := range remoteBranchHeads {
		remote, _ := ParseBranchname(branchName)
		if remote == remoteOrigin {
			LogInfof("Branch '%s' matches remote '%s'", branchName, remoteOrigin)
			originBranches[branchName] = branchHead
		} else {
			LogInfof("Branch '%s' does not match remote '%s', not adding", branchName, remoteOrigin)
		}
	}

	// Early exit if no branches to check
	if len(originBranches) == 0 {
		LogInfo("No branches found for the specified origin")
		return mergedBranches, nil
	}

	// Create a hash set for faster lookups - Handle multiple branches with same hash
	branchHashes := make(map[string][]string, len(originBranches))
	for branchName, branchHead := range originBranches {
		hashStr := branchHead.String()
		branchHashes[hashStr] = append(branchHashes[hashStr], branchName)
	}

	// Iterate through commits with early termination
	foundBranches := make(map[string]bool)
	err = masterCommits.ForEach(func(commit *object.Commit) error {
		// Early termination: stop when all branches are found
		if len(foundBranches) == len(originBranches) {
			return fmt.Errorf("all branches found") // Use error to break early
		}

		commitHash := commit.Hash.String()
		if branchNames, exists := branchHashes[commitHash]; exists {
			for _, branchName := range branchNames {
				if !foundBranches[branchName] {
					LogInfof("Branch %s head (%s) was found in master, so has been merged!", branchName, commitHash)
					mergedBranches = append(mergedBranches, branchName)
					foundBranches[branchName] = true
				}
			}
		}
		return nil
	})

	// Handle early termination error (which is expected)
	if err != nil && err.Error() != "all branches found" {
		LogFatalError("looking for merged commits failed:", err)
		return nil, err
	}

	sort.Strings(mergedBranches)
	return mergedBranches, nil
}
