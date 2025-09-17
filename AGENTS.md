# Repository Guidelines

## Project Structure & Module Organization
The CLI entrypoint lives in `main.go`, while reusable logic is grouped under `internal/` (e.g., `githelpers.go`, `prompthelpers.go`). Integration specs reside in `features/` with Ruby step definitions in `features/step_definitions/`. Generated binaries land in `bin/`, release artifacts in `dist/`, and the Homebrew formula sits in `Formula/`. Supporting scripts such as `install.sh` and `test-goreleaser.sh` configure installs and release checks.

## Build, Test, and Development Commands
Use `make build` to compile an optimized binary into `bin/gitsweeper`. `make test` runs `go test ./...` across the module, while `make cucumber` executes the Cucumber features via `bundle exec cucumber`. `make lint` verifies `golangci-lint` (>=2.3.0) and lints the codebase, and `make cover` produces race-enabled coverage profiles. For multi-platform snapshots, prefer `make build-all` (goreleaser snapshot build).

## Coding Style & Naming Conventions
All Go code must be `gofmt`ed; the default tab indentation and `MixedCaps` identifiers follow standard Go style. Keep exported symbols minimal—prefer `internal/` packages for helpers. Run `golangci-lint run` before sending a PR; align with existing helper naming patterns such as `loghelpers` and `slicehelpers`. Shell scripts should pass `shellcheck`; mirror the defensive checks already present in `install.sh`.

## Testing Guidelines
Place unit tests alongside source files using the `_test.go` suffix (`internal/githelpers_test.go` as reference). Use table-driven tests for new logic and ensure race detection succeeds (`go test -race ./...`). Cucumber scenarios live in `features/*.feature`; name steps in plain language mirroring CLI behavior and run them with `bundle exec cucumber`. Keep coverage reports (`profile.out`) out of version control.

## Commit & Pull Request Guidelines
Follow the dominant Conventional Commit style: prefixes like `feat:` and `fix:` appear throughout (`fix: add explicit permissions...`, `feat: implement Windows-specific zip builds`). Include context-rich messages and reference pull requests or issues in parentheses when relevant. Before opening a PR, run `make lint test cucumber` and note the results. Provide a concise summary, expected CLI output (e.g., `gitsweeper preview`), and link related issues. Screenshots or logs are encouraged when behavior changes.
