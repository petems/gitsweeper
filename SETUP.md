# Setup Guide for Gitsweeper

This guide covers various ways to install and set up gitsweeper for different use cases.

## Quick Installation

### Option 1: Install Script (Recommended)

Our install script provides the most flexible installation method:

```bash
# Install latest version to ./bin
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh

# Install specific version
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s v0.1.0

# Install to system location (requires sudo)
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s -- -b /usr/local/bin

# Force compilation from source (requires Go toolchain)
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh -s -- -f
```

The install script will:
- Detect your platform automatically (Linux, macOS, Windows)
- Download pre-compiled binaries when available
- Fall back to source compilation when needed
- Verify checksums for security

### Option 2: Homebrew (macOS)

If you're on macOS and use Homebrew:

```bash
# Add our tap
brew tap petems/gitsweeper

# Install gitsweeper (uses binary when possible)
brew install gitsweeper

# Force compilation from source
brew install --build-from-source gitsweeper

# Or install with source option
brew install gitsweeper --with-source
```

The Homebrew formula will:
- Use pre-compiled binaries for Intel Macs when available
- Automatically compile from source on Apple Silicon
- Allow forcing source compilation via options

## Manual Installation

### Pre-compiled Binaries

1. Visit the [releases page](https://github.com/petems/gitsweeper/releases)
2. Download the appropriate archive for your platform:
   - **Linux x86_64**: `gitsweeper-vX.Y.Z-linux-amd64.tar.gz`
   - **Linux ARM64**: `gitsweeper-vX.Y.Z-linux-arm64.tar.gz`
   - **macOS Intel**: `gitsweeper-vX.Y.Z-darwin-amd64.tar.gz`
   - **macOS Apple Silicon**: `gitsweeper-vX.Y.Z-darwin-arm64.tar.gz`
   - **Windows x86_64**: `gitsweeper-vX.Y.Z-windows-amd64.zip`

3. Extract and install:
   ```bash
   # Linux/macOS
   tar -xzf gitsweeper-vX.Y.Z-your-platform.tar.gz
   sudo mv gitsweeper /usr/local/bin/
   chmod +x /usr/local/bin/gitsweeper
   
   # Windows
   unzip gitsweeper-vX.Y.Z-windows-amd64.zip
   # Move gitsweeper.exe to a directory in your PATH
   ```

### Build from Source

#### Prerequisites
- Go 1.21 or later
- Git

#### Steps
```bash
# Option 1: Using go install
go install github.com/petems/gitsweeper@latest

# Option 2: Clone and build
git clone https://github.com/petems/gitsweeper.git
cd gitsweeper
make build

# Option 3: Build specific version
git clone https://github.com/petems/gitsweeper.git
cd gitsweeper
git checkout v0.1.0
go build -ldflags "-X main.gitCommit=$(git rev-parse --short HEAD)" .
```

## Development Setup

### Quick Development Environment

Use our Brewfile to set up all development dependencies:

```bash
# Clone the repository
git clone https://github.com/petems/gitsweeper.git
cd gitsweeper

# Install all development dependencies (macOS)
brew bundle

# Build and test
make
```

### Manual Development Setup

#### Required Tools
- Go 1.21+
- Git
- Make
- golangci-lint (for linting)
- Ruby (for integration tests)
- Docker (for integration tests)

#### Installation
```bash
# Install Go
# Visit https://golang.org/dl/ or use your package manager

# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Install Ruby (for Cucumber tests)
# Visit https://www.ruby-lang.org/en/downloads/ or use your package manager

# Install bundler and gems
gem install bundler
bundle install

# Install Docker
# Visit https://docs.docker.com/get-docker/
```

#### Development Commands
```bash
# Build
make build

# Run tests
make test

# Run linting
make lint

# Run integration tests
make cucumber  # or bundle exec cucumber

# Build for all platforms
make build-all

# Create release archives
make release-archives

# Clean build artifacts
make clean
```

## Configuration

### Environment Variables

- `BINDIR`: Installation directory for the install script (default: `./bin`)
- `FORCE_SOURCE`: Force source compilation (default: `false`)

### Build Flags

When building from source, you can customize the build:

```bash
# Build with version information
go build -ldflags "-X main.gitCommit=$(git rev-parse --short HEAD)" .

# Build with additional flags
go build -ldflags "-X main.gitCommit=$(git rev-parse --short HEAD) -s -w" .
```

## Troubleshooting

### Install Script Issues

**Script fails with "platform not supported"**
- Try force building from source: `curl ... | sh -s -- -f`
- Check if Go is installed: `go version`

**Permission denied errors**
- Use a custom install directory: `curl ... | sh -s -- -b ~/bin`
- Make sure the target directory is in your PATH

### Homebrew Issues

**Formula not found**
- Make sure you've added the tap: `brew tap petems/gitsweeper`
- Update Homebrew: `brew update`

**Build fails on Apple Silicon**
- The formula should automatically build from source on ARM64
- Try explicitly: `brew install --build-from-source gitsweeper`

### General Issues

**Binary not found after installation**
- Check if the installation directory is in your PATH: `echo $PATH`
- Verify the binary exists: `ls -la /usr/local/bin/gitsweeper`
- Try running with full path: `/usr/local/bin/gitsweeper version`

**Git repository errors**
- Make sure you're in a Git repository: `git status`
- Check if you have proper Git remotes: `git remote -v`

## Security

### Checksum Verification

All releases include SHA256 checksums. When downloading manually:

```bash
# Download the archive and checksum
wget https://github.com/petems/gitsweeper/releases/download/v0.1.0/gitsweeper-v0.1.0-linux-amd64.tar.gz
wget https://github.com/petems/gitsweeper/releases/download/v0.1.0/gitsweeper_0.1.0_checksums.txt

# Verify the checksum
sha256sum -c gitsweeper_0.1.0_checksums.txt --ignore-missing
```

The install script automatically verifies checksums when downloading binaries.

## Contributing

See the main README.md for contribution guidelines. The development setup above should get you started with making changes to gitsweeper.