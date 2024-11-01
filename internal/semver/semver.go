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
	Build             *Build
}

// String returns the version string
func (v *SemVersion) String() string {
	version := v.Version.String()
	if v.PreReleaseVersion != nil && v.PreReleaseVersion.String() != "" {
		version += "-" + v.PreReleaseVersion.String()
	}
	if v.Build != nil && v.Build.String() != "" {
		version += "+" + v.Build.String()
	}
	return version
}

// Bump returns a new SemVersion instance after incrementing the specified part
func (v *SemVersion) Bump(part VersionPart, preReleaseLabels []string, buildLabel string) (*SemVersion, error) {
	var version *Version
	var preReleaseVersion *PreReleaseVersion
	var build *Build
	versionPart := versionPartInt(part)
	if versionPart >= vMajor && versionPart <= vPatch {
		// bump the root version
		version = v.Version.Bump(versionPart)

		// reset all pre-release versions
		preReleaseVersion = NewPrereleaseVersion("", 0, 0, 0)
	} else if versionPart >= prNext && versionPart <= prPatch {
		version = NewVersion(v.Version.major, v.Version.minor, v.Version.patch)
		preReleaseVersion, _ = v.PreReleaseVersion.Bump(versionPart, preReleaseLabels)

	} else if versionPart == prBuild {
		version = NewVersion(v.Version.major, v.Version.minor, v.Version.patch)
		preReleaseVersion = NewPrereleaseVersion(v.PreReleaseVersion.Label, v.PreReleaseVersion.Version.major, v.PreReleaseVersion.Version.minor, v.PreReleaseVersion.Version.patch)
		if v.Build != nil {
			build = v.Build.Bump()
		} else {
			build = NewBuild(buildLabel, 1)
		}
	} else {
		return nil, fmt.Errorf("invalid version part: %d", versionPart)
	}
	return &SemVersion{
		Version:           version,
		PreReleaseVersion: preReleaseVersion,
		Build:             build,
	}, nil

}

// Compare compares two SemVersion instances.
// Returns -1 if v is less than other, 1 if v is greater than other, and 0 if they are equal.
func (v *SemVersion) Compare(other *SemVersion) int {
	if v.Version.major != other.Version.major {
		if v.Version.major < other.Version.major {
			return -1
		}
		return 1
	}

	if v.Version.minor != other.Version.minor {
		if v.Version.minor < other.Version.minor {
			return -1
		}
		return 1
	}

	if v.Version.patch != other.Version.patch {
		if v.Version.patch < other.Version.patch {
			return -1
		}
		return 1
	}

	if v.PreReleaseVersion != nil && other.PreReleaseVersion != nil {
		preReleaseComparison := v.PreReleaseVersion.Compare(other.PreReleaseVersion)
		if preReleaseComparison != 0 {
			return preReleaseComparison
		}
	} else if v.PreReleaseVersion != nil {
		return -1
	} else if other.PreReleaseVersion != nil {
		return 1
	}

	if v.Build != nil && other.Build != nil {
		buildComparison := v.Build.Compare(other.Build)
		if buildComparison != 0 {
			return buildComparison
		}
	} else if v.Build != nil {
		return 1
	} else if other.Build != nil {
		return -1
	}

	return 0
}

// ParseSemVersion parses a semantic version string and returns a new SemVersion instance
func ParseSemVersion(versionStr string) (*SemVersion, error) {
	isPreRelease := strings.Contains(versionStr, "-")
	isBuild := strings.Contains(versionStr, "+")

	var rootPart string
	var preReleasePart string
	var buildPart string
	if !isPreRelease && !isBuild {
		rootPart = versionStr
	} else if isPreRelease && isBuild {
		parts := strings.Split(versionStr, "-")
		rootPart = parts[0]
		prAndBuildParts := strings.Split(parts[1], "+")
		preReleasePart = prAndBuildParts[0]
		buildPart = prAndBuildParts[1]
	} else if isBuild && !isPreRelease {
		parts := strings.Split(versionStr, "+")
		rootPart = parts[0]
		buildPart = parts[1]
	} else if !isBuild && isPreRelease { // pre-release and no build
		parts := strings.Split(versionStr, "-")
		rootPart = parts[0]
		preReleasePart = parts[1]
	}

	version, err := ParseVersion(rootPart)
	if err != nil {
		return nil, err
	}
	preReleaseVersion, err := ParsePrereleaseVersion(preReleasePart)
	if err != nil {
		return nil, err
	}

	build, err := ParseBuild(buildPart)
	if err != nil {
		return nil, err
	}
	return &SemVersion{
		Version:           version,
		PreReleaseVersion: preReleaseVersion,
		Build:             build,
	}, nil
}

func ValidateSemVersion(versionStr string) bool {
	_, err := ParseSemVersion(versionStr)
	return err == nil
}
