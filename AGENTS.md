# Repository Guidelines

## Go Version Requirement
**This project requires Go 1.23 or later.** The minimum version is enforced by security dependencies in `go-git` v5.14.0+ which require `golang.org/x/crypto@v0.35.0+` to mitigate CVE-2025-22869 (GO-2025-3487), a denial of service vulnerability in SSH implementations. All CI workflows and local development should use Go 1.23+.

## Project Structure & Module Organization
The CLI entrypoint lives in `main.go`, while reusable logic is grouped under `internal/` (e.g., `githelpers.go`, `prompthelpers.go`). Integration specs reside in `features/` with Ruby step definitions in `features/step_definitions/`. Generated binaries land in `bin/`, release artifacts in `dist/`, and the Homebrew formula sits in `Formula/`. Supporting scripts such as `install.sh` and `test-goreleaser.sh` configure installs and release checks.

## Git Operations & Authentication Strategy
The codebase uses [go-git](https://github.com/go-git/go-git) for read operations (repository analysis, branch detection, commit traversal) but **shells out to the system `git` command for branch deletion** (`git push --delete`). This hybrid approach is intentional: go-git excels at read operations, but authenticated push operations are extremely complex due to the wide variety of authentication methods in production environments (SSH keys with passphrases, SSH agents, credential helpers, tokens, deploy keys, etc.). Rather than reimplementing complex authentication logic, we leverage the user's existing Git configuration by shelling out for deletions. See https://github.com/go-git/go-git/issues/28 for context on go-git's authentication challenges.

## Build, Test, and Development Commands
Use `make build` to compile an optimized binary into `bin/gitsweeper`. `make test` runs `go test ./...` across the module, while `make cucumber` executes the Cucumber features via `bundle exec cucumber`. `make lint` verifies `golangci-lint` (>=2.3.0) and lints the codebase, and `make cover` produces race-enabled coverage profiles. For multi-platform snapshots, prefer `make build-all` (goreleaser snapshot build).

## Coding Style & Naming Conventions
All Go code must be `gofmt`ed; the default tab indentation and `MixedCaps` identifiers follow standard Go style. Keep exported symbols minimal—prefer `internal/` packages for helpers. Run `golangci-lint run` before sending a PR; align with existing helper naming patterns such as `loghelpers` and `slicehelpers`. Shell scripts should pass `shellcheck`; mirror the defensive checks already present in `install.sh`.

## Testing Guidelines
Place unit tests alongside source files using the `_test.go` suffix (`internal/githelpers_test.go` as reference). Use table-driven tests for new logic and ensure race detection succeeds (`go test -race ./...`).

### Acceptance Tests (Critical for CLI Changes)
**IMPORTANT: Always run acceptance tests (`make cucumber`) before committing any changes that affect:**
- Error messages or error handling
- CLI output format or content
- Command-line argument parsing
- User interaction flows

Cucumber scenarios live in `features/*.feature`; name steps in plain language mirroring CLI behavior. The acceptance tests build a Docker container with the compiled binary and test real-world scenarios. When modifying error messages or output, you MUST update the corresponding feature file expectations. Run them with `make cucumber` or `bundle exec cucumber features/specific.feature` for targeted testing. Keep coverage reports (`profile.out`) out of version control.

## Commit & Pull Request Guidelines
Follow the dominant Conventional Commit style: prefixes like `feat:` and `fix:` appear throughout (`fix: add explicit permissions...`, `feat: implement Windows-specific zip builds`). Include context-rich messages and reference pull requests or issues in parentheses when relevant. Before opening a PR, run `make lint test cucumber` and note the results. Provide a concise summary, expected CLI output (e.g., `gitsweeper preview`), and link related issues. Screenshots or logs are encouraged when behavior changes.
