# `gitsweeper`

A CLI tool for cleaning up git repositories.

**Note**: This version uses ultra-optimized code for improved performance and reduced binary size (7.8MB). The older unoptimized implementations have been removed.

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

### Quick Install (Recommended)

Use our install script for the easiest installation:

```bash
# Install latest version
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh

# Install specific version
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s v0.1.0

# Install to custom location
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s -- -b /usr/local/bin

# Force build from source (requires Go)
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s -- -f
```

### Homebrew (macOS)

```bash
# Install directly from the formula in this repository
brew install --formula Formula/gitsweeper.rb

# Or if you prefer using a tap (requires separate homebrew-gitsweeper repository)
brew tap petems/gitsweeper
brew install gitsweeper

# Force compilation from source (auto-detects platform, prefers binaries for speed)
brew install --build-from-source --formula Formula/gitsweeper.rb
```

> See [Formula/README.md](Formula/README.md) for more details about the Homebrew formula.

### Pre-built Binaries

Download the latest release for your platform from the [GitHub releases page](https://github.com/petems/gitsweeper/releases):

- **Linux (x86_64)**: `gitsweeper-vX.Y.Z-linux-amd64.tar.gz`
- **Linux (ARM64)**: `gitsweeper-vX.Y.Z-linux-arm64.tar.gz`
- **macOS (Intel)**: `gitsweeper-vX.Y.Z-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `gitsweeper-vX.Y.Z-darwin-arm64.tar.gz`
- **Windows (x86_64)**: `gitsweeper-vX.Y.Z-windows-amd64.zip`

#### Manual Installation Steps:

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

### Development Environment

If you're contributing to gitsweeper, you can set up a development environment using Homebrew:

```bash
# Install all development dependencies
brew bundle

# Build and test
make
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
