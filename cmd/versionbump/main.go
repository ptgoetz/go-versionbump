package main

import (
	"bufio"
	"flag"
	"fmt"
	vbc "github.com/ptgoetz/go-versionbump/internal/config"
	"github.com/ptgoetz/go-versionbump/internal/git"
	vbu "github.com/ptgoetz/go-versionbump/internal/utils"
	vbv "github.com/ptgoetz/go-versionbump/internal/version"
	"os"
	"path"
	"strings"
)


var opts vbc.Options

func main() {
	// Define command-line flags
	// TODO Migrate to Cobra
	flag.BoolVar(&opts.ShowVersion, "V", false, "Show the version of Config and exit.")
	flag.StringVar(&opts.ConfigPath, "config", "versionbump.yaml", "The path to the configuration file")
	flag.BoolVar(&opts.DryRun, "dry-run", false, "Dry run. Don't change anything, just report what would "+
		"be changed")
	flag.BoolVar(&opts.NoPrompt, "no-prompt", false, "Don't prompt the user for confirmation before making changes.")
	flag.BoolVar(&opts.Quiet, "quiet", false, "Don't print verbose output.")
	flag.StringVar(&opts.ResetVersion, "reset", "", "Reset the version to the specified value.")
	flag.BoolVar(&opts.NoGit, "no-git", false, "Don't perform any git operations.")
	flag.Parse()
	args := flag.Args()

	// TODO: The trash below badly needs refactoring

	// print the version and exit
	if opts.ShowVersion {
		fmt.Println(vbv.VersionBumpVersion)
		os.Exit(0)
	}

	if len(args) != 1 && opts.ResetVersion == "" {
		fmt.Println("ERROR: no version part specified.")
		flag.Usage()
		os.Exit(1)
	}

	// Log the version and configuration path
	logVerbose(vbv.VersionBumpVersion)
	logVerbose(fmt.Sprintf("Config path: %s", opts.ConfigPath))

	// Load the configuration file
	config, root, err := vbc.LoadConfig(opts.ConfigPath)
	if err != nil {
		fmt.Printf("Error loading configuration file: '%s' %v\n", opts.ConfigPath, err)
		os.Exit(1)
	}

	// Perform a git pre-flight check
	gitPreFlight(root, config)

	// Log the pre-flight information
	logVerbose(fmt.Sprintf("Project root: %s", root))
	logVerbose(fmt.Sprintf("Current Version: %s", config.Version))

	// Log the files that will be updated
	logVerbose("Tracked Files:")
	for _, file := range config.Files {
		logVerbose(fmt.Sprintf("  - %s", file.Path))
	}

	if opts.ResetVersion != "" {
		if !vbv.ValidateVersion(opts.ResetVersion) {
			fmt.Printf("ERROR: Invalid semantic version for reset: %s\n", opts.ResetVersion)
			os.Exit(1)
		}
		parsedVersion, _ := vbv.ParseVersion(opts.ResetVersion)

		performReset(config, root, parsedVersion)
		os.Exit(0)
	}

	var bumpMetadata *vbc.GitMeta
	if len(args) == 1 { // We got a version bump request.
		if !vbv.ValidateVersionPart(args[0]) {
			fmt.Printf("ERROR: Invalid version bump part: %s\n", args[0])
			os.Exit(1)
		}
		bumpMetadata, err = config.BumpAndGetMetaData(args[0])
		if err != nil {
			fmt.Printf("ERROR: Unable to bump version: %v\n", err)
			os.Exit(1)
		}
		changePreFlight(root, config, args)
	} else {
		fmt.Println("ERROR: no version bump part specified.")
		os.Exit(1)
	}

	if !opts.DryRun {
		proceed := true
		if !opts.NoPrompt {
			fmt.Println("The following files will be updated:")
			for _, file := range config.Files {
				fmt.Printf("  - %s\n", file.Path)
			}
			proceed = promptUserConfirmation("Do you want to proceed with the changes?")
		}
		if proceed {
			changePreFlight(root, config, args)
			makeChanges(root, config, args[0])
			gitCommit(bumpMetadata, root, config)
		}
	}
}

func performReset(config *vbc.Config, root string, resetVersion *vbv.Version) {
	// Log the reset information
	logVerbose(fmt.Sprintf("Resetting version to: %s", resetVersion.String()))
	if promptUserConfirmation("Do you want to proceed with the reset?") {
		makeChanges(root, config, "")
		logVerbose("\nVersion reset complete. No git operations were performed.")
		logVerbose("Please verify you have the correct version set and commit the changes manually.")
	}
}

