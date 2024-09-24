package config

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"github.com/ptgoetz/go-versionbump/internal/version"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"sort"
)

const (
	DefaultGitCommitTemplate     = "Bump version {old} --> {new}"
	DefaultGitTagTemplate        = "v{new}"
	DefaultGitTagMessageTemplate = "Release version {new}"
)

// Config represents the version bump configuration.
type Config struct {
	Version               string          `yaml:"version"`
	PreReleaseVersion     string          `yaml:"prerelease-version"`
	PreReleaseLabels      []string        `yaml:"prerelease-labels"`
	GitCommit             bool            `yaml:"git-commit"`
	GitCommitTemplate     string          `yaml:"git-commit-template"`
	GitSign               bool            `yaml:"git-sign"`
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
	ConfigPath      string
	Quiet           bool
	NoPrompt        bool
	ShowVersion     bool
	ResetVersion    string
	NoGit           bool
	NoColor         bool
	BumpPart        string
	PreReleaseLabel string
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

// HasLabel returns true if a given label is in the list of pre-release labels
func (v Config) HasLabel(label string) bool {
	for _, l := range v.PreReleaseLabels {
		if l == label {
			return true
		}
	}
	return false
}

// GetSortedLabels returns a sorted slice of pre-release labels
func (v Config) GetSortedLabels() []string {
	// Make a copy of the input slice to avoid modifying the original
	sortedStrings := make([]string, len(v.PreReleaseLabels))
	copy(sortedStrings, v.PreReleaseLabels)

	// Sort the strings lexically
	sort.Strings(sortedStrings)

	return sortedStrings
}

// VersionedFile represents the file to be updated with the new version.
type VersionedFile struct {
	Path    string   `yaml:"path"`
	Replace []string `yaml:"replace"`
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

	// make sure we can resolve the parent directory
	root, err := utils.ParentDirAbsolutePath(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("error getting parent directory: %w", err)
	}

	// validate the version string is not empty
	if config.Version == "" {
		return nil, "", fmt.Errorf("version string is required")
	}

	if !version.ValidateVersion(config.Version) {
		return nil, "", fmt.Errorf("invalid version string: %s", config.Version)
	}

	configPtr := &config
	// include the config file as a file to update
	configPtr.Files = append(configPtr.Files, VersionedFile{Path: configFile, Replace: []string{"version: \"{version}\""}})

	return configPtr, root, nil
}
