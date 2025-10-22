package internal

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBranchname(t *testing.T) {
	remote, branchShort := ParseBranchname("origin/cool-branch-name")

	assert.Equal(t, "origin", remote)
	assert.Equal(t, "cool-branch-name", branchShort)
}

func TestParseBranchname_OneSlash(t *testing.T) {
	remote, branchShort := ParseBranchname("origin/janedoe/cool-branch-name")

	assert.Equal(t, "origin", remote)
	assert.Equal(t, "janedoe/cool-branch-name", branchShort)
}

func TestParseBranchname_TwoSlashes(t *testing.T) {
	remote, branchShort := ParseBranchname("origin/janedoe/cool-branch-name/and-another-branch-after")

	assert.Equal(t, "origin", remote)
	assert.Equal(t, "janedoe/cool-branch-name/and-another-branch-after", branchShort)
}

func TestParseBranchname_ThreeSlashes(t *testing.T) {
	remote, branchShort := ParseBranchname("origin/janedoe/cool-branch-name/and-another-branch-after/one-more-for-luck")

	assert.Equal(t, "origin", remote)
	assert.Equal(t, "janedoe/cool-branch-name/and-another-branch-after/one-more-for-luck", branchShort)
}

func TestGetRemoteBranches(t *testing.T) {
	testCases := []struct {
		name                string
		masterBranchName    string
		remoteBranches      map[string]string
		skipBranches        map[string]bool
		expectedBranches    []string
		expectedBranchCount int
	}{
		{
			name:             "skips custom master branch",
			masterBranchName: "main",
			remoteBranches: map[string]string{
				"refs/remotes/origin/main":           "1111111111111111111111111111111111111111",
				"refs/remotes/origin/feature-branch": "2222222222222222222222222222222222222222",
			},
			skipBranches:        map[string]bool{},
			expectedBranches:    []string{"origin/feature-branch"},
			expectedBranchCount: 1,
		},
		{
			name:             "no other branches",
			masterBranchName: "main",
			remoteBranches: map[string]string{
				"refs/remotes/origin/main": "1111111111111111111111111111111111111111",
			},
			skipBranches:        map[string]bool{},
			expectedBranches:    []string{},
			expectedBranchCount: 0,
		},
		{
			name:             "skips branch from skip list",
			masterBranchName: "master",
			remoteBranches: map[string]string{
				"refs/remotes/origin/master":         "1111111111111111111111111111111111111111",
				"refs/remotes/origin/feature-branch": "2222222222222222222222222222222222222222",
				"refs/remotes/origin/dont-delete":    "3333333333333333333333333333333333333333",
			},
			skipBranches:        map[string]bool{"dont-delete": true},
			expectedBranches:    []string{"origin/feature-branch"},
			expectedBranchCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, err := git.Init(memory.NewStorage(), memfs.New())
			require.NoError(t, err)

			for branchName, hashStr := range tc.remoteBranches {
				hash := plumbing.NewHash(hashStr)
				setErr := repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(branchName), hash))
				require.NoError(t, setErr)
			}

			branches, err := getRemoteBranches(repo, "origin", tc.masterBranchName, tc.skipBranches)
			require.NoError(t, err)
			assert.Len(t, branches, tc.expectedBranchCount)

			branchNames := make([]string, len(branches))
			for i, b := range branches {
				branchNames[i] = b.Name
			}

			assert.ElementsMatch(t, tc.expectedBranches, branchNames)
		})
	}
}

// Additional comprehensive tests for ParseBranchname edge cases.
func TestParseBranchname_EmptyString(t *testing.T) {
	remote, branchShort := ParseBranchname("")

	assert.Empty(t, remote)
	assert.Empty(t, branchShort)
}

func TestParseBranchname_NoSlash(t *testing.T) {
	remote, branchShort := ParseBranchname("originonly")

	assert.Equal(t, "originonly", remote)
	assert.Empty(t, branchShort)
}

func TestParseBranchname_LeadingSlash(t *testing.T) {
	// Leading slash should parse the entire string as remote since idx would be 0
	remote, branchShort := ParseBranchname("/branch-name")

	assert.Equal(t, "/branch-name", remote)
	assert.Empty(t, branchShort)
}

func TestParseBranchname_TrailingSlash(t *testing.T) {
	remote, branchShort := ParseBranchname("origin/feature-branch/")

	assert.Equal(t, "origin", remote)
	assert.Equal(t, "feature-branch/", branchShort)
}

