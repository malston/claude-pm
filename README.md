# claude-pm

A comprehensive CLI tool for managing Claude Code plugins, marketplaces, and MCP servers.

## Overview

`claude-pm` provides visibility into and control over your Claude Code installation, including:

- Installed plugins and their state (enable/disable individual plugins)
- Marketplace repositories
- MCP server configuration (granular control over individual MCP servers)
- Plugin updates and maintenance

This is a community-built tool that provides a unified interface for managing Claude Code plugins and MCP servers.

## Installation

### From Source

```bash
go install github.com/malston/claude-pm@latest
```

### From GitHub Releases

Download the appropriate binary for your platform from the [releases page](https://github.com/malston/claude-pm/releases).

### Building Locally

```bash
git clone https://github.com/malston/claude-pm.git
cd claude-pm
go build -o bin/claude-pm ./cmd/claude-pm
```

## Usage

### Quick Start

```bash
# Show overview of your Claude installation
claude-pm status

# List all installed plugins
claude-pm plugins list

# Enable/disable plugins
claude-pm disable compound-engineering@every-marketplace
claude-pm enable compound-engineering@every-marketplace

# List installed marketplaces
claude-pm marketplaces

# List MCP servers by plugin
claude-pm mcp list

# Enable/disable specific MCP servers
claude-pm mcp disable superpowers-chrome@superpowers-marketplace:chrome
claude-pm mcp enable superpowers-chrome@superpowers-marketplace:chrome
```

### Commands

#### Status

Show an overview of your Claude Code installation:

```bash
claude-pm status
```

Example output:
```
╔════════════════════════════════════════╗
║           claude-pm Status             ║
╚════════════════════════════════════════╝

Marketplaces (7)
  ✓ superpowers-marketplace
  ✓ claude-code-plugins
  ✓ every-marketplace
  ...

Plugins (27 total)
  ✓ 8 enabled

MCP Servers
  → Run 'claude-pm mcp list' for details

Issues Detected
  ⚠ 19 plugins have stale paths
  → Run 'claude-pm doctor' for details
```

#### Plugins

List all installed plugins with detailed information:

```bash
claude-pm plugins list
```

Shows version, status, installation path, and type (local/cached) for each plugin.

#### Marketplaces

List all installed marketplace repositories:

```bash
claude-pm marketplaces
```

Shows source, repository, location, and last update time for each marketplace.

#### MCP Servers

List MCP servers grouped by the plugin that provides them:

```bash
claude-pm mcp list
```

Shows command, arguments, and environment variables for each MCP server.

#### Enable/Disable Plugins

Enable or disable individual plugins:

```bash
# Disable a plugin
claude-pm disable hookify@claude-code-plugins

# Re-enable a plugin
claude-pm enable hookify@claude-code-plugins
```

When you disable a plugin:
- It's removed from `installed_plugins.json`
- Its metadata is saved to `~/.claude-pm/config.json`
- All commands, agents, skills, and MCP servers become unavailable

When you re-enable it:
- The metadata is restored from config
- The plugin becomes available again immediately
- No need to reinstall

#### Enable/Disable MCP Servers

Control individual MCP servers without disabling the entire plugin:

```bash
# Disable a specific MCP server
claude-pm mcp disable compound-engineering@every-marketplace:playwright

# Re-enable it
claude-pm mcp enable compound-engineering@every-marketplace:playwright
```

This is useful for:
- Reducing MCP context usage
- Disabling heavy servers (like Playwright) when not needed
- Keeping plugin features (commands, skills) while removing MCP tools

**Note:** MCP server changes require restarting Claude Code to take effect.

## Configuration

### Global Config File

`claude-pm` stores its configuration in `~/.claude-pm/config.json`:

```json
{
  "disabledPlugins": {
    "hookify@claude-code-plugins": {
      "version": "1.0.0",
      "installPath": "/Users/you/.claude/plugins/...",
      ...
    }
  },
  "disabledMcpServers": [
    "compound-engineering@every-marketplace:playwright"
  ],
  "claudeDir": "/Users/you/.claude",
  "preferences": {
    "autoUpdate": false,
    "verboseOutput": false
  }
}
```

The config file is created automatically on first use.

### Global Flags

- `--claude-dir` - Override the Claude installation directory (default: `~/.claude`)

Example:
```bash
claude-pm --claude-dir /custom/path status
```

## Roadmap

### Phase 1: Core Status & Discovery ✅ (Complete)

- ✅ `claude-pm status` - Overview of installation
- ✅ `claude-pm plugins list` - List all plugins
- ✅ `claude-pm marketplaces` - List marketplaces
- ✅ `claude-pm mcp list` - List MCP servers

### Phase 2: Enable/Disable Control ✅ (Complete)

- ✅ `claude-pm enable <plugin>` - Enable a plugin
- ✅ `claude-pm disable <plugin>` - Disable a plugin
- ✅ `claude-pm mcp disable <plugin>:<server>` - Disable specific MCP server
- ✅ `claude-pm mcp enable <plugin>:<server>` - Re-enable MCP server
- ✅ Global config file for tracking disabled plugins and MCP servers

### Phase 3: Update & Maintenance

- `claude-pm update` - Update marketplaces and plugins
- `claude-pm cleanup` - Clean stale plugins
- `claude-pm fix-paths` - Fix plugin path issues
- `claude-pm doctor` - Diagnose common issues

### Phase 4: Project-Level Config & Polish

- `claude-pm init` - Initialize project config
- `.claude-pm.json` - Project-specific configuration
- Config merging (global + project)
- `claude-pm maintain` - Interactive TUI mode
- `claude-pm config show` - Show effective merged config

## Architecture

The project is organized as follows:

```
claude-pm/
├── cmd/claude-pm/          # Main entry point
├── internal/
│   ├── claude/             # Claude data structures (plugins, marketplaces)
│   ├── mcp/                # MCP server discovery and management
│   ├── commands/           # CLI commands
│   ├── config/             # Configuration management (future)
│   └── ui/                 # UI components (future)
```

## Development

### Prerequisites

- Go 1.21 or later

### Building

```bash
go build -o bin/claude-pm ./cmd/claude-pm
```

### Testing

```bash
# Test against your actual Claude installation
./bin/claude-pm status
./bin/claude-pm plugins list
./bin/claude-pm mcp list
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

This tool was created to improve upon the existing bash scripts in the Claude Code ecosystem. It aims to provide a more robust, cross-platform solution for managing Claude Code installations.
