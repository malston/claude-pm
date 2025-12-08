// ABOUTME: Embedded default profiles for claudeup
// ABOUTME: Uses Go embed to bundle profiles in the binary
package profile

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
)

//go:embed profiles/*.json
var embeddedProfiles embed.FS

// EnsureDefaultProfiles extracts embedded profiles to the profiles directory
// if they don't already exist
func EnsureDefaultProfiles(profilesDir string) error {
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return err
	}

	entries, err := embeddedProfiles.ReadDir("profiles")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		destPath := filepath.Join(profilesDir, name)

		// Skip if file already exists
		if _, err := os.Stat(destPath); err == nil {
			continue
		}

		// Read embedded file
		data, err := embeddedProfiles.ReadFile("profiles/" + name)
		if err != nil {
			return err
		}

		// Write to profiles directory
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

// GetEmbeddedProfile returns an embedded profile by name
func GetEmbeddedProfile(name string) (*Profile, error) {
	data, err := embeddedProfiles.ReadFile("profiles/" + name + ".json")
	if err != nil {
		return nil, err
	}

	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// ListEmbeddedProfiles returns all embedded profiles
func ListEmbeddedProfiles() ([]*Profile, error) {
	entries, err := embeddedProfiles.ReadDir("profiles")
	if err != nil {
		return nil, err
	}

	var profiles []*Profile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if len(name) < 6 || name[len(name)-5:] != ".json" {
			continue
		}

		profileName := name[:len(name)-5]
		p, err := GetEmbeddedProfile(profileName)
		if err != nil {
			continue
		}
		profiles = append(profiles, p)
	}

	return profiles, nil
}
