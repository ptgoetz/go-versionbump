package main

import (
	"flag"
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal"
	vbc "github.com/ptgoetz/go-versionbump/internal/config"
	vbv "github.com/ptgoetz/go-versionbump/internal/version"
	"log"
	"os"
)

func main() {
	// Define command-line flags
	var opts vbc.Options
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
	//fmt.Printf("OS.ARGS: %V\n", args)

	// print the version and exit
	if opts.ShowVersion {
		fmt.Println(vbv.VersionBumpVersion)
		os.Exit(0)
	}

	if len(args) > 0 {
		opts.BumpPart = args[0]
	}
	if len(args) == 0 && opts.ResetVersion == "" {
		fmt.Println("ERROR: no version part specified.")
		flag.Usage()
		os.Exit(1)
	}

	// Create a new VersionBump instance
	vb, err := internal.NewVersionBump(opts)
	if err != nil {
		log.Fatal(err)
	}
	vb.Run()
}
