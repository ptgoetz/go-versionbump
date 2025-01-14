package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig tests the LoadConfig function
func TestLoadConfig(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "loadConfigTest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create a temporary YAML config file
	filePath := filepath.Join(dir, "config.yaml")
	yamlContent := `
version: "1.0.0"
git-commit: true
git-tag: false
git-sign: true
files:
  - path: "version.go"
    replace: 
      - "v{version}"
  - path: "README.md"
    replace: 
      - "version: {version}"
`
	if err := os.WriteFile(filePath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write to YAML config file: %v", err)
	}

	// Load the config file
	config, root, err := LoadConfig(filePath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify the parsed configuration
	if config.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', but got '%s'", config.Version)
	}
	if !config.GitCommit {
		t.Errorf("Expected git-commit to be true, but got false")
	}
	if config.GitTag {
		t.Errorf("Expected git-tag to be false, but got true")
	}
	if len(config.Files) != 3 {
		t.Fatalf("Expected 2 files, but got %d", len(config.Files))
	}
	if config.Files[0].Path != "version.go" || config.Files[0].Replace[0] != "v{version}" {
		t.Errorf("Unexpected file config for 'version.go': %+v", config.Files[0])
	}
	if config.Files[1].Path != "README.md" || config.Files[1].Replace[0] != "version: {version}" {
		t.Errorf("Unexpected file config for 'README.md': %+v", config.Files[1])
	}

	// Verify the root directory
	expectedRoot, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	if root != expectedRoot {
		t.Errorf("Expected root directory '%s', but got '%s'", expectedRoot, root)
	}
}

// TestLoadConfigFileNotFound tests the LoadConfig function with a non-existent file
func TestLoadConfigFileNotFound(t *testing.T) {
	_, _, err := LoadConfig("non_existent.yaml")
	if err == nil {
		t.Fatal("Expected an error when loading a non-existent file, but got none")
	}
}

// TestLoadConfigInvalidYAML tests the LoadConfig function with invalid YAML content
func TestLoadConfigInvalidYAML(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "loadConfigInvalidYAMLTest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create a temporary invalid YAML config file
	filePath := filepath.Join(dir, "config.yaml")
	invalidYamlContent := `
version: 1.0.0"
git-commit: true
git-tag: false
git-sigrue
files:
  - path: "version.go"
    replace: "v{version}"
  - path: "README.md"
    replace: "version: {version"
`
	if err := os.WriteFile(filePath, []byte(invalidYamlContent), 0644); err != nil {
		t.Fatalf("Failed to write to YAML config file: %v", err)
	}

	// Load the config file
	_, _, err = LoadConfig(filePath)
	if err == nil {
		t.Fatal("Expected an error when loading an invalid YAML file, but got none")
	}
}
