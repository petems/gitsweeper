# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gitsweeper` is a CLI tool for cleaning up merged git branches, written in Go. It's a modern rewrite of the Python-based `git-sweep` tool, designed to identify and delete remote branches that have been merged into the main branch.

## Architecture

- **`main.go`**: Entry point with command-line argument parsing using Go's standard `flag` package. Supports `preview` and `cleanup` commands with flags for debug, origin, master branch name, skip patterns, and force mode.
- **`internal/`**: Core functionality split into focused helper modules:
  - `githelpers.go`: Git operations using `go-git` library for repository operations, branch detection, and deletion
  - `prompthelpers.go`: User interaction utilities for confirmation prompts
  - `loghelpers.go`: Lightweight logging setup
  - `slicehelpers.go`: Utility functions for slice operations

The application uses the `go-git` library for Git operations rather than shelling out to Git commands, making it more portable and reliable.

## Development Commands

### Build and Test
```bash
# Full development cycle (clean, build, format, lint, test, install)
make

# Build optimized binary
make build

# Run tests
make test

# Run cucumber integration tests (requires Ruby/Bundle setup)
make cucumber

# Run tests with coverage
make cover
```

### Code Quality
```bash
# Lint code (requires golangci-lint v2.3.0+)
make lint

# Auto-fix linting issues
make lint-fix

# Format code
make fmt
```

### Release and Distribution
```bash
# Build for all platforms using goreleaser
make build-all

# Create snapshot release (no publishing)
make release-snapshot

# Create full release (requires proper Git tags)
make release
```

### Development Dependencies
Install development tools via Homebrew:
```bash
brew bundle
```

## Key Implementation Details

- Uses ultra-optimized build flags (`-s -w -trimpath`) for minimal binary size (7.8MB)
- Git operations are performed through the `go-git` library for better cross-platform compatibility
- Branch detection logic identifies remote branches merged into the specified master branch
- Support for skipping branches via comma-separated patterns
- Progress indication for large branch deletion operations
- Confirmation prompts can be bypassed with `--force` flag

## Testing

The project includes both unit tests (Go) and integration tests (Cucumber/Ruby). The integration tests require a Ruby environment with Bundle.

Run unit tests: `make test`  
Run integration tests: `make cucumber`  
Generate coverage report: `make cover` then `make cover_html`

## Configuration Files

- **`.golangci.yml`**: Comprehensive linting configuration with 40+ enabled linters
- **`.goreleaser.yml`**: Multi-platform build and release automation
- **`Makefile`**: Development workflow automation
- **`go.mod`**: Go 1.23+ with `go-git` as primary dependency