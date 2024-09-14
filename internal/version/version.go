package version

import (
	"fmt"
	"strconv"
	"strings"
)

const VersionBumpVersion = "VersionBump v0.3.0"

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
