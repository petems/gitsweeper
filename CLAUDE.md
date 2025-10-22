# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gitsweeper` is a CLI tool for cleaning up merged git branches, written in Go. It's a modern rewrite of the Python-based `git-sweep` tool, designed to identify and delete remote branches that have been merged into the main branch.

**Go Version Requirement: This project requires Go 1.23 or later.** This is due to security dependencies in `go-git` v5.14.0+ which require `golang.org/x/crypto@v0.35.0+` to mitigate CVE-2025-22869 (GO-2025-3487), a denial of service vulnerability in SSH implementations. Users unable to upgrade to Go 1.23 should use a version of gitsweeper built with go-git v5.13.x or earlier.

## Architecture

- **`main.go`**: Entry point with command-line argument parsing using Go's standard `flag` package. Supports `preview` and `cleanup` commands with flags for debug, origin, master branch name, skip patterns, and force mode.
- **`internal/`**: Core functionality split into focused helper modules:
  - `githelpers.go`: Git operations using `go-git` library for repository operations and branch detection; shells out to git for branch deletion
  - `prompthelpers.go`: User interaction utilities for confirmation prompts
  - `loghelpers.go`: Lightweight logging setup
  - `slicehelpers.go`: Utility functions for slice operations

The application uses the `go-git` library for most Git operations (reading, analysis, branch detection) rather than shelling out to Git commands, making it more portable and reliable.

### Authentication Handling

**Branch deletion uses shell commands (`git push --delete`) instead of go-git's push operations.** While go-git is excellent for read operations, it has significant complexity and limitations when dealing with authenticated push operations. There's a huge variety of authentication methods in the wild (SSH keys with passphrases, SSH agents, various credential helpers, tokens, deploy keys, etc.), and trying to handle them all through go-git's authentication API is overly complex and error-prone.

By shelling out to the system's `git` command for deletion, we leverage the user's existing Git configuration and authentication setup automatically. The system git already knows how to work with SSH agents, credential helpers, and other authentication mechanisms configured by the user.

See the go-git project's long-standing authentication complexity issues: https://github.com/go-git/go-git/issues/28

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

The project includes both unit tests (Go) and acceptance tests (Cucumber/Ruby). The acceptance tests require a Ruby environment with Bundle and use Docker for integration testing.

**IMPORTANT: Always run acceptance tests before committing changes that affect CLI output or error messages.**

Run unit tests: `make test`
Run acceptance tests: `make cucumber`
Generate coverage report: `make cover` then `make cover_html`

### Running Acceptance Tests

Acceptance tests use Cucumber with Ruby step definitions. They build a Docker container with the compiled binary and test real-world CLI scenarios:

```bash
# Run all acceptance tests
make cucumber

# Run specific feature
bundle exec cucumber features/preview.feature
```

The acceptance tests verify:
- CLI output format and messaging
- Error handling and error messages
- Branch detection and cleanup behavior
- User interaction flows

**When modifying error messages or CLI output, you MUST update the corresponding Cucumber feature expectations and run the tests to verify.**

## Configuration Files

- **`.golangci.yml`**: Comprehensive linting configuration with 40+ enabled linters
- **`.goreleaser.yml`**: Multi-platform build and release automation
- **`Makefile`**: Development workflow automation
- **`go.mod`**: Go 1.23+ with `go-git` as primary dependency