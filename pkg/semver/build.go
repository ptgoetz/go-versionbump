package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"strconv"
	"strings"
)

type build struct {
	index int
	label string
}

// newBuild creates a new build instance
func newBuild(label string, index int) *build {
	return &build{
		index: index,
		label: label,
	}
}

// parseBuild parses a build version string and returns a new build instance
func parseBuild(buildStr string) (*build, error) {
	if buildStr == "" {
		return nil, nil
	}
	// build specific logic
	// build.1
	vals := strings.Split(buildStr, ".")
	if len(vals) != 2 {
		return nil, fmt.Errorf("invalid build version: %s", buildStr)
	}
	if vals[1] == "" {
		return nil, fmt.Errorf("invalid build version, build version is required: %s", buildStr)
	}
	// build label must be alphabetic
	if !utils.IsAllAlphabetic(vals[0]) {
		return nil, fmt.Errorf("invalid build version, build label must not contain digits: %s", buildStr)
	}
	buildNum, err := strconv.Atoi(vals[1])
	if err != nil {
		return nil, fmt.Errorf("invalid build version, build number must be an integer: %s", buildStr)
	}
	return &build{
		index: buildNum,
		label: vals[0],
	}, nil
}

// String returns the build version string
func (b *build) String() string {
	if b.index > 0 {
		return fmt.Sprintf("%s.%d", b.label, b.index)
	} else {
		return ""
	}
}

// Compare compares two build instances.
// Returns -1 if b is less than other, 1 if b is greater than other, and 0 if they are equal.
func (b *build) Compare(other *build) int {
	if b.label != other.label {
		if b.label < other.label {
			return -1
		}
		return 1
	}

	if b.index != other.index {
		if b.index < other.index {
			return -1
		}
		return 1
	}

	return 0
}

func (b *build) bump() *build {
	return newBuild(b.label, b.index+1)
}
