package main

import (
	"fmt"
	"os"

	"github.com/ptgoetz/go-versionbump/internal"
	vbc "github.com/ptgoetz/go-versionbump/internal/config"
	"github.com/ptgoetz/go-versionbump/pkg/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var opts vbc.Options

var rootCmd = &cobra.Command{
	Use:   "versionbump",
	Short: `VersionBump is a command-line tool designed to automate version string management in projects.`,
	Long:  `VersionBump is a command-line tool designed to automate version string management in projects.`,
	RunE:  runRootCmd, // Use RunE for better error handling
}

var majorCmd = &cobra.Command{
	Use:   semver.Major.String(),
	Short: `Bump the major version number (e.g. 1.2.3 -> 2.0.0).`,
	Long:  `Bump the major version number (e.g. 1.2.3 -> 2.0.0).`,
	RunE:  bumpMajor, // Use RunE for better error handling
}

var minorCmd = &cobra.Command{
	Use:   semver.Minor.String(),
	Short: `Bump the minor version number (e.g. 1.2.3 -> 1.3.0).`,
	Long:  `Bump the minor version number (e.g. 1.2.3 -> 1.3.0).`,
	RunE:  bumpMinor, // Use RunE for better error handling
}

var patchCmd = &cobra.Command{
	Use:   semver.Patch.String(),
	Short: `Bump the patch version number (e.g. 1.2.3 -> 1.2.4).`,
	Long:  `Bump the patch version number (e.g. 1.2.3 -> 1.2.4).`,
	RunE:  bumpPatch, // Use RunE for better error handling
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: `Show the effective configuration of the project.`,
	Long:  `Show the effective configuration of the project.`,
	RunE:  runConfigCmd, // Use RunE for better error handling
}

var resetCmd = &cobra.Command{
	Use:   "set <version>",
	Short: `Set the project version to the specified value.`,
	Long:  `Set the project version to the specified value. Value can be any valid semantic version string.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runResetCmd, // Use RunE for better error handling
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: `Initialize a new versionbump configuration file.`,
	Long:  `Initialize a new versionbump configuration file.`,
	RunE:  runInitCmd, // Use RunE for better error handling
}

var showCmd = &cobra.Command{
	Use:   "show [version]",
	Short: `Show potential versioning paths for the project version or a specific version.`,
	Long:  `Show potential versioning paths for the project version or a specific version.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vb, err := internal.NewVersionBump(opts)
		if err != nil {
			return err
		}
		versionStr := ""
		if len(args) > 0 {
			versionStr = args[0]
		}

		return vb.Show(versionStr)
	},
}

var showVersionCmd = &cobra.Command{
	Use:   "show-version",
	Short: `Show the current project version.`,
	Long:  `Show the current project version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vb, err := internal.NewVersionBump(opts)
		if err != nil {
			return err
		}
		vb.ShowVersion()
		return nil
	},
}

var showLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: `Show the latest project release version based on git tags.`,
	Long:  `Show the latest project release version based on git tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vb, err := internal.NewVersionBump(opts)
		if err != nil {
			return err
		}
		return vb.LatestVersion()
	},
}

var gitTagHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: `Show the sorted version history based on git tags.`,
	Long:  `Show the sorted version history based on git tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vb, err := internal.NewVersionBump(opts)
		if err != nil {
			return err
		}
		err = vb.GitTagHistory()
		return err
	},
}

var preReleaseNextCmd = &cobra.Command{
	Use:   semver.PreRelease.String(),
	Short: `Bump the next pre-release version label (e.g. 1.2.3-alpha -> 1.2.3-beta).`,
	Long:  `Bump the next pre-release version label (e.g. 1.2.3-alpha -> 1.2.3-beta).`,
	RunE:  bumpPreReleaseNext, // Use RunE for better error handling
}

var preReleaseMajorCmd = &cobra.Command{
	Use:   semver.PreReleaseMajor.String(),
	Short: `Bump the pre-release major version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.1).`,
	Long:  `Bump the pre-release major version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.1).`,
	RunE:  bumpPreReleaseMajor, // Use RunE for better error handling
}

var preReleaseMinorCmd = &cobra.Command{
	Use:   semver.PreReleaseMinor.String(),
	Short: `Bump the pre-release minor version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.1).`,
	Long:  `Bump the pre-release minor version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.1).`,
	RunE:  bumpPreReleaseMinor, // Use RunE for better error handling
}

var preReleasePatchCmd = &cobra.Command{
	Use:   semver.PreReleasePatch.String(),
	Short: `Bump the pre-release patch version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.0.1).`,
	Long:  `Bump the pre-release patch version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.0.1).`,
	RunE:  bumpPreReleasePatch, // Use RunE for better error handling
}

var preReleaseNewMajorCmd = &cobra.Command{
	Use:   semver.PreReleaseNewMajor.String(),
	Short: `Bump the major version and apply the first pre-release label (e.g. 1.2.3 -> 2.0.0-alpha).`,
	Long:  `Bump the major version and apply the first pre-release label (e.g. 1.2.3 -> 2.0.0-alpha).`,
	RunE:  bumpNewPreReleaseMajor, // Use RunE for better error handling
}

var preReleaseNewMinorCmd = &cobra.Command{
	Use:   semver.PreReleaseNewMinor.String(),
	Short: `Bump the minor version and apply the first pre-release label (e.g. 1.2.3 -> 1.3.0-alpha).`,
	Long:  `Bump the minor version and apply the first pre-release label (e.g. 1.2.3 -> 1.3.0-alpha).`,
	RunE:  bumpNewPreReleaseMinor, // Use RunE for better error handling
}

var preReleaseNewPatchCmd = &cobra.Command{
	Use:   semver.PreReleaseNewPatch.String(),
	Short: `Bump the patch version and apply the first pre-release label (e.g. 1.2.3 -> 1.2.4-alpha).`,
	Long:  `Bump the patch version and apply the first pre-release label (e.g. 1.2.3 -> 1.2.4-alpha).`,
	RunE:  bumpNewPreReleasePatch, // Use RunE for better error handling
}

var releaseCmd = &cobra.Command{
	Use:   semver.Release.String(),
	Short: `Bump the pre-release version to a release version (e.g. 1.2.3-alpha -> 1.2.3).`,
	Long:  `Bump the pre-release version to a release version (e.g. 1.2.3-alpha -> 1.2.3).`,
	RunE:  bumpRelease, // Use RunE for better error handling
}

var preReleaseBuildCmd = &cobra.Command{
	Use:   semver.PreReleaseBuild.String(),
	Short: `Bump the pre-release build version number (e.g. 1.2.3 -> 1.2.3+build.1).`,
	Long:  `Bump the pre-release build version number (e.g. 1.2.3 -> 1.2.3+build.1).`,
	RunE:  bumpPreReleaseBuild, // Use RunE for better error handling
}

func init() {
	rootCmd.Flags().BoolVarP(&opts.ShowVersion, "version", "V", false, "Show the VersionBump version and exit.")

	commonFlags := pflag.NewFlagSet("common", pflag.ExitOnError)
	commonFlags.StringVarP(&opts.ConfigPath, "config", "c", "versionbump.yaml", "The path to the configuration file")
	commonFlags.BoolVar(&opts.NoPrompt, "no-prompt", false, "Don't prompt the user for confirmation before making changes.")
	commonFlags.BoolVarP(&opts.Quiet, "quiet", "q", false, "Don't print verbose output.")
	commonFlags.BoolVar(&opts.NoGit, "no-git", false, "Don't perform any git operations.")
	commonFlags.BoolVar(&opts.NoColor, "no-color", false, "Disable color output.")

	configColorFlags := pflag.NewFlagSet("config-color", pflag.ExitOnError)
	configColorFlags.StringVarP(&opts.ConfigPath, "config", "c", "versionbump.yaml", "The path to the configuration file")
	configColorFlags.BoolVar(&opts.NoColor, "no-color", false, "Disable color output.")

	commonFlags.AddFlagSet(configColorFlags)

	initFlags := pflag.NewFlagSet("init", pflag.ExitOnError)
	initFlags.StringVarP(&opts.InitOpts.File, "file", "f", "versionbump.yaml", "The name of the configuration file to create.")
	initFlags.BoolVar(&opts.InitOpts.NoInteractive, "no-interactive", false, "Don't prompt for interactive input.")
	initCmd.Flags().AddFlagSet(initFlags)

	prereleaserFlags := pflag.NewFlagSet("prelease", pflag.ExitOnError)
	prereleaserFlags.AddFlagSet(commonFlags)

	preReleaseNextCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleaseMajorCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleaseMinorCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleasePatchCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleaseBuildCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleaseNewMajorCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleaseNewMinorCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleaseNewPatchCmd.Flags().AddFlagSet(prereleaserFlags)
	releaseCmd.Flags().AddFlagSet(prereleaserFlags)

	showCmd.Flags().AddFlagSet(configColorFlags)
	showVersionCmd.Flags().AddFlagSet(commonFlags)
	showLatestCmd.Flags().AddFlagSet(commonFlags)
	configCmd.Flags().AddFlagSet(configColorFlags)

	majorCmd.Flags().AddFlagSet(commonFlags)
	minorCmd.Flags().AddFlagSet(commonFlags)
	patchCmd.Flags().AddFlagSet(commonFlags)
	resetCmd.Flags().AddFlagSet(commonFlags)

	rootCmd.AddCommand(majorCmd)
	rootCmd.AddCommand(minorCmd)
	rootCmd.AddCommand(patchCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(preReleaseNextCmd)
	rootCmd.AddCommand(preReleaseMajorCmd)
	rootCmd.AddCommand(preReleaseMinorCmd)
	rootCmd.AddCommand(preReleasePatchCmd)
	rootCmd.AddCommand(preReleaseBuildCmd)
	rootCmd.AddCommand(preReleaseNewMajorCmd)
	rootCmd.AddCommand(preReleaseNewMinorCmd)
	rootCmd.AddCommand(preReleaseNewPatchCmd)
	rootCmd.AddCommand(releaseCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(showVersionCmd)
	rootCmd.AddCommand(showLatestCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(gitTagHistoryCmd)
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	if opts.ShowVersion {
		fmt.Println(internal.Version)
		return nil
	}
	return cmd.Help()
}

func bumpMajor(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.Major)
}

func bumpMinor(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.Minor)
}

func bumpPatch(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.Patch)
}

func bumpPreReleaseNext(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreRelease)
}

func bumpPreReleaseMajor(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreReleaseMajor)
}

func bumpPreReleaseMinor(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreReleaseMinor)
}

func bumpPreReleasePatch(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreReleasePatch)
}

func bumpPreReleaseBuild(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreReleaseBuild)
}

func bumpNewPreReleaseMajor(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreReleaseNewMajor)
}

func bumpNewPreReleaseMinor(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreReleaseNewMinor)
}

func bumpNewPreReleasePatch(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.PreReleaseNewPatch)
}

func bumpRelease(cmd *cobra.Command, args []string) error {
	return runVersionBump(semver.Release)
}

func runResetCmd(cmd *cobra.Command, args []string) error {
	opts.ResetVersion = args[0]
	vb, err := internal.NewVersionBump(opts)
	if err != nil {
		return err
	}

	vb.Run()
	return nil
}

func runInitCmd(cmd *cobra.Command, args []string) error {
	return internal.InitVersionBumpProject(opts)
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	vb, err := internal.NewVersionBump(opts)
	if err != nil {
		return err
	}

	return vb.ShowEffectiveConfig()
}

// runVersionBump contains the logic for executing the version bump process
func runVersionBump(bumpPart semver.BumpStrategy) error {
	opts.BumpPart = bumpPart

	vb, err := internal.NewVersionBump(opts)
	if err != nil {
		return err
	}

	// Run the version bump process
	vb.Run()
	return nil
}
