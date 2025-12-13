# claudeup

CLI tool for managing Claude Code configurations, profiles, and sandboxed environments.

## Project Structure

- `cmd/claudeup/` - Main entry point
- `internal/commands/` - Cobra command implementations
- `internal/profile/` - Profile management (save, load, apply, snapshot)
- `internal/claude/` - Claude Code configuration file handling
- `internal/sandbox/` - Docker-based sandboxed execution
- `internal/secrets/` - Secret resolution (env, 1Password, keychain)
- `test/integration/` - Integration tests

## Plans and Documentation

**Design documents and implementation plans go in a separate repository:**

```
https://github.com/claudeup/claudeup-superpowers.git
```

Clone locally as `claudeup-feature-plans` or similar. When brainstorming features or creating implementation plans, save them there - not in this repository.

## Development

```bash
# Run tests
go test ./...

# Build
go build -o bin/claudeup ./cmd/claudeup

# Run integration tests
go test ./test/integration/... -v
```

## Worktrees

Feature development uses git worktrees in `.worktrees/` directory (already in .gitignore).

## Embedded Profiles

Built-in profiles are embedded from `internal/profile/profiles/*.json` using Go's embed directive.
