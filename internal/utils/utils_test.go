package utils

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
	count, err := CountStringsInFile(filePath, "Hello")
	if err != nil {
		t.Fatalf("countStringOccurrences failed: %v", err)
	}

	// Verify the count
	expectedCount := 3
	if count != expectedCount {
		t.Errorf("Expected count %d, but got %d", expectedCount, count)
	}
}

// TestGetParentDirAbsolutePath tests the GetParentDirAbsolutePath function
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
	parentDir, err := ParentDirAbsolutePath(filePath)
	if err != nil {
		t.Fatalf("GetParentDirAbsolutePath failed: %v", err)
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

func TestIsAllAlphabetic(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"alpha", true},
		{"alpha1", false},
		{"alpha-", false},
		{"", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual := IsAllAlphabetic(test.input)
			if actual != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, actual)
			}
		})
	}
}

func TestStartsWithDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1alpha", true},
		{"alpha1", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual := StartsWithDigit(test.input)
			if actual != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, actual)
			}
		})
	}
}
