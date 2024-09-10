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
