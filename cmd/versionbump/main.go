package main

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal"
	vbc "github.com/ptgoetz/go-versionbump/internal/config"
	vbv "github.com/ptgoetz/go-versionbump/internal/version"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var opts vbc.Options

func main() {
	var rootCmd = &cobra.Command{
		Use:   "versionbump [bump-part]",
		Short: "VersionBump is a tool for managing version bumps",
		Long:  `VersionBump is a tool for managing version bumps with optional git integration.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.ShowVersion {
				fmt.Println(vbv.VersionBumpVersion)
				os.Exit(0)
			}

			if len(args) > 0 {
				opts.BumpPart = args[0]
			}
			if len(args) == 0 && opts.ResetVersion == "" {
				fmt.Println("ERROR: no version part specified.")
				cmd.Usage()
				os.Exit(1)
			}

			vb, err := internal.NewVersionBump(opts)
			if err != nil {
				log.Fatal(err)
			}
			vb.Run()
		},
	}

	rootCmd.Flags().BoolVarP(&opts.ShowVersion, "version", "V", false, "Show the version of Config and exit.")
	rootCmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "versionbump.yaml", "The path to the configuration file")
	rootCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Dry run. Don't change anything, just report what would be changed")
	rootCmd.Flags().BoolVar(&opts.NoPrompt, "no-prompt", false, "Don't prompt the user for confirmation before making changes.")
	rootCmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "Don't print verbose output.")
	rootCmd.Flags().StringVar(&opts.ResetVersion, "reset", "", "Reset the version to the specified value.")
	rootCmd.Flags().BoolVar(&opts.NoGit, "no-git", false, "Don't perform any git operations.")
	rootCmd.Flags().BoolVar(&opts.NoColor, "no-color", false, "Disable color output.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
