# claude-pm

A comprehensive CLI tool for managing Claude Code plugins, marketplaces, and MCP servers.

## Overview

`claude-pm` provides visibility into and control over your Claude Code installation, including:

- Installed plugins and their state
- Marketplace repositories
- MCP server configuration
- Plugin updates and maintenance

This tool is designed to replace and improve upon the existing bash scripts in `~/.claude/scripts/` with a more robust, cross-platform solution.

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

# List installed marketplaces
claude-pm marketplaces

# List MCP servers by plugin
claude-pm mcp list
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

## Configuration

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

### Phase 2: Enable/Disable Control (Next)

- `claude-pm enable <plugin>` - Enable a plugin
- `claude-pm disable <plugin>` - Disable a plugin
- `claude-pm mcp disable <plugin>:<server>` - Disable specific MCP server
- `claude-pm mcp enable <plugin>:<server>` - Re-enable MCP server
- Global config file for tracking disabled MCP servers

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
