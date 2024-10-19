package semver

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	VersionMajor = iota
	VersionMinor
	VersionPatch
	PreReleaseNext
	PreReleaseMajor
	PreReleaseMinor
	PreReleasePatch
	PreReleaseBuild
)

type SemVersion struct {
	Version                *Version
	PreReleaseVersionLabel string
	PreReleaseVersion      *PreReleaseVersion
	Build                  int
}

// String returns the version string
func (v *SemVersion) String() string {
	version := v.Version.String()
	preReleaseStr := v.PreReleaseVersion.String()
	if v.PreReleaseVersion != nil && preReleaseStr != "" {
		version += "-" + v.PreReleaseVersion.String()
	}
	if v.Build != 0 {
		version += "+" + strconv.Itoa(v.Build)
	}
	return version
}

// ParseSemVersion parses a semantic version string and returns a new SemVersion instance
func ParseSemVersion(versionStr string) (*SemVersion, error) {
	isPreRelease := strings.Index(versionStr, "-") != -1
	isBuild := strings.Index(versionStr, "+") != -1
	isPreReleaseBuild := isPreRelease && isBuild
	fmt.Printf("isPreRelease: %v, isBuild: %v, isPreReleaseBuild: %v\n", isPreRelease, isBuild, isPreReleaseBuild)

	var rootPart string
	var preReleasePart string
	var buildPart string
	if !isPreRelease && !isBuild {
		rootPart = versionStr
	}
	if isPreRelease /* && !isBuild */ {
		parts := strings.Split(versionStr, "-")
		rootPart = parts[0]
		preReleasePart = parts[1]
	}
	if !isPreRelease && isBuild {
		parts := strings.Split(versionStr, "+")
		rootPart = parts[0]
		buildPart = parts[1]
	}
	fmt.Printf("rootPart: %s, preReleasePart: %s, buildPart: %s\n", rootPart, preReleasePart, buildPart)

	version, err := ParseVersion(rootPart)
	if err != nil {
		return nil, err
	}
	preReleaseVersion, err := ParsePrereleaseVersion(preReleasePart)
	if err != nil {
		return nil, err
	}
	return &SemVersion{
		Version:           version,
		PreReleaseVersion: preReleaseVersion,
	}, nil
}
