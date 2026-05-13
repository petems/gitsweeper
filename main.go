// Package main implements the gitsweeper CLI for cleaning up merged remote git branches.
//
// gitsweeper exposes two subcommands, preview and cleanup. Both detect remote branches
// that have already been merged into the configured main branch using a two-pass strategy:
// commit-hash reachability via go-git, followed by a git cherry / git patch-id check to
// catch squash merges and rebases. cleanup additionally deletes the matched branches by
// shelling out to git push --delete so that the user's existing git authentication is
// reused. See https://gitsweeper.readthedocs.io/ for full documentation.
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

// cliFlags holds the parsed command-line flag values.
type cliFlags struct {
	debug      bool
	origin     string
	master     string
	skip       string
	force      bool
	maxCommits int
	noDeep     bool
}

func main() {
	// Define command-line flags using standard library
	var (
		debug      = flag.Bool("debug", false, "Enable debug mode")
		version    = flag.Bool("version", false, "Show version")
		help       = flag.Bool("help", false, "Show help")
		origin     = flag.String("origin", "origin", "The name of the remote you wish to clean up")
		master     = flag.String("master", "master", "The name of what you consider the master branch")
		skip       = flag.String("skip", "", "Comma-separated list of branches to skip")
		force      = flag.Bool("force", false, "Do not ask, cleanup immediately")
		maxCommits = flag.Int("max-commits", 0, "Maximum number of commits to check (0 = default 10000)")
		noDeep     = flag.Bool("no-deep-check", false, "Disable git cherry squash-merge detection")
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

	// Merge command-position flags into the top-level values
	flags := &cliFlags{
		debug:      *debug,
		origin:     *origin,
		master:     *master,
		skip:       *skip,
		force:      *force,
		maxCommits: *maxCommits,
		noDeep:     *noDeep,
	}
	if flag.NArg() > 1 {
		mergeCommandFlags(flags, flag.Args()[1:])
	}

	// Setup lightweight logger
	hlpr.SetupLightLogger(flags.debug)

	opts := hlpr.MergeDetectionOptions{
		MaxCommits:    flags.maxCommits,
		DisableCherry: flags.noDeep,
	}

	switch command {
	case "preview":
		handlePreview(flags.origin, flags.master, flags.skip, opts)
	case "cleanup":
		handleCleanup(flags.origin, flags.master, flags.skip, flags.force, opts)
	case "version":
		fmt.Printf("%s %s\n", Version, gitCommit)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}
}

// mergeCommandFlags parses flags that appear after the command name and
// merges them into the already-parsed top-level flag values.
func mergeCommandFlags(flags *cliFlags, args []string) {
	cmdFlags := flag.NewFlagSet("", flag.ExitOnError)
	force := cmdFlags.Bool("force", false, "Do not ask, cleanup immediately")
	debug := cmdFlags.Bool("debug", false, "Enable debug mode")
	origin := cmdFlags.String("origin", "", "The name of the remote you wish to clean up")
	master := cmdFlags.String("master", "", "The name of what you consider the master branch")
	skip := cmdFlags.String("skip", "", "Comma-separated list of branches to skip")
	maxCommits := cmdFlags.Int("max-commits", 0, "Maximum number of commits to check (0 = default 10000)")
	noDeep := cmdFlags.Bool("no-deep-check", false, "Disable git cherry squash-merge detection")

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing command flags: %s\n", err)
		os.Exit(1)
	}

	if *force {
		flags.force = true
	}
	if *debug {
		flags.debug = true
	}
	if *origin != "" {
		flags.origin = *origin
	}
	if *master != "" {
		flags.master = *master
	}
	if *skip != "" {
		flags.skip = *skip
	}
	if *maxCommits != 0 {
		flags.maxCommits = *maxCommits
	}
	if *noDeep {
		flags.noDeep = true
	}
}

func handlePreview(origin, master, skipBranches string, opts hlpr.MergeDetectionOptions) {
	repo, err := hlpr.GetCurrentDirAsGitRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: This is not a Git repository\n")
		os.Exit(1)
	}

	mergedBranches, err := hlpr.GetMergedBranches(repo, origin, master, skipBranches, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when looking for branches: %s\n", err)
		os.Exit(1)
	}

	if len(mergedBranches) == 0 {
		fmt.Println("No remote branches are available for cleaning up")
	} else {
		fmt.Println("\nThese branches have been merged into master:")
		for _, result := range mergedBranches {
			fmt.Printf("  %s\n", result.Name)
		}
		fmt.Println("\nTo delete them, run again with `gitsweeper cleanup`")
	}
}

func handleCleanup(origin, master, skipBranches string, force bool, opts hlpr.MergeDetectionOptions) {
	repo, err := hlpr.GetCurrentDirAsGitRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: This is not a Git repository\n")
		os.Exit(1)
	}

	mergedBranches, err := hlpr.GetMergedBranches(repo, origin, master, skipBranches, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when looking for branches: %s\n", err)
		os.Exit(1)
	}

	if len(mergedBranches) == 0 {
		fmt.Println("No remote branches are available for cleaning up")
		return
	}

	fmt.Println("\nThese branches have been merged into master:")
	for _, result := range mergedBranches {
		fmt.Printf("  %s\n", result.Name)
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
	for i, result := range mergedBranches {
		remote, branchShort := hlpr.ParseBranchName(result.Name)
		if total > 10 {
			fmt.Printf("  [%d/%d] deleting %s", i+1, total, result.Name)
		} else {
			fmt.Printf("  deleting %s", result.Name)
		}

		err := hlpr.DeleteBranch(repo, remote, branchShort)
		if err != nil {
			fmt.Printf(" - (failed: %s)\n", err)
		} else {
			fmt.Printf(" - (done)\n")
		}
	}
}
