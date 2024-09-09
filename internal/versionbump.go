package internal

import (
	"bufio"
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/config"
	"github.com/ptgoetz/go-versionbump/internal/git"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	vbu "github.com/ptgoetz/go-versionbump/internal/utils"
	vbv "github.com/ptgoetz/go-versionbump/internal/version"

	"os"
	"path"
	"strings"
)

type VersionBump struct {
	Config     config.Config
	Options    config.Options
	ParentDir  string
	OldVersion string
	NewVersion string
}

// NewVersionBump creates a new VersionBump instance.
// It loads the configuration file and determines/validates the old and new versions.
// If the reset version option is set, the new version is set to the reset version.
func NewVersionBump(options config.Options) (*VersionBump, error) {
	// Log the version and configuration path

	cfg, parentDir, err := config.LoadConfig(options.ConfigPath)
	if err != nil {
		return nil, err
	}

	// determine the old and new versions
	var oldVersion string
	var newVersion string

	// get the old version from the config
	vTemp, err := vbv.ParseVersion(cfg.Version)
	if err != nil {
		logFatal(fmt.Sprintf("Failed to parse semantic version string for old version: %s", err))
	}
	oldVersion = vTemp.String()

	// get the new or reset version
	if options.IsResetVersion() {
		if !vbv.ValidateVersion(options.ResetVersion) {
			logFatal(fmt.Sprintf("Failed to parse semantic version reset string: %s\n", options.ResetVersion))
		}
		// set new version to the reset version
		vTemp, _ = vbv.ParseVersion(options.ResetVersion)
		newVersion = vTemp.String()
	} else {
		if !vbv.ValidateVersionPart(options.BumpPart) {
			logFatal(fmt.Sprintf("Invalid version part: %s", options.BumpPart))
		}
		_ = vTemp.StringBump(options.BumpPart)
		newVersion = vTemp.String()
	}

	vb := &VersionBump{
		Config:     *cfg,
		Options:    options,
		ParentDir:  parentDir,
		OldVersion: oldVersion,
		NewVersion: newVersion,
	}

	return vb, nil
}

func (vb *VersionBump) GitMetadata() (*config.GitMeta, error) {
	var commitMessageTemplate string
	if vb.Config.GitCommitTemplate != "" {
		commitMessageTemplate = vb.Config.GitCommitTemplate
	} else {
		commitMessageTemplate = config.DefaultGitCommitTemplate
	}
	commitMessage := utils.ReplaceInString(commitMessageTemplate, "{old}", vb.OldVersion)
	commitMessage = utils.ReplaceInString(commitMessage, "{new}", vb.NewVersion)

	var tagTemplate string
	if vb.Config.GitTagTemplate != "" {
		tagTemplate = vb.Config.GitTagTemplate
	} else {
		tagTemplate = config.DefaultGitTagTemplate
	}
	tagName := utils.ReplaceInString(tagTemplate, "{old}", vb.OldVersion)
	tagName = utils.ReplaceInString(tagName, "{new}", vb.NewVersion)

	var tagMessageTemplate string
	if vb.Config.GitTagMessageTemplate != "" {
		tagMessageTemplate = vb.Config.GitTagMessageTemplate
	} else {
		tagMessageTemplate = config.DefaultGitTagMessageTemplate
	}
	tagMessage := utils.ReplaceInString(tagMessageTemplate, "{old}", vb.OldVersion)
	tagMessage = utils.ReplaceInString(tagMessage, "{new}", vb.NewVersion)

	return &config.GitMeta{
		CommitMessage: commitMessage,
		TagMessage:    tagMessage,
		TagName:       tagName,
	}, nil
}

// -------------------------

// -------------------------

func (vb *VersionBump) Run() {
	vb.preamble()
	vb.gitPreFlight()
	vb.logTrackedFiles()
	vb.bumpPreflight()
	if vb.promptProceedWithChanges() {
		vb.makeChanges()
		vb.gitCommit()
	}
}

func (vb *VersionBump) gitPreFlight() {
	if vb.Options.NoGit {
		return
	}

	// make sure the `git` command is available
	if vb.Config.IsGitRequired() {
		isGitAvalable, version := git.IsGitAvailable()
		if !isGitAvalable {
			logFatal("Git is required by the configuration but is not available. " +
				"VersionBump requires Git to be installed and available in the system PATH in order †o perform Giit " +
				"operations")
		} else {
			logVerbose(vb.Options, fmt.Sprintf("Git version: %s", strings.TrimSpace(version)[12:]))
		}
	}

	// check if the parent directory is a Git repository
	isGitRepo, err := git.IsGitRepository(vb.ParentDir)
	if err != nil {
		logFatal(fmt.Sprintf("Error checking for git repository: %v\n", err))
	}
	if !isGitRepo {

		if vb.Options.NoPrompt {
			logFatal("The project root is not a Git repository, but Git options are enabled in the " +
				"configuration file.")
		}
		if promptUserConfirmation("Do you want to initialize a Git repository in this directory?") {
			err := git.InitializeGitRepo(vb.ParentDir)
			if err != nil {
				logFatal(fmt.Sprintf("Unable to initialize Git repository: %v\n", err))
			}
		}
	}

	// check if the Git repository has pending changes
	isDirty, _ := git.HasPendingChanges(vb.ParentDir)
	if isDirty {
		fmt.Println("ERROR: The Git repository has pending changes. Please commit or stash them before proceeding.")
		os.Exit(1)
	}
}

// gitPreFlight performs a pre-flight check for Git operations.
func (vb *VersionBump) preamble() {
	logVerbose(vb.Options, vbv.VersionBumpVersion)
	logVerbose(vb.Options, fmt.Sprintf("Configuration file: %s", vb.Options.ConfigPath))
	logVerbose(vb.Options, fmt.Sprintf("Project root directory: %s", vb.ParentDir))
}