func TestParseBranchname_MultipleConsecutiveSlashes(t *testing.T) {
	remote, branchShort := ParseBranchname("origin//feature//branch")

	assert.Equal(t, "origin", remote)
	assert.Equal(t, "/feature//branch", branchShort)
}

func TestParseBranchname_OnlySlash(t *testing.T) {
	remote, branchShort := ParseBranchname("/")

	assert.Equal(t, "/", remote)
	assert.Empty(t, branchShort)
}

func TestParseBranchname_SpecialCharacters(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedRemote string
		expectedBranch string
	}{
		{
			name:           "hyphenated branch",
			input:          "origin/feature-branch-name",
			expectedRemote: "origin",
			expectedBranch: "feature-branch-name",
		},
		{
			name:           "underscored branch",
			input:          "origin/feature_branch_name",
			expectedRemote: "origin",
			expectedBranch: "feature_branch_name",
		},
		{
			name:           "dots in branch",
			input:          "origin/release/1.2.3",
			expectedRemote: "origin",
			expectedBranch: "release/1.2.3",
		},
		{
			name:           "uppercase remote",
			input:          "ORIGIN/branch",
			expectedRemote: "ORIGIN",
			expectedBranch: "branch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			remote, branchShort := ParseBranchname(tc.input)
			assert.Equal(t, tc.expectedRemote, remote)
			assert.Equal(t, tc.expectedBranch, branchShort)
		})
	}
}

// Tests for RemoteBranches function.
func TestRemoteBranches_EmptyRepository(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	iter, err := RemoteBranches(repo.Storer)
	require.NoError(t, err)
	require.NotNil(t, iter)

	// Count references
	count := 0
	err = iter.ForEach(func(_ *plumbing.Reference) error {
		count++
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Empty repository should have no remote branches")
}

func TestRemoteBranches_OnlyLocalBranches(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	// Add local branches only
	localBranches := []string{
		"refs/heads/master",
		"refs/heads/develop",
		"refs/heads/feature",
	}

	for _, branchName := range localBranches {
		hash := plumbing.NewHash("1111111111111111111111111111111111111111")
		setErr := repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(branchName), hash))
		require.NoError(t, setErr)
	}

	iter, err := RemoteBranches(repo.Storer)
	require.NoError(t, err)

	// Count remote references (should be 0)
	count := 0
	err = iter.ForEach(func(_ *plumbing.Reference) error {
		count++
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Repository with only local branches should have no remote branches")
}

func TestRemoteBranches_MultipleRemoteBranches(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	// Add multiple remote branches
	remoteBranches := map[string]string{
		"refs/remotes/origin/master":    "1111111111111111111111111111111111111111",
		"refs/remotes/origin/develop":   "2222222222222222222222222222222222222222",
		"refs/remotes/upstream/master":  "3333333333333333333333333333333333333333",
		"refs/remotes/upstream/feature": "4444444444444444444444444444444444444444",
		"refs/heads/local-branch":       "5555555555555555555555555555555555555555",
	}

	for branchName, hashStr := range remoteBranches {
		hash := plumbing.NewHash(hashStr)
		setErr := repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(branchName), hash))
		require.NoError(t, setErr)
	}

	iter, err := RemoteBranches(repo.Storer)
	require.NoError(t, err)

	// Collect remote branch names
	var foundBranches []string
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		foundBranches = append(foundBranches, ref.Name().String())
		return nil
	})
	require.NoError(t, err)

	// Should only include remote branches (4 total), not the local branch
	assert.Len(t, foundBranches, 4)
	assert.Contains(t, foundBranches, "refs/remotes/origin/master")
	assert.Contains(t, foundBranches, "refs/remotes/origin/develop")
	assert.Contains(t, foundBranches, "refs/remotes/upstream/master")
	assert.Contains(t, foundBranches, "refs/remotes/upstream/feature")
	assert.NotContains(t, foundBranches, "refs/heads/local-branch")
}

