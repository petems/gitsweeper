package internal

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
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
