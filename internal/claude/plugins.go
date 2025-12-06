// ABOUTME: Data structures and functions for managing Claude Code plugins
// ABOUTME: Handles reading and writing installed_plugins.json
package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PluginRegistry represents the installed_plugins.json file structure
type PluginRegistry struct {
	Version int                       `json:"version"`
	Plugins map[string]PluginMetadata `json:"plugins"`
}

// PluginMetadata represents metadata for an installed plugin
type PluginMetadata struct {
	Version      string `json:"version"`
	InstalledAt  string `json:"installedAt"`
	LastUpdated  string `json:"lastUpdated"`
	InstallPath  string `json:"installPath"`
	GitCommitSha string `json:"gitCommitSha"`
	IsLocal      bool   `json:"isLocal"`
}

// LoadPlugins reads and parses the installed_plugins.json file
func LoadPlugins(claudeDir string) (*PluginRegistry, error) {
	pluginsPath := filepath.Join(claudeDir, "plugins", "installed_plugins.json")

	data, err := os.ReadFile(pluginsPath)
	if err != nil {
		return nil, err
	}

	var registry PluginRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}

	return &registry, nil
}

// SavePlugins writes the plugin registry back to installed_plugins.json
func SavePlugins(claudeDir string, registry *PluginRegistry) error {
	pluginsPath := filepath.Join(claudeDir, "plugins", "installed_plugins.json")

	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pluginsPath, data, 0644)
}

// PathExists checks if a plugin's install path actually exists
func (p *PluginMetadata) PathExists() bool {
	if p.InstallPath == "" {
		return false
	}
	_, err := os.Stat(p.InstallPath)
	return err == nil
}

// DisablePlugin removes a plugin from the registry
func (r *PluginRegistry) DisablePlugin(pluginName string) bool {
	if _, exists := r.Plugins[pluginName]; !exists {
		return false // Plugin not found
	}
	delete(r.Plugins, pluginName)
	return true
}

// EnablePlugin adds a plugin back to the registry
// Note: This requires having the plugin metadata available
func (r *PluginRegistry) EnablePlugin(pluginName string, metadata PluginMetadata) {
	r.Plugins[pluginName] = metadata
}

// PluginExists checks if a plugin is in the registry
func (r *PluginRegistry) PluginExists(pluginName string) bool {
	_, exists := r.Plugins[pluginName]
	return exists
}