func TestRemoteBranches_FiltersCorrectly(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	// Add various types of references
	references := map[string]string{
		"refs/remotes/origin/master":  "1111111111111111111111111111111111111111",
		"refs/heads/master":           "2222222222222222222222222222222222222222",
		"refs/tags/v1.0.0":            "3333333333333333333333333333333333333333",
		"refs/remotes/origin/feature": "4444444444444444444444444444444444444444",
	}

	for refName, hashStr := range references {
		hash := plumbing.NewHash(hashStr)
		setErr := repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(refName), hash))
		require.NoError(t, setErr)
	}

	iter, err := RemoteBranches(repo.Storer)
	require.NoError(t, err)

	// Collect and verify only remote branches are returned
	var foundBranches []string
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		foundBranches = append(foundBranches, ref.Name().String())
		// Verify each reference is actually a remote branch
		assert.True(t, ref.Name().IsRemote(), "Reference %s should be a remote branch", ref.Name())
		return nil
	})
	require.NoError(t, err)

	assert.Len(t, foundBranches, 2)
	assert.Contains(t, foundBranches, "refs/remotes/origin/master")
	assert.Contains(t, foundBranches, "refs/remotes/origin/feature")
	assert.NotContains(t, foundBranches, "refs/heads/master")
	assert.NotContains(t, foundBranches, "refs/tags/v1.0.0")
}

// Tests for DeleteBranch function.
// Note: Since DeleteBranch now shells out to git command, we can't fully test it
// without a real git environment, but we can test the function signature and error handling

func TestDeleteBranch_NilRepository(t *testing.T) {
	// The function now ignores the repository parameter (marked with _)
	// so it should work even with nil
	err := DeleteBranch(nil, "origin", "test-branch")

	// We expect an error because git command will fail (no actual git repo)
	// but the function itself should handle nil repo gracefully
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete branch")
}

func TestDeleteBranch_ErrorMessageFormat(t *testing.T) {
	// Test that error messages contain the expected information
	err := DeleteBranch(nil, "test-remote", "test-branch-name")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "test-branch-name")
	assert.Contains(t, err.Error(), "test-remote")
	assert.Contains(t, err.Error(), "failed to delete branch")
}

func TestDeleteBranch_SpecialCharactersInBranchName(t *testing.T) {
	// Test various branch name formats
	testCases := []struct {
		name       string
		remote     string
		branchName string
	}{
		{
			name:       "hyphenated branch",
			remote:     "origin",
			branchName: "feature-branch-name",
		},
		{
			name:       "underscored branch",
			remote:     "origin",
			branchName: "feature_branch_name",
		},
		{
			name:       "slash in branch",
			remote:     "origin",
			branchName: "feature/sub-feature",
		},
		{
			name:       "dots in branch",
			remote:     "origin",
			branchName: "release/1.2.3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := DeleteBranch(nil, tc.remote, tc.branchName)

			// All should fail (no real git), but should handle the names without panic
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.branchName)
			assert.Contains(t, err.Error(), tc.remote)
		})
	}
}

// Tests for RemoteBranchesToStrings helper function.
func TestRemoteBranchesToStrings_EmptyArray(t *testing.T) {
	result := RemoteBranchesToStrings([]*git.Remote{})

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestRemoteBranchesToStrings_SingleRemote(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	require.NoError(t, err)

	result := RemoteBranchesToStrings([]*git.Remote{remote})

	assert.Len(t, result, 1)
	assert.Equal(t, "origin", result[0])
}

func TestRemoteBranchesToStrings_MultipleRemotes(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	origin, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	require.NoError(t, err)

	upstream, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{"https://github.com/upstream/repo.git"},
	})
	require.NoError(t, err)

	result := RemoteBranchesToStrings([]*git.Remote{origin, upstream})

	assert.Len(t, result, 2)
	assert.Contains(t, result, "origin")
	assert.Contains(t, result, "upstream")
}

// Tests for minInt helper function.
func TestMinInt_FirstSmaller(t *testing.T) {
	assert.Equal(t, 5, minInt(5, 10))
}

func TestMinInt_SecondSmaller(t *testing.T) {
	assert.Equal(t, 3, minInt(10, 3))
}

func TestMinInt_Equal(t *testing.T) {
	assert.Equal(t, 7, minInt(7, 7))
}

func TestMinInt_NegativeNumbers(t *testing.T) {
	assert.Equal(t, -10, minInt(-5, -10))
	assert.Equal(t, -10, minInt(-10, -5))
	assert.Equal(t, -5, minInt(5, -5))
}

func TestMinInt_Zero(t *testing.T) {
	assert.Equal(t, 0, minInt(0, 10))
	assert.Equal(t, 0, minInt(10, 0))
	assert.Equal(t, 0, minInt(0, 0))
	assert.Equal(t, -5, minInt(0, -5))
}

