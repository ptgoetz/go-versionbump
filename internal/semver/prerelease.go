package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"strconv"
	"strings"
)

// PreReleaseVersion is similar to Version, but its string representation is reduced by removing trailing ".0"
type PreReleaseVersion struct {
	Version    *Version
	Label      string
	Build      int
	BuildLabel string
}

// NewPrereleaseVersion creates a new immutable PreReleaseVersion instance
func NewPrereleaseVersion(label string, major int, minor int, patch int, buildLabel string, buildValue int) *PreReleaseVersion {
	version := NewVersion(major, minor, patch)
	return &PreReleaseVersion{
		Label:      label,
		Version:    version,
		Build:      buildValue,
		BuildLabel: buildLabel,
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
		return NewPrereleaseVersion(version, 0, 0, 0, "", 0), nil
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

	buildLabel := ""
	buildVal := 0
	if len(parts) > 1 {
		build := parts[1]
		vals := strings.Split(build, ".")

		if !utils.StartsWithDigit(build) {
			buildLabel = vals[0]
			valStr := vals[1:]
			buildVal, err = strconv.Atoi(valStr[0])
			if err != nil {
				return nil, fmt.Errorf("invalid build version: %s", valStr[0])
			}
		}
	}

	return NewPrereleaseVersion(label, major, minor, patch, buildLabel, buildVal), nil
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
	if v.BuildLabel != "" && v.Build >= 0 {
		retval = fmt.Sprintf("%s+%s.%d", retval, v.BuildLabel, v.Build)
	}
	return retval
}

// Bump returns a new PreReleaseVersion instance after incrementing the specified part
func (v *PreReleaseVersion) Bump(versionPart int) *PreReleaseVersion {
	switch versionPart {
	// TODO: Implement bumping for prerelease and build versions
	case PreReleaseMajor:
		return NewPrereleaseVersion(v.Label, v.Version.major+1, 0, 0, "", 0)
	case PreReleaseMinor:
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor+1, 0, "", 0)
	case PreReleasePatch:
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor, v.Version.patch+1, "", 0)
	case PreReleaseBuild:
		return NewPrereleaseVersion(v.Label, v.Version.major, v.Version.minor, v.Version.patch /* TODO: don't hard-code */, "build", v.Build+1)
	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
}
