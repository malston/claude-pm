// ABOUTME: Integration tests for enable/disable commands
// ABOUTME: Tests plugin and MCP server enable/disable workflows
package integration

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/claudeup/claudeup/internal/claude"
	"github.com/claudeup/claudeup/internal/config"
)

var _ = Describe("PluginDisableEnable", func() {
	var env *TestEnv

	BeforeEach(func() {
		env = SetupTestEnv()

		env.CreateMarketplace("test-marketplace", "test/repo")
		env.CreatePlugin("test-plugin", "test-marketplace", "1.0.0", nil)

		pluginPath := filepath.Join(env.ClaudeDir, "plugins", "marketplaces", "test-marketplace", "plugins", "test-plugin")
		env.CreatePluginRegistry(map[string]claude.PluginMetadata{
			"test-plugin@test-marketplace": {
				Version:      "1.0.0",
				InstallPath:  pluginPath,
				GitCommitSha: "abc123",
				IsLocal:      true,
			},
		})
	})

	AfterEach(func() {
		env.Cleanup()
	})

	It("disables and re-enables a plugin", func() {
		Expect(env.PluginExists("test-plugin@test-marketplace")).To(BeTrue(), "Plugin should exist initially")
		Expect(env.PluginCount()).To(Equal(1))

		registry := env.LoadPluginRegistry()
		metadata, _ := registry.GetPlugin("test-plugin@test-marketplace")

		cfg := config.DefaultConfig()
		cfg.DisablePlugin("test-plugin@test-marketplace", config.DisabledPlugin{
			Version:      metadata.Version,
			InstallPath:  metadata.InstallPath,
			GitCommitSha: metadata.GitCommitSha,
			IsLocal:      metadata.IsLocal,
		})

		registry.DisablePlugin("test-plugin@test-marketplace")
		err := claude.SavePlugins(env.ClaudeDir, registry)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.PluginExists("test-plugin@test-marketplace")).To(BeFalse(), "Plugin should not exist after disable")
		Expect(env.PluginCount()).To(Equal(0))
		Expect(cfg.IsPluginDisabled("test-plugin@test-marketplace")).To(BeTrue(), "Plugin should be in disabled list")

		disabledMeta, exists := cfg.EnablePlugin("test-plugin@test-marketplace")
		Expect(exists).To(BeTrue(), "Plugin should exist in disabled list")

		registry = env.LoadPluginRegistry()
		registry.EnablePlugin("test-plugin@test-marketplace", claude.PluginMetadata{
			Version:      disabledMeta.Version,
			InstallPath:  disabledMeta.InstallPath,
			GitCommitSha: disabledMeta.GitCommitSha,
			IsLocal:      disabledMeta.IsLocal,
		})
		err = claude.SavePlugins(env.ClaudeDir, registry)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.PluginExists("test-plugin@test-marketplace")).To(BeTrue(), "Plugin should exist after re-enable")
		Expect(env.PluginCount()).To(Equal(1))
		Expect(cfg.IsPluginDisabled("test-plugin@test-marketplace")).To(BeFalse(), "Plugin should not be in disabled list after enable")
	})
})

var _ = Describe("MCPServerDisableEnable", func() {
	It("disables and re-enables an MCP server", func() {
		cfg := config.DefaultConfig()
		serverRef := "test-plugin@test-marketplace:test-server"

		Expect(cfg.IsMCPServerDisabled(serverRef)).To(BeFalse(), "MCP server should not be disabled initially")

		Expect(cfg.DisableMCPServer(serverRef)).To(BeTrue(), "DisableMCPServer should return true for first disable")
		Expect(cfg.IsMCPServerDisabled(serverRef)).To(BeTrue(), "MCP server should be disabled")

		Expect(cfg.DisableMCPServer(serverRef)).To(BeFalse(), "DisableMCPServer should return false for already disabled server")

		Expect(cfg.EnableMCPServer(serverRef)).To(BeTrue(), "EnableMCPServer should return true")
		Expect(cfg.IsMCPServerDisabled(serverRef)).To(BeFalse(), "MCP server should not be disabled after enable")

		Expect(cfg.EnableMCPServer(serverRef)).To(BeFalse(), "EnableMCPServer should return false for already enabled server")
	})
})

var _ = Describe("MultiplePluginsDisable", func() {
	var env *TestEnv

	BeforeEach(func() {
		env = SetupTestEnv()

		env.CreateMarketplace("test-marketplace", "test/repo")
		env.CreatePlugin("plugin1", "test-marketplace", "1.0.0", nil)
		env.CreatePlugin("plugin2", "test-marketplace", "2.0.0", nil)
		env.CreatePlugin("plugin3", "test-marketplace", "3.0.0", nil)

		plugin1Path := filepath.Join(env.ClaudeDir, "plugins", "marketplaces", "test-marketplace", "plugins", "plugin1")
		plugin2Path := filepath.Join(env.ClaudeDir, "plugins", "marketplaces", "test-marketplace", "plugins", "plugin2")
		plugin3Path := filepath.Join(env.ClaudeDir, "plugins", "marketplaces", "test-marketplace", "plugins", "plugin3")

		env.CreatePluginRegistry(map[string]claude.PluginMetadata{
			"plugin1@test-marketplace": {
				Version:     "1.0.0",
				InstallPath: plugin1Path,
			},
			"plugin2@test-marketplace": {
				Version:     "2.0.0",
				InstallPath: plugin2Path,
			},
			"plugin3@test-marketplace": {
				Version:     "3.0.0",
				InstallPath: plugin3Path,
			},
		})
	})

	AfterEach(func() {
		env.Cleanup()
	})

	It("disables multiple plugins", func() {
		Expect(env.PluginCount()).To(Equal(3))

		registry := env.LoadPluginRegistry()
		cfg := config.DefaultConfig()

		for _, name := range []string{"plugin1@test-marketplace", "plugin3@test-marketplace"} {
			metadata, _ := registry.GetPlugin(name)
			cfg.DisablePlugin(name, config.DisabledPlugin{
				Version:     metadata.Version,
				InstallPath: metadata.InstallPath,
			})
			registry.DisablePlugin(name)
		}

		err := claude.SavePlugins(env.ClaudeDir, registry)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.PluginCount()).To(Equal(1))
		Expect(env.PluginExists("plugin2@test-marketplace")).To(BeTrue(), "plugin2 should still exist")
		Expect(env.PluginExists("plugin1@test-marketplace")).To(BeFalse(), "plugin1 should not exist after disable")
		Expect(env.PluginExists("plugin3@test-marketplace")).To(BeFalse(), "plugin3 should not exist after disable")

		Expect(cfg.IsPluginDisabled("plugin1@test-marketplace")).To(BeTrue(), "plugin1 should be in disabled list")
		Expect(cfg.IsPluginDisabled("plugin3@test-marketplace")).To(BeTrue(), "plugin3 should be in disabled list")
	})
})