// Tests for BranchInfo struct usage.
func TestBranchInfo_StructFields(t *testing.T) {
	hash := plumbing.NewHash("1111111111111111111111111111111111111111")

	branchInfo := BranchInfo{
		Name:   "origin/feature-branch",
		Hash:   hash,
		Remote: "origin",
		Short:  "feature-branch",
	}

	assert.Equal(t, "origin/feature-branch", branchInfo.Name)
	assert.Equal(t, hash, branchInfo.Hash)
	assert.Equal(t, "origin", branchInfo.Remote)
	assert.Equal(t, "feature-branch", branchInfo.Short)
}

func TestBranchInfo_EmptyValues(t *testing.T) {
	branchInfo := BranchInfo{}

	assert.Empty(t, branchInfo.Name)
	assert.Equal(t, plumbing.ZeroHash, branchInfo.Hash)
	assert.Empty(t, branchInfo.Remote)
	assert.Empty(t, branchInfo.Short)
}

// Test constants.
func TestConstants(t *testing.T) {
	assert.Equal(t, 10000, MaxCommitsToCheck)
	assert.Equal(t, 4, ConcurrentWorkers)
	assert.Equal(t, 100, BatchSize)

	// Ensure constants are positive
	assert.Positive(t, MaxCommitsToCheck)
	assert.Positive(t, ConcurrentWorkers)
	assert.Positive(t, BatchSize)
}

// Integration test for the filtering in getRemoteBranches.
func TestGetRemoteBranches_ComplexFiltering(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	// Set up a complex scenario with multiple remotes and branches
	branches := map[string]string{
		"refs/remotes/origin/master":       "1111111111111111111111111111111111111111",
		"refs/remotes/origin/develop":      "2222222222222222222222222222222222222222",
		"refs/remotes/origin/feature/test": "3333333333333333333333333333333333333333",
		"refs/remotes/origin/skip-this":    "4444444444444444444444444444444444444444",
		"refs/remotes/upstream/master":     "5555555555555555555555555555555555555555",
		"refs/remotes/upstream/develop":    "6666666666666666666666666666666666666666",
		"refs/heads/local-only":            "7777777777777777777777777777777777777777",
	}

	for branchName, hashStr := range branches {
		hash := plumbing.NewHash(hashStr)
		setErr := repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(branchName), hash))
		require.NoError(t, setErr)
	}

	// Test filtering for origin with skip list
	skipSet := map[string]bool{
		"skip-this": true,
	}

	result, err := getRemoteBranches(repo, "origin", "master", skipSet)
	require.NoError(t, err)

	// Should have 2 branches: develop and feature/test (master is filtered, skip-this is skipped, upstream branches are filtered)
	assert.Len(t, result, 2)

	branchNames := make([]string, len(result))
	for i, b := range result {
		branchNames[i] = b.Name
	}

	assert.Contains(t, branchNames, "origin/develop")
	assert.Contains(t, branchNames, "origin/feature/test")
	assert.NotContains(t, branchNames, "origin/master")
	assert.NotContains(t, branchNames, "origin/skip-this")
	assert.NotContains(t, branchNames, "upstream/master")
}

func TestGetRemoteBranches_DifferentRemotes(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	branches := map[string]string{
		"refs/remotes/origin/main":        "1111111111111111111111111111111111111111",
		"refs/remotes/origin/feature-1":   "2222222222222222222222222222222222222222",
		"refs/remotes/upstream/main":      "3333333333333333333333333333333333333333",
		"refs/remotes/upstream/feature-2": "4444444444444444444444444444444444444444",
	}

	for branchName, hashStr := range branches {
		hash := plumbing.NewHash(hashStr)
		setErr := repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(branchName), hash))
		require.NoError(t, setErr)
	}

	// Test upstream remote
	result, err := getRemoteBranches(repo, "upstream", "main", map[string]bool{})
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "upstream/feature-2", result[0].Name)
	assert.Equal(t, "upstream", result[0].Remote)
	assert.Equal(t, "feature-2", result[0].Short)
}

func TestGetRemoteBranches_AllBranchesSkipped(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	branches := map[string]string{
		"refs/remotes/origin/main":      "1111111111111111111111111111111111111111",
		"refs/remotes/origin/feature-1": "2222222222222222222222222222222222222222",
		"refs/remotes/origin/feature-2": "3333333333333333333333333333333333333333",
	}

	for branchName, hashStr := range branches {
		hash := plumbing.NewHash(hashStr)
		setErr := repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(branchName), hash))
		require.NoError(t, setErr)
	}

	// Skip all non-master branches
	skipSet := map[string]bool{
		"feature-1": true,
		"feature-2": true,
	}

	result, err := getRemoteBranches(repo, "origin", "main", skipSet)
	require.NoError(t, err)
	assert.Empty(t, result)
}