package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestBumpPreRelease ensures that bumping the version correctly returns a new subversion instance
func TestBumpPreRelease(t *testing.T) {
	preReleaseLabels := []string{"alpha", "beta", "rc"}
	tests := []struct {
		input      string
		bumpType   int
		expected   string
		shouldFail bool
	}{
		{"foo", prNext, "", true},
		{"rc", prNext, "", true},
		{"beta", prNext, "rc", false},
		{"2.5.1", prNext, "alpha", false},
		{"alpha", prNext, "beta", false},
		{"alpha", prMinor, "alpha.0.1", false},
		{"0.0.0", prMajor, "alpha", false},
		{"1.0.0", prMinor, "alpha.1.1", false},
		{"1.1.0", prPatch, "alpha.1.1.1", false},
		{"2.5.1", prMajor, "alpha.2", false},
		{"2.5.1", prMinor, "alpha.2.6", false},
		{"2.5.1", prPatch, "alpha.2.5.2", false},
	}

	for _, test := range tests {
		subv, err := ParsePrereleaseVersion(test.input)
		if err != nil {
			t.Fatalf("Unexpected error for input %s: %v", test.input, err)
		}

		bumped, err := subv.Bump(test.bumpType, preReleaseLabels)

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
		{"1.2", "1.2", false},
		{"1", "1", false},
		{"alpha", "alpha", false},
		{"alpha.1", "alpha.1", false},
		{"alpha.0.1", "alpha.0.1", false},
		{"alpha.0.0.1", "alpha.0.0.1", false},
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
