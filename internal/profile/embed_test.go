// ABOUTME: Tests for embedded profile functionality
// ABOUTME: Validates that default profiles are properly embedded and extracted
package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetEmbeddedProfile(t *testing.T) {
	p, err := GetEmbeddedProfile("default")
	if err != nil {
		t.Fatalf("Failed to get embedded default profile: %v", err)
	}

	if p.Name != "default" {
		t.Errorf("Expected name 'default', got %q", p.Name)
	}

	if p.Description == "" {
		t.Error("Expected description to be set")
	}
}

func TestGetEmbeddedProfileNotFound(t *testing.T) {
	_, err := GetEmbeddedProfile("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent profile")
	}
}

func TestEnsureDefaultProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	profilesDir := filepath.Join(tmpDir, "profiles")

	err := EnsureDefaultProfiles(profilesDir)
	if err != nil {
		t.Fatalf("EnsureDefaultProfiles failed: %v", err)
	}

	// Check that default.json was created
	defaultPath := filepath.Join(profilesDir, "default.json")
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Error("default.json was not created")
	}

	// Load and verify
	p, err := Load(profilesDir, "default")
	if err != nil {
		t.Fatalf("Failed to load extracted profile: %v", err)
	}

	if p.Name != "default" {
		t.Errorf("Expected name 'default', got %q", p.Name)
	}
}

func TestEnsureDefaultProfilesDoesNotOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	profilesDir := filepath.Join(tmpDir, "profiles")
	os.MkdirAll(profilesDir, 0755)

	// Create a custom default.json
	customContent := `{"name": "default", "description": "Custom description"}`
	defaultPath := filepath.Join(profilesDir, "default.json")
	os.WriteFile(defaultPath, []byte(customContent), 0644)

	// Run ensure - should not overwrite
	err := EnsureDefaultProfiles(profilesDir)
	if err != nil {
		t.Fatalf("EnsureDefaultProfiles failed: %v", err)
	}

	// Verify custom content is preserved
	p, err := Load(profilesDir, "default")
	if err != nil {
		t.Fatalf("Failed to load profile: %v", err)
	}

	if p.Description != "Custom description" {
		t.Errorf("Profile was overwritten, got description: %q", p.Description)
	}
}
