# Installation

`gitsweeper` ships pre-built binaries for Linux, macOS, and Windows. The fastest way to install is the one-line install script; Homebrew and manual download options are also available.

## Quick install (recommended)

The install script auto-detects your platform, downloads the appropriate binary, and verifies its checksum:

```bash
# Install latest version to ./bin
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh

# Install a specific version
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s v0.1.0

# Install to a system directory (may need sudo)
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s -- -b /usr/local/bin

# Force compilation from source (requires Go)
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s -- -f
```

The script will:

- Detect your platform (Linux, macOS, Windows)
- Download a pre-compiled binary when available
- Fall back to source compilation when needed
- Verify SHA256 checksums

## Homebrew (macOS)

Install directly from the formula in this repository:

```bash
brew install --formula Formula/gitsweeper.rb
```

Or via a tap (requires a separate `homebrew-gitsweeper` repository):

```bash
brew tap petems/gitsweeper
brew install gitsweeper
```

Force a source build (auto-detects platform, prefers binaries for speed):

```bash
brew install --build-from-source --formula Formula/gitsweeper.rb
```

See the [Homebrew formula reference](../reference/homebrew-formula.md) for details on how the formula picks between binaries and source.

## Pre-built binaries

Download the right archive for your platform from the [GitHub releases page](https://github.com/petems/gitsweeper/releases):

| Platform | Archive |
| --- | --- |
| Linux x86_64 | `gitsweeper-vX.Y.Z-linux-amd64.tar.gz` |
| Linux ARM64 | `gitsweeper-vX.Y.Z-linux-arm64.tar.gz` |
| macOS Intel | `gitsweeper-vX.Y.Z-darwin-amd64.tar.gz` |
| macOS Apple Silicon | `gitsweeper-vX.Y.Z-darwin-arm64.tar.gz` |
| Windows x86_64 | `gitsweeper-vX.Y.Z-windows-amd64.zip` |

### Extract and install

=== "Linux / macOS"

    ```bash
    tar -xzf gitsweeper-vX.Y.Z-your-platform.tar.gz
    sudo mv gitsweeper /usr/local/bin/
    chmod +x /usr/local/bin/gitsweeper
    ```

=== "Windows"

    ```powershell
    Expand-Archive gitsweeper-vX.Y.Z-windows-amd64.zip
    # Move gitsweeper.exe to a directory in your PATH
    ```

### Verify the checksum

All releases publish a `gitsweeper_X.Y.Z_checksums.txt` file alongside the archives:

```bash
wget https://github.com/petems/gitsweeper/releases/download/v0.1.0/gitsweeper-v0.1.0-linux-amd64.tar.gz
wget https://github.com/petems/gitsweeper/releases/download/v0.1.0/gitsweeper_0.1.0_checksums.txt
sha256sum -c gitsweeper_0.1.0_checksums.txt --ignore-missing
```

The install script does this automatically.

## Build from source

Requires Go 1.23 or later — see the [architecture notes](../explanation/architecture.md#go-version-requirement) for why.

```bash
# Easiest
go install github.com/petems/gitsweeper@latest

# Or clone and build
git clone https://github.com/petems/gitsweeper.git
cd gitsweeper
make build

# Or build a specific version
git checkout v0.1.0
go build -ldflags "-X main.gitCommit=$(git rev-parse --short HEAD)" .
```

## Troubleshooting

??? failure "Script reports `platform not supported`"
    Try forcing a source build: `curl ... | sh -s -- -f`. Make sure Go is installed (`go version`).

??? failure "Permission denied on default `./bin` path"
    Pick a writable location: `curl ... | sh -s -- -b ~/bin`. Make sure that directory is on your `PATH`.

??? failure "Brew formula not found"
    If you used `brew tap`, run `brew update` first. Otherwise install directly from the formula file: `brew install --formula Formula/gitsweeper.rb`.

??? failure "`gitsweeper: command not found` after install"
    Check `echo $PATH` and confirm the install directory is on it. Try the full path: `/usr/local/bin/gitsweeper --help`.

??? failure "`git repository not found` when running"
    `gitsweeper` runs against the current directory. Confirm you're inside a checkout with `git status` and that you have remotes configured with `git remote -v`.
