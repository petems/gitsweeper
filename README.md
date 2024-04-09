# `gitsweeper`

A CLI tool for cleaning up git repositories.

## Usage

### List branches merged into master

```bash
$ gitsweeper preview
Fetching from the remote...

These branches have been merged into master:
  origin/merged_already_to_master

To delete them, run again with `gitsweeper cleanup`
```

### Cleanup branches merged into master

```bash
$ gitsweeper cleanup
Fetching from the remote...

These branches have been merged into master:
  origin/merged_already_to_master
```

## Installation

```bash
go install github.com/petems/gitsweeper@latest
```

Eventually I'll configure Travis to build binaries and setup a `brew tap` for OSX and Linux.

## Background

`gitsweeper` is a tribute to a tool I've been using for a long time, [git-sweep](b.com/arc90/git-sweep). git-sweep is a great tool written in Python.

However, since then it seems to have been abandoned. It's not had a commit pushed [since 2016](https://github.com/arc90/git-sweep/commit/d7522b4de1dbc85570ec36b82bc155a4fa371b5e), seems to be [broken with Python 3](https://github.com/arc90/git-sweep/issues/44).

I've been trying to learn more Go recently, and Go has some excellent CLI library tools as well as the ability to build a self-contained binary for distribution, rather than having to make sure it works with various versions of go etc.

`gitsweeper` matches the output matches the original tool quite a lot:

```
$ git-sweep preview
Fetching from the remote
These branches have been merged into master:

  merged_already_to_master

To delete them, run again with `git-sweep cleanup`
```

but has a few changes that are tweaked toward my requirements.
