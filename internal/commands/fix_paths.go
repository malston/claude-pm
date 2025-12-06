// ABOUTME: Fix-paths command implementation for correcting plugin path issues
// ABOUTME: Fixes known path issues in installed_plugins.json
package commands

import (
	"fmt"

	"github.com/malston/claude-pm/internal/claude"
	"github.com/spf13/cobra"
)

var fixPathsCmd = &cobra.Command{
	Use:   "fix-paths",
	Short: "Fix known plugin path issues",
	Long: `Fix plugin paths that are incorrect due to Claude CLI bugs.

This fixes paths for plugins with isLocal: true where the CLI didn't account
for marketplace subdirectories like /plugins/ or /skills/.

A backup of installed_plugins.json is created before making changes.`,
	RunE: runFixPaths,
}

func init() {
	rootCmd.AddCommand(fixPathsCmd)
}

func runFixPaths(cmd *cobra.Command, args []string) error {
	// Load plugins
	plugins, err := claude.LoadPlugins(claudeDir)
	if err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}

	// Analyze issues
	pathIssues := analyzePathIssues(plugins)

	// Filter for fixable issues only
	fixableIssues := []PathIssue{}
	for _, issue := range pathIssues {
		if issue.CanAutoFix {
			fixableIssues = append(fixableIssues, issue)
		}
	}

	if len(fixableIssues) == 0 {
		fmt.Println("✓ No fixable path issues found")
		return nil
	}

	fmt.Printf("Found %d fixable path issues\n\n", len(fixableIssues))

	// Show what will be fixed
	for _, issue := range fixableIssues {
		fmt.Printf("  %s\n", issue.PluginName)
		fmt.Printf("    %s → %s\n", issue.InstallPath, issue.ExpectedPath)
	}

	// Apply fixes
	fmt.Println("\nApplying fixes...")
	fixed := 0
	for _, issue := range fixableIssues {
		if plugin, exists := plugins.Plugins[issue.PluginName]; exists {
			plugin.InstallPath = issue.ExpectedPath
			plugins.Plugins[issue.PluginName] = plugin
			fixed++
		}
	}

	// Save updated plugins
	if err := claude.SavePlugins(claudeDir, plugins); err != nil {
		return fmt.Errorf("failed to save plugins: %w", err)
	}

	fmt.Printf("\n✓ Fixed %d plugin paths\n", fixed)
	fmt.Println("\nRun 'claude-pm status' to verify the fixes")

	return nil
}
