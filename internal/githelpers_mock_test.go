package internal

import (
	"io"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/stretchr/testify/assert"
)

// Mock Storer
type mockStorer struct {
	storer.Storer
	refs []*plumbing.Reference
}

func (s *mockStorer) IterReferences() (storer.ReferenceIter, error) {
	return storer.NewReferenceSliceIter(s.refs), nil
}

func (s *mockStorer) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	for _, ref := range s.refs {
		if ref.Name() == name {
			return ref, nil
		}
	}
	return nil, plumbing.ErrReferenceNotFound
}

// Mock CommitIter
type mockCommitIter struct {
	object.CommitIter
	commits []*object.Commit
	idx     int
}

func (i *mockCommitIter) Next() (*object.Commit, error) {
	if i.idx >= len(i.commits) {
		return nil, io.EOF
	}
	commit := i.commits[i.idx]
	i.idx++
	return commit, nil
}

func (i *mockCommitIter) ForEach(f func(*object.Commit) error) error {
	for {
		c, err := i.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := f(c); err != nil {
			if err.Error() == "all branches found" || err.Error() == "max commits reached" {
				return nil
			}
			return err
		}
	}
	return nil
}

func (i *mockCommitIter) Close() {}

// Mock Repository
type mockRepository struct {
	storer  storer.ReferenceStorer
	commits []*object.Commit
	remotes []*git.Remote
}

func (r *mockRepository) Storer() storer.ReferenceStorer {
	return r.storer
}

func (r *mockRepository) Log(opts *git.LogOptions) (object.CommitIter, error) {
	return &mockCommitIter{commits: r.commits}, nil
}

func (r *mockRepository) Remotes() ([]*git.Remote, error) {
	return r.remotes, nil
}

func TestGetMergedBranches_Mock(t *testing.T) {
	// Commits
	commit1 := &object.Commit{Hash: plumbing.NewHash("1111111111111111111111111111111111111111")}
	commit2 := &object.Commit{Hash: plumbing.NewHash("2222222222222222222222222222222222222222")}
	commit3 := &object.Commit{Hash: plumbing.NewHash("3333333333333333333333333333333333333333")}
	commit4 := &object.Commit{Hash: plumbing.NewHash("4444444444444444444444444444444444444444")}

	masterCommits := []*object.Commit{commit1, commit2, commit3}

	testCases := []struct {
		name           string
		branches       []*plumbing.Reference
		masterCommits  []*object.Commit
		skipBranches   string
		expectedResult []string
		expectedErr    error
	}{
		{
			name: "no merged branches",
			branches: []*plumbing.Reference{
				plumbing.NewHashReference(plumbing.NewBranchReferenceName("master"), commit1.Hash),
				plumbing.NewHashReference(plumbing.NewRemoteReferenceName("origin", "feature1"), commit4.Hash),
			},
			masterCommits:  masterCommits,
			expectedResult: []string{},
		},
		{
			name: "one merged branch",
			branches: []*plumbing.Reference{
				plumbing.NewHashReference(plumbing.NewBranchReferenceName("master"), commit1.Hash),
				plumbing.NewHashReference(plumbing.NewRemoteReferenceName("origin", "feature1"), commit2.Hash),
			},
			masterCommits:  masterCommits,
			expectedResult: []string{"origin/feature1"},
		},
		{
			name: "multiple merged branches",
			branches: []*plumbing.Reference{
				plumbing.NewHashReference(plumbing.NewBranchReferenceName("master"), commit1.Hash),
				plumbing.NewHashReference(plumbing.NewRemoteReferenceName("origin", "feature1"), commit2.Hash),
				plumbing.NewHashReference(plumbing.NewRemoteReferenceName("origin", "feature2"), commit3.Hash),
			},
			masterCommits:  masterCommits,
			expectedResult: []string{"origin/feature1", "origin/feature2"},
		},
		{
			name: "with skipped branches",
			branches: []*plumbing.Reference{
				plumbing.NewHashReference(plumbing.NewBranchReferenceName("master"), commit1.Hash),
				plumbing.NewHashReference(plumbing.NewRemoteReferenceName("origin", "feature1"), commit2.Hash),
				plumbing.NewHashReference(plumbing.NewRemoteReferenceName("origin", "feature2"), commit3.Hash),
			},
			masterCommits:  masterCommits,
			skipBranches:   "feature1",
			expectedResult: []string{"origin/feature2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockRepository{
				storer:  &mockStorer{refs: tc.branches},
				commits: tc.masterCommits,
				remotes: []*git.Remote{git.NewRemote(nil, &config.RemoteConfig{Name: "origin"})},
			}
			result, err := getMergedBranches(mockRepo, "origin", "master", tc.skipBranches)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expectedResult, result)
			}
		})
	}
}
