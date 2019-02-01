package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	log "github.com/sirupsen/logrus"

	"gopkg.in/alecthomas/kingpin.v2"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	debug = kingpin.Flag("debug", "Enable debug mode.").Bool()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	log.SetLevel(log.FatalLevel)
	log.SetOutput(os.Stdout)

	if *debug {
		log.SetLevel(log.InfoLevel)
	}

	log.Info("Getting current working directory")
	dir, err := os.Getwd()
	log.Infof("Current working directory is %s", dir)
	if err != nil {
		log.Fatalf("Error opening current directory: %s", err)
		return
	}

	log.Infof("Attempting to open Git drectory at %s", dir)
	repo, err := git.PlainOpen(dir)

	if err != nil {
		log.Fatalf("Error reading %s as Git repo: %s", dir, err)
		return
	}

	log.Info("Attempting to list branches from repo")
	branchRefs, err := repo.Branches()
	if err != nil {
		log.Fatalf("list branches failed: %s", err)
	}

	branchHeads := make(map[string]plumbing.Hash)

	err = branchRefs.ForEach(func(reference *plumbing.Reference) error {
		branchName := strings.TrimPrefix(reference.Name().String(), "refs/heads/")
		branchHead := reference.Hash()
		branchHeads[branchName] = branchHead
		return nil
	})

	if err != nil {
		log.Fatalf("list branches failed: %s", err)
		return
	}

	fmt.Printf("There are %d branches\n", len(branchHeads))

	for branchName, branchHead := range branchHeads {
		fmt.Printf("Branch %s head is: %s\n", branchName, branchHead)
	}

	nonMasterBranchRefs := branchHeads

	delete(nonMasterBranchRefs, "master")

	masterCommits, err := repo.Log(&git.LogOptions{From: branchHeads["master"]})

	err = masterCommits.ForEach(func(commit *object.Commit) error {
		for branchName, branchHead := range nonMasterBranchRefs {
			if branchHead.String() == commit.Hash.String() {
				fmt.Printf("Branch %s head (%s) was found in master, so has been merged!\n", branchName, branchHead)
			}
		}
		return nil
	})

}
