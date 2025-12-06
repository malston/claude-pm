// ABOUTME: Cleanup command implementation for removing stale plugin entries
// ABOUTME: Removes plugins from installed_plugins.json where paths don't exist
package commands

import (
	"fmt"

	"github.com/malston/claude-pm/internal/claude"
	"github.com/spf13/cobra"
)

var (
	cleanupReinstall bool
	cleanupDryRun    bool
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove stale plugin entries",
	Long: `Remove plugin entries where the installation path no longer exists.

This cleans up orphaned entries in installed_plugins.json after marketplace
structure changes or manual plugin deletions.`,
	RunE: runCleanup,
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().BoolVar(&cleanupReinstall, "reinstall", false, "Offer to reinstall removed plugins")
	cleanupCmd.Flags().BoolVar(&cleanupDryRun, "dry-run", false, "Show what would be removed without making changes")
}

func runCleanup(cmd *cobra.Command, args []string) error {
	// Load plugins
	plugins, err := claude.LoadPlugins(claudeDir)
	if err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}

	// Analyze issues
	pathIssues := analyzePathIssues(plugins)

	// Filter for non-fixable issues (truly missing)
	stalePlugins := []PathIssue{}
	for _, issue := range pathIssues {
		if !issue.CanAutoFix {
			stalePlugins = append(stalePlugins, issue)
		}
	}

	if len(stalePlugins) == 0 {
		fmt.Println("✓ No stale plugins found")
		return nil
	}

	if cleanupDryRun {
		fmt.Printf("Would remove %d stale plugin entries:\n\n", len(stalePlugins))
	} else {
		fmt.Printf("Found %d stale plugin entries\n\n", len(stalePlugins))
	}

	// Show what will be removed
	for _, issue := range stalePlugins {
		fmt.Printf("  • %s\n", issue.PluginName)
		fmt.Printf("    Path: %s\n", issue.InstallPath)
	}

	if cleanupDryRun {
		fmt.Println("\nRun without --dry-run to remove these entries")
		return nil
	}

	fmt.Println("\nThese plugins will be removed from installed_plugins.json")
	fmt.Print("Continue? [y/N]: ")

	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Cancelled")
		return nil
	}

	// Remove stale entries
	removed := 0
	for _, issue := range stalePlugins {
		if plugins.DisablePlugin(issue.PluginName) {
			removed++
		}
	}

	// Save updated plugins
	if err := claude.SavePlugins(claudeDir, plugins); err != nil {
		return fmt.Errorf("failed to save plugins: %w", err)
	}

	fmt.Printf("\n✓ Removed %d stale plugin entries\n", removed)

	if cleanupReinstall && removed > 0 {
		fmt.Println("\nTo reinstall these plugins, use:")
		for _, issue := range stalePlugins {
			fmt.Printf("  claude plugin install %s\n", issue.PluginName)
		}
	}

	return nil
}
