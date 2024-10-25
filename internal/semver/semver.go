package semver

import (
	"fmt"
	"strings"
)

const (
	vMajor = iota
	vMinor
	vPatch
	prNext
	prMajor
	prMinor
	prPatch
	prBuild
)

type VersionPart string

const (
	Major           VersionPart = "major"
	Minor           VersionPart = "minor"
	Patch           VersionPart = "patch"
	PreReleaseNext  VersionPart = "prerelease-next"
	PreReleaseMajor VersionPart = "prerelease-major"
	PreReleaseMinor VersionPart = "prerelease-minor"
	PreReleasePatch VersionPart = "prerelease-patch"
	PreReleaseBuild VersionPart = "prerelease-build"
)

func versionPartInt(part VersionPart) int {
	switch part {
	case Major:
		return vMajor
	case Minor:
		return vMinor
	case Patch:
		return vPatch
	case PreReleaseNext:
		return prNext
	case PreReleaseMajor:
		return prMajor
	case PreReleaseMinor:
		return prMinor
	case PreReleasePatch:
		return prPatch
	case PreReleaseBuild:
		return prBuild
	default:
		panic(fmt.Sprintf("invalid version part: %s", part))
	}
}

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
func (v *SemVersion) Bump(part VersionPart, preReleaseLabels []string, buildLabel string) (*SemVersion, error) {
	versionPart := versionPartInt(part)
	if versionPart >= vMajor && versionPart <= vPatch {
		// bump the root version
		v.Version = v.Version.Bump(versionPart)

		// reset all pre-release versions
		v.PreReleaseVersion = NewPrereleaseVersion("", 0, 0, 0, nil)
	} else if versionPart >= prNext && versionPart <= prBuild {
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

	var rootPart string
	var preReleasePart string
	if !isPreRelease && !isBuild {
		rootPart = versionStr
	}
	if isPreRelease {
		parts := strings.Split(versionStr, "-")
		rootPart = parts[0]
		preReleasePart = parts[1]
	}
	if !isPreRelease && isBuild {
		parts := strings.Split(versionStr, "+")
		rootPart = parts[0]
	}

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
