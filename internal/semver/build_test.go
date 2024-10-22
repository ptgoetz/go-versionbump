package semver

import "testing"

func TestBuild(t *testing.T) {
	tests := []struct {
		input    string
		bumpType int
		expected string
	}{
		{"build.1", PreReleaseBuild, "build.2"},
		{"", PreReleaseBuild, "build.1"},
		{"", PreReleaseBuild, "build.1"},
		{"foo.1", PreReleaseBuild, "foo.2"},
	}
	for _, test := range tests {
		build, err := ParseBuild(test.input)
		if err != nil {
			t.Fatalf("Unexpected error for input %s: %v", test.input, err)
		}

		bumped := build.Bump()
		if result := bumped.String(); result != test.expected {
			t.Errorf("For input '%s' and bumpType %d, expected %s, but got %s", test.input, test.bumpType, test.expected, result)
		}
	}
}
