package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ptgoetz/go-versionbump/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewVersionBump(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "versionBumpTest")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create a temporary YAML config file
	filePath := filepath.Join(dir, "versionbump.yaml")
	yamlContent := `
  version: "1.0.0"
  git-commit: true
  git-tag: false
  files:
    - path: "version.go"
      replace: "v{version}"
    - path: "README.md"
      replace: "Version: {version}"
`
	if err := os.WriteFile(filePath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write to YAML config file: %v", err)
	}

	options := config.Options{
		ConfigPath: filePath,
		BumpPart:   "patch",
	}

	vb, err := NewVersionBump(options)
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", vb.GetOldVersion())
	assert.Equal(t, "1.0.1", vb.GetNewVersion())
	assert.Equal(t, dir, vb.ParentDir)
}

func TestGitMetadata(t *testing.T) {
	vb := &VersionBump{
		Config: config.Config{
			Version:               "1.0.0",
			GitCommitTemplate:     "Commit {old} to {new}",
			GitTagTemplate:        "v{new}",
			GitTagMessageTemplate: "Tagging version {new}",
		},
		Options: config.Options{
			BumpPart: "patch",
		},
	}

	gitMeta, err := vb.GitMetadata()
	assert.NoError(t, err)
	assert.Equal(t, "Commit 1.0.0 to 1.0.1", gitMeta.CommitMessage)
	assert.Equal(t, "v1.0.1", gitMeta.TagName)
	assert.Equal(t, "Tagging version 1.0.1", gitMeta.TagMessage)
}
