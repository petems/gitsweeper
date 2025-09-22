package internal

import (
	"testing"

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
