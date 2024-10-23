package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"sort"
	"strconv"
	"strings"
)

// PreReleaseVersion is similar to Version, but its string representation is reduced by removing trailing ".0"
type PreReleaseVersion struct {
	Version *Version
	Label   string
	Build   *Build
}

// NewPrereleaseVersion creates a new immutable PreReleaseVersion instance
func NewPrereleaseVersion(label string, major int, minor int, patch int, build *Build) *PreReleaseVersion {
	version := NewVersion(major, minor, patch)
	return &PreReleaseVersion{
		Label:   label,
		Version: version,
		Build:   build,
	}
}

// ParsePrereleaseVersion parses a version string and returns a new PreReleaseVersion instance.
// It handles versions with 1, 2, or 3 parts. E.g., "1" becomes "1.0.0", "1.2" becomes "1.2.0".
func ParsePrereleaseVersion(versionStr string) (*PreReleaseVersion, error) {
	// alpha+build.1
	// alpha.1
	parts := strings.Split(versionStr, "+")

	version := parts[0]

	if utils.IsAllAlphabetic(version) && len(parts) == 1 {
		return NewPrereleaseVersion(version, 0, 0, 0, nil), nil
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

	// Label specific logic

	// Build specific logic
	var build *Build
	if len(parts) > 1 {
		buildStr := parts[1]
		build, err = ParseBuild(buildStr)
		if err != nil {
			return nil, err
		}

	}

	return NewPrereleaseVersion(label, major, minor, patch, build), nil
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
	if v.Build != nil && v.Build.Index >= 0 {
		retval = fmt.Sprintf("%s+%s.%d", retval, v.Build.Label, v.Build.Index)
	}
	return retval
}

// Bump returns a new PreReleaseVersion instance after incrementing the specified part
func (v *PreReleaseVersion) Bump(versionPart int, preReleaseLabels []string, buildLabel string) (*PreReleaseVersion, error) {
	switch versionPart {
	// TODO: Implement bumping for prerelease and build versions
	case PreReleaseMajor:
		return NewPrereleaseVersion(v.Label, v.Version.major+1, 0, 0, nil), nil
	case PreReleaseMinor:
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor+1, 0, nil), nil
	case PreReleasePatch:
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor, v.Version.patch+1, nil), nil
	case PreReleaseBuild:
		// TODO: Move to build.go
		if v.Build == nil {
			v.Build = NewBuild(buildLabel, 0)
		}
		v.Build = v.Build.Bump()
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor, v.Version.patch, v.Build), nil
	case PreReleaseNext:
		// sort pre-release labels
		sort.Strings(preReleaseLabels)
		// find the index of the current label
		idx := sort.SearchStrings(preReleaseLabels, v.Label)
		fmt.Printf("Found label %s at index %d\n", v.Label, idx)

	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
	return nil, nil
}
