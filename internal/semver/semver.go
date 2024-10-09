package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"strconv"
	"strings"
)

const (
	VersionMajor = iota
	VersionMinor
	VersionPatch
	PrereleaseNext
	PrereleaseMajor
	PrereleaseMinor
	PrereleasePatch
	Build
)

// Version represents a semantic version
type Version struct {
	major int
	minor int
	patch int
}

// PreReleaseVersion is similar to Version, but its string representation is reduced by removing trailing ".0"
type PreReleaseVersion struct {
	Version *Version
	Label   string
}

type SemVersion struct {
	Version                *Version
	PreReleaseVersionLabel string
	PreReleaseVersion      *PreReleaseVersion
	Build                  int
}

// String returns the version string
func (v *SemVersion) String() string {
	version := v.Version.String()
	if v.PreReleaseVersion != nil && v.PreReleaseVersion.Version.String() == "" {
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
	if isPreRelease && !isBuild {
		parts := strings.Split(versionStr, "-")
		rootPart = parts[0]
		preReleasePart = parts[1]
	}
	if !isPreRelease && isBuild {
		parts := strings.Split(versionStr, "+")
		rootPart = parts[0]
		buildPart = parts[1]
	}
	if isPreRelease && isBuild {
		parts := strings.Split(versionStr, "-")
		rootPart = parts[0]
		parts = strings.Split(parts[1], "+")
		preReleasePart = parts[0]
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

// NewVersion creates a new immutable Version instance
func NewVersion(major int, minor int, patch int) *Version {
	return &Version{major: major, minor: minor, patch: patch}
}

// ParseVersion parses a version string and returns a new Version instance
func ParseVersion(version string) (*Version, error) {
	vals := strings.Split(version, ".")
	if len(vals) != 3 {
		return nil, fmt.Errorf("invalid semantic version string: %s", version)
	}
	major, err := strconv.Atoi(vals[VersionMajor])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(vals[VersionMinor])
	if err != nil {
		return nil, err
	}
	patch, err := strconv.Atoi(vals[VersionPatch])
	if err != nil {
		return nil, err
	}
	return NewVersion(major, minor, patch), nil
}

// String returns the version string
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

// Bump returns a new Version instance after incrementing the specified part
func (v *Version) Bump(versionPart int) *Version {
	switch versionPart {
	case VersionMajor:
		return NewVersion(v.major+1, 0, 0)
	case VersionMinor:
		return NewVersion(v.major, v.minor+1, 0)
	case VersionPatch:
		return NewVersion(v.major, v.minor, v.patch+1)
	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
}

// ValidateVersion checks if the provided version string is a valid semantic version
func ValidateVersion(version string) bool {
	_, err := ParseVersion(version)
	return err == nil
}

// NewPrereleaseVersion creates a new immutable PreReleaseVersion instance
func NewPrereleaseVersion(label string, major int, minor int, patch int) *PreReleaseVersion {
	version := NewVersion(major, minor, patch)
	return &PreReleaseVersion{
		Label:   label,
		Version: version,
	}
}

// ParsePrereleaseVersion parses a version string and returns a new PreReleaseVersion instance.
// It handles versions with 1, 2, or 3 parts. E.g., "1" becomes "1.0.0", "1.2" becomes "1.2.0".
func ParsePrereleaseVersion(version string) (*PreReleaseVersion, error) {

	if utils.IsAllAlphabetic(version) {
		return NewPrereleaseVersion(version, 0, 0, 0), nil
	}

	vals := strings.Split(version, ".")
	var major, minor, patch int
	var err error
	var label string
	if !utils.StartsWithDigit(version) {
		label = vals[0]
		vals = vals[1:]
	}

	switch len(vals) {
	case 0:
		// No version parts provided, e.g. "alpha"
		major, minor, patch = 0, 0, 0
	case 1:
		// Only major part provided, e.g. "1"
		major, err = strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("invalid major version: %s", vals[0])
		}
		minor, patch = 0, 0
	case 2:
		// Major and minor parts provided, e.g. "1.2"
		major, err = strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("invalid major version: %s", vals[0])
		}
		minor, err = strconv.Atoi(vals[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", vals[1])
		}
		patch = 0
	case 3:
		// Full semantic version provided, e.g. "1.2.3"
		major, err = strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("invalid major version: %s", vals[0])
		}
		minor, err = strconv.Atoi(vals[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", vals[1])
		}
		patch, err = strconv.Atoi(vals[2])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", vals[2])
		}
	}

	return NewPrereleaseVersion(label, major, minor, patch), nil
}

// String returns the reduced version string by removing trailing ".0" parts
func (v *PreReleaseVersion) String() string {
	var retval string
	if v.Version.patch != 0 {
		retval = fmt.Sprintf("%d.%d.%d", v.Version.major, v.Version.minor, v.Version.patch)
	} else if v.Version.minor != 0 {
		retval = fmt.Sprintf("%d.%d", v.Version.major, v.Version.minor)
	} else {
		retval = fmt.Sprintf("%d", v.Version.major)
	}

	// alpha.0.0.0 -> alpha
	if retval == "0" {
		retval = fmt.Sprintf("%s", v.Label)
	} else if v.Label != "" {
		retval = fmt.Sprintf("%s.%s", v.Label, retval)
	}
	if v.Label != "" {
		retval = fmt.Sprintf("%s", retval)
	}
	return retval
}

// Bump returns a new PreReleaseVersion instance after incrementing the specified part
func (v *PreReleaseVersion) Bump(versionPart int) *PreReleaseVersion {
	switch versionPart {
	case VersionMajor:
		return NewPrereleaseVersion(v.Label, v.Version.major+1, 0, 0)
	case VersionMinor:
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor+1, 0)
	case VersionPatch:
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor, v.Version.patch+1)
	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
}
