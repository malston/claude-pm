// ABOUTME: Tests for snapshot functionality
// ABOUTME: Validates creating a profile from current Claude state
package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotFromState(t *testing.T) {
	// Create temp directories
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	pluginsDir := filepath.Join(claudeDir, "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create mock installed_plugins.json
	pluginsData := map[string]interface{}{
		"version": 2,
		"plugins": map[string]interface{}{
			"superpowers@superpowers-marketplace": []map[string]interface{}{{
				"scope":       "user",
				"version":     "1.0.0",
				"installPath": "/path/to/plugin",
			}},
			"frontend-design@claude-code-plugins": []map[string]interface{}{{
				"scope":       "user",
				"version":     "2.0.0",
				"installPath": "/path/to/plugin2",
			}},
		},
	}
	writeJSON(t, filepath.Join(pluginsDir, "installed_plugins.json"), pluginsData)

	// Create mock known_marketplaces.json
	marketplacesData := map[string]interface{}{
		"superpowers-marketplace": map[string]interface{}{
			"source": map[string]interface{}{
				"source": "github",
				"repo":   "anthropics/superpowers-marketplace",
			},
		},
		"claude-code-plugins": map[string]interface{}{
			"source": map[string]interface{}{
				"source": "github",
				"repo":   "anthropics/claude-code-plugins",
			},
		},
	}
	writeJSON(t, filepath.Join(pluginsDir, "known_marketplaces.json"), marketplacesData)

	// Create mock ~/.claude.json with MCP servers
	claudeJSON := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"context7": map[string]interface{}{
				"type":    "stdio",
				"command": "npx",
				"args":    []string{"-y", "@upstash/context7-mcp"},
				"env":     map[string]string{},
			},
		},
	}
	claudeJSONPath := filepath.Join(tmpDir, ".claude.json")
	writeJSON(t, claudeJSONPath, claudeJSON)

	// Create snapshot
	p, err := Snapshot("test-snapshot", claudeDir, claudeJSONPath)
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	// Verify profile
	if p.Name != "test-snapshot" {
		t.Errorf("Name mismatch: got %q", p.Name)
	}

	if len(p.Plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(p.Plugins))
	}

	if len(p.Marketplaces) != 2 {
		t.Errorf("Expected 2 marketplaces, got %d", len(p.Marketplaces))
	}

	if len(p.MCPServers) != 1 {
		t.Errorf("Expected 1 MCP server, got %d", len(p.MCPServers))
	}

	// Verify MCP server details
	if len(p.MCPServers) > 0 {
		mcp := p.MCPServers[0]
		if mcp.Name != "context7" {
			t.Errorf("MCP name mismatch: got %q", mcp.Name)
		}
		if mcp.Command != "npx" {
			t.Errorf("MCP command mismatch: got %q", mcp.Command)
		}
	}
}

func TestSnapshotEmptyState(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	claudeJSONPath := filepath.Join(tmpDir, ".claude.json")

	// Don't create any files - test with empty state
	p, err := Snapshot("empty", claudeDir, claudeJSONPath)
	if err != nil {
		t.Fatalf("Snapshot failed on empty state: %v", err)
	}

	if p.Name != "empty" {
		t.Errorf("Name mismatch: got %q", p.Name)
	}

	if len(p.Plugins) != 0 {
		t.Errorf("Expected 0 plugins, got %d", len(p.Plugins))
	}
}

func writeJSON(t *testing.T, path string, data interface{}) {
	t.Helper()
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, bytes, 0644); err != nil {
		t.Fatal(err)
	}
}
