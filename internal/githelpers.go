package internal

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	gitshell "code.gitea.io/git"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"

	log "github.com/sirupsen/logrus"
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

	deleteBranchShortName := fmt.Sprintf(":%s", branchShortName)

	err := gitshell.Push(".", gitshell.PushOptions{Remote: remote, Branch: deleteBranchShortName})

	return err
}

func RemoteBranchesToStrings(gitRemoteArray []*git.Remote) []string {
	stringArray := make([]string, len(gitRemoteArray))
	for i, v := range gitRemoteArray {
		stringArray[i] = v.Config().Name
	}
	return stringArray
}

func GetCurrentDirAsGitRepo() (*git.Repository, error) {
	log.Info("Getting current working directory")

	dir, err := os.Getwd()

	log.Infof("Current working directory is %s", dir)

	if err != nil {
		FatalError("Error opening current directory:", err)
		return nil, err
	}

	log.Infof("Attempting to open Git drectory at %s", dir)

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

	log.Info("Attempting to get master information from branches from repo")

	branchRefs, err := repo.Branches()
	if err != nil {
		FatalError("list branches failed:", err)
	}

	branchHeads := make(map[string]plumbing.Hash)

	fmt.Println("Fetching from the remote...")

	listRemotes, err := repo.Remotes()

	if err != nil {
		FatalError("Error looking for remotes", err)
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
		FatalError("list branches failed:", err)
		return nil, err
	}

	masterCommits, err := repo.Log(&git.LogOptions{From: branchHeads[masterBranchName]})

	if err != nil {
		FatalError("get commits from master failed:", err)
		return nil, err
	}

	remoteBranches, err := RemoteBranches(repo.Storer)

	if err != nil {
		FatalError("list remote branches failed:", err)
		return nil, err
	}

	remoteBranchHeads := make(map[string]plumbing.Hash)

	err = remoteBranches.ForEach(func(branch *plumbing.Reference) error {
		remoteBranchName := strings.TrimPrefix(branch.Name().String(), "refs/remotes/")
		remoteBranchHead := branch.Hash()
		remoteBranchHeads[remoteBranchName] = remoteBranchHead
		return nil
	})

	if err != nil {
		FatalError("iterating remote branches failed:", err)
		return nil, err
	}

	for branchName, branchHead := range remoteBranchHeads {
		log.Infof("Remote Branch %s head is: %s", branchName, branchHead)
	}

	mergedBranches := make([]string, 0)

	masterBranchRemote := fmt.Sprintf("%s/%s", remoteOrigin, masterBranchName)

	delete(remoteBranchHeads, masterBranchRemote)

	log.Infof("Origin has been set to '%s', restricting branches to preview to that origin", remoteOrigin)

	err = masterCommits.ForEach(func(commit *object.Commit) error {
		for branchName, branchHead := range remoteBranchHeads {
			remote, _ := ParseBranchname(branchName)
			if remote == remoteOrigin {
				log.Infof("Branch '%s' matches remote '%s'", branchName, remoteOrigin)
				if branchHead.String() == commit.Hash.String() {
					log.Infof("Branch %s head (%s) was found in master, so has been merged!\n", branchName, branchHead)
					mergedBranches = append(mergedBranches, branchName)
				}
			} else {
				log.Infof("Branch '%s' does not match remote '%s', not adding to", branchName, remoteOrigin)
			}
		}
		return nil
	})

	if err != nil {
		FatalError("looking for merged commits failed:", err)
		return nil, err
	}

	sort.Strings(mergedBranches)

	return mergedBranches, nil
}
