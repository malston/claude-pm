// ABOUTME: Tests for project detection logic
// ABOUTME: Validates matching projects to profiles based on detect rules
package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMatchesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create package.json
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatal(err)
	}

	profile := &Profile{
		Name: "nodejs",
		Detect: DetectRules{
			Files: []string{"package.json"},
		},
	}

	match, err := Detect(tmpDir, profile)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !match {
		t.Error("Expected profile to match, but it didn't")
	}
}

func TestDetectMatchesContains(t *testing.T) {
	tmpDir := t.TempDir()

	// Create package.json with react
	content := `{"name": "test", "dependencies": {"react": "^18.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	profile := &Profile{
		Name: "react",
		Detect: DetectRules{
			Files:    []string{"package.json"},
			Contains: map[string]string{"package.json": "react"},
		},
	}

	match, err := Detect(tmpDir, profile)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !match {
		t.Error("Expected profile to match react dependency")
	}
}

func TestDetectNoMatchMissingFile(t *testing.T) {
	tmpDir := t.TempDir()

	profile := &Profile{
		Name: "go",
		Detect: DetectRules{
			Files: []string{"go.mod"},
		},
	}

	match, err := Detect(tmpDir, profile)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if match {
		t.Error("Expected no match for missing go.mod")
	}
}

func TestDetectNoMatchContentMismatch(t *testing.T) {
	tmpDir := t.TempDir()

	// Create package.json WITHOUT react
	content := `{"name": "test", "dependencies": {"vue": "^3.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	profile := &Profile{
		Name: "react",
		Detect: DetectRules{
			Files:    []string{"package.json"},
			Contains: map[string]string{"package.json": "react"},
		},
	}

	match, err := Detect(tmpDir, profile)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if match {
		t.Error("Expected no match - vue project should not match react profile")
	}
}

func TestDetectEmptyRulesNoMatch(t *testing.T) {
	tmpDir := t.TempDir()

	profile := &Profile{
		Name:   "empty",
		Detect: DetectRules{},
	}

	match, err := Detect(tmpDir, profile)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if match {
		t.Error("Expected no match for profile with no detect rules")
	}
}

func TestFindMatchingProfiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	profiles := []*Profile{
		{Name: "nodejs", Detect: DetectRules{Files: []string{"package.json"}}},
		{Name: "go", Detect: DetectRules{Files: []string{"go.mod"}}},
		{Name: "rust", Detect: DetectRules{Files: []string{"Cargo.toml"}}},
	}

	matches := FindMatchingProfiles(tmpDir, profiles)

	if len(matches) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matches))
	}

	if matches[0].Name != "go" {
		t.Errorf("Expected 'go' profile to match, got %q", matches[0].Name)
	}
}
