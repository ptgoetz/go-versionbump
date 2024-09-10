package git

import (
	"bytes"
	"fmt"
	vbc "github.com/ptgoetz/go-versionbump/internal/config"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitAvailable checks if the 'git' command is available on the system and returns the Git version if available.
func IsGitAvailable() (bool, string) {
	// Attempt to run the 'git --version' command
	out, _, err := runGitCommand("", "--version")
	if err != nil {
		return false, ""
	}
	return true, out
}

// IsRepository checks if the given directory is a Git repository.
func IsRepository(dirPath string) (bool, error) {
	out, _, err := runGitCommand(dirPath, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, nil
	}
	// Check if the output is "true\n"
	if string(out) == "true\n" {
		return true, nil
	}
	return false, nil
}

// HasPendingChanges checks if the given directory has pending changes (uncommitted changes) in the Git repository.
func HasPendingChanges(dirPath string) (bool, error) {
	out, _, err := runGitCommand(dirPath, "status", "--porcelain", "--untracked-files=no")
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}
	// If the output is not empty, there are pending changes
	if len(out) > 0 {
		return true, nil
	}
	return false, nil
}

// InitializeGitRepo initializes a new Git repository in the specified directory path.
func InitializeGitRepo(dirPath string) error {
	_, _, err := runGitCommand(dirPath, "init")
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}
	return nil
}

// AddFiles adds the specified files to the staging area of the git repository.
func AddFiles(dirPath string, files ...vbc.VersionedFile) error {
	for _, file := range files {
		_, _, err := runGitCommand(dirPath, "add", file.Path)
		if err != nil {
			return fmt.Errorf("failed to add file to git staging area: %w", err)
		}
	}
	return nil
}

// IsSigningEnabled checks if GPG signing is enabled for commits in the git repository.
func IsSigningEnabled(dirPath string) (bool, error) {
	out, _, err := runGitCommand(dirPath, "config", "--get", "commit.gpgsign")
	if err != nil {
		return false, fmt.Errorf("failed to get git config 'commit.gpgsign': %w", err)
	}
	return strings.TrimSpace(out) == "true", nil
}

// GetSigningKey returns the GPG signing key used for signing commits and tags.
func GetSigningKey(dirPath string) (string, error) {
	out, _, err := runGitCommand(dirPath, "config", "--get", "user.signingkey")
	if err != nil {
		return "", fmt.Errorf("failed to get git signing key: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// CommitChanges commits any pending staged changes to the git repository.
// It performsa a `git commit -am <message>`.
func CommitChanges(dirPath string, commitMessage string, sign bool) error {
	args := []string{"commit", "-am", commitMessage}
	if sign {
		args = append(args, "-S")
	}
	_, _, err := runGitCommand(dirPath, args...)
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	return nil
}

// TagChanges creates a new tag in the git repository with the specified name and message.
func TagChanges(root string, name string, message string, sign bool) error {
	args := []string{"tag", "-a", name, "-m", message}
	if sign {
		args = append(args, "-s")
	}
	_, _, err := runGitCommand(root, args...)
	if err != nil {
		return fmt.Errorf("failed to tag changes: %w", err)
	}
	return nil
}

// GetTags returns a list of git tags for the given project directory
func GetTags(projectDir string) ([]string, error) {
	out, _, err := runGitCommand(projectDir, "tag", "--list")
	if err != nil {
		return nil, fmt.Errorf("failed to get git tags: %w", err)
	}
	// Convert the output to a slice of strings, one per line
	tags := strings.Split(strings.TrimSpace(out), "\n")
	return tags, nil
}

// TagExists checks if the given tag exists in the git repository of the project directory
func TagExists(projectDir string, tagName string) (bool, error) {
	tags, err := GetTags(projectDir)
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
	out, _, err := runGitCommand(projectDir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// runGitCommand runs a git command in the specified directory and returns the output and error messages.
func runGitCommand(root string, args ...string) (string, string, error) {
	absPath, err := filepath.Abs(root)
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = absPath

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	// Some git commends will return a non-zero exit code even if they succeed
	// For example `git conig --get user.email` will return 1 if the email is not set
	// We don't want to treat this as an error
	_ = cmd.Run()

	// Run the command
	if stdErr.String() != "" {
		return "", "", fmt.Errorf("git command failed: %s", stdErr.String())
	}

	return stdOut.String(), stdErr.String(), nil
}
