package internal

import (
	"bufio"
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/config"
	"github.com/ptgoetz/go-versionbump/internal/git"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	vbu "github.com/ptgoetz/go-versionbump/internal/utils"
	vbv "github.com/ptgoetz/go-versionbump/internal/version"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

type VersionBump struct {
	Config    config.Config
	Options   config.Options
	ParentDir string
}

// NewVersionBump creates a new VersionBump instance.
func NewVersionBump(options config.Options) (*VersionBump, error) {

	cfg, parentDir, err := config.LoadConfig(options.ConfigPath)
	if err != nil {
		logFatal(options, fmt.Sprintf("Error loading configuration file: %v", err))
	}

	vb := &VersionBump{
		Config:    *cfg,
		Options:   options,
		ParentDir: parentDir,
	}

	return vb, nil
}

func (vb *VersionBump) GetOldVersion() string {
	if !vbv.ValidateVersion(vb.Config.Version) {
		logFatal(vb.Options, fmt.Sprintf("Failed to parse semantic version string for old version: %s", vb.Config.Version))
	}
	oldVersion, _ := vbv.ParseVersion(vb.Config.Version)
	return oldVersion.String()
}

func (vb *VersionBump) GetNewVersion() string {
	if vb.Options.IsResetVersion() {
		v, err := vbv.ParseVersion(vb.Options.ResetVersion)
		if err != nil {
			logFatal(vb.Options, fmt.Sprintf("Failed to parse semantic version string for reset version: %s", vb.Options.ResetVersion))
		}
		return v.String()
	}
	oldVersionStr := vb.GetOldVersion()
	oldVersion, _ := vbv.ParseVersion(oldVersionStr)
	newVersion := oldVersion.StringBump(vb.Options.BumpPart)
	return newVersion.String()
}

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

func (vb *VersionBump) Show(versionStr string) error {
	var curVersionStr string
	isProject := false
	if versionStr != "" {
		curVersionStr = versionStr
	} else {
		curVersionStr = vb.Config.Version
		isProject = true
	}
	curVersion, err := vbv.ParseVersion(curVersionStr)
	if err != nil {
		return err
	}

	if !isProject {
		logVerbose(vb.Options, fmt.Sprintf("Potential versioning paths for version: %s",
			curVersion.String()))
	} else {
		logVerbose(vb.Options, fmt.Sprintf("Potential versioning paths for project version: %s",
			curVersion.String()))
	}
	// we now know we have a valid version
	majorVersion := curVersion.StringBump(vbv.VersionMajorStr)
	minorVersion := curVersion.StringBump(vbv.VersionMinorStr)
	patchVersion := curVersion.StringBump(vbv.VersionPatchStr)

	padLen := len(curVersion.String()) - 2
	padding := utils.PaddingString(padLen, " ")

	tree := fmt.Sprintf(
		`%s ── bump ─┬─ major ─ %s
        %s    ├─ minor ─ %s
        %s    ╰─ patch ─ %s
`,
		curVersion.String(),
		majorVersion.String(),
		padding,
		minorVersion.String(),
		padding,
		patchVersion.String())

	printColorOpts(vb.Options, tree, ColorBlue)
	return nil
}

func (vb *VersionBump) ShowEffectiveConfig() error {
	logVerbose(vb.Options, fmt.Sprintf("Config file: %s", vb.Options.ConfigPath))
	logVerbose(vb.Options, fmt.Sprintf("Project root: %s", vb.ParentDir))
	logVerbose(vb.Options, "Effective Configuration YAML:")

	conf := &vb.Config

	// set defaults if not overridden
	if conf.GitCommitTemplate == "" {
		conf.GitCommitTemplate = config.DefaultGitCommitTemplate
	}
	if conf.GitTagTemplate == "" {
		conf.GitTagTemplate = config.DefaultGitTagTemplate
	}
	if conf.GitTagMessageTemplate == "" {
		conf.GitTagMessageTemplate = config.DefaultGitTagMessageTemplate
	}

	b, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	printColorOpts(vb.Options, string(b), ColorBlue)
	return nil
}

func (vb *VersionBump) GitMetadata() (*config.GitMeta, error) {
	var commitMessageTemplate string
	if vb.Config.GitCommitTemplate != "" {
		commitMessageTemplate = vb.Config.GitCommitTemplate
	} else {
		commitMessageTemplate = config.DefaultGitCommitTemplate
	}
	commitMessage := utils.ReplaceInString(commitMessageTemplate, "{old}", vb.GetOldVersion())
	commitMessage = utils.ReplaceInString(commitMessage, "{new}", vb.GetNewVersion())

	var tagTemplate string
	if vb.Config.GitTagTemplate != "" {
		tagTemplate = vb.Config.GitTagTemplate
	} else {
		tagTemplate = config.DefaultGitTagTemplate
	}
	tagName := utils.ReplaceInString(tagTemplate, "{old}", vb.GetOldVersion())
	tagName = utils.ReplaceInString(tagName, "{new}", vb.GetNewVersion())

	var tagMessageTemplate string
	if vb.Config.GitTagMessageTemplate != "" {
		tagMessageTemplate = vb.Config.GitTagMessageTemplate
	} else {
		tagMessageTemplate = config.DefaultGitTagMessageTemplate
	}
	tagMessage := utils.ReplaceInString(tagMessageTemplate, "{old}", vb.GetOldVersion())
	tagMessage = utils.ReplaceInString(tagMessage, "{new}", vb.GetNewVersion())

	return &config.GitMeta{
		CommitMessage: commitMessage,
		TagMessage:    tagMessage,
		TagName:       tagName,
	}, nil
}

func (vb *VersionBump) gitPreFlight() {
	if vb.Options.NoGit {
		return
	}

	logVerbose(vb.Options, "Checking git configuration...")

	// make sure the `git` command is available
	if vb.Config.IsGitRequired() {
		isGitAvalable, version := git.IsGitAvailable()
		if !isGitAvalable {
			logFatal(vb.Options, "Git is required by the configuration but is not available. "+
				"VersionBump requires Git to be installed and available in the system PATH in order †o perform Giit "+
				"operations")
		} else {
			logVerbose(vb.Options, fmt.Sprintf("Git version: %s", strings.TrimSpace(version)[12:]))
		}
	}

	// check if the parent directory is a Git repository
	isGitRepo, err := git.IsRepository(vb.ParentDir)
	if err != nil {
		logFatal(vb.Options, fmt.Sprintf("Error checking for git repository: %v\n", err))
	}
	if !isGitRepo {

		if vb.Options.NoPrompt {
			logFatal(vb.Options, "The project root is not a Git repository, but Git options are enabled in the "+
				"configuration file.")
		}
		if promptUserConfirmation("The project directory is not a git repository.\nDo you want to initialize a git repository in the project directory?") {
			err := git.InitializeGitRepo(vb.ParentDir)
			if err != nil {
				logFatal(vb.Options, fmt.Sprintf("Unable to initialize Git repository: %v\n", err))
			}
			logVerbose(vb.Options, "Initialized Git repository.\nAdding tracked files...")
			vb.logTrackedFiles()
			err = git.AddFiles(vb.ParentDir, vb.Config.Files...)
			if err != nil {
				logFatal(vb.Options, fmt.Sprintf("Error adding files to the Git staging area: %v\n", err))
			}
			logVerbose(vb.Options, "Performing initial commit.")
			err = git.CommitChanges(vb.ParentDir, "Initial commit", vb.Config.GitSign)
			if err != nil {
				logFatal(vb.Options, fmt.Sprintf("Error committing initial changes: %v\n", err))
			}
		} else {
			os.Exit(0)
		}
	}

	branch, err := git.GetCurrentBranch(vb.ParentDir)
	if err != nil {
		logFatal(vb.Options, fmt.Sprintf("Error getting current branch: %v\n", err))
	}
	logVerbose(vb.Options, fmt.Sprintf("Current branch: %s", branch))

	if vb.Config.GitTag {
		// check to see if the tag already exists
		logVerbose(vb.Options, "Checking for existing tag...")
		gitMeta, err := vb.GitMetadata()
		if err != nil {
			logFatal(vb.Options, fmt.Sprintf("Unable to get Git metadata: %v\n", err))
		}
		tagExists, err := git.TagExists(vb.ParentDir, gitMeta.TagName)
		if err != nil {
			logFatal(vb.Options, fmt.Sprintf("Error checking for existing tag: %v\n", err))
		}
		if tagExists {
			logFatal(vb.Options, fmt.Sprintf("Tag '%s' already exists in the git repository. "+
				"Please bump to a different version or remove the existing tag.\n", gitMeta.TagName))
		}
	}

	// check if the Git repository has pending changes
	isDirty, _ := git.HasPendingChanges(vb.ParentDir)
	if isDirty {
		logFatal(vb.Options, "The Git repository has pending changes. Please commit or stash them before proceeding.")
	}

	// check if GPG signing is enabled for commits
	signKey, err := git.GetSigningKey(vb.ParentDir)
	if err != nil {
		logFatal(vb.Options, fmt.Sprintf("Error checking for GPG signing key: %v\n", err))
	}
	signByDefault, err := git.IsSigningEnabled(vb.ParentDir)
	if err != nil {
		logFatal(vb.Options, fmt.Sprintf("Error checking if GPG signing is enabled: %v\n", err))
	}
	if signByDefault || vb.Config.GitSign {
		logVerbose(vb.Options, "GPG signing of git commits is enabled. Checking configuration...")
	}
	// sanity check signing key
	if (signByDefault || vb.Config.GitSign) && signKey == "" {
		logFatal(vb.Options, "GPG signing of git commits is enabled but no signing key is configured. "+
			"Please configure a signing key in git.")
	}
	if signByDefault && !vb.Config.GitSign {
		logWarning(vb.Options, "GPG signing of git commits is enabled by default in the git configuration. "+
			"Consider enabling GPG signing in the VersionBump configuration.")
	}
	if signByDefault || vb.Config.GitSign {
		logVerbose(vb.Options, fmt.Sprintf("Git commits will be signed with GPG key: %s", signKey))
	}

}

// preamble prints the version bump preamble.
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
		logFatal(vb.Options, fmt.Sprintf("Unable to get Git metadata: %v\n", err))
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
		err := git.CommitChanges(vb.ParentDir, gitMeta.CommitMessage, vb.Config.GitSign)
		if err != nil {
			fmt.Printf("Error committing changes: %v\n", err)
			os.Exit(1)
		}
		logVerbose(vb.Options, fmt.Sprintf("Committed changes with message: %s", gitMeta.CommitMessage))
	}
	if vb.Config.GitTag {
		logVerbose(vb.Options, "Tagging changes...")
		err := git.TagChanges(vb.ParentDir, gitMeta.TagName, gitMeta.TagMessage, vb.Config.GitSign)
		if err != nil {
			fmt.Printf("Error tagging changes: %v\n", err)
			os.Exit(1)
		}
		logVerbose(vb.Options,
			fmt.Sprintf(
				"Tag '%s' created with message: %s",
				gitMeta.TagName,
				gitMeta.TagMessage))
	}
}

