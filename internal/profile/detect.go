// ABOUTME: Project detection logic for profile suggestions
// ABOUTME: Matches project files and content against profile detect rules
package profile

import (
	"os"
	"path/filepath"
	"strings"
)

// Detect checks if a profile's detect rules match the given directory
func Detect(dir string, p *Profile) (bool, error) {
	rules := p.Detect

	// No rules means no match
	if len(rules.Files) == 0 && len(rules.Contains) == 0 {
		return false, nil
	}

	// Check all required files exist
	for _, file := range rules.Files {
		filePath := filepath.Join(dir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return false, nil
		}
	}

	// Check file content patterns
	for file, pattern := range rules.Contains {
		filePath := filepath.Join(dir, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return false, nil // File doesn't exist or can't be read
		}

		if !strings.Contains(string(content), pattern) {
			return false, nil
		}
	}

	return true, nil
}

// FindMatchingProfiles returns all profiles that match the given directory
func FindMatchingProfiles(dir string, profiles []*Profile) []*Profile {
	var matches []*Profile

	for _, p := range profiles {
		match, err := Detect(dir, p)
		if err == nil && match {
			matches = append(matches, p)
		}
	}

	return matches
}

// SuggestProfile finds the best matching profile for a directory
// Returns nil if no profiles match
func SuggestProfile(dir string, profiles []*Profile) *Profile {
	matches := FindMatchingProfiles(dir, profiles)

	if len(matches) == 0 {
		return nil
	}

	// Return the first match (profiles should be ordered by priority)
	return matches[0]
}
