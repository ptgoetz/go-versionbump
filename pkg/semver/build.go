package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"strconv"
	"strings"
)

type BuildVersion struct {
	number int
	label  string
}

// newBuild creates a new BuildVersion instance
func newBuild(label string, index int) *BuildVersion {
	return &BuildVersion{
		number: index,
		label:  label,
	}
}

func (b *BuildVersion) Number() int {
	return b.number
}

func (b *BuildVersion) Label() string {
	return b.label
}

// parseBuild parses a BuildVersion version string and returns a new BuildVersion instance
func parseBuild(buildStr string) (*BuildVersion, error) {
	if buildStr == "" {
		return nil, nil
	}
	// BuildVersion specific logic
	// BuildVersion.1
	vals := strings.Split(buildStr, ".")
	if len(vals) != 2 {
		return nil, fmt.Errorf("invalid BuildVersion version: %s", buildStr)
	}
	if vals[1] == "" {
		return nil, fmt.Errorf("invalid BuildVersion version, BuildVersion version is required: %s", buildStr)
	}
	// BuildVersion label must be alphabetic
	if !utils.IsAllAlphabetic(vals[0]) {
		return nil, fmt.Errorf("invalid BuildVersion version, BuildVersion label must not contain digits: %s", buildStr)
	}
	buildNum, err := strconv.Atoi(vals[1])
	if err != nil {
		return nil, fmt.Errorf("invalid BuildVersion version, BuildVersion number must be an integer: %s", buildStr)
	}
	return &BuildVersion{
		number: buildNum,
		label:  vals[0],
	}, nil
}

// String returns the BuildVersion version string
func (b *BuildVersion) String() string {
	if b.number > 0 {
		return fmt.Sprintf("%s.%d", b.label, b.number)
	} else {
		return ""
	}
}

// Compare compares two BuildVersion instances.
// Returns -1 if buildVersion is less than other, 1 if buildVersion is greater than other, and 0 if they are equal.
func (b *BuildVersion) Compare(other *BuildVersion) int {
	if b.label != other.label {
		if b.label < other.label {
			return -1
		}
		return 1
	}

	if b.number != other.number {
		if b.number < other.number {
			return -1
		}
		return 1
	}

	return 0
}

func (b *BuildVersion) bump() *BuildVersion {
	return newBuild(b.label, b.number+1)
}
