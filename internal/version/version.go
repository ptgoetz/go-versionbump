package version

import (
	"fmt"
	"strconv"
	"strings"
)

const VersionBumpVersion = "VersionBump v0.4.1"

const (
	VersionMajor = iota
	VersionMinor
	VersionPatch

	VersionMajorStr = "major"
	VersionMinorStr = "minor"
	VersionPatchStr = "patch"
)

type Version struct {
	major int
	minor int
	patch int
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

// StringBump returns a new Version instance after incrementing the specified part (as a string)
func (v *Version) StringBump(versionPart string) *Version {
	switch versionPart {
	case VersionMajorStr:
		return v.Bump(VersionMajor)
	case VersionMinorStr:
		return v.Bump(VersionMinor)
	case VersionPatchStr:
		return v.Bump(VersionPatch)
	default:
		panic(fmt.Sprintf("invalid version part: %s. Call `ValidateVersionPart()` to prevent this error.\n", versionPart))
	}
}

// ValidateVersionPart checks if the provided version part string is valid
func ValidateVersionPart(part string) bool {
	switch part {
	case VersionMajorStr, VersionMinorStr, VersionPatchStr:
		return true
	default:
		return false
	}
}

// ValidateVersion checks if the provided version string is a valid semantic version
func ValidateVersion(version string) bool {
	_, err := ParseVersion(version)
	return err == nil
}

// Subversion is similar to Version, but its string representation is reduced by removing trailing ".0"
type Subversion struct {
	major int
	minor int
	patch int
}

// NewSubversion creates a new immutable Subversion instance
func NewSubversion(major int, minor int, patch int) *Subversion {
	return &Subversion{major: major, minor: minor, patch: patch}
}

// ParseSubversion parses a version string and returns a new Subversion instance.
// It handles versions with 1, 2, or 3 parts. E.g., "1" becomes "1.0.0", "1.2" becomes "1.2.0".
func ParseSubversion(version string) (*Subversion, error) {
	vals := strings.Split(version, ".")
	var major, minor, patch int
	var err error

	switch len(vals) {
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
	default:
		return nil, fmt.Errorf("invalid version string: %s", version)
	}

	return NewSubversion(major, minor, patch), nil
}

// String returns the reduced version string by removing trailing ".0" parts
func (v *Subversion) String() string {
	if v.patch != 0 {
		return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	} else if v.minor != 0 {
		return fmt.Sprintf("%d.%d", v.major, v.minor)
	}
	return fmt.Sprintf("%d", v.major)
}

// Bump returns a new Subversion instance after incrementing the specified part
func (v *Subversion) Bump(versionPart int) *Subversion {
	switch versionPart {
	case VersionMajor:
		return NewSubversion(v.major+1, 0, 0)
	case VersionMinor:
		return NewSubversion(v.major, v.minor+1, 0)
	case VersionPatch:
		return NewSubversion(v.major, v.minor, v.patch+1)
	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
}

// StringBump returns a new Subversion instance after incrementing the specified part (as a string)
func (v *Subversion) StringBump(versionPart string) *Subversion {
	switch versionPart {
	case VersionMajorStr:
		return v.Bump(VersionMajor)
	case VersionMinorStr:
		return v.Bump(VersionMinor)
	case VersionPatchStr:
		return v.Bump(VersionPatch)
	default:
		panic(fmt.Sprintf("invalid version part: %s. Call `ValidateVersionPart()` to prevent this error.\n", versionPart))
	}
}
