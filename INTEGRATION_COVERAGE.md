# Integration Test Coverage

This document explains how to use the new integration test coverage functionality introduced in Go 1.20+. This feature allows you to collect code coverage from running actual compiled binaries, providing a more realistic view of code coverage than traditional unit tests.

## Overview

Integration test coverage works by:
1. Building a coverage-instrumented binary using `go build -cover`
2. Running acceptance tests against this instrumented binary
3. Collecting coverage data files during execution
4. Processing the coverage data to generate reports

This approach provides coverage metrics for the actual production code paths, including main functions and command-line interfaces that are typically excluded from unit test coverage.

## Requirements

- Go 1.20 or later
- The `go tool covdata` command (included with Go 1.20+)

## Available Commands

### Basic Integration Coverage

```bash
make integration-cover
```

This command:
- Builds a coverage-instrumented binary
- Runs all acceptance tests against the binary
- Generates a summary coverage report

Example output:
```
Building coverage-instrumented binary...
Running integration tests with coverage...
=== RUN   TestVersionCommand
=== RUN   TestVersionCommand/Version_with_no_flags
=== RUN   TestVersionCommand/Version_with_--debug_flag
--- PASS: TestVersionCommand (0.01s)
    --- PASS: TestVersionCommand/Version_with_no_flags (0.00s)
    --- PASS: TestVersionCommand/Version_with_--debug_flag (0.00s)
...
Processing coverage data...
        command-line-arguments          coverage: 34.9% of statements
        github.com/petems/gitsweeper/internal           coverage: 30.8% of statements
```

### HTML Coverage Report

```bash
make integration-cover-html
```

Generates an HTML coverage report that you can open in a browser:
- Creates `integration-coverage.html`
- Shows line-by-line coverage with color coding
- Includes clickable file navigation

### Function-Level Coverage

```bash
make integration-cover-func
```

Shows coverage statistics for individual functions:

```
/workspace/main.go:37:                                          main            34.9%
github.com/petems/gitsweeper/internal/githelpers.go:21:         RemoteBranches  0.0%
github.com/petems/gitsweeper/internal/githelpers.go:32:         ParseBranchname 0.0%
github.com/petems/gitsweeper/internal/githelpers.go:39:         DeleteBranch    0.0%
...
total:                                                          (statements)    31.9%
```

### Coverage Data Merging

```bash
make integration-cover-merge
```

Merges multiple coverage data files into a single dataset:
- Useful for combining results from multiple test runs
- Reduces file count and size
- Enables analysis of coverage across different scenarios

## How It Works

### Environment Variables

The system uses two key environment variables:

1. **GOCOVERDIR**: Specifies where coverage data files are written
2. **GITSWEEPER_TEST_BINARY**: Points to the coverage-instrumented binary

### Test Helper Integration

The `TestHelper` in `acceptance_test.go` automatically:
- Detects when a coverage-instrumented binary is provided
- Passes through the `GOCOVERDIR` environment variable to child processes
- Uses absolute paths to ensure proper execution

### Coverage Data Files

Coverage data is stored in the `covdatafiles/` directory:
- `covmeta.*`: Metadata about instrumented packages
- `covcounters.*`: Execution counts for each instrumented statement
- Multiple files are created for different test executions

## Comparison with Traditional Coverage

| Feature | Unit Test Coverage | Integration Coverage |
|---------|-------------------|---------------------|
| **Scope** | Individual packages | Entire application |
| **Main Function** | ❌ Not covered | ✅ Covered |
| **CLI Interface** | ❌ Not covered | ✅ Covered |
| **Real Execution** | ❌ Mock/stub heavy | ✅ Actual binary |
| **Performance** | ⚡ Fast | 🐌 Slower |
| **Accuracy** | 📊 Package-level | 🎯 End-to-end |

## Coverage Results Analysis

### Current Coverage (Example)

Based on the integration tests, we see:
- **Main function**: 34.9% coverage
- **Internal packages**: 30.8% coverage
- **Overall**: ~32% coverage

### Understanding the Numbers

- **High coverage functions**: `SetupLogger` (100%), `GetCurrentDirAsGitRepo` (81.8%)
- **Low coverage functions**: Many helper functions (0%) - these may not be exercised by current acceptance tests
- **Opportunities**: Error handling paths, interactive features, cleanup operations

## Best Practices

### 1. Combine with Unit Tests

Integration coverage complements but doesn't replace unit test coverage:

```bash
# Run all coverage types
make all-cover-html       # Unit test coverage
make integration-cover-html  # Integration coverage
```

### 2. Use for CI/CD Validation

Integration coverage helps validate that your acceptance tests actually exercise the code:

```bash
# In CI pipeline
make integration-cover
# Fail if coverage drops below threshold
```

### 3. Focus on End-to-End Scenarios

Design acceptance tests to cover:
- Complete user workflows
- Error conditions
- Edge cases in CLI parsing
- Integration between components

### 4. Monitor Coverage Trends

Track integration coverage over time:
- Set minimum coverage thresholds
- Monitor for coverage regression
- Identify untested code paths

## Cleanup

To remove all coverage files:

```bash
make clean-coverage
```

This removes:
- `*.out`, `*.html`, `*.txt` coverage files
- `covdatafiles/` directory
- `merged-covdata/` directory

## Troubleshooting

### Binary Not Found

If you see "binary not found" errors:
- Ensure the binary path is absolute
- Check that `go build -cover` succeeded
- Verify the `GITSWEEPER_TEST_BINARY` environment variable

### No Coverage Data

If no coverage data is generated:
- Ensure Go 1.20+ is installed
- Check that `GOCOVERDIR` is set and writable
- Verify the binary was built with `-cover`

### Low Coverage Numbers

If coverage seems too low:
- Review your acceptance test scenarios
- Add tests for error conditions
- Consider interactive command testing
- Check if main function initialization is covered

## References

- [Go Blog: Code coverage for Go integration tests](https://go.dev/blog/integration-test-coverage)
- [Go 1.20 Release Notes](https://go.dev/doc/go1.20)
- [`go tool covdata` documentation](https://pkg.go.dev/cmd/covdata)