package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a semantic version
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
	major, err := strconv.Atoi(vals[vMajor])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(vals[vMinor])
	if err != nil {
		return nil, err
	}
	patch, err := strconv.Atoi(vals[vPatch])
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
	case vMajor:
		return NewVersion(v.major+1, 0, 0)
	case vMinor:
		return NewVersion(v.major, v.minor+1, 0)
	case vPatch:
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
