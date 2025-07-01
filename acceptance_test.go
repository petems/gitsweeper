package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelper provides utilities for acceptance testing
type TestHelper struct {
	t       *testing.T
	workDir string
	binPath string
}

// NewTestHelper creates a new test helper instance
func NewTestHelper(t *testing.T) *TestHelper {
	workDir, err := os.MkdirTemp("", "gitsweeper-test-*")
	require.NoError(t, err)

	binPath := filepath.Join(workDir, "gitsweeper-test")

	helper := &TestHelper{
		t:       t,
		workDir: workDir,
		binPath: binPath,
	}

	// Build the test binary
	helper.buildTestBinary()

	return helper
}

// Cleanup removes temporary files
func (h *TestHelper) Cleanup() {
	// Remove temporary directory
	os.RemoveAll(h.workDir)
}

// buildTestBinary builds the gitsweeper binary for testing
func (h *TestHelper) buildTestBinary() {
	cmd := exec.Command("go", "build", "-o", h.binPath, "main.go")
	err := cmd.Run()
	require.NoError(h.t, err, "Failed to build test binary")
}

// RunCommand executes the gitsweeper command and returns the result
func (h *TestHelper) RunCommand(args ...string) *CommandResult {
	return h.RunCommandInDir(h.workDir, args...)
}

// RunCommandInDir executes the gitsweeper command in a specific directory
func (h *TestHelper) RunCommandInDir(dir string, args ...string) *CommandResult {
	cmd := exec.Command(h.binPath, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	result := &CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
	}

	return result
}

// RunCommandInteractive executes a command with input simulation
func (h *TestHelper) RunCommandInteractive(input string, args ...string) *CommandResult {
	return h.RunCommandInteractiveInDir(h.workDir, input, args...)
}

// RunCommandInteractiveInDir executes a command with input simulation in a specific directory
func (h *TestHelper) RunCommandInteractiveInDir(dir, input string, args ...string) *CommandResult {
	cmd := exec.Command(h.binPath, args...)
	cmd.Dir = dir

	stdin, err := cmd.StdinPipe()
	require.NoError(h.t, err)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Start()
	require.NoError(h.t, err)

	// Send input
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input+"\n")
	}()

	err = cmd.Wait()
	
	result := &CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
	}

	return result
}

// CommandResult holds the result of a command execution
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// ContainsOutput checks if the output contains the expected string
func (r *CommandResult) ContainsOutput(expected string) bool {
	return strings.Contains(r.Stdout, expected) || strings.Contains(r.Stderr, expected)
}

// MatchesOutput checks if the output matches exactly
func (r *CommandResult) MatchesOutput(expected string) bool {
	return strings.TrimSpace(r.Stdout) == strings.TrimSpace(expected)
}

// CreateTestGitRepo creates a test git repository
func (h *TestHelper) CreateTestGitRepo(repoName string) string {
	repoPath := filepath.Join(h.workDir, repoName)
	err := os.MkdirAll(repoPath, 0755)
	require.NoError(h.t, err, "Failed to create repo directory")

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	err = cmd.Run()
	require.NoError(h.t, err, "Failed to initialize git repo")

	return repoPath
}

