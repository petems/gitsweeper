package internal

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// initTestRepo creates a git repo in a temp dir with an initial commit on master.
func initTestRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(out))
	}

	run("init", "-b", "master")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")

	// Create initial commit.
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0o644))
	run("add", "README.md")
	run("commit", "-m", "initial commit")

	return dir
}

func TestCherryCheckBranch_SquashMerged(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	dir := initTestRepo(t)

	run := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(out))
	}

	// Create a feature branch with commits.
	run("checkout", "-b", "feature-branch")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "feature.txt"), []byte("feature content\n"), 0o644))
	run("add", "feature.txt")
	run("commit", "-m", "add feature")

	// Go back to master and squash-merge.
	run("checkout", "master")
	run("merge", "--squash", "feature-branch")
	run("commit", "-m", "squash merge feature-branch")

	merged, err := CherryCheckBranch(dir, "master", "feature-branch", 0)
	require.NoError(t, err)
	assert.True(t, merged, "squash-merged branch should be detected as merged")
}

func TestCherryCheckBranch_NotMerged(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	dir := initTestRepo(t)

	run := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(out))
	}

	// Create a feature branch with commits but don't merge.
	run("checkout", "-b", "unmerged-branch")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "unmerged.txt"), []byte("unmerged content\n"), 0o644))
	run("add", "unmerged.txt")
	run("commit", "-m", "add unmerged feature")

	run("checkout", "master")

	merged, err := CherryCheckBranch(dir, "master", "unmerged-branch", 0)
	require.NoError(t, err)
	assert.False(t, merged, "unmerged branch should not be detected as merged")
}

func TestCherryCheckBranch_SameAsMaster(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	dir := initTestRepo(t)

	run := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(out))
	}

	// Create branch at same point as master (no extra commits).
	run("checkout", "-b", "same-as-master")
	run("checkout", "master")

	merged, err := CherryCheckBranch(dir, "master", "same-as-master", 0)
	require.NoError(t, err)
	assert.True(t, merged, "branch at same point as master should be detected as merged")
}

func TestCherryCheckBranch_MultipleCommitsSquashed(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	dir := initTestRepo(t)

	run := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(out))
	}

	// Create a feature branch with multiple commits.
	run("checkout", "-b", "multi-commit")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("content1\n"), 0o644))
	run("add", "file1.txt")
	run("commit", "-m", "add file1")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("content2\n"), 0o644))
	run("add", "file2.txt")
	run("commit", "-m", "add file2")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file3.txt"), []byte("content3\n"), 0o644))
	run("add", "file3.txt")
	run("commit", "-m", "add file3")

	// Go back to master and squash-merge.
	run("checkout", "master")
	run("merge", "--squash", "multi-commit")
	run("commit", "-m", "squash merge multi-commit")

	merged, err := CherryCheckBranch(dir, "master", "multi-commit", 0)
	require.NoError(t, err)
	assert.True(t, merged, "squash-merged branch with multiple commits should be detected as merged")
}

func TestCherryCheckBranch_InvalidRef(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	dir := initTestRepo(t)

	// Invalid branch ref should return false (not merged), not error the whole run.
	merged, err := CherryCheckBranch(dir, "master", "nonexistent-branch", 0)
	assert.False(t, merged)
	assert.NoError(t, err)
}

func TestCherryCheckBranch_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name      string
		repoPath  string
		upstream  string
		branchRef string
		errMsg    string
	}{
		{
			name:      "empty repo path",
			repoPath:  "",
			upstream:  "master",
			branchRef: "feature",
			errMsg:    "repo path cannot be empty",
		},
		{
			name:      "empty upstream",
			repoPath:  "/tmp",
			upstream:  "",
			branchRef: "feature",
			errMsg:    "upstream ref cannot be empty",
		},
		{
			name:      "empty branch ref",
			repoPath:  "/tmp",
			upstream:  "master",
			branchRef: "",
			errMsg:    "branch ref cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			merged, err := CherryCheckBranch(tc.repoPath, tc.upstream, tc.branchRef, 0)
			assert.False(t, merged)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestCherryCheckBranches_Concurrent(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	dir := initTestRepo(t)

	run := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(out))
	}

	// Branch 1: squash-merged.
	run("checkout", "-b", "merged-1")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "m1.txt"), []byte("m1\n"), 0o644))
	run("add", "m1.txt")
	run("commit", "-m", "merged-1 commit")
	run("checkout", "master")
	run("merge", "--squash", "merged-1")
	run("commit", "-m", "squash merge merged-1")

	// Branch 2: squash-merged.
	run("checkout", "-b", "merged-2")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "m2.txt"), []byte("m2\n"), 0o644))
	run("add", "m2.txt")
	run("commit", "-m", "merged-2 commit")
	run("checkout", "master")
	run("merge", "--squash", "merged-2")
	run("commit", "-m", "squash merge merged-2")

	// Branch 3: not merged.
	run("checkout", "-b", "unmerged-1")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "u1.txt"), []byte("u1\n"), 0o644))
	run("add", "u1.txt")
	run("commit", "-m", "unmerged-1 commit")
	run("checkout", "master")

	branches := []BranchInfo{
		{Name: "merged-1", Short: "merged-1"},
		{Name: "merged-2", Short: "merged-2"},
		{Name: "unmerged-1", Short: "unmerged-1"},
	}

	results, err := CherryCheckBranches(dir, "master", branches, 0)
	require.NoError(t, err)

	assert.Len(t, results, 2)
	assert.Contains(t, results, "merged-1")
	assert.Contains(t, results, "merged-2")
	assert.NotContains(t, results, "unmerged-1")
}
