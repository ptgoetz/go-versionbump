package config

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

const (
	DefaultGitCommitTemplate     = "Bump version {old} --> {new}"
	DefaultGitTagTemplate        = "v{new}"
	DefaultGitTagMessageTemplate = "Release version {new}"
)

// Config represents the version bump configuration.
type Config struct {
	Version               string          `yaml:"version"`
	GitCommit             bool            `yaml:"git-commit"`
	GitCommitTemplate     string          `yaml:"git-commit-template"`
	GitTag                bool            `yaml:"git-tag"`
	GitTagTemplate        string          `yaml:"git-tag-template"`
	GitTagMessageTemplate string          `yaml:"git-tag-message-template"`
	Files                 []VersionedFile `yaml:"files"`
}

type GitMeta struct {
	OldVersion    string
	NewVersion    string
	CommitMessage string
	TagMessage    string
	TagName       string
}

type Options struct {
	ConfigPath   string
	DryRun       bool
	Quiet        bool
	NoPrompt     bool
	ShowVersion  bool
	ResetVersion string
	NoGit        bool
	BumpPart     string
}

func (o Options) IsResetVersion() bool {
	return o.ResetVersion != ""
}

func (vbm *GitMeta) String() string {
	return fmt.Sprintf("Commit Message: %s\nTag Message: %s\nTag Name: %s",
		vbm.CommitMessage, vbm.TagMessage, vbm.TagName)
}

// IsGitRequired returns true if any of the Git options are enabled.
func (v Config) IsGitRequired() bool {
	return v.GitCommit || v.GitTag
}

// VersionedFile represents the file to be updated with the new version.
type VersionedFile struct {
	Path    string `yaml:"path"`
	Replace string `yaml:"replace"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(filePath string) (*Config, string, error) {
	// Open the YAML file
	configFile := path.Base(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("error opening config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing config file: %v", err)
		}
	}(file)

	// Parse the YAML file into the Config struct
	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing config file: %w", err)
	}

	root, err := utils.ParentDirAbsolutePath(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("error getting parent directory: %w", err)
	}

	configPtr := &config
	// include the config file as a file to update
	configPtr.Files = append(configPtr.Files, VersionedFile{Path: configFile, Replace: "version: \"{version}\""})

	return configPtr, root, nil
}