func (vb *VersionBump) logTrackedFiles() {
	// Log the files that will be updated
	logVerbose(vb.Options, "Tracked Files:")
	for _, file := range vb.Config.Files {
		logVerbose(vb.Options, fmt.Sprintf("  - %s", file.Path))
	}
}

// gitCommit conditionally commits the changes to the Git repository.
func (vb *VersionBump) gitCommit() {
	if vb.Options.NoGit || !vb.Config.IsGitRequired() {
		return
	}
	gitMeta, err := vb.GitMetadata()
	if err != nil {
		logFatal(fmt.Sprintf("Unable to get Git metadata: %v\n", err))
	}
	if vb.Config.IsGitRequired() && !vb.Options.NoPrompt {

		logVerbose(vb.Options, fmt.Sprintf("Commit Message: %s\nTag Message: %s\nTag Name: %s",
			gitMeta.CommitMessage,
			gitMeta.TagMessage,
			gitMeta.TagName))
		proceed := promptUserConfirmation("Do you want to commit the changes to the git repository?")
		if !proceed {
			os.Exit(1)
		}
	}

	// commit changes
	if vb.Config.GitCommit {
		logVerbose(vb.Options, "Committing changes...")
		err := git.CommitChanges(vb.ParentDir, gitMeta.CommitMessage)
		if err != nil {
			fmt.Printf("Error committing changes: %v\n", err)
			os.Exit(1)
		}
		logVerbose(vb.Options, fmt.Sprintf("Committed changes with message: %s", gitMeta.CommitMessage))
	}
	if vb.Config.GitTag {
		logVerbose(vb.Options, "Tagging changes...")
		err := git.TagChanges(vb.ParentDir, gitMeta.TagName, gitMeta.TagMessage)
		if err != nil {
			fmt.Printf("Error tagging changes: %v\n", err)
			os.Exit(1)
		}
		logVerbose(vb.Options,
			fmt.Sprintf(
				"Tagged '%s' created with message: %s",
				gitMeta.TagName,
				gitMeta.TagMessage))
	}
}

// bumpPreflight performs a pre-flight check for the version bump operation.
func (vb *VersionBump) bumpPreflight() {
	if vb.Options.ResetVersion == "" {
		logVerbose(vb.Options, fmt.Sprintf("Bumping version part: %s", vb.Options.BumpPart))
	} else {
		logVerbose(vb.Options, fmt.Sprintf("Resetting version to: %s", vb.NewVersion))
	}
	logVerbose(vb.Options, fmt.Sprintf("Will bump version %s --> %s", vb.OldVersion, vb.NewVersion))

	// log what changes will be made to each file
	for _, file := range vb.Config.Files {
		find := vbu.ReplaceInString(file.Replace, "{version}", vb.OldVersion)
		replace := vbu.ReplaceInString(file.Replace, "{version}", vb.NewVersion)

		logVerbose(vb.Options, file.Path)
		logVerbose(vb.Options, fmt.Sprintf("     Find: \"%s\"", find))
		logVerbose(vb.Options, fmt.Sprintf("  Replace: \"%s\"", replace))
		count, err := vbu.CountStringsInFile(path.Join(vb.ParentDir, file.Path), find)
		if err != nil {
			fmt.Println(fmt.Errorf("error getting replacement count: a%v", err))
			os.Exit(1)
		}
		if count > 0 {
			logVerbose(vb.Options, fmt.Sprintf("    Found %d replacement(s)", count))
		} else {
			fmt.Println("ERROR: No replacements found in file: ", file.Path)
			os.Exit(1)
		}
	}
}

// makeChanges updates the version in the files.
func (vb *VersionBump) makeChanges() {
	// at this point we have already checked the config and there are no errors
	for _, file := range vb.Config.Files {
		find := vbu.ReplaceInString(file.Replace, "{version}", vb.OldVersion)
		replace := vbu.ReplaceInString(file.Replace, "{version}", vb.NewVersion)

		if !vb.Options.DryRun {
			var resolvedPath string
			if path.IsAbs(file.Path) {
				resolvedPath = file.Path
			} else {
				resolvedPath = path.Join(vb.ParentDir, file.Path)
			}
			err := vbu.ReplaceInFile(resolvedPath, find, replace)
			if err != nil {
				fmt.Println(fmt.Errorf("error updating file %s: a%v", file.Path, err))
				os.Exit(1)
			}
			logVerbose(vb.Options, fmt.Sprintf("Updated file: %s", file.Path))
		}
	}
}

// promptProceedWithChanges prompts the user to proceed with the changes.
func (vb *VersionBump) promptProceedWithChanges() bool {
	if !vb.Options.NoPrompt {
		if !promptUserConfirmation("Proceed with the changes?") {
			logVerbose(vb.Options, "Cancelled by user.")
			os.Exit(0)
		}
	}
	return true
}

// promptUserConfirmation prompts the user with the given prompt string and expects 'y' or 'n' input.
// It returns true for 'y' and false for 'n'.
func promptUserConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Print the prompt and read the user's input
		fmt.Printf("%s [y/N]: ", prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		// Trim the input and convert it to lowercase
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "" {
			input = "n"
		}

		// Check if the input is 'y' or 'n'
		if input == "y" {
			return true
		} else if input == "n" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
}

func logFatal(msg string) {
	fmt.Printf("ERROR: %s", msg)
	os.Exit(1)
}

func logVerbose(opts config.Options, msg string) {
	if opts.DryRun || !opts.Quiet {
		fmt.Println(msg)
	}
}
