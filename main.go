package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	git "gopkg.in/src-d/go-git.v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("gitsweeper", "A command-line tool for cleaning up merged branches.")

	debug = app.Flag("debug", "Enable debug mode.").Bool()

	preview = app.Command("preview", "Show the branches to delete.")

	cleanup = app.Command("cleanup", "Delete the remote branches.")
)

func setupLogger() {
	log.SetOutput(os.Stderr)
	textFormatter := new(prefixed.TextFormatter)
	textFormatter.FullTimestamp = true
	textFormatter.TimestampFormat = "01 Jan 2019 15:04:05"
	log.SetFormatter(textFormatter)
	log.SetLevel(log.FatalLevel)
}

func fatalError(msg string, err error) {
	if err != nil {
		log.WithError(err).Fatal(msg)
	} else {
		log.Fatal(msg)
	}
	os.Exit(1)
}

func remoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}

func main() {
	kingpin.Version("0.0.1")

	setupLogger()

	if *debug {
		log.SetLevel(log.InfoLevel)
		log.Info("--debug setting detected - Info level logs enabled")
	}

	log.Info("Getting current working directory")

	dir, err := os.Getwd()

	log.Infof("Current working directory is %s", dir)

	if err != nil {
		fatalError("Error opening current directory:", err)
		return
	}

	log.Infof("Attempting to open Git drectory at %s", dir)

	repo, err := git.PlainOpen(dir)

	if err != nil {
		fatalError(fmt.Sprintf("Error reading %s as Git repo", dir), err)
		return
	}

	log.Info("Attempting to get master information from branches from repo")

	branchRefs, err := repo.Branches()
	if err != nil {
		fatalError("list branches failed:", err)
	}

	branchHeads := make(map[string]plumbing.Hash)

	fmt.Println("Fetching from the remote...")

	err = branchRefs.ForEach(func(reference *plumbing.Reference) error {
		branchName := strings.TrimPrefix(reference.Name().String(), "refs/heads/")
		branchHead := reference.Hash()
		branchHeads[branchName] = branchHead
		return nil
	})

	if err != nil {
		fatalError("list branches failed:", err)
		return
	}

	masterCommits, err := repo.Log(&git.LogOptions{From: branchHeads["master"]})

	remoteBranches, err := remoteBranches(repo.Storer)

	remoteBranchHeads := make(map[string]plumbing.Hash)

	err = remoteBranches.ForEach(func(branch *plumbing.Reference) error {
		remoteBranchName := strings.TrimPrefix(branch.Name().String(), "refs/remotes/")
		remoteBranchHead := branch.Hash()
		remoteBranchHeads[remoteBranchName] = remoteBranchHead
		return nil
	})

	for branchName, branchHead := range remoteBranchHeads {
		log.Infof("Remote Branch %s head is: %s", branchName, branchHead)
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case preview.FullCommand():

		mergedBranches := make([]string, 0)

		err = masterCommits.ForEach(func(commit *object.Commit) error {
			for branchName, branchHead := range remoteBranchHeads {
				if (branchHead.String() == commit.Hash.String()) && (branchName != "origin/master") {
					log.Infof("Branch %s head (%s) was found in master, so has been merged!\n", branchName, branchHead)
					mergedBranches = append(mergedBranches, branchName)
				}
			}
			return nil
		})

		if len(mergedBranches) == 0 {
			fmt.Println("No branches already merged into master!")
		} else {
			sort.Strings(mergedBranches)
			fmt.Println("\nThese branches have been merged into master:")
			for _, branchName := range mergedBranches {
				fmt.Printf("  %s\n", branchName)
			}
		}
	case cleanup.FullCommand():
		fmt.Println("Cleanup has not been implemented yet!")
		os.Exit(1)
	}

}
