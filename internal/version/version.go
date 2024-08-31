package version

import (
	"fmt"
	"strconv"
	"strings"
)

const VersionBumpVersion = "VersionBump v0.0.2"

const (
	VersionMajor = iota
	VersionMinor
	VersionPatch

	VersionMajorStr = "major"
	VersionMinorStr = "minor"
	VersionPatchStr = "patch"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func NewVersion(major int, minor int, patch int) *Version {
	return &Version{Major: major, Minor: minor, Patch: patch}
}

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

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) Bump(versionPart int) error {
	switch versionPart {
	case VersionMajor:
		v.Major++
		v.Minor = 0
		v.Patch = 0
		return nil
	case VersionMinor:
		v.Minor++
		v.Patch = 0
		return nil
	case VersionPatch:
		v.Patch++
		return nil
	default:
		return fmt.Errorf("invalid version part: %d", versionPart)
	}
}

func (v *Version) StringBump(versionPart string) error {
	switch versionPart {
	case VersionMajorStr:
		return v.Bump(VersionMajor)
	case VersionMinorStr:
		return v.Bump(VersionMinor)
	case VersionPatchStr:
		return v.Bump(VersionPatch)
	default:
		return fmt.Errorf("invalid version part: %s", versionPart)
	}
}
