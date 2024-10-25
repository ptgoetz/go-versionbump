package semver

import (
	"fmt"
	"github.com/ptgoetz/go-versionbump/internal/utils"
	"strconv"
	"strings"
)

type Build struct {
	Index int
	Label string
}

// NewBuild creates a new Build instance
func NewBuild(label string, index int) *Build {
	return &Build{
		Index: index,
		Label: label,
	}
}

// ParseBuild parses a build version string and returns a new Build instance
func ParseBuild(buildStr string) (*Build, error) {
	if buildStr == "" {
		return nil, nil
	}
	// Build specific logic
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
	return &Build{
		Index: buildNum,
		Label: vals[0],
	}, nil
}

// String returns the build version string
func (b *Build) String() string {
	if b.Index > 0 {
		return fmt.Sprintf("%s.%d", b.Label, b.Index)
	} else {
		return ""
	}
}

func (b *Build) Bump() *Build {
	return NewBuild(b.Label, b.Index+1)
}
