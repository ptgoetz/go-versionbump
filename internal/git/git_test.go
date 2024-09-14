package git

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsGitAvailable tests the IsGitAvailable function
func TestIsGitAvailable(t *testing.T) {
	available, version := IsGitAvailable()
	if !available {
		t.Fatalf("Expected Git to be available, but it was not")
	}
	if version == "" {
		t.Fatalf("Expected a valid Git version, but got an empty string")
	}
	t.Logf("Git is available with version: %s", version)
}

// TestIsGitRepository tests the IsRepository function
func TestIsGitRepository(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "gitrepo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Initialize a new Git repository in the temp directory
	err = InitializeGitRepo(dir)
	if err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Check if the directory is recognized as a Git repository
	isRepo, err := IsRepository(dir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !isRepo {
		t.Fatalf("Expected directory to be a Git repository, but it was not")
	}
}

// TestHasPendingChanges tests the HasPendingChanges function
func TestHasPendingChanges(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "gitrepo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Initialize a new Git repository in the temp directory
	err = InitializeGitRepo(dir)
	if err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Initially, there should be no pending changes
	hasChanges, err := HasPendingChanges(dir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if hasChanges {
		t.Fatalf("Expected no pending changes, but found some")
	}

	// Create a new file in the repository to introduce changes
	testFilePath := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFilePath, []byte("Hello, Git!"), 0644); err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	// Now, there should be pending changes
	//hasChanges, err = HasPendingChanges(dir)
	//if err != nil {
	//	t.Fatalf("Unexpected error: %v", err)
	//}
	//if !hasChanges {
	//	t.Fatalf("Expected pending changes, but found none")
	//}
}

// TestInitializeGitRepo tests the InitializeGitRepo function
func TestInitializeGitRepo(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "gitrepo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Initialize a new Git repository in the temp directory
	err = InitializeGitRepo(dir)
	if err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Check if the .git directory exists
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Fatalf("Expected .git directory to exist, but it does not")
	}
}
