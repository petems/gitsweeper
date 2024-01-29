package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	hlpr "github.com/petems/gitsweeper/internal"
)

// Version is what is returned by the `-v` flag
const Version = "0.1.0"

// gitCommit is the gitcommit its built from
var gitCommit = "development"

var (
	app = kingpin.New("gitsweeper", "A command-line tool for cleaning up merged branches.")

	debug = app.Flag("debug", "Enable debug mode.").Bool()

	preview       = app.Command("preview", "Show the branches to delete.")
	previewOrigin = preview.Flag("origin", "The name of the remote you wish to clean up").Default("origin").String()
	previewMaster = preview.Flag("master", "The name of what you consider the master branch").Default("master").String()
	previewSkip   = preview.Flag("skip", "Comma-separated list of branches to skip").String()

	cleanup       = app.Command("cleanup", "Delete the remote branches.")
	cleanupForce  = cleanup.Flag("force", "Do not ask, cleanup immediately").Default("false").Bool()
	cleanupOrigin = cleanup.Flag("origin", "The name of the remote you wish to clean up").Default("origin").String()
	cleanupMaster = cleanup.Flag("master", "The name of what you consider the master branch").Default("master").String()
	cleanupSkip   = cleanup.Flag("skip", "Comma-separated list of branches to skip").String()

	version = app.Command("version", "Show the version.")
)

func main() {
	kingpin.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case preview.FullCommand():

		hlpr.SetupLogger(*debug)

		_, err := hlpr.GetCurrentDirAsGitRepo()

		if err != nil {
			kingpin.Fatalf("This is not a Git repository")
		}

		mergedBranches, err := hlpr.GetMergedBranches(*previewOrigin, *previewMaster, *previewSkip)

		if err != nil {
			kingpin.Fatalf("Error when looking for branches: %s", err)
		}

		if len(mergedBranches) == 0 {
			fmt.Println("No remote branches are available for cleaning up")
		} else {
			fmt.Println("\nThese branches have been merged into master:")
			for _, branchName := range mergedBranches {
				fmt.Printf("  %s\n", branchName)
			}
			fmt.Println("\nTo delete them, run again with `gitsweeper cleanup`")
		}
	case cleanup.FullCommand():

		hlpr.SetupLogger(*debug)

		mergedBranches, err := hlpr.GetMergedBranches(*cleanupOrigin, *cleanupMaster, *cleanupSkip)

		if err != nil {
			kingpin.Fatalf("Error when looking for branches %s", err)
		}

		if len(mergedBranches) == 0 {
			fmt.Println("No remote branches are available for cleaning up")
		} else {
			fmt.Println("\nThese branches have been merged into master:")
			for _, branchName := range mergedBranches {
				fmt.Printf("  %s\n", branchName)
			}
			if !*(cleanupForce) {
				confirmDeleteBranches, err := hlpr.AskForConfirmation("Delete these branches?", os.Stdin)
				if err != nil {
					hlpr.FatalError("\nError when awaiting input", err)
				} else {
					if !confirmDeleteBranches {
						fmt.Printf("OK, aborting.\n")
						os.Exit(0)
					}
				}
			}
			fmt.Printf("\n")
			for _, branchName := range mergedBranches {
				remote, branchShort := hlpr.ParseBranchname(branchName)
				fmt.Printf("  deleting %s", branchName)
				repo, _ := hlpr.GetCurrentDirAsGitRepo()
				err := hlpr.DeleteBranch(repo, remote, branchShort)
				if err != nil {
					hlpr.FatalError("\nCould not delete branch", err)
				} else {
					fmt.Printf(" - (done)\n")
				}
			}
		}
	case version.FullCommand():
		hlpr.SetupLogger(*debug)
		fmt.Printf("%s %s\n", Version, gitCommit)
	default:
		os.Exit(0)
	}

}
