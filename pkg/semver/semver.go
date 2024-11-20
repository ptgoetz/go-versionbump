// Package `semver` provides utilities for parsing, validating, and manipulating semantic versions.
// It supports operations such as version comparison, version bumping, and sorting of versions.
// The package also includes support for pre-release and build metadata as defined in the Semantic Versioning 2.0.0 specification.
package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"sort"
	"strings"
)

const (
	vMajor = iota
	vMinor
	vPatch
	prNext
	prMajor
	prNewMajor
	prNewMinor
	prNewPatch
	prMinor
	prPatch
	prBuild
)

type BumpStrategy string

const (
	Major              BumpStrategy = "major"
	Minor              BumpStrategy = "minor"
	Patch              BumpStrategy = "patch"
	PreRelease         BumpStrategy = "pre"
	PreReleaseMajor    BumpStrategy = "pre-major"
	PreReleaseMinor    BumpStrategy = "pre-minor"
	PreReleasePatch    BumpStrategy = "pre-patch"
	PreReleaseBuild    BumpStrategy = "pre-build"
	PreReleaseNewMajor BumpStrategy = "pre-new-major"
	PreReleaseNewMinor BumpStrategy = "pre-new-minor"
	PreReleaseNewPatch BumpStrategy = "pre-new-patch"
)

func versionPartInt(part BumpStrategy) int {
	switch part {
	case Major:
		return vMajor
	case Minor:
		return vMinor
	case Patch:
		return vPatch
	case PreRelease:
		return prNext
	case PreReleaseMajor:
		return prMajor
	case PreReleaseMinor:
		return prMinor
	case PreReleasePatch:
		return prPatch
	case PreReleaseBuild:
		return prBuild
	case PreReleaseNewMajor:
		return prNewMajor
	case PreReleaseNewMinor:
		return prNewMinor
	case PreReleaseNewPatch:
		return prNewPatch
	default:
		panic(fmt.Sprintf("invalid rootVersion part: %s", part))
	}
}

// SemanticVersion represents a parsed and validated semantic version
type SemanticVersion struct {
	rootVersion       *Version
	preReleaseVersion *PreReleaseVersion
	buildVersion      *BuildVersion
}

// RootVersion returns the root version part of the SemanticVersion
func (v *SemanticVersion) RootVersion() *Version {
	return v.rootVersion
}

// PreReleaseVersion returns the pre-release version part of the SemanticVersion
func (v *SemanticVersion) PreReleaseVersion() *PreReleaseVersion {
	return v.preReleaseVersion
}

// BuildVersion returns the BuildVersion part of the SemanticVersion
func (v *SemanticVersion) BuildVersion() *BuildVersion {
	return v.buildVersion
}

// String returns the string representation of the SemanticVersion instance
func (v *SemanticVersion) String() string {
	if v == nil {
		return ""
	}
	version := v.rootVersion.String()
	if v.preReleaseVersion != nil && v.preReleaseVersion.String() != "" {
		version += "-" + v.preReleaseVersion.String()
	}
	if v.buildVersion != nil && v.buildVersion.String() != "" {
		version += "+" + v.buildVersion.String()
	}
	return version
}

// Bump returns a new SemanticVersion instance after incrementing the specified part.
// If the part is a pre-release part, `preReleaseLabels` must be provided. If the part is a build version part,
// `buildLabel` must be provided. If the part is a root version part, preReleaseLabels and buildLabel are ignored.
func (v *SemanticVersion) Bump(strategy BumpStrategy, preReleaseLabels []string, buildLabel string) (*SemanticVersion, error) {
	var version *Version
	var preReleaseVersion *PreReleaseVersion
	var build *BuildVersion
	var err error
	versionPart := versionPartInt(strategy)

	switch {
	case versionPart >= vMajor && versionPart <= vPatch:
		// bump the root version
		version = v.rootVersion.bump(versionPart)
		// reset all pre-release versions
		preReleaseVersion = newPrereleaseVersion("", 0, 0, 0)
	case versionPart >= prNewMajor && versionPart <= prNewPatch:
		// bump the root version
		version = v.rootVersion.bump(versionPart - 5)
		// reset all pre-release versions
		preReleaseVersion = newPrereleaseVersion(preReleaseLabels[0], 0, 0, 0)
	case versionPart >= prNext && versionPart <= prPatch:
		version = newVersion(v.rootVersion.major, v.rootVersion.minor, v.rootVersion.patch)
		preReleaseVersion, err = v.preReleaseVersion.bump(versionPart, preReleaseLabels)
		if err != nil {
			return nil, err
		}
	case versionPart == prBuild:
		version = newVersion(v.rootVersion.major, v.rootVersion.minor, v.rootVersion.patch)
		preReleaseVersion = newPrereleaseVersion(v.preReleaseVersion.label, v.preReleaseVersion.version.major, v.preReleaseVersion.version.minor, v.preReleaseVersion.version.patch)
		if v.buildVersion != nil {
			build = v.buildVersion.bump()
		} else {
			build = newBuild(buildLabel, 1)
		}
	default:
		return nil, fmt.Errorf("invalid version strategy: %d", versionPart)
	}

	return &SemanticVersion{
		rootVersion:       version,
		preReleaseVersion: preReleaseVersion,
		buildVersion:      build,
	}, nil
}

