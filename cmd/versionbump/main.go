package main

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal"
	vbc "github.com/ptgoetz/go-versionbump/internal/config"
	"github.com/ptgoetz/go-versionbump/internal/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
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
	Use:   "major",
	Short: `Bump the major version number (e.g. 1.2.3 -> 2.0.0).`,
	Long:  `Bump the major version number (e.g. 1.2.3 -> 2.0.0).`,
	RunE:  bumpMajor, // Use RunE for better error handling
}

var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: `Bump the minor version number (e.g. 1.2.3 -> 1.3.0).`,
	Long:  `Bump the minor version number (e.g. 1.2.3 -> 1.3.0).`,
	RunE:  bumpMinor, // Use RunE for better error handling
}

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: `Bump the patch version number (e.g. 1.2.3 -> 1.2.4).`,
	Long:  `Bump the patch version number (e.g. 1.2.3 -> 1.2.4).`,
	RunE:  bumpPatch, // Use RunE for better error handling
}

var preReleaseNextcCmd = &cobra.Command{
	Use:   "patch",
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
	Use:   "reset <version>",
	Short: `Reset the project version to the specified value.`,
	Long:  `Reset the project version to the specified value. Value can be any valid semantic version string.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runResetCmd, // Use RunE for better error handling
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

var preReleaseMajorCmd = &cobra.Command{
	Use:   "prerelease-major",
	Short: `Bump the pre-release major version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.1).`,
	Long:  `ump the pre-release major version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.1).`,
	RunE:  bumpPreReleaseMajor, // Use RunE for better error handling
}

var preReleaseMinorCmd = &cobra.Command{
	Use:   "prerelease-minor",
	Short: `Bump the pre-release minor version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.1).`,
	Long:  `ump the pre-release minor version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.1).`,
	RunE:  bumpPreReleaseMinor, // Use RunE for better error handling
}

var preReleasePatchCmd = &cobra.Command{
	Use:   "prerelease-patch",
	Short: `Bump the pre-release patch version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.0.1).`,
	Long:  `ump the pre-release patch version number (e.g. 1.2.3-alpha -> 1.2.3-alpha.0.0.1).`,
	RunE:  bumpPreReleasePatch, // Use RunE for better error handling
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

	prereleaserFlags := pflag.NewFlagSet("prelease", pflag.ExitOnError)
	prereleaserFlags.AddFlagSet(commonFlags)

	preReleaseMajorCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleaseMinorCmd.Flags().AddFlagSet(prereleaserFlags)
	preReleasePatchCmd.Flags().AddFlagSet(prereleaserFlags)

	showCmd.Flags().AddFlagSet(configColorFlags)
	configCmd.Flags().AddFlagSet(configColorFlags)

	majorCmd.Flags().AddFlagSet(commonFlags)
	minorCmd.Flags().AddFlagSet(commonFlags)
	patchCmd.Flags().AddFlagSet(commonFlags)
	resetCmd.Flags().AddFlagSet(commonFlags)

	rootCmd.AddCommand(majorCmd)
	rootCmd.AddCommand(minorCmd)
	rootCmd.AddCommand(patchCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(preReleaseMajorCmd)
	rootCmd.AddCommand(preReleaseMinorCmd)
	rootCmd.AddCommand(preReleasePatchCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(configCmd)

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
	return runVersionBump(semver.PreReleaseNext)
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

func runResetCmd(cmd *cobra.Command, args []string) error {
	opts.ResetVersion = args[0]
	vb, err := internal.NewVersionBump(opts)
	if err != nil {
		return err
	}

	vb.Run()
	return nil
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	vb, err := internal.NewVersionBump(opts)
	if err != nil {
		return err
	}

	return vb.ShowEffectiveConfig()
}

// runVersionBump contains the logic for executing the version bump process
func runVersionBump(bumpPart semver.VersionPart) error {
	opts.BumpPart = bumpPart

	vb, err := internal.NewVersionBump(opts)
	if err != nil {
		return err
	}

	// Run the version bump process
	vb.Run()
	return nil
}
