package semver

import (
	"github.com/stretchr/testify/assert"
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
		{"1.0.0-alpha+build.1", "1.0.0-alpha+build.1", false},
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

				// Assert the parsed version matches the expected output
				assert.Equal(t, test.expected, semVer.String(), "expected %s, got %s", test.expected, semVer.String())
			}
		})
	}
}

func TestParsePrereleaseVersion(t *testing.T) {
	tests := []struct {
		versionStr string
		expected   string
		shouldFail bool
	}{
		{"alpha", "alpha", false},
		{"alpha.1", "alpha.1", false},
	}

	for _, test := range tests {
		t.Run(test.versionStr, func(t *testing.T) {
			prereleaseVersion, err := ParsePrereleaseVersion(test.versionStr)

			if test.shouldFail {
				// Assert an error was returned
				assert.Error(t, err, "expected an error for version %s", test.versionStr)
			} else {
				// Assert no error was returned
				assert.NoError(t, err, "unexpected error for version %s", test.versionStr)

				// Assert the parsed version matches the expected output
				assert.Equal(t, test.expected, prereleaseVersion.String(), "expected %s, got %s", test.expected, prereleaseVersion.String())
			}
		})
	}

}

// TestBumpPreRelease ensures that bumping the version correctly returns a new subversion instance
func TestBumpPreRelease(t *testing.T) {
	tests := []struct {
		input    string
		bumpType int
		expected string
	}{
		{"0.0.0", VersionMajor, "1"},
		{"1.0.0", VersionMinor, "1.1"},
		{"1.1.0", VersionPatch, "1.1.1"},
		{"2.5.1", VersionMajor, "3"},
		{"2.5.1", VersionMinor, "2.6"},
		{"2.5.1", VersionPatch, "2.5.2"},
	}

	for _, test := range tests {
		subv, err := ParsePrereleaseVersion(test.input)
		if err != nil {
			t.Fatalf("Unexpected error for input %s: %v", test.input, err)
		}

		bumped := subv.Bump(test.bumpType)
		if result := bumped.String(); result != test.expected {
			t.Errorf("For input %s and bumpType %d, expected %s, but got %s", test.input, test.bumpType, test.expected, result)
		}
	}
}