// Compare compares two SemanticVersion instances.
// Returns -1 if v is less than other, 1 if v is greater than other, and 0 if they are equal.
func (v *SemanticVersion) Compare(other *SemanticVersion) int {
	if v.rootVersion.major != other.rootVersion.major {
		if v.rootVersion.major < other.rootVersion.major {
			return -1
		}
		return 1
	}

	if v.rootVersion.minor != other.rootVersion.minor {
		if v.rootVersion.minor < other.rootVersion.minor {
			return -1
		}
		return 1
	}

	if v.rootVersion.patch != other.rootVersion.patch {
		if v.rootVersion.patch < other.rootVersion.patch {
			return -1
		}
		return 1
	}

	if v.preReleaseVersion != nil && other.preReleaseVersion != nil {
		preReleaseComparison := v.preReleaseVersion.Compare(other.preReleaseVersion)
		if preReleaseComparison != 0 {
			return preReleaseComparison
		}
	} else if v.preReleaseVersion != nil {
		return -1
	} else if other.preReleaseVersion != nil {
		return 1
	}

	if v.buildVersion != nil && other.buildVersion != nil {
		buildComparison := v.buildVersion.Compare(other.buildVersion)
		if buildComparison != 0 {
			return buildComparison
		}
	} else if v.buildVersion != nil {
		return 1
	} else if other.buildVersion != nil {
		return -1
	}

	return 0
}

// ParseSemVersion parses a semantic rootVersion string and returns a new SemanticVersion instance
func ParseSemVersion(versionStr string) (*SemanticVersion, error) {
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
	} else if !isBuild && isPreRelease { // pre-release and no BuildVersion
		parts := strings.Split(versionStr, "-")
		rootPart = parts[0]
		preReleasePart = parts[1]
	}

	version, err := parseVersion(rootPart)
	if err != nil {
		return nil, err
	}
	preReleaseVersion, err := parsePrereleaseVersion(preReleasePart)
	if err != nil {
		return nil, err
	}

	build, err := parseBuild(buildPart)
	if err != nil {
		return nil, err
	}
	return &SemanticVersion{
		rootVersion:       version,
		preReleaseVersion: preReleaseVersion,
		buildVersion:      build,
	}, nil
}

// ValidateSemVersion checks if the provided string is a valid semantic version
func ValidateSemVersion(versionStr string) bool {
	_, err := ParseSemVersion(versionStr)
	return err == nil
}

// ValidatePreReleaseLabels checks if the provided pre-release labels are valid
func ValidatePreReleaseLabels(preReleaseLabels []string) bool {
	for _, label := range preReleaseLabels {
		if !utils.IsAllAlphabetic(label) {
			return false
		}
	}
	return true
}

// ValidateBuildLabel checks if the provided pre-release labels are valid
func ValidateBuildLabel(buildLabel string) bool {
	return utils.IsAllAlphanumeric(buildLabel)
}

// ValidatePreReleaseLabelsString checks if the provided pre-release labels string is valid.
// The labels must be comma-separated.
func ValidatePreReleaseLabelsString(preReleaseLabels string) bool {
	labels := strings.Split(preReleaseLabels, ",")
	return ValidatePreReleaseLabels(labels)
}

// SortVersions sorts a slice of SemanticVersion instances in descending order (latest to oldest)
func SortVersions(versions []*SemanticVersion) {
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) > 0
	})
}
