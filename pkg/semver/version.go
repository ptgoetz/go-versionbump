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

// newVersion creates a new immutable Version instance
func newVersion(major int, minor int, patch int) *Version {
	return &Version{major: major, minor: minor, patch: patch}
}

// parseVersion parses a version string and returns a new Version instance
func parseVersion(version string) (*Version, error) {
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
	return newVersion(major, minor, patch), nil
}

// String returns the version string
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

// bump returns a new Version instance after incrementing the specified part
func (v *Version) bump(versionPart int) *Version {
	switch versionPart {
	case vMajor:
		return newVersion(v.major+1, 0, 0)
	case vMinor:
		return newVersion(v.major, v.minor+1, 0)
	case vPatch:
		return newVersion(v.major, v.minor, v.patch+1)
	default:
		panic(fmt.Sprintf("invalid version part: %d.\n", versionPart))
	}
}
