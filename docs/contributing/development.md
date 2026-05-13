# Development

This page covers setting up a `gitsweeper` development environment and running the day-to-day commands you'll need.

## Prerequisites

- **Go 1.23+** — required by `go-git` security dependencies. See [architecture notes](../explanation/architecture.md#go-version-requirement) for the CVE background.
- **Git**
- **Make**
- **golangci-lint** (>= 2.3.0) — for linting.
- **Ruby + Bundler** — for the Cucumber acceptance tests.
- **Docker** — the acceptance tests build a container with the compiled binary.

## Quick setup (macOS)

All development dependencies are listed in the repo's `Brewfile`:

```bash
git clone https://github.com/petems/gitsweeper.git
cd gitsweeper
brew bundle
make
```

`make` (the default target) does the full development cycle: clean, build, format, lint, test, install.

## Manual setup (other platforms)

```bash
# 1. Install Go 1.23+ from https://golang.org/dl/

# 2. Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b $(go env GOPATH)/bin

# 3. Install Ruby and Bundler (e.g. via your package manager or rbenv)
gem install bundler
bundle install

# 4. Install Docker (https://docs.docker.com/get-docker/)

# 5. Clone and build
git clone https://github.com/petems/gitsweeper.git
cd gitsweeper
make build
```

## Daily commands

```bash
# Build optimized binary into bin/gitsweeper
make build

# Run unit tests
make test

# Generate a coverage profile
make cover
make cover_html        # open the HTML report

# Lint
make lint
make lint-fix          # auto-fix what can be auto-fixed

# Format
make fmt

# Run acceptance tests (builds a Docker image)
make cucumber

# Build for all platforms with goreleaser
make build-all

# Snapshot release (no publish)
make release-snapshot
```

## Project layout

```text
gitsweeper/
├── main.go                     CLI entry point + flag parsing
├── internal/                   All non-CLI helpers (intentionally not importable)
│   ├── githelpers.go           go-git operations, hash-match pass
│   ├── cherrycheck.go          git cherry / patch-id shell-outs
│   ├── prompthelpers.go        Confirmation prompts
│   ├── loghelpers.go           Lightweight logger
│   └── slicehelpers.go         Slice/set utilities
├── features/                   Cucumber specs
│   └── step_definitions/       Ruby step definitions
├── Formula/                    Homebrew formula + placeholder vars
├── docs/                       This site
├── .github/workflows/          CI: lint, test, release, codeql, docs build
├── .goreleaser.yml             Multi-platform build + release automation
├── .golangci.yml               40+ enabled linters
├── Makefile                    Development workflow automation
└── go.mod                      Go 1.23+ with go-git as primary dependency
```

See [Architecture](../explanation/architecture.md) for the design rationale behind the `go-git` + shell-out hybrid.

## Coding style

- Code must be `gofmt`ed. Use tabs for indentation and `MixedCaps` for identifiers.
- Keep exported symbols minimal — prefer `internal/` packages for helpers. This is a CLI binary, not a reusable library.
- Run `golangci-lint run` before sending a PR. Use `make lint`.
- Shell scripts should pass `shellcheck`; mirror the defensive checks already in `install.sh`.

## Commit and PR conventions

- Follow [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `chore:`, `docs:`, etc. Use `*` for bullet points in commit bodies.
- Run `make lint test cucumber` before opening a PR.
- Include a concise summary, expected CLI output if behaviour changed, and a link to any related issues.

## Configuration variables

| Variable | Used by | Purpose |
| --- | --- | --- |
| `BINDIR` | `install.sh` | Override install directory (default `./bin`) |
| `FORCE_SOURCE` | `install.sh` | Force source compilation when set |
| `GIT_TERMINAL_PROMPT` | runtime | Set to `0` automatically when shelling out for `git push --delete` |
