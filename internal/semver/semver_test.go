package semver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	sv "golang.org/x/mod/semver"
	"testing"
)

func TestParseSemVersion(t *testing.T) {
	tests := []struct {
		versionStr string
		expected   string
		shouldFail bool
	}{
		{"1.2.3", "1.2.3", false},
		{"1.0.0-alpha", "1.0.0-alpha", false},
		{"1.0.0-alpha.1", "1.0.0-alpha.1", false},
		{"1.0.0-alpha+build.1", "1.0.0-alpha+build.1", false},
		{"1.0.0+build.1", "1.0.0+build.1", false},
		{"2.0", "", true},     // Should fail as it's not a valid semantic version
		{"", "", true},        // Empty version string should fail
		{"1.2.3.4", "", true}, // Invalid semver format
	}

	for _, test := range tests {
		t.Run(test.versionStr, func(t *testing.T) {
			semVer, err := ParseSemVersion(test.versionStr)

			if test.shouldFail {
				// Assert an error was returned
				assert.Error(t, err, "expected an error for version %s", test.versionStr)
			} else {
				// Assert no error was returned
				assert.NoError(t, err, "unexpected error for version %s", test.versionStr)

				semVerStr := semVer.String()
				// Assert the parsed version matches the expected output
				assert.Equal(t, sv.IsValid("v"+semVerStr), true, "expected a valid semver version")
				assert.Equal(t, test.expected, semVer.String(), "expected %s, got %s", test.expected, semVerStr)
			}
		})
	}
}

// TestBumpPreRelease ensures that bumping the version correctly returns a new subversion instance
func TestSemVersion_Bump(t *testing.T) {
	buildLabel := "ptgoetz"
	preReleaseLabels := []string{"alpha", "beta", "rc"}
	tests := []struct {
		input      string
		bumpType   VersionPart
		expected   string
		shouldFail bool
	}{
		{"0.0.0", Patch, "0.0.1", false},
		{"0.0.1", Minor, "0.1.0", false},
		{"0.1.0", Major, "1.0.0", false},
		{"1.0.0", PreReleaseNext, "1.0.0-alpha", false},
		{"1.0.0-alpha", PreReleaseBuild, "1.0.0-alpha+ptgoetz.1", false},
		{"1.0.0-alpha+ptgoetz.1", PreReleaseBuild, "1.0.0-alpha+ptgoetz.2", false},
		{"1.0.0-alpha+ptgoetz.1", PreReleaseNext, "1.0.0-beta", false},
		{"1.0.0-beta", PreReleasePatch, "1.0.0-beta.0.0.1", false},
		{"1.0.0-beta", PreReleaseMinor, "1.0.0-beta.0.1", false},
		{"1.0.0-beta", PreReleaseMajor, "1.0.0-beta.1", false},
	}

	for _, test := range tests {
		version, err := ParseSemVersion(test.input)
		if err != nil {
			t.Fatalf("Unexpected error for input %s: %v", test.input, err)
		}
		assert.Equal(t, sv.IsValid("v"+version.String()), true, "expected a valid semver version")

		bumped, err := version.Bump(test.bumpType, preReleaseLabels, buildLabel)

		if test.shouldFail {
			// Assert an error was returned
			assert.Error(t, err, "expected an error for version %s", test.input)
		} else {

			// Assert no error was returned
			assert.NoError(t, err, "unexpected error for version %s", test.input)
			result := bumped.String()
			assert.Equal(t, sv.IsValid("v"+result), true, "expected a valid semver version")
			assert.Equal(t, test.expected, result, "expected %s, got %s", test.expected, result)
		}

	}
}

func TestSemVersion_Compare(t *testing.T) {
	tests := []struct {
		version1 string
		version2 string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"1.1.0", "1.0.0", 1},
		{"1.0.0", "1.1.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0", "2.0.0", -1},
		{"1.0.0-alpha", "1.0.0-beta", -1},
		{"1.0.0-beta", "1.0.0-alpha", 1},
		{"1.0.0-alpha", "1.0.0-alpha", 0},
		{"1.0.0-alpha+build.1", "1.0.0-alpha+build.2", -1},
		{"1.0.0-alpha+build.2", "1.0.0-alpha+build.1", 1},
		{"1.0.0-alpha+build.1", "1.0.0-alpha+build.1", 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s vs %s", test.version1, test.version2), func(t *testing.T) {
			semVer1, err := ParseSemVersion(test.version1)
			assert.NoError(t, err, "unexpected error for version %s", test.version1)

			semVer2, err := ParseSemVersion(test.version2)
			assert.NoError(t, err, "unexpected error for version %s", test.version2)

			result := semVer1.Compare(semVer2)
			assert.Equal(t, test.expected, result, "expected %d, got %d", test.expected, result)
		})
	}
}
