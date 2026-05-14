# Motivation

## A tribute to `git-sweep`

`gitsweeper` is a tribute to a tool that's been part of my workflow for years: [`git-sweep`](https://github.com/arc90/git-sweep) by arc90. `git-sweep` is a great little Python utility — point it at a repo, and it tells you which remote branches have already been merged so you can clean them up.

Unfortunately, `git-sweep` has been effectively abandoned:

- The last commit was [pushed in 2016](https://github.com/arc90/git-sweep/commit/d7522b4de1dbc85570ec36b82bc155a4fa371b5e).
- It's [broken on Python 3](https://github.com/arc90/git-sweep/issues/44).
- The Python packaging ecosystem has moved on around it.

## Why a Go rewrite

I've been learning Go and wanted a project I'd actually use. Go fits this problem well:

- **Self-contained binaries.** A 7.8 MB binary works on any Linux/macOS/Windows box without worrying about Python versions, virtualenvs, or system packages.
- **First-class CLI ergonomics.** The `flag` package in the standard library is enough for a small utility like this.
- **A mature git library.** [`go-git`](https://github.com/go-git/go-git) handles read operations (branch listing, commit traversal) entirely in-process — no shelling out for the read path.

The output is intentionally close to the original `git-sweep` so muscle memory transfers:

```text
$ git-sweep preview
Fetching from the remote
These branches have been merged into master:

  merged_already_to_master

To delete them, run again with `git-sweep cleanup`
```

```text
$ gitsweeper preview
Fetching from the remote...

These branches have been merged into master:
  origin/merged_already_to_master

To delete them, run again with `gitsweeper cleanup`
```

## What's new vs. `git-sweep`

A few changes were made along the way, all motivated by real-world pain:

- **Squash-merge detection.** The biggest one. `git-sweep` (and `gitsweeper`'s earliest versions) only detected branches whose tip commit appeared in the main branch — that misses every branch merged via GitHub's *Squash and merge* button. See [Squash-Merge Detection](../usage/squash-merge-detection.md).
- **Forge-independent.** `gitsweeper` uses only vanilla git commands. No GitHub/GitLab API tokens needed, no rate limits to worry about, works on self-hosted git just as well as GitHub.
- **Cross-platform binaries.** Pre-built binaries for Linux (amd64/arm64), macOS (Intel/Apple Silicon), and Windows.
- **Configurable history depth.** `--max-commits` lets you tune how far back to scan; the default of 10 000 covers most projects.
- **Skip lists.** Comma-separated `--skip` excludes specific branches you don't want considered.

## What this tool is not

- **A replacement for `git branch -D`.** `gitsweeper` only deletes *remote* branches via `git push --delete`. Your local branches are untouched.
- **A history rewriter.** `gitsweeper` never touches the main branch or rewrites commits. It only deletes refs.
- **A reviewer.** `gitsweeper` tells you which branches *look* merged. It cannot tell you whether your colleague had uncommitted experimental work on a branch they squash-merged. Use `--force` carefully.

If you've used `git-sweep` for years and just want the same workflow with squash-merge support and a single binary — that's what `gitsweeper` is.
