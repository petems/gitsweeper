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

### Pre-built Binaries (Recommended)

Download the latest release for your platform from the [GitHub releases page](https://github.com/petems/gitsweeper/releases):

- **Linux (x86_64)**: `gitsweeper-vX.Y.Z-linux-amd64.tar.gz`
- **Linux (ARM64)**: `gitsweeper-vX.Y.Z-linux-arm64.tar.gz`
- **macOS (Intel)**: `gitsweeper-vX.Y.Z-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `gitsweeper-vX.Y.Z-darwin-arm64.tar.gz`
- **Windows (x86_64)**: `gitsweeper-vX.Y.Z-windows-amd64.zip`

#### Installation Steps:

1. Download the appropriate archive for your platform
2. Extract the archive:
   ```bash
   # For Linux/macOS
   tar -xzf gitsweeper-vX.Y.Z-your-platform.tar.gz
   
   # For Windows
   unzip gitsweeper-vX.Y.Z-windows-amd64.zip
   ```
3. Move the binary to a directory in your PATH:
   ```bash
   # Linux/macOS
   sudo mv gitsweeper /usr/local/bin/
   chmod +x /usr/local/bin/gitsweeper
   
   # Windows: Move gitsweeper.exe to a directory in your PATH
   ```

### Build from Source

```bash
go install github.com/petems/gitsweeper@latest
```

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