// bumpPreflight performs a pre-flight check for the version bump operation.
func (vb *VersionBump) bumpPreflight() {
	if !vb.Options.IsResetVersion() {
		logVerbose(vb.Options, fmt.Sprintf("Bumping version part: %s", vb.Options.BumpPart))
	} else {
		logVerbose(vb.Options, fmt.Sprintf("Resetting version to: %s", vb.GetNewVersion()))
	}
	logVerbose(vb.Options, fmt.Sprintf("Will bump version %s --> %s", vb.GetOldVersion(), vb.GetNewVersion()))

	// log what changes will be made to each file
	for _, file := range vb.Config.Files {
		for _, replace := range file.Replace {
			find := vbu.ReplaceInString(replace, "{version}", vb.GetOldVersion())
			replace := vbu.ReplaceInString(replace, "{version}", vb.GetNewVersion())

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
				logFatal(vb.Options, fmt.Sprintf("No replacements found in file: %s\n", file.Path))
			}
		}
	}
}

// makeChanges updates the version in the files.
func (vb *VersionBump) makeChanges() {
	// at this point we have already checked the config and there are no errors
	for _, file := range vb.Config.Files {
		for _, replace := range file.Replace {
			find := vbu.ReplaceInString(replace, "{version}", vb.GetOldVersion())
			replace := vbu.ReplaceInString(replace, "{version}", vb.GetNewVersion())

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
		printColor(fmt.Sprintf("%s [y/N]: ", prompt), ColorBlue)
		input, err := reader.ReadString('\n')
		if err != nil {
			printColor("Error reading input. Please try again.", ColorYellow)
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
			logVerbose(config.Options{}, "Operation canceled by user.")
			return false
		} else {
			printColor("Invalid input. Please enter 'y' or 'n'.", ColorYellow)
		}
	}
}

func logWarning(opts config.Options, msg string) {
	printColorOpts(opts, fmt.Sprintf("WARNING: %s\n", msg), ColorYellow)
}

func logFatal(opts config.Options, msg string) {
	printColorOpts(opts, fmt.Sprintf("ERROR: %s\n", msg), ColorRed)
	os.Exit(1)
}

func logVerbose(opts config.Options, msg string) {
	if !opts.Quiet {
		printColorOpts(opts, fmt.Sprintf("%s\n", msg), ColorLightGray)
	}
}

const (
	ColorRed       = "red"
	ColorGreen     = "green"
	ColorYellow    = "yellow"
	ColorBlue      = "blue"
	ColorMagenta   = "magenta"
	ColorCyan      = "cyan"
	ColorWhite     = "white"
	ColorLightGray = "lightgray"
)

// printColor prints the given text in different colors based on the color code
func printColor(text string, color string) {

	var colorCode string

	// Define ANSI color codes
	switch color {
	case ColorRed:
		colorCode = "\033[31m"
	case ColorGreen:
		colorCode = "\033[32m"
	case ColorYellow:
		colorCode = "\033[33m"
	case ColorBlue:
		colorCode = "\033[34m"
	case ColorMagenta:
		colorCode = "\033[35m"
	case ColorCyan:
		colorCode = "\033[36m"
	case ColorWhite:
		colorCode = "\033[37m"
	case ColorLightGray:
		colorCode = "\033[38;5;246m"
	default:
		colorCode = "\033[0m" // Default color (white)
	}

	// Print the text with the chosen color
	fmt.Printf("%s%s\033[0m", colorCode, text)

}

// printColorOpts conditionally prints the given text in different colors based on the color code
func printColorOpts(opts config.Options, text string, color string) {
	if opts.NoColor {
		fmt.Print(text)
		return
	} else {
		printColor(text, color)
	}
}
