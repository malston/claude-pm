// ABOUTME: MCP command implementation for managing MCP servers
// ABOUTME: Lists and shows information about MCP servers provided by plugins
package commands

import (
	"fmt"
	"sort"

	"github.com/malston/claude-pm/internal/claude"
	"github.com/malston/claude-pm/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Manage MCP servers",
	Long:  `List and manage MCP servers provided by Claude Code plugins.`,
}

var mcpListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all MCP servers",
	Long:  `Display MCP servers grouped by the plugin that provides them.`,
	RunE:  runMCPList,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
	mcpCmd.AddCommand(mcpListCmd)
}

func runMCPList(cmd *cobra.Command, args []string) error {
	// Load plugins
	plugins, err := claude.LoadPlugins(claudeDir)
	if err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}

	// Discover MCP servers
	mcpServers, err := mcp.DiscoverMCPServers(plugins)
	if err != nil {
		return fmt.Errorf("failed to discover MCP servers: %w", err)
	}

	if len(mcpServers) == 0 {
		fmt.Println("No MCP servers found in installed plugins.")
		return nil
	}

	// Sort by plugin name for consistent output
	sort.Slice(mcpServers, func(i, j int) bool {
		return mcpServers[i].PluginName < mcpServers[j].PluginName
	})

	// Print header
	fmt.Println("=== MCP Servers by Plugin ===\n")

	// Count total servers
	totalServers := 0
	for _, pluginServers := range mcpServers {
		totalServers += len(pluginServers.Servers)
	}

	// Print each plugin's MCP servers
	for _, pluginServers := range mcpServers {
		fmt.Printf("✓ %s\n", pluginServers.PluginName)

		// Sort server names
		serverNames := make([]string, 0, len(pluginServers.Servers))
		for name := range pluginServers.Servers {
			serverNames = append(serverNames, name)
		}
		sort.Strings(serverNames)

		// Print each server
		for _, serverName := range serverNames {
			server := pluginServers.Servers[serverName]
			fmt.Printf("   ✓ %s\n", serverName)
			fmt.Printf("      Command: %s\n", server.Command)
			if len(server.Args) > 0 {
				fmt.Printf("      Args:    %v\n", server.Args)
			}
			if len(server.Env) > 0 {
				fmt.Printf("      Env:     %d variables\n", len(server.Env))
			}
		}
		fmt.Println()
	}

	fmt.Printf("Total: %d MCP servers from %d plugins\n", totalServers, len(mcpServers))

	return nil
}
