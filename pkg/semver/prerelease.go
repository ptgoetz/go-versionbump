package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"sort"
	"strconv"
	"strings"
)

// preReleaseVersion is similar to Version, but its string representation is reduced by removing trailing ".0"
type preReleaseVersion struct {
	version *Version
	label   string
}

// newPrereleaseVersion creates a new immutable preReleaseVersion instance
func newPrereleaseVersion(label string, major int, minor int, patch int) *preReleaseVersion {
	version := newVersion(major, minor, patch)
	return &preReleaseVersion{
		label:   label,
		version: version,
	}
}

// parsePrereleaseVersion parses a version string and returns a new preReleaseVersion instance.
// It handles versions with 1, 2, or 3 parts. E.g., "1" becomes "1.0.0", "1.2" becomes "1.2.0".
func parsePrereleaseVersion(versionStr string) (*preReleaseVersion, error) {
	// alpha+build.1
	// alpha.1
	parts := strings.Split(versionStr, "+")

	version := parts[0]

	if utils.IsAllAlphabetic(version) && len(parts) == 1 {
		return newPrereleaseVersion(version, 0, 0, 0), nil
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

	return newPrereleaseVersion(label, major, minor, patch), nil
}

// String returns the reduced version string by removing trailing ".0" parts
func (v *preReleaseVersion) String() string {
	var retval string
	if v.version.patch != 0 {
		retval = fmt.Sprintf("%d.%d.%d", v.version.major, v.version.minor, v.version.patch)
	} else if v.version.minor != 0 {
		retval = fmt.Sprintf("%d.%d", v.version.major, v.version.minor)
	} else {
		retval = fmt.Sprintf("%d", v.version.major)
	}

	// alpha.0.0.0 -> alpha
	if retval == "0" {
		retval = fmt.Sprint(v.label)
	} else if v.label != "" {
		retval = fmt.Sprintf("%s.%s", v.label, retval)
	}
	if v.label != "" {
		retval = fmt.Sprint(retval)
	}
	return retval
}

// Compare compares two preReleaseVersion instances.
// Returns -1 if v is less than other, 1 if v is greater than other, and 0 if they are equal.
func (v *preReleaseVersion) Compare(other *preReleaseVersion) int {
	if v.label != other.label {
		if v.label < other.label {
			return -1
		}
		return 1
	}

	if v.version.major != other.version.major {
		if v.version.major < other.version.major {
			return -1
		}
		return 1
	}

	if v.version.minor != other.version.minor {
		if v.version.minor < other.version.minor {
			return -1
		}
		return 1
	}

	if v.version.patch != other.version.patch {
		if v.version.patch < other.version.patch {
			return -1
		}
		return 1
	}

	return 0
}

// bump returns a new preReleaseVersion instance after incrementing the specified part
func (v *preReleaseVersion) bump(versionPart int, preReleaseLabels []string) (*preReleaseVersion, error) {
	if len(preReleaseLabels) == 0 {
		panic("preReleaseVersion.bump(): preReleaseLabels cannot be empty")
	}
	// sort pre-release labels
	sort.Strings(preReleaseLabels)
	// if the label is empty, this is the first pre-release version, so return the first label
	label := v.label
	if v.label == "" {
		label = preReleaseLabels[0]
	}

	switch versionPart {
	case prMajor:
		if v.label == "" {
			return newPrereleaseVersion(label, v.version.major, 0, 0), nil
		} else {
			return newPrereleaseVersion(label, v.version.major+1, 0, 0), nil
		}

	case prMinor:
		return newPrereleaseVersion(label, v.version.major, v.version.minor+1, 0), nil
	case prPatch:
		return newPrereleaseVersion(label, v.version.major, v.version.minor, v.version.patch+1), nil
	case prNext:
		// find the index of the current label
		idx := indexOf(label, preReleaseLabels)
		// if the version being bumped has no label, return the first label
		offset := 1
		if v.label == "" {
			offset = 0
		}
		if idx == -1 {
			return nil, fmt.Errorf("label %s not found in preReleaseLabels: %v", v.label, preReleaseLabels)
		} else if idx == len(preReleaseLabels)-1 {
			return nil, fmt.Errorf("cannot bump beyond the last label %s", v.label)
		} else {
			return newPrereleaseVersion(preReleaseLabels[idx+offset], 0, 0, 0), nil
		}
	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
}

// TODO: move to utils package
func indexOf(s string, arr []string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}
