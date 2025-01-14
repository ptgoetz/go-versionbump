package internal

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ptgoetz/go-versionbump/internal/config"
	"github.com/ptgoetz/go-versionbump/internal/git"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	vbu "github.com/ptgoetz/go-versionbump/internal/utils"
	"github.com/ptgoetz/go-versionbump/pkg/semver"
	"gopkg.in/yaml.v2"
)

const Version = "1.0.0-alpha"

// VersionBump represents the Version bump operation.
type VersionBump struct {
	Config    config.Config
	Options   config.Options
	ParentDir string
}

// NewVersionBump creates a new VersionBump instance.
func NewVersionBump(options config.Options) (*VersionBump, error) {

	cfg, parentDir, err := config.LoadConfig(options.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error loading configuration file: %v", err)
	}

	vb := &VersionBump{
		Config:    *cfg,
		Options:   options,
		ParentDir: parentDir,
	}

	return vb, nil
}

func (vb *VersionBump) GetOldVersion() string {
	if !semver.ValidateSemVersion(vb.Config.Version) {
		logFatal(vb.Options, fmt.Sprintf("Failed to parse semantic version string for old version: %s", vb.Config.Version))
	}
	oldVersion, _ := semver.ParseSemVersion(vb.Config.Version)
	return oldVersion.String()
}

