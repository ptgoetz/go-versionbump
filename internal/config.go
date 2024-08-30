package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

const (
	DefaultGitCommitTemplate     = "Bump version {old} --> {new}"
	DefaultGitTagTemplate        = "v{new}"
	DefaultGitTagMessageTemplate = "Release version {new}"
)

// VersionBump represents the version bump configuration.
type VersionBump struct {
	Version               string          `yaml:"version"`
	GitCommit             bool            `yaml:"git-commit"`
	GitCommitTemplate     string          `yaml:"git-commit-template"`
	GitTag                bool            `yaml:"git-tag"`
	GitTagTemplate        string          `yaml:"git-tag-template"`
	GitTagMessageTemplate string          `yaml:"git-tag-message-template"`
	Files                 []VersionedFile `yaml:"files"`
}

type VersionBumpMetadata struct {
	OldVersion    string
	NewVersion    string
	CommitMessage string
	TagMessage    string
	TagName       string
}

func (vbm *VersionBumpMetadata) String() string {
	return fmt.Sprintf("Commit Message: %s\nTag Message: %s\nTag Name: %s",
		vbm.CommitMessage, vbm.TagMessage, vbm.TagName)
}

// IsGitRequired returns true if any of the Git options are enabled.
func (v VersionBump) IsGitRequired() bool {
	return v.GitCommit || v.GitTag
}

func (v VersionBump) BumpAndGetMetaData(versionPart string) (*VersionBumpMetadata, error) {
	version, err := ParseVersion(v.Version)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse semantic version string: %s\n", v.Version)

	}
	oldVersionStr := version.String()
	// bump version
	err = version.StringBump(versionPart)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %v\n", err)
	}
	newVersionStr := version.String()

	// check for template overrides
	var commitMessageTemplate string
	if v.GitCommitTemplate != "" {
		commitMessageTemplate = v.GitCommitTemplate
	} else {
		commitMessageTemplate = DefaultGitCommitTemplate
	}
	commitMessage := ReplaceInString(commitMessageTemplate, "{old}", oldVersionStr)
	commitMessage = ReplaceInString(commitMessage, "{new}", newVersionStr)

	var tagTemplate string
	if v.GitTagTemplate != "" {
		tagTemplate = v.GitTagTemplate
	} else {
		tagTemplate = DefaultGitTagTemplate
	}
	tagName := ReplaceInString(tagTemplate, "{old}", oldVersionStr)
	tagName = ReplaceInString(tagName, "{new}", newVersionStr)

	var tagMessageTemplate string
	if v.GitTagMessageTemplate != "" {
		tagMessageTemplate = v.GitTagMessageTemplate
	} else {
		tagMessageTemplate = DefaultGitTagMessageTemplate
	}
	tagMessage := ReplaceInString(tagMessageTemplate, "{old}", oldVersionStr)
	tagMessage = ReplaceInString(tagMessage, "{new}", newVersionStr)

	return &VersionBumpMetadata{
		OldVersion:    oldVersionStr,
		NewVersion:    newVersionStr,
		CommitMessage: commitMessage,
		TagMessage:    tagMessage,
		TagName:       tagName,
	}, nil
}

// VersionedFile represents the file to be updated with the new version.
type VersionedFile struct {
	Path    string `yaml:"path"`
	Replace string `yaml:"replace"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(filePath string) (*VersionBump, string, error) {
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

	// Parse the YAML file into the VersionBump struct
	var config VersionBump
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing config file: %w", err)
	}

	root, err := getParentDirAbsolutePath(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("error getting parent directory: %w", err)
	}

	configPtr := &config
	// include the config file as a file to update
	configPtr.Files = append(configPtr.Files, VersionedFile{Path: configFile, Replace: "version: \"{version}\""})

	return configPtr, root, nil
}
