# Profiles and Setup Command Design

## Overview

Add profile-based configuration management to claude-pm, enabling:
- First-time setup with `claude-pm setup`
- Reset to known-good-state with `claude-pm profile use <name>`
- Portable, shareable configuration profiles

## Profile Schema

Profiles are stored in `~/.claude-pm/profiles/<name>.json`:

```json
{
  "name": "default",
  "description": "Base setup for general development",
  "mcpServers": [
    {
      "name": "context7",
      "command": "npx",
      "args": ["-y", "@context7/mcp"],
      "scope": "user",
      "secrets": {}
    }
  ],
  "marketplaces": [
    {"source": "github", "repo": "anthropics/claude-code-plugins"},
    {"source": "github", "repo": "anthropics/claude-code-templates"}
  ],
  "plugins": [
    "superpowers@superpowers-marketplace",
    "episodic-memory@superpowers-marketplace"
  ],
  "detect": {
    "files": ["package.json"],
    "contains": {"package.json": ["react", "vue"]}
  }
}
```

### Secret References

Secrets use an abstract reference format supporting multiple backends:

```json
"secrets": {
  "GITHUB_TOKEN": {
    "description": "GitHub personal access token with repo scope",
    "sources": [
      {"type": "env", "key": "GITHUB_TOKEN"},
      {"type": "1password", "ref": "op://Private/GitHub PAT/credential"},
      {"type": "keychain", "service": "github-pat", "account": "claude"}
    ]
  }
}
```

Resolution tries each source in order. First success wins.

## Commands

### `claude-pm setup [--profile <name>] [--yes]`

First-time setup or reset:
1. Install Claude CLI if missing
2. Create profiles directory if missing
3. Select profile (flag, interactive picker, or default)
4. Apply profile
5. Run `claude-pm doctor` to verify

TTY behavior: Shows summary and offers profile picker unless `--yes`.

### `claude-pm profile list`

Lists available profiles with descriptions. Marks active profile.

### `claude-pm profile use <name> [--yes]`

Applies a profile using **replace strategy**:
- Remove plugins/MCP servers not in profile
- Install everything in profile
- Confirms changes unless `--yes`

### `claude-pm profile create <name>`

Snapshots current state as a new profile.

### `claude-pm profile show <name>`

Displays profile contents in readable format.

### `claude-pm profile suggest`

Analyzes current directory against profile `detect` rules, suggests best match.

## Secret Resolution

### Supported Backends (v1)

| Type | Implementation | Platform |
|------|----------------|----------|
| `env` | `os.Getenv(key)` | All |
| `1password` | `op read <ref>` | All (if `op` installed) |
| `keychain` | `security find-generic-password` | macOS |

### Resolution Algorithm

1. Try each source in order
2. First successful read wins
3. If all fail and TTY: prompt user, offer to save
4. If all fail and non-interactive: error with message

### User Preferences

```json
// ~/.claude-pm/config.json
{
  "preferences": {
    "secretBackend": "1password"
  }
}
```

## File Structure

```
~/.claude-pm/
├── config.json              # global prefs (existing)
└── profiles/
    ├── default.json
    ├── frontend.json
    └── go-backend.json
```

## Code Organization

New packages:

```
internal/
├── profile/
│   ├── profile.go           # Profile struct, Load/Save
│   ├── apply.go             # Apply profile (replace strategy)
│   ├── detect.go            # Project detection logic
│   └── snapshot.go          # Create profile from current state
├── secrets/
│   ├── resolver.go          # Resolution chain interface
│   ├── env.go               # Env var backend
│   ├── onepassword.go       # 1Password backend
│   └── keychain.go          # macOS Keychain backend
└── commands/
    ├── setup.go             # claude-pm setup
    └── profile.go           # claude-pm profile *
```

### Key Interfaces

```go
type SecretResolver interface {
    Name() string
    Available() bool
    Resolve(ref string) (string, error)
    Store(key, value, description string) error
}

type Profile struct {
    Name         string        `json:"name"`
    Description  string        `json:"description"`
    MCPServers   []MCPServer   `json:"mcpServers"`
    Marketplaces []Marketplace `json:"marketplaces"`
    Plugins      []string      `json:"plugins"`
    Detect       DetectRules   `json:"detect,omitempty"`
}
```

## Apply Logic

When `profile use` runs:

1. **Load current state** from Claude's config files
2. **Compute diff** (what to remove, what to install)
3. **Resolve secrets** for all MCP servers (fail fast)
4. **Confirm** if TTY and not `--yes`
5. **Execute changes** via `claude` CLI commands
6. **Record active profile** in config.json

## Bootstrap

`claude-pm setup` handles the chicken-and-egg problem:

1. Check for Claude CLI → install if missing
2. Check for profiles dir → create with bundled default
3. **Detect existing installation** (see below)
4. Apply selected profile

Default profile is embedded in binary using Go's `embed` package.

### Existing Installation Protection

If `~/.claude` exists with plugins, MCP servers, or marketplaces configured, setup offers to preserve the current state as a profile before proceeding:

```
$ claude-pm setup

Existing Claude Code installation detected:
  → 3 MCP servers, 2 marketplaces, 8 plugins

Options:
  [s] Save current setup as a profile, then continue
  [c] Continue anyway (will replace current setup)
  [a] Abort

Choice: s
Profile name [current]: my-old-setup
✓ Saved as 'my-old-setup'

Proceeding with setup using profile: default
...
```

This allows users to switch back with `claude-pm profile use my-old-setup`.

With `--yes` flag: skips the prompt and continues without saving (for scripted installs where the user knows what they're doing).

## One-liner Install

```bash
curl -fsSL https://raw.githubusercontent.com/malston/claude-pm/main/install.sh | bash && claude-pm setup
```

The install.sh downloads the platform-appropriate binary. All logic lives in Go.
