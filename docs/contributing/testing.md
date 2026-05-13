# Testing

`gitsweeper` has two test suites: Go unit tests and Cucumber acceptance tests. **Acceptance tests must be run before committing any change that touches CLI output, error messages, or argument parsing.**

## Unit tests

Go unit tests live alongside their source files (`internal/githelpers_test.go` etc.).

```bash
make test
```

This runs `go test ./...` across the whole module with race detection enabled.

### Coverage

```bash
make cover           # produces a coverage profile
make cover_html      # opens an HTML report in your browser
```

### Conventions

- Place tests in `_test.go` files next to the code they cover.
- Use table-driven tests for any logic with more than one branch.
- Run with `-race` (the Makefile does this for you) so concurrency bugs in `findMergedBranchesConcurrent` are caught.
- Keep `profile.out` and other coverage artifacts out of version control.

## Acceptance tests (Cucumber)

Acceptance tests live under `features/` with Ruby step definitions in `features/step_definitions/`. They build a Docker container with the compiled binary and exercise real CLI scenarios.

```bash
# Run the full suite
make cucumber

# Or target a single feature
bundle exec cucumber features/preview.feature
```

The suite verifies:

- CLI output format and exact wording.
- Error handling and error message content.
- Branch detection and cleanup behaviour against a real git fixture.
- User interaction flows (confirmation prompts).

### When you **must** run acceptance tests

You must run `make cucumber` before committing any change that affects:

- Error messages or error handling paths.
- The format or content of CLI output (including info logs that users might see).
- Command-line argument parsing.
- Confirmation prompts or other user interaction flows.

If a feature file's expectation no longer matches the new output, **update the feature file in the same commit**. Skipping this leaves acceptance tests broken on master.

### Why this matters

Unit tests verify code correctness. Acceptance tests verify the user-facing contract: "When I run `gitsweeper preview` against a repo with no merged branches, what exactly do I see?" Those expectations are documented in plain-language Gherkin steps, which means the suite doubles as living documentation.

## CI

GitHub Actions runs on every PR:

- `.github/workflows/golang.yml` — `make lint test`.
- `.github/workflows/aruba.yml` — `make cucumber`.
- `.github/workflows/codeql.yml` — CodeQL security analysis.
- `.github/workflows/install-script-validation.yml` — exercises `install.sh` on Linux, macOS, Windows.
- `.github/workflows/docs.yml` — `mkdocs build --strict` (catches broken nav/links before Read the Docs does).

All five jobs must pass before a PR can be merged.