func (vb *VersionBump) GetNewVersion() string {
	if vb.Options.IsResetVersion() {
		v, err := semver.ParseSemVersion(vb.Options.ResetVersion)
		if err != nil {
			logFatal(vb.Options, fmt.Sprintf("Failed to parse semantic version string for reset version: %s", vb.Options.ResetVersion))
		}
		return v.String()
	}
	oldVersionStr := vb.GetOldVersion()
	oldVersion, _ := semver.ParseSemVersion(oldVersionStr)
	newVersion, _ := oldVersion.Bump(vb.Options.BumpPart, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
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

func (vb *VersionBump) ShowVersion() {
	fmt.Println(vb.Config.Version)
}

func checkBumpError(vb *VersionBump, v *semver.SemanticVersion, err error) string {
	if err != nil {
		logWarning(vb.Options, err.Error())
		return "❌"
	} else {
		return v.String()
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
	curVersion, err := semver.ParseSemVersion(curVersionStr)
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
	majorVersion, err := curVersion.Bump(semver.Major, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	majorVersionStr := checkBumpError(vb, majorVersion, err)
	minorVersion, err := curVersion.Bump(semver.Minor, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	minorVersionStr := checkBumpError(vb, minorVersion, err)
	patchVersion, err := curVersion.Bump(semver.Patch, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	patchVersionStr := checkBumpError(vb, patchVersion, err)

	releaseVersion, err := curVersion.Bump(semver.Release, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	releaseVersionStr := checkBumpError(vb, releaseVersion, err)

	prNextVersion, err := curVersion.Bump(semver.PreRelease, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prNextVersionStr := checkBumpError(vb, prNextVersion, err)

	prMajorVersion, err := curVersion.Bump(semver.PreReleaseMajor, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prMajorVersionStr := checkBumpError(vb, prMajorVersion, err)
	prMinorVersion, err := curVersion.Bump(semver.PreReleaseMinor, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prMinorVersionStr := checkBumpError(vb, prMinorVersion, err)
	prPatchVersion, err := curVersion.Bump(semver.PreReleasePatch, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prPatchVersionStr := checkBumpError(vb, prPatchVersion, err)

	prNewMajorVersion, err := curVersion.Bump(semver.PreReleaseNewMajor, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prNewMajorVersionStr := checkBumpError(vb, prNewMajorVersion, err)
	prNewMinorVersion, err := curVersion.Bump(semver.PreReleaseNewMinor, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prNewMinorVersionStr := checkBumpError(vb, prNewMinorVersion, err)
	prNewPatchVersion, err := curVersion.Bump(semver.PreReleaseNewPatch, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prNewPatchVersionStr := checkBumpError(vb, prNewPatchVersion, err)

	prBuildVersion, err := curVersion.Bump(semver.PreReleaseBuild, vb.Config.PreReleaseLabels, vb.Config.BuildLabel)
	prBuildVersionStr := checkBumpError(vb, prBuildVersion, err)

	if prNextVersionStr == "" {
		prNextVersionStr = "❌"
	}

	padLen := len(curVersion.String())
	padding := utils.PaddingString(padLen, " ")

	tree := fmt.Sprintf(
		`%s ─┬─ major ─ %s
  %s├─ minor ─ %s
  %s├─ patch ─ %s
  %s├─ release ─ %s
  %s├─ new-pre-major ─ %s
  %s├─ new-pre-minor ─ %s
  %s├─ new-pre-patch ─ %s
  %s├─ pre ─ %s
  %s├─ pre-major ─ %s
  %s├─ pre-minor ─ %s
  %s├─ pre-patch ─ %s
  %s╰─ pre-build ─ %s
`,
		curVersion.String(),
		majorVersionStr,
		padding,
		minorVersionStr,
		padding,
		patchVersionStr,
		padding,
		releaseVersionStr,
		padding,
		prNewMajorVersionStr,
		padding,
		prNewMinorVersionStr,
		padding,
		prNewPatchVersionStr,
		padding,
		prNextVersionStr,
		padding,
		prMajorVersionStr,
		padding,
		prMinorVersionStr,
		padding,
		prPatchVersionStr,
		padding,
		prBuildVersionStr)

	printColorOpts(vb.Options, tree, ColorLightBlue)
	return nil
}

func (vb *VersionBump) GitTagHistory() error {
	if vb.Options.NoGit {
		return nil
	}
	logVerbose(vb.Options, "version History:")
	versions, err := vb.GetSortedVersions()
	if err != nil {
		return err
	}
	for _, version := range versions {
		logVerbose(vb.Options, fmt.Sprintf("  - %s", version.String()))
	}
	return nil
}

func (vb *VersionBump) LatestVersion() error {
	versions, err := vb.GetSortedVersions()
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		return fmt.Errorf("no versions found")
	}
	fmt.Println(versions[0].String())
	return nil
}

func (vb *VersionBump) GetSortedVersions() ([]*semver.SemanticVersion, error) {
	tags, err := git.GetTags(vb.ParentDir)
	if err != nil {
		return nil, err
	}
	versions := make([]*semver.SemanticVersion, 0)
	for _, tag := range tags {
		vStr, err := ExtractVersion(vb.Config.GitTagTemplate, tag)
		if err != nil {
			logVerbose(vb.Options, fmt.Sprintf("Error extracting version from tag: %s", err.Error()))
			continue
		}
		v, err := semver.ParseSemVersion(vStr)
		if err == nil {
			versions = append(versions, v)
		} else {
			logVerbose(vb.Options, fmt.Sprintf("Error parsing tag: %s", tag))
		}

	}
	semver.SortVersions(versions)
	return versions, nil
}

func ExtractVersion(template, value string) (string, error) {
	version := "{new}"
	idx := strings.Index(template, version)
	idx2 := idx + len(version)
	right := template[idx2:]
	left := template[:idx]
	if !strings.Contains(value, left) || !strings.Contains(value, right) {
		return "", fmt.Errorf("value '%s' does not match template '%s'", value, template)
	}

	proposed := value[idx : len(value)-len(right)]
	newValue := strings.ReplaceAll(template, version, proposed)
	if len(newValue) != len(value) {
		return "", fmt.Errorf("value '%s' does not match template '%s'", value, template)
	}

	return proposed, nil
}

func (vb *VersionBump) ShowEffectiveConfig() error {
	logVerbose(vb.Options, fmt.Sprintf("Config file: %s", vb.Options.ConfigPath))
	logVerbose(vb.Options, fmt.Sprintf("Project root: %s", vb.ParentDir))
	logVerbose(vb.Options, "Effective Configuration YAML:")

	conf := &vb.Config
	b, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	printColorOpts(vb.Options, string(b), ColorLightBlue)
	return nil
}

func InitVersionBumpProject(opts config.Options) error {
	// check to see if a configuration file already exists
	if utils.FileExists(opts.InitOpts.File) {
		return fmt.Errorf("configuration file already exists: %s", opts.InitOpts.File)
	}

	conf := config.NewConfig()

	// prompt the user for pre-release labels
	prLabels := promptUserForValue(
		"Enter pre-release labels (comma-separated)",
		strings.Join(conf.PreReleaseLabels, ","),
		semver.ValidatePreReleaseLabelsString)
	conf.PreReleaseLabels = strings.Split(prLabels, ",")

	// prompt the user for build label
	buildLabel := promptUserForValue(
		"Enter build label",
		conf.BuildLabel,
		semver.ValidateBuildLabel)
	conf.BuildLabel = buildLabel

	// prompt the user for the initial version
	initVersionStr := promptUserForValue("Enter the initial version", "0.0.0", semver.ValidateSemVersion)
	conf.Version = initVersionStr

	gitAvail, _ := git.IsGitAvailable()
	if gitAvail {
		if promptUserConfirm("Git is installed on this system. \nDo you want to enable Git features?") {
			conf.GitCommit = promptUserConfirm("Do you want to enable Git commit feature?")
			if conf.GitCommit {
				conf.GitTag = promptUserConfirm("Do you want to enable the Git tag feature?")
			}
		}
	}

	//
	tmpl, err := template.New("yaml").Parse(config.DefaultConfigTemplate)
	if err != nil {
		panic(err)
	}

	// create the configuration file
	f, err := os.Create(opts.InitOpts.File)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, conf)
	if err != nil {
		panic(err)
	}
	//
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
		if promptUserConfirm("The project directory is not a git repository.\nDo you want to initialize a git repository in the project directory?") {
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

// preamble prints the Version bump preamble.
func (vb *VersionBump) preamble() {
	logVerbose(vb.Options, fmt.Sprintf("VersionBump %s", Version))
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
		proceed := promptUserConfirm("Do you want to commit the changes to the git repository?")
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

// bumpPreflight performs a pre-flight check for the Version bump operation.
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

// makeChanges updates the Version in the files.
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
		if !promptUserConfirm("Proceed with the changes?") {
			os.Exit(0)
		}
	}
	return true
}

// promptUserConfirm prompts the user with the given prompt string and expects 'y' or 'n' input.
// It returns true for 'y' and false for 'n'.
func promptUserConfirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Print the prompt and read the user's input
		printColor(fmt.Sprintf("%s [y/N]: ", prompt), ColorLightBlue)
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
			printColor("Invalid input. Please enter 'y' or 'n'.\n", ColorYellow)
		}
	}
}

// promptUserForValue prompts the user for a value with the given prompt string.
// It returns the default value if the user input is empty.
// The validator function is used to validate the user input.
func promptUserForValue(prompt string, defaultValue string, validator func(string) bool) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Print the prompt and read the user's input
		printColor(fmt.Sprintf("%s [%s]: ", prompt, defaultValue), ColorLightBlue)
		input, err := reader.ReadString('\n')
		if err != nil {
			printColor("Error reading input. Please try again.\n", ColorYellow)
			continue
		}

		// Trim the input
		input = strings.TrimSpace(input)
		if input == "" {
			return defaultValue
		}

		// Validate the input
		if validator(input) {
			return input
		} else {
			printColor("Invalid input. Please try again.\n", ColorYellow)
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
	ColorLightBlue = "light-blue"
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
	case ColorLightBlue:
		colorCode = "\033[38;5;117m" // Light blue using extended 256 color code
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
