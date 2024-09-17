package version

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	version := "1.2.3"
	expected := Version{major: 1, minor: 2, patch: 3}
	actual, _ := ParseVersion(version)
	if *actual != expected {
		t.Errorf("Expected %v but got %v", expected, *actual)
	}
}

func TestVersion_String(t *testing.T) {
	version := Version{major: 1, minor: 2, patch: 3}
	expected := "1.2.3"
	actual := version.String()
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

// TestBumpSubversion ensures that bumping the version correctly returns a new subversion instance
func TestBumpSubversion(t *testing.T) {
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
		subv, err := ParseSubversion(test.input)
		if err != nil {
			t.Fatalf("Unexpected error for input %s: %v", test.input, err)
		}

		bumped := subv.Bump(test.bumpType)
		if result := bumped.String(); result != test.expected {
			t.Errorf("For input %s and bumpType %d, expected %s, but got %s", test.input, test.bumpType, test.expected, result)
		}
	}
}
