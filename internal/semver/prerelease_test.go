package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestBumpPreRelease ensures that bumping the version correctly returns a new subversion instance
func TestBumpPreRelease(t *testing.T) {
	buildLabel := "build"
	preReleaseLabels := []string{"alpha", "beta", "rc"}
	tests := []struct {
		input      string
		bumpType   int
		expected   string
		shouldFail bool
	}{
		{"foo", PreReleaseNext, "", true},
		{"rc", PreReleaseNext, "", true},
		{"beta", PreReleaseNext, "rc", false},
		{"2.5.1", PreReleaseNext, "alpha", false},
		{"alpha", PreReleaseNext, "beta", false},
		{"alpha", PreReleaseMinor, "alpha.0.1", false},
		{"2.5.1", PreReleaseBuild, "2.5.1+build.1", false},
		{"0.0.0", PreReleaseMajor, "1", false},
		{"1.0.0", PreReleaseMinor, "1.1", false},
		{"1.1.0", PreReleasePatch, "1.1.1", false},
		{"2.5.1", PreReleaseMajor, "3", false},
		{"2.5.1", PreReleaseMinor, "2.6", false},
		{"2.5.1", PreReleasePatch, "2.5.2", false},
		{"2.5.1+build.1", PreReleaseBuild, "2.5.1+build.2", false},
		{"2.5.1", PreReleaseBuild, "2.5.1+build.1", false},
		// TODO: Add test cases for pre-release and build versions

	}

	for _, test := range tests {
		subv, err := ParsePrereleaseVersion(test.input)
		if err != nil {
			t.Fatalf("Unexpected error for input %s: %v", test.input, err)
		}

		bumped, err := subv.Bump(test.bumpType, preReleaseLabels, buildLabel)

		if test.shouldFail {
			// Assert an error was returned
			assert.Error(t, err, "expected an error for version %s", test.input)
		} else {
			// Assert no error was returned
			assert.NoError(t, err, "unexpected error for version %s", test.input)
			result := bumped.String()
			assert.Equal(t, test.expected, result, "expected %s, got %s", test.expected, result)
		}

	}
}

func TestParsePrereleaseVersion(t *testing.T) {
	tests := []struct {
		versionStr string
		expected   string
		shouldFail bool
	}{
		{"1.2.3", "1.2.3", false},
		{"alpha", "alpha", false},
		{"alpha.1", "alpha.1", false},
		{"alpha+build.1", "alpha+build.1", false},
	}

	for _, test := range tests {
		t.Run(test.versionStr, func(t *testing.T) {
			semVer, err := ParsePrereleaseVersion(test.versionStr)

			if test.shouldFail {
				// Assert an error was returned
				assert.Error(t, err, "expected an error for version %s", test.versionStr)
			} else {
				// Assert no error was returned
				assert.NoError(t, err, "unexpected error for version %s", test.versionStr)

				semVerStr := semVer.String()
				// Assert the parsed version matches the expected output
				assert.Equal(t, test.expected, semVer.String(), "expected %s, got %s", test.expected, semVerStr)
			}
		})
	}
}