// CloneRepo clones a git repository
func (h *TestHelper) CloneRepo(repoURL string) string {
	cmd := exec.Command("git", "clone", repoURL)
	cmd.Dir = h.workDir
	err := cmd.Run()
	require.NoError(h.t, err, "Failed to clone repository")

	// Extract repo name from URL
	parts := strings.Split(repoURL, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	
	return filepath.Join(h.workDir, repoName)
}

// CreateBareRepo creates a bare git repository
func (h *TestHelper) CreateBareRepo(repoName string) string {
	repoPath := filepath.Join(h.workDir, repoName)
	cmd := exec.Command("git", "init", "--bare", repoPath)
	err := cmd.Run()
	require.NoError(h.t, err, "Failed to create bare repository")
	
	return repoPath
}

// AddRemote adds a new remote to a git repository
func (h *TestHelper) AddRemote(repoDir, remoteName, remoteURL string) {
	cmd := exec.Command("git", "remote", "add", remoteName, remoteURL)
	cmd.Dir = repoDir
	err := cmd.Run()
	require.NoError(h.t, err, "Failed to add remote")

	cmd = exec.Command("git", "fetch", remoteName)
	cmd.Dir = repoDir
	err = cmd.Run()
	require.NoError(h.t, err, "Failed to fetch remote")
}

// CreateDirectory creates a directory
func (h *TestHelper) CreateDirectory(dirName string) string {
	dirPath := filepath.Join(h.workDir, dirName)
	err := os.MkdirAll(dirPath, 0755)
	require.NoError(h.t, err, "Failed to create directory")
	return dirPath
}

// CheckCommandExists checks if a command is available in the system
func (h *TestHelper) CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// RequireCommand skips the test if the command is not available
func (h *TestHelper) RequireCommand(command string) {
	if !h.CheckCommandExists(command) {
		h.t.Skipf("Command %s not available", command)
	}
}

// CheckPortFree checks if a port is free
func (h *TestHelper) CheckPortFree(port string) bool {
	cmd := exec.Command("lsof", "-i", "TCP:"+port)
	err := cmd.Run()
	return err != nil // If lsof fails, port is likely free
}

// RequirePortFree skips the test if the port is not free
func (h *TestHelper) RequirePortFree(port string) {
	if !h.CheckPortFree(port) {
		h.t.Skipf("Port %s is not free", port)
	}
}

// Acceptance Tests

func TestVersionCommand(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("Version with no flags", func(t *testing.T) {
		result := helper.RunCommand("version")
		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.Stdout, "0.1.0 development")
	})

	t.Run("Version with --debug flag", func(t *testing.T) {
		result := helper.RunCommand("--debug", "version")
		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.Stdout, "0.1.0 development")
		// Note: Debug output might go to stderr
		output := result.Stdout + result.Stderr
		assert.Contains(t, output, "--debug setting detected - Info level logs enabled")
	})
}

func TestNoArgumentCommand(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("No arguments shows help", func(t *testing.T) {
		result := helper.RunCommand()
		// The application should show help or usage information
		// Based on kingpin behavior, it should exit with 0 and show help
		assert.Equal(t, 0, result.ExitCode)
	})
}

func TestCleanupCommand(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// Skip if required commands are not available
	helper.RequireCommand("go")
	helper.RequireCommand("git")

	t.Run("In a non-git repo", func(t *testing.T) {
		nonGitDir := helper.CreateDirectory("not-a-git-repo")
		
		result := helper.RunCommandInDir(nonGitDir, "cleanup")
		
		assert.Equal(t, 1, result.ExitCode)
		// The error message should indicate repository does not exist
		output := result.Stdout + result.Stderr
		assert.Contains(t, output, "repository does not exist")
	})

	t.Run("In a git repo with no remotes", func(t *testing.T) {
		repoDir := helper.CreateTestGitRepo("test-repo")
		
		result := helper.RunCommandInDir(repoDir, "cleanup", "--force")
		
		// The application will fail if there are no remotes
		assert.Equal(t, 1, result.ExitCode)
		output := result.Stdout + result.Stderr
		assert.Contains(t, output, "Error when looking for branches")
	})


}

func TestPreviewCommand(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// Skip if required commands are not available
	helper.RequireCommand("go")
	helper.RequireCommand("git")

	t.Run("Preview in non-git repo", func(t *testing.T) {
		nonGitDir := helper.CreateDirectory("not-a-git-repo")
		
		result := helper.RunCommandInDir(nonGitDir, "preview")
		
		assert.Equal(t, 1, result.ExitCode)
		output := result.Stdout + result.Stderr
		assert.Contains(t, output, "This is not a Git repository")
	})

	t.Run("Preview in git repo with no remotes", func(t *testing.T) {
		repoDir := helper.CreateTestGitRepo("test-repo")
		
		result := helper.RunCommandInDir(repoDir, "preview")
		
		assert.Equal(t, 1, result.ExitCode)
		output := result.Stdout + result.Stderr
		assert.Contains(t, output, "Error when looking for branches")
	})

	t.Run("Preview with custom master branch", func(t *testing.T) {
		repoDir := helper.CreateTestGitRepo("test-repo")
		
		result := helper.RunCommandInDir(repoDir, "preview", "--master=main")
		
		// This will fail because there are no remotes
		assert.Equal(t, 1, result.ExitCode)
		output := result.Stdout + result.Stderr
		assert.Contains(t, output, "Error when looking for branches")
	})
}