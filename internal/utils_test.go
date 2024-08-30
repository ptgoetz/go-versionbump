package internal

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCountStringOccurrences tests the countStringOccurrences function
func TestCountStringOccurrences(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "countStringOccurrencesTest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create a temporary file
	filePath := filepath.Join(dir, "test.txt")
	content := "Hello, world!\nHello, Go!\nHello, world!\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	// Count occurrences of "Hello"
	count, err := CountStringOccurrences(filePath, "Hello")
	if err != nil {
		t.Fatalf("countStringOccurrences failed: %v", err)
	}

	// Verify the count
	expectedCount := 3
	if count != expectedCount {
		t.Errorf("Expected count %d, but got %d", expectedCount, count)
	}
}

// TestGetParentDirAbsolutePath tests the getParentDirAbsolutePath function
func TestGetParentDirAbsolutePath(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "getParentDirAbsolutePathTest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create a temporary file path
	filePath := filepath.Join(dir, "subdir", "test.txt")
	parentDirPath := filepath.Join(dir, "subdir")

	// Make sure the subdir exists
	if err := os.MkdirAll(parentDirPath, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Get the parent directory absolute path
	parentDir, err := getParentDirAbsolutePath(filePath)
	if err != nil {
		t.Fatalf("getParentDirAbsolutePath failed: %v", err)
	}

	// Verify the parent directory path
	expectedParentDir, err := filepath.Abs(parentDirPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	if parentDir != expectedParentDir {
		t.Errorf("Expected parent directory %s, but got %s", expectedParentDir, parentDir)
	}
}
