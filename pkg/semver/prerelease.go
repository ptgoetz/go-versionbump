package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"sort"
	"strconv"
	"strings"
)

// PreReleaseVersion is similar to Version, but its string representation is reduced by removing trailing ".0"
// All PreReleaseVersion instances have an associated `label`.
type PreReleaseVersion struct {
	version *Version
	label   string
}

// newPrereleaseVersion creates a new immutable PreReleaseVersion instance
func newPrereleaseVersion(label string, major int, minor int, patch int) *PreReleaseVersion {
	version := newVersion(major, minor, patch)
	return &PreReleaseVersion{
		label:   label,
		version: version,
	}
}

// parsePrereleaseVersion parses a rootVersion string and returns a new PreReleaseVersion instance.
// It handles versions with 1, 2, or 3 parts. E.g., "1" becomes "1.0.0", "1.2" becomes "1.2.0".
func parsePrereleaseVersion(versionStr string) (*PreReleaseVersion, error) {
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
		// No rootVersion parts provided, e.g. "alpha"
		major, minor, patch = 0, 0, 0
	case 1:
		// Only major part provided, e.g. "1"
		major, err = strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("invalid major rootVersion: %s", vals[0])
		}
		minor, patch = 0, 0
	case 2:
		// Major and minor parts provided, e.g. "1.2"
		major, err = strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("invalid major rootVersion: %s", vals[0])
		}
		minor, err = strconv.Atoi(vals[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor rootVersion: %s", vals[1])
		}
		patch = 0
	case 3:
		// Full semantic rootVersion provided, e.g. "1.2.3"
		major, err = strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("invalid major rootVersion: %s", vals[0])
		}
		minor, err = strconv.Atoi(vals[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor rootVersion: %s", vals[1])
		}
		patch, err = strconv.Atoi(vals[2])
		if err != nil {
			return nil, fmt.Errorf("invalid patch rootVersion: %s", vals[2])
		}
	}

	return newPrereleaseVersion(label, major, minor, patch), nil
}

// Label returns the pre-release label
func (v *PreReleaseVersion) Label() string {
	return v.label
}

// Version returns the pre-release Version
func (v *PreReleaseVersion) Version() *Version {
	return v.version
}

// String returns the reduced rootVersion string by removing trailing ".0" parts
func (v *PreReleaseVersion) String() string {
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

// Compare compares two PreReleaseVersion instances.
// Returns -1 if v is less than other, 1 if v is greater than other, and 0 if they are equal.
func (v *PreReleaseVersion) Compare(other *PreReleaseVersion) int {
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

// bump returns a new PreReleaseVersion instance after incrementing the specified part
func (v *PreReleaseVersion) bump(versionPart int, preReleaseLabels []string) (*PreReleaseVersion, error) {
	if len(preReleaseLabels) == 0 {
		panic("PreReleaseVersion.bump(): preReleaseLabels cannot be empty")
	}
	// sort pre-release labels
	sort.Strings(preReleaseLabels)
	// if the label is empty, this is the first pre-release rootVersion, so return the first label
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
		// find the number of the current label
		idx := indexOf(label, preReleaseLabels)
		// if the rootVersion being bumped has no label, return the first label
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
		panic(fmt.Sprintf("invalid rootVersion part: %d.\n", versionPart))
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
