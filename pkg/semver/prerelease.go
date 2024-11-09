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
func ParsePrereleaseVersion(versionStr string) (*PreReleaseVersion, error) {
	// alpha+build.1
	// alpha.1
	parts := strings.Split(versionStr, "+")

	version := parts[0]

	if utils.IsAllAlphabetic(version) && len(parts) == 1 {
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
		retval = fmt.Sprint(v.Label)
	} else if v.Label != "" {
		retval = fmt.Sprintf("%s.%s", v.Label, retval)
	}
	if v.Label != "" {
		retval = fmt.Sprint(retval)
	}
	return retval
}

// Compare compares two PreReleaseVersion instances.
// Returns -1 if v is less than other, 1 if v is greater than other, and 0 if they are equal.
func (v *PreReleaseVersion) Compare(other *PreReleaseVersion) int {
	if v.Label != other.Label {
		if v.Label < other.Label {
			return -1
		}
		return 1
	}

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

	return 0
}

// Bump returns a new PreReleaseVersion instance after incrementing the specified part
func (v *PreReleaseVersion) Bump(versionPart int, preReleaseLabels []string) (*PreReleaseVersion, error) {
	if len(preReleaseLabels) == 0 {
		panic("PreReleaseVersion.Bump(): preReleaseLabels cannot be empty")
	}
	// sort pre-release labels
	sort.Strings(preReleaseLabels)
	// if the label is empty, this is the first pre-release version, so return the first label
	label := v.Label
	if v.Label == "" {
		label = preReleaseLabels[0]
	}

	switch versionPart {
	case prMajor:
		if v.Label == "" {
			return NewPrereleaseVersion(label, v.Version.major, 0, 0), nil
		} else {
			return NewPrereleaseVersion(label, v.Version.major+1, 0, 0), nil
		}

	case prMinor:
		return NewPrereleaseVersion(label, v.Version.major, v.Version.minor+1, 0), nil
	case prPatch:
		return NewPrereleaseVersion(label, v.Version.major, v.Version.minor, v.Version.patch+1), nil
	case prNext:
		// find the index of the current label
		idx := indexOf(label, preReleaseLabels)
		// if the version being bumped has no label, return the first label
		offset := 1
		if v.Label == "" {
			offset = 0
		}
		if idx == -1 {
			return nil, fmt.Errorf("label %s not found in preReleaseLabels: %v", v.Label, preReleaseLabels)
		} else if idx == len(preReleaseLabels)-1 {
			return nil, fmt.Errorf("cannot bump beyond the last label %s", v.Label)
		} else {
			return NewPrereleaseVersion(preReleaseLabels[idx+offset], 0, 0, 0), nil
		}
	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
}

func indexOf(s string, arr []string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}
