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
		{"2.0", "", true},     // Should fail as it's not a valid semantic rootVersion
		{"", "", true},        // Empty rootVersion string should fail
		{"1.2.3.4", "", true}, // Invalid semver format
	}

	for _, test := range tests {
		t.Run(test.versionStr, func(t *testing.T) {
			semVer, err := ParseSemVersion(test.versionStr)

			if test.shouldFail {
				// Assert an error was returned
				assert.Error(t, err, "expected an error for rootVersion %s", test.versionStr)
			} else {
				// Assert no error was returned
				assert.NoError(t, err, "unexpected error for rootVersion %s", test.versionStr)

				semVerStr := semVer.String()
				// Assert the parsed rootVersion matches the expected output
				assert.Equal(t, sv.IsValid("v"+semVerStr), true, "expected a valid semver rootVersion")
				assert.Equal(t, test.expected, semVer.String(), "expected %s, got %s", test.expected, semVerStr)
			}
		})
	}
}

// TestBumpPreRelease ensures that bumping the rootVersion correctly returns a new subversion instance
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
		{"1.0.0", PreRelease, "1.0.0-alpha", false},
		{"1.0.0", PreReleaseMajor, "1.0.0-alpha", false},
		{"1.0.0-alpha", PreReleaseBuild, "1.0.0-alpha+ptgoetz.1", false},
		{"1.0.0-alpha+ptgoetz.1", PreReleaseBuild, "1.0.0-alpha+ptgoetz.2", false},
		{"1.0.0-alpha+ptgoetz.1", PreRelease, "1.0.0-beta", false},
		{"1.0.0-beta", PreReleasePatch, "1.0.0-beta.0.0.1", false},
		{"1.0.0-beta", PreReleaseMinor, "1.0.0-beta.0.1", false},
		{"1.0.0-beta", PreReleaseMajor, "1.0.0-beta.1", false},
	}

	for _, test := range tests {
		version, err := ParseSemVersion(test.input)
		if err != nil {
			t.Fatalf("Unexpected error for input %s: %v", test.input, err)
		}
		assert.Equal(t, sv.IsValid("v"+version.String()), true, "expected a valid semver rootVersion")

		bumped, err := version.Bump(test.bumpType, preReleaseLabels, buildLabel)

		if test.shouldFail {
			// Assert an error was returned
			assert.Error(t, err, "expected an error for rootVersion %s", test.input)
		} else {

			// Assert no error was returned
			assert.NoError(t, err, "unexpected error for rootVersion %s", test.input)
			result := bumped.String()
			assert.Equal(t, sv.IsValid("v"+result), true, "expected a valid semver rootVersion")
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
			assert.NoError(t, err, "unexpected error for rootVersion %s", test.version1)

			semVer2, err := ParseSemVersion(test.version2)
			assert.NoError(t, err, "unexpected error for rootVersion %s", test.version2)

			result := semVer1.Compare(semVer2)
			assert.Equal(t, test.expected, result, "expected %d, got %d", test.expected, result)
		})
	}
}

func TestAccessors(t *testing.T) {
	v, err := ParseSemVersion("1.2.3-alpha.4.5.6+build.1")
	assert.NoError(t, err, "unexpected error for version string  '%s'", "1.2.3-alpha.4.5.6+build.1")
	assert.Equal(t, 1, v.RootVersion().Major(), "expected major version to be 1")
	assert.Equal(t, 2, v.RootVersion().Minor(), "expected minor version to be 2")
	assert.Equal(t, 3, v.RootVersion().Patch(), "expected major version to be 3")
	assert.Equal(t, "alpha", v.PreReleaseVersion().Label(), "expected pre-release label to be 'alpha'")
	assert.Equal(t, 4, v.PreReleaseVersion().Version().Major(), "expected pre-release major version to be 4")
	assert.Equal(t, 5, v.PreReleaseVersion().Version().Minor(), "expected pre-release minor version to be 5")
	assert.Equal(t, 6, v.PreReleaseVersion().Version().Patch(), "expected pre-release patch version to be 6")
	assert.Equal(t, "build.1", v.BuildVersion().String(), "expected build version to be 'build.1'")
	assert.Equal(t, 1, v.BuildVersion().Number(), "expected build number to be 1")
}
