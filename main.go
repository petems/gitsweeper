package main

import (
	"flag"
	"fmt"
	"os"

	hlpr "github.com/petems/gitsweeper/internal"
)

// Version is what is returned by the `-v` flag.
const Version = "0.1.0"

// gitCommit is the gitcommit its built from.
var gitCommit = "development"

func main() {
	// Define command-line flags using standard library
	var (
		debug   = flag.Bool("debug", false, "Enable debug mode")
		version = flag.Bool("version", false, "Show version")
		help    = flag.Bool("help", false, "Show help")
		origin  = flag.String("origin", "origin", "The name of the remote you wish to clean up")
		master  = flag.String("master", "master", "The name of what you consider the master branch")
		skip    = flag.String("skip", "", "Comma-separated list of branches to skip")
		force   = flag.Bool("force", false, "Do not ask, cleanup immediately")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: gitsweeper [<flags>] <command> [<args> ...]\n\n")
		fmt.Fprintf(os.Stderr, "A command-line tool for cleaning up merged branches.\n")
	}

	// Parse flags before the command
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("%s %s\n", Version, gitCommit)
		return
	}

	// Handle help flag
	if *help || flag.NArg() == 0 {
		flag.Usage()
		return
	}

	command := flag.Arg(0)

	// Parse remaining flags after the command
	if flag.NArg() > 1 {
		// Create a new flag set for the remaining arguments
		cmdFlags := flag.NewFlagSet("", flag.ExitOnError)
		cmdFlags.Bool("force", false, "Do not ask, cleanup immediately")
		cmdFlags.Bool("debug", false, "Enable debug mode")
		cmdFlags.String("origin", "origin", "The name of the remote you wish to clean up")
		cmdFlags.String("master", "master", "The name of what you consider the master branch")
		cmdFlags.String("skip", "", "Comma-separated list of branches to skip")

		// Parse the remaining arguments
		if err := cmdFlags.Parse(flag.Args()[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing command flags: %s\n", err)
			os.Exit(1)
		}

		// Update the original flags with command-specific flags
		if cmdFlags.Lookup("force") != nil && cmdFlags.Lookup("force").Value.String() == "true" {
			*force = true
		}
		if cmdFlags.Lookup("debug") != nil && cmdFlags.Lookup("debug").Value.String() == "true" {
			*debug = true
		}
		if cmdFlags.Lookup("origin") != nil && cmdFlags.Lookup("origin").Value.String() != "" {
			*origin = cmdFlags.Lookup("origin").Value.String()
		}
		if cmdFlags.Lookup("master") != nil && cmdFlags.Lookup("master").Value.String() != "" {
			*master = cmdFlags.Lookup("master").Value.String()
		}
		if cmdFlags.Lookup("skip") != nil && cmdFlags.Lookup("skip").Value.String() != "" {
			*skip = cmdFlags.Lookup("skip").Value.String()
		}
	}

	// Setup lightweight logger
	hlpr.SetupLightLogger(*debug)

	switch command {
	case "preview":
		handlePreview(*origin, *master, *skip)
	case "cleanup":
		handleCleanup(*origin, *master, *skip, *force)
	case "version":
		fmt.Printf("%s %s\n", Version, gitCommit)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}
}

func handlePreview(origin, master, skipBranches string) {
	repo, err := hlpr.GetCurrentDirAsGitRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: This is not a Git repository\n")
		os.Exit(1)
	}

	mergedBranches, err := hlpr.GetMergedBranches(repo, origin, master, skipBranches)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when looking for branches: %s\n", err)
		os.Exit(1)
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
}

func handleCleanup(origin, master, skipBranches string, force bool) {
	repo, err := hlpr.GetCurrentDirAsGitRepo()
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"gitsweeper-int-test: error: Error when looking for branches repository does not exist\n",
		)
		os.Exit(1)
	}

	mergedBranches, err := hlpr.GetMergedBranches(repo, origin, master, skipBranches)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gitsweeper-int-test: error: Error when looking for branches %s\n", err)
		os.Exit(1)
	}

	if len(mergedBranches) == 0 {
		fmt.Println("No remote branches are available for cleaning up")
		return
	}

	fmt.Println("\nThese branches have been merged into master:")
	for _, branchName := range mergedBranches {
		fmt.Printf("  %s\n", branchName)
	}

	if !force {
		confirmDeleteBranches, confirmErr := hlpr.AskForConfirmation("Delete these branches?", os.Stdin)
		if confirmErr != nil {
			hlpr.LogFatalError("\nError when awaiting input", confirmErr)
		}
		if !confirmDeleteBranches {
			fmt.Printf("OK, aborting.\n")
			return
		}
	}

	fmt.Printf("\n")

	// Process deletions with progress indication for large sets
	total := len(mergedBranches)
	for i, branchName := range mergedBranches {
		remote, branchShort := hlpr.ParseBranchName(branchName)
		if total > 10 {
			fmt.Printf("  [%d/%d] deleting %s", i+1, total, branchName)
		} else {
			fmt.Printf("  deleting %s", branchName)
		}

		err := hlpr.DeleteBranch(repo, remote, branchShort)
		if err != nil {
			fmt.Printf(" - (failed: %s)\n", err)
		} else {
			fmt.Printf(" - (done)\n")
		}
	}
}
