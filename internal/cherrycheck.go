package internal

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// DetectionMethod indicates how a merged branch was detected.
type DetectionMethod string

const (
	// DetectionHashMatch means the branch tip was found in master's commit history.
	DetectionHashMatch DetectionMethod = "hash-match"
	// DetectionCherry means git cherry determined all commits were applied upstream.
	DetectionCherry DetectionMethod = "cherry"
)

// MergedBranchResult holds a merged branch name and how it was detected.
type MergedBranchResult struct {
	Name   string
	Method DetectionMethod
}

// MergeDetectionOptions controls merge detection behavior.
type MergeDetectionOptions struct {
	MaxCommits    int  // 0 = use default (MaxCommitsToCheck)
	DisableCherry bool // --no-deep-check sets this to true
}

// CherryCheckBranch checks if a branch has been fully merged into upstream.
// It first tries `git cherry` (works for single-commit squash merges, rebases,
// and cherry-picks). If that finds unapplied commits, it falls back to comparing
// the combined branch diff via `git patch-id` (works for multi-commit squash merges).
func CherryCheckBranch(repoPath, upstream, branchRef string, maxCommits int) (bool, error) {
	if repoPath == "" {
		return false, errors.New("repo path cannot be empty")
	}
	if upstream == "" {
		return false, errors.New("upstream ref cannot be empty")
	}
	if branchRef == "" {
		return false, errors.New("branch ref cannot be empty")
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return false, fmt.Errorf("git command not found in PATH: %w", err)
	}

	// Pass 1: git cherry (handles single-commit squash, rebase, cherry-pick)
	merged, cherryErr := cherryCheck(repoPath, gitPath, upstream, branchRef)
	if cherryErr != nil {
		LogInfof("git cherry check failed for %s, skipping: %s", branchRef, cherryErr)
		return false, nil
	}
	if merged {
		return true, nil
	}

	// Pass 2: patch-id comparison (handles multi-commit squash merges)
	return patchIDCheck(repoPath, gitPath, upstream, branchRef, maxCommits)
}

// cherryCheck runs `git cherry` and returns true if all commits are applied.
func cherryCheck(repoPath, gitPath, upstream, branchRef string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, gitPath, "cherry", upstream, branchRef)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	trimmedOutput := strings.TrimSpace(string(output))

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, fmt.Errorf("timeout running git cherry for %s after 30s: %w", branchRef, err)
		}
		LogInfof("git cherry failed for %s: %s (output: %s)", branchRef, err, trimmedOutput)
		return false, err
	}

	// Empty output means no commits between merge-base and branch.
	if trimmedOutput == "" {
		return true, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(trimmedOutput))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "+") {
			return false, nil
		}
	}

	return true, nil
}

// patchIDCheck compares the combined diff of a branch against recent upstream
// commits using git patch-id. This detects multi-commit squash merges where
// multiple branch commits were combined into a single upstream commit.
func patchIDCheck(repoPath, gitPath, upstream, branchRef string, maxCommits int) (bool, error) {
	if maxCommits <= 0 {
		maxCommits = MaxCommitsToCheck
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get the combined diff of the branch relative to upstream.
	branchPatchID, err := getCombinedPatchID(ctx, repoPath, gitPath, upstream, branchRef)
	if err != nil || branchPatchID == "" {
		LogInfof("Could not get patch-id for branch %s: %v", branchRef, err)
		return false, nil
	}

	// Get all upstream patch-ids in a single pass by piping git log -p
	// into git patch-id. This uses only 2 processes instead of 2 per commit.
	maxCountFlag := fmt.Sprintf("--max-count=%d", maxCommits)
	//nolint:gosec // upstream/branchRef are validated caller inputs; maxCountFlag is derived from an int
	logCmd := exec.CommandContext(
		ctx, gitPath, "log", "-p",
		upstream, "--not", branchRef, maxCountFlag,
	)
	logCmd.Dir = repoPath

	patchIDCmd := exec.CommandContext(ctx, gitPath, "patch-id", "--stable")
	patchIDCmd.Dir = repoPath

	pipe, logPipeErr := logCmd.StdoutPipe()
	if logPipeErr != nil {
		return false, nil
	}
	patchIDCmd.Stdin = pipe

	if startErr := logCmd.Start(); startErr != nil {
		return false, nil
	}

	patchOutput, patchErr := patchIDCmd.Output()
	_ = logCmd.Wait() //nolint:errcheck // best-effort cleanup after pipe consumer finished

	if patchErr != nil {
		LogInfof("Could not get upstream patch-ids: %v", patchErr)
		return false, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(patchOutput)))
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) >= 2 && parts[0] == branchPatchID {
			LogInfof("Branch %s matched upstream commit %s via patch-id", branchRef, parts[1])
			return true, nil
		}
	}

	return false, nil
}

// getCombinedPatchID returns the patch-id of the combined diff of branchRef relative to upstream.
func getCombinedPatchID(ctx context.Context, repoPath, gitPath, upstream, branchRef string) (string, error) {
	diffRange := upstream + "..." + branchRef
	//nolint:gosec // upstream and branchRef are validated caller inputs, not user-controlled shell args
	diffCmd := exec.CommandContext(ctx, gitPath, "diff", diffRange)
	diffCmd.Dir = repoPath

	patchIDCmd := exec.CommandContext(ctx, gitPath, "patch-id", "--stable")
	patchIDCmd.Dir = repoPath

	pipe, err := diffCmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	patchIDCmd.Stdin = pipe

	if startErr := diffCmd.Start(); startErr != nil {
		return "", startErr
	}

	patchOutput, patchErr := patchIDCmd.Output()
	if patchErr != nil {
		_ = diffCmd.Wait() //nolint:errcheck // best-effort cleanup after pipe consumer failed
		return "", patchErr
	}
	_ = diffCmd.Wait() //nolint:errcheck // diff already fully consumed by patch-id cmd

	return extractPatchID(string(patchOutput)), nil
}

// extractPatchID extracts the patch-id hash from `git patch-id` output.
// Output format is: "<patch-id> <commit-id>".
func extractPatchID(output string) string {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return ""
	}
	parts := strings.Fields(trimmed)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// CherryCheckBranches runs CherryCheckBranch concurrently for multiple branches
// and returns the names of branches that are fully merged.
func CherryCheckBranches(repoPath, upstream string, branches []BranchInfo, maxCommits int) ([]string, error) {
	if len(branches) == 0 {
		return nil, nil
	}

	LogInfof("Running git cherry check on %d branches", len(branches))

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []string
	)

	// Bounded worker pool.
	sem := make(chan struct{}, ConcurrentWorkers)

	for _, branch := range branches {
		wg.Add(1)
		go func(b BranchInfo) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			merged, err := CherryCheckBranch(repoPath, upstream, b.Name, maxCommits)
			if err != nil {
				LogInfof("Cherry check error for %s: %s", b.Name, err)
				return
			}

			if merged {
				LogInfof("Branch %s detected as merged via git cherry", b.Name)
				mu.Lock()
				results = append(results, b.Name)
				mu.Unlock()
			}
		}(branch)
	}

	wg.Wait()
	return results, nil
}
