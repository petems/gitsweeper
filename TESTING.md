# Testing Guide

This project uses native Go testing for both unit and acceptance tests, replacing the previous Ruby/Aruba-based testing setup.

## Test Structure

### Unit Tests
- Located in `internal/githelpers_test.go`
- Test individual functions and components
- Run with: `make test` or `go test ./...`

### Acceptance Tests
- Located in `acceptance_test.go`
- Test the complete application behavior end-to-end
- Use a custom `TestHelper` to manage test environments
- Run with: `make acceptance-test`

## Running Tests

### All Tests
```bash
make test-all
```

### Unit Tests Only
```bash
make test
```

### Acceptance Tests Only
```bash
make acceptance-test
```

### Specific Test
```bash
go test -v -run TestVersionCommand
```

## Coverage Reports

### Acceptance Test Coverage
```bash
# Run acceptance tests with coverage
make acceptance-cover

# Generate HTML coverage report for acceptance tests
make acceptance-cover-html

# Show function-level coverage for acceptance tests
make acceptance-cover-func
```

### Combined Coverage
```bash
# Run all tests with combined coverage report
make all-cover

# Generate HTML coverage report for all tests
make all-cover-html
```

### Coverage Cleanup
```bash
# Clean up all coverage files
make clean-coverage
```

**Note**: Acceptance tests typically show 0% code coverage because they test the application as a black box by executing the compiled binary. This is expected behavior and indicates proper separation between unit tests (which test internal functions) and acceptance tests (which test end-to-end functionality).

## Test Dependencies

### Required for Acceptance Tests
- **Docker**: Required for running dummy git server containers
- **Git**: Required for git operations in tests
- **Go**: Required for building test binaries

### Optional Dependencies
- **lsof**: Used for port checking (tests will skip if not available)

## Test Helper Features

The `TestHelper` provides utilities for acceptance testing:

### Environment Management
- Creates temporary directories for each test
- Builds test binaries automatically
- Cleans up resources after tests

### Command Execution
- `RunCommand()`: Execute commands and capture output
- `RunCommandInteractive()`: Simulate user input
- `RunCommandInDir()`: Execute commands in specific directories

### Git Operations
- `CloneRepo()`: Clone repositories for testing
- `CreateBareRepo()`: Create bare git repositories
- `AddRemote()`: Add git remotes

### Docker Integration
- `StartDummyGitServer()`: Start containerized git servers
- Automatic container cleanup

### Dependency Checking
- `RequireCommand()`: Skip tests if commands are missing
- `RequirePortFree()`: Skip tests if ports are in use

## Test Scenarios

The acceptance tests cover the same scenarios as the previous Cucumber tests:

### Version Command
- Basic version output
- Version with debug flag

### Cleanup Command
- Force cleanup of merged branches
- Interactive cleanup (yes/no prompts)
- Error handling for non-git repositories
- Custom remote specification
- Branch skipping functionality
- Debug output verification

### Preview Command
- Preview branches to be deleted
- Custom remote and master branch options
- Skip functionality

## Migration from Ruby/Aruba

The new Go-based testing approach provides several advantages over the previous Ruby/Aruba setup:

### Benefits
- **Native Go**: No Ruby dependencies or external tools
- **Faster**: No process spawning overhead for simple tests
- **Better Integration**: Direct access to Go testing tools and IDE support
- **Maintainable**: Single language for application and tests
- **Portable**: Fewer external dependencies

### Removed Files
The following Ruby-based testing files are no longer needed:
- `Gemfile`
- `Gemfile.lock`
- `cucumber.yml`
- `features/` directory and all `.feature` files
- `features/step_definitions/` directory
- `features/support/` directory

### Docker Usage
Docker is still used for integration testing to provide:
- Dummy git servers for realistic testing scenarios
- Isolated test environments
- Consistent test data

## Writing New Tests

### Adding Unit Tests
Add test functions to existing `*_test.go` files in the `internal/` package:

```go
func TestNewFunction(t *testing.T) {
    result := NewFunction("input")
    assert.Equal(t, "expected", result)
}
```

### Adding Acceptance Tests
Add test functions to `acceptance_test.go`:

```go
func TestNewCommand(t *testing.T) {
    helper := NewTestHelper(t)
    defer helper.Cleanup()

    result := helper.RunCommand("new-command", "--flag")
    assert.Equal(t, 0, result.ExitCode)
    assert.Contains(t, result.Stdout, "expected output")
}
```

### Test Naming Convention
- Unit tests: `TestFunctionName`
- Acceptance tests: `TestCommandName`
- Sub-tests use `t.Run()` with descriptive names

## Troubleshooting

### Docker Issues
If Docker tests fail:
1. Ensure Docker is running
2. Check if the `petems/dummy-git-repo` image is available
3. Verify port 8008 is free
4. Tests will skip automatically if Docker is unavailable

### Port Conflicts
If port 8008 is in use:
1. Stop the conflicting service
2. Or modify the test to use a different port
3. Tests will skip automatically if the port is not free

### Missing Dependencies
Tests will automatically skip if required commands are not available:
- `docker` - For container-based tests
- `git` - For git operations
- `lsof` - For port checking (optional)