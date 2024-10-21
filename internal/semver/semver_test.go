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
		{"1.0.0-alpha.1", "1.0.0-alpha.1", false},
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

				semVerStr := semVer.String()
				// Assert the parsed version matches the expected output
				assert.Equal(t, test.expected, semVer.String(), "expected %s, got %s", test.expected, semVerStr)
			}
		})
	}
}
