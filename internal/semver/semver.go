package semver

import (
	"fmt"
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
	Version           *Version
	PreReleaseVersion *PreReleaseVersion
}

// String returns the version string
func (v *SemVersion) String() string {
	version := v.Version.String()
	preReleaseStr := v.PreReleaseVersion.String()
	if v.PreReleaseVersion != nil && preReleaseStr != "" {
		version += "-" + v.PreReleaseVersion.String()
	}
	return version
}

// Bump returns a new SemVersion instance after incrementing the specified part
func (v *SemVersion) Bump(versionPart int, preReleaseLabels []string, buildLabel string) (*SemVersion, error) {
	if versionPart >= VersionMajor && versionPart <= VersionPatch {
		// bump the root version
		v.Version = v.Version.Bump(versionPart)

		// reset all pre-release versions
		v.PreReleaseVersion = NewPrereleaseVersion("", 0, 0, 0, nil)
	} else if versionPart >= PreReleaseNext && versionPart <= PreReleaseBuild {
		v.PreReleaseVersion, _ = v.PreReleaseVersion.Bump(versionPart, preReleaseLabels, buildLabel)

	} else {
		return nil, fmt.Errorf("invalid version part: %d", versionPart)
	}
	return nil, nil

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