func gitCommit(bumpMetadata *vbc.GitMeta, root string, config *vbc.Config) {
	if opts.NoGit {
		return
	}
	if config.IsGitRequired() && !opts.NoPrompt {
		logVerbose(bumpMetadata.String())
		if !opts.NoPrompt {
			proceed := promptUserConfirmation("Do you want to commit the changes to the git repository?")
			if !proceed {
				os.Exit(1)
			}
		}
	}
	// commit changes
	if config.GitCommit {
		logVerbose("Committing changes...")
		err := git.CommitChanges(root, bumpMetadata.CommitMessage)
		if err != nil {
			fmt.Printf("Error committing changes: %v\n", err)
			os.Exit(1)
		}
		logVerbose(fmt.Sprintf("Committed changes with message: %s", bumpMetadata.CommitMessage))
	}
	if config.GitTag {
		logVerbose("Tagging changes...")
		err := git.TagChanges(root, bumpMetadata.TagName, bumpMetadata.TagMessage)
		if err != nil {
			fmt.Printf("Error tagging changes: %v\n", err)
			os.Exit(1)
		}
		logVerbose(fmt.Sprintf("Tagged '%s' created with message: %s", bumpMetadata.TagName, bumpMetadata.TagMessage))
	}
}

func gitPreFlight(root string, config *vbc.Config) {
	if opts.NoGit {
		return
	}
	if config.IsGitRequired() {
		isGitAvalable, version := git.IsGitAvailable()
		if !isGitAvalable {
			fmt.Printf("ERROR: Git is required by the configuration but not available. " +
				"Config requires Git to be installed and available in the system PATH.")
			os.Exit(1)
		} else {
			logVerbose(fmt.Sprintf("Git version: %s", strings.TrimSpace(version)[12:]))
		}
	}
	isGitRepo, err := git.IsGitRepository(root)
	if err != nil {
		fmt.Printf("Error checking for git repository: %v\n", err)
		os.Exit(1)
	}
	if !isGitRepo {
		fmt.Println("ERROR: The project root is not a Git repository, but Git options are enabled in the " +
			"configuration file.")
		if opts.NoPrompt {
			os.Exit(1)
		}
		if promptUserConfirmation("Do you want to initialize a Git repository in this directory?") {
			err := git.InitializeGitRepo(root)
			if err != nil {
				fmt.Printf("Error initializing Git repository: %v\n", err)
				os.Exit(1)
			}
		}
	}
	isDirty, _ := git.HasPendingChanges(root)
	if isDirty {
		fmt.Println("ERROR: The Git repository has pending changes. Please commit or stash them before proceeding.")
		os.Exit(1)
	}
}

func changePreFlight(root string, config *vbc.Config, args []string) {
	// parse the current version:
	curVersion, err := vbv.ParseVersion(config.Version)
	if err != nil {
		fmt.Printf("Failed to parse semantic version string: %s\n", config.Version)
		os.Exit(1)
	}
	currentVersionStr := curVersion.String()
	// bump version
	err = curVersion.StringBump(args[0])
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	nextVersionStr := curVersion.String()
	logVerbose(fmt.Sprintf("Bumping version part: %s", args[0]))
	logVerbose(fmt.Sprintf("Will bump version %s --> %s", currentVersionStr, nextVersionStr))

	// log what changes will be made to each file
	for _, file := range config.Files {
		find := vbu.ReplaceInString(file.Replace, "{version}", currentVersionStr)
		replace := vbu.ReplaceInString(file.Replace, "{version}", nextVersionStr)

		logVerbose(file.Path)
		logVerbose(fmt.Sprintf("     Find: \"%s\"", find))
		logVerbose(fmt.Sprintf("  Replace: \"%s\"", replace))
		count, err := vbu.CountStringsInFile(path.Join(root, file.Path), find)
		if err != nil {
			fmt.Println(fmt.Errorf("error getting replacement count: a%v", err))
			os.Exit(1)
		}
		if count > 0 {
			logVerbose(fmt.Sprintf("    Found %d replacement(s)", count))
		} else {
			fmt.Println("ERROR: No replacements found in file: ", file.Path)
			os.Exit(1)
		}
	}
}

func makeChanges(root string, config *vbc.Config, versionPart string) {
	// at this point we have already checked the config and there are no errors
	var currentVersionStr, nextVersionStr string
	if versionPart != "" {
		curVersion, _ := vbv.ParseVersion(config.Version)
		currentVersionStr = curVersion.String()
		// bump version
		_ = curVersion.StringBump(versionPart)
		nextVersionStr = curVersion.String()
	} else { // reset
		cv, _ := vbv.ParseVersion(config.Version)
		currentVersionStr = cv.String()
		nv, _ := vbv.ParseVersion(opts.ResetVersion)
		nextVersionStr = nv.String()
	}

	for _, file := range config.Files {
		find := vbu.ReplaceInString(file.Replace, "{version}", currentVersionStr)
		replace := vbu.ReplaceInString(file.Replace, "{version}", nextVersionStr)

		if !opts.DryRun {
			var resolvedPath string
			if path.IsAbs(file.Path) {
				resolvedPath = file.Path
			} else {
				resolvedPath = path.Join(root, file.Path)
			}
			err := vbu.ReplaceInFile(resolvedPath, find, replace)
			if err != nil {
				fmt.Println(fmt.Errorf("error updating file %s: a%v", file.Path, err))
				os.Exit(1)
			}
			logVerbose(fmt.Sprintf("Updated file: %s", file.Path))
		}
	}
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

func logVerbose(msg string) {
	if opts.DryRun || !opts.Quiet {
		fmt.Println(msg)
	}
}
