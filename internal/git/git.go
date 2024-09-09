package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitAvailable checks if the 'git' command is available on the system and returns the Git version if available.
func IsGitAvailable() (bool, string) {
	// Attempt to run the 'git --version' command
	cmd := exec.Command("git", "--version")

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command and check if it executes successfully
	if err := cmd.Run(); err != nil {
		return false, ""
	}

	// Return true and the output (which contains the Git version)
	return true, out.String()
}

// IsGitRepository checks if the given directory is a Git repository.
func IsGitRepository(dirPath string) (bool, error) {
	// Ensure the path is an absolute path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Run the 'git rev-parse --is-inside-work-tree' command in the specified directory
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = absPath

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return false, nil // Not a git repository if command fails
	}

	// Check if the output is "true\n"
	if string(output) == "true\n" {
		return true, nil
	}

	return false, nil
}

// HasPendingChanges checks if the given directory has pending changes (uncommitted changes) in the Git repository.
func HasPendingChanges(dirPath string) (bool, error) {
	// Ensure the path is an absolute path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Run the 'git status --porcelain' command in the specified directory
	cmd := exec.Command("git", "status", "--porcelain", "--untracked-files=no")
	cmd.Dir = absPath

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	// If the output is not empty, there are pending changes
	if len(output) > 0 {
		return true, nil
	}

	return false, nil
}

// InitializeGitRepo initializes a new Git repository in the specified directory path.
func InitializeGitRepo(dirPath string) error {
	// Ensure the directory path is an absolute path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Ensure the directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Run the 'git init' command in the specified directory
	cmd := exec.Command("git", "init")
	cmd.Dir = absPath

	// Execute the command and check for errors
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	return nil
}

// CommitChanges commits any pending staged changes to the git repository.
// It performsa a `git commit -am <message>`.
func CommitChanges(dirPath string, commitMessage string) error {
	// Ensure the path is an absolute path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cmd := exec.Command("git", "commit", "-am", commitMessage)
	cmd.Dir = absPath

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println(out.String())
		return fmt.Errorf("git commit failed: %w", err)
	}
	// get the exit status of the command
	if exitStatus := cmd.ProcessState.ExitCode(); exitStatus != 0 {
		return fmt.Errorf("git commit failed with exit status %d: %s", exitStatus, out.String())
	}

	return nil
}

func TagChanges(root string, name string, message string) interface{} {
	// Ensure the path is an absolute path
	absPath, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cmd := exec.Command("git", "tag", "-a", name, "-m", message)
	cmd.Dir = absPath

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git tag failed: %w", err)
	}

	// get the exit status of the command
	if exitStatus := cmd.ProcessState.ExitCode(); exitStatus != 0 {
		return fmt.Errorf("git tag failed with exit status %d: %s", exitStatus, output)
	}

	return nil
}

// GetTags returns a list of git tags for the given project directory
func GetGitTags(projectDir string) ([]string, error) {
	// Prepare the git command
	cmd := exec.Command("git", "tag")
	cmd.Dir = projectDir

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Convert the output to a slice of strings, one per line
	tags := strings.Split(strings.TrimSpace(out.String()), "\n")

	return tags, nil
}

// TagExists checks if the given tag exists in the git repository of the project directory
func TagExists(projectDir string, tagName string) (bool, error) {
	tags, err := GetGitTags(projectDir)
	if err != nil {
		return false, err
	}
	for _, tag := range tags {
		if tag == tagName {
			return true, nil
		}
	}
	return false, nil
}

// GetCurrentBranch returns the current branch of the git repository in the given project directory
func GetCurrentBranch(projectDir string) (string, error) {
	// Prepare the git command
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = projectDir

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
