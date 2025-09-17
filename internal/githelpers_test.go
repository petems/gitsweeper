package internal

import (
	"testing"

	memfs "github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
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

func TestGetRemoteBranchesOptimizedSkipsCustomMaster(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	assert.NoError(t, err)

	masterHash := plumbing.NewHash("1111111111111111111111111111111111111111")
	featureHash := plumbing.NewHash("2222222222222222222222222222222222222222")

	err = repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName("refs/remotes/origin/main"), masterHash))
	assert.NoError(t, err)
	err = repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName("refs/remotes/origin/feature-branch"), featureHash))
	assert.NoError(t, err)

	branches, err := getRemoteBranchesOptimized(repo, "origin", "main", map[string]bool{})
	assert.NoError(t, err)
	assert.Len(t, branches, 1)
	assert.Equal(t, "origin/feature-branch", branches[0].Name)
}
