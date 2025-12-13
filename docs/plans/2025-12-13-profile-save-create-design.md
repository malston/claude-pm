# Profile Save and Create Commands Design

## Overview

Refactor profile management commands to clarify the distinction between saving current state and creating profiles from existing ones.

**Key insight:** `save` captures reality, `create` copies profiles.

## Command Structure

### `claudeup profile save [name]`

Saves current Claude Code state to a profile file.

| Scenario | Behavior |
|----------|----------|
| `profile save foo` (new) | Creates `foo.json` from current state |
| `profile save foo` (exists) | Prompts "Overwrite foo?" then saves |
| `profile save foo -y` (exists) | Overwrites without prompting |
| `profile save` (no name) | Saves to active profile, or errors if none active |

### `claudeup profile create <name>`

Creates a new profile by copying an existing one.

| Scenario | Behavior |
|----------|----------|
| `profile create foo` | Prompts "Base profile?" with list of existing profiles |
| `profile create foo --from bar` | Copies `bar` to new profile `foo` |
| `profile create foo -y` | Copies active profile to `foo` (errors if none active) |
| `profile create foo` (exists) | Error: "Profile foo already exists. Use 'profile save' to update it." |

## Flags and Interactive Prompts

### `profile save`

```
Flags:
  -y, --yes    Skip overwrite confirmation
```

No special prompts needed.

### `profile create`

```
Flags:
  -y, --yes           Use active profile as base without prompting
  --from <profile>    Specify base profile explicitly
```

Interactive prompt (when no `--from` and no `-y`):

```
$ claudeup profile create my-python-setup

Which profile should "my-python-setup" be based on?

  1) default         Minimal Claude Code configuration
  2) developer       Full development setup with plugins
  3) untrusted       Sandboxed profile for untrusted code

Enter number or name: _
```

### Error Cases

- `profile create foo -y` with no active profile: `Error: no active profile. Use --from <profile> to specify base.`
- `profile create foo` when `foo` exists: `Error: profile "foo" already exists. Use 'profile save foo' to update it.`
- `profile create foo --from bar` when `bar` doesn't exist: `Error: profile "bar" not found.`

## Implementation Changes

### Files to modify

`internal/commands/profile_cmd.go`:
- Rename `profileCreateCmd` to `profileSaveCmd`
- Update `runProfileCreate` to `runProfileSave`
- Add new `profileCreateCmd` with copy/fork logic
- Add `--from` flag to create command
- Add interactive profile picker helper function

### New helper function

```go
func promptProfileSelection(profilesDir string, newName string) (*profile.Profile, error)
```

Lists available profiles (user + embedded), displays numbered menu, accepts number or name input, returns selected profile.

### Logic flow for `runProfileCreate` (new behavior)

1. Check if target profile name already exists - error if so
2. If `--from` specified - load that profile
3. Else if `-y` and active profile exists - load active profile
4. Else if `-y` and no active profile - error
5. Else - prompt with `promptProfileSelection()`
6. Clone the selected profile with new name
7. Save to profiles directory

### Logic flow for `runProfileSave`

1. If no name given - use active profile name (error if none)
2. If profile exists and no `-y` - prompt overwrite confirmation
3. Call `profile.Snapshot()` (unchanged)
4. Save profile

## Testing Strategy

### Unit tests in `internal/commands/profile_cmd_test.go`

**profile save:**
- Save to new profile name creates file
- Save to existing profile with confirmation overwrites
- Save with no name + active profile saves to active
- Save with no name + no active profile errors

**profile create:**
- Create with `--from` valid profile copies correctly
- Create with `--from` invalid profile errors
- Create with `-y` + active profile copies active
- Create with `-y` + no active profile errors
- Create when target name exists errors
- Create clones all fields (plugins, MCP servers, marketplaces)

### Integration tests in `test/integration/`

- End-to-end: save profile, modify Claude state, create variant, verify both profiles exist with correct contents
- Verify embedded profiles work as `--from` source

## Manual Testing Scripts

### `scripts/test/profile-save-test.sh`

```bash
#!/bin/bash
# Tests profile save functionality

set -e
echo "=== Profile Save Tests ==="

# Test 1: Save to new profile
echo "[1] Testing: save to new profile name"
claudeup profile save test-save-new
claudeup profile show test-save-new
echo "✓ Pass"

# Test 2: Save to existing profile (will prompt - tester confirms)
echo "[2] Testing: save to existing profile (answer 'y' when prompted)"
claudeup profile save test-save-new
echo "✓ Pass"

# Test 3: Save with -y skips prompt
echo "[3] Testing: save with -y flag"
claudeup profile save test-save-new -y
echo "✓ Pass"

# Test 4: Save with no name uses active profile
echo "[4] Testing: save with no name (uses active profile)"
claudeup profile use test-save-new -y
claudeup profile save
echo "✓ Pass"

# Cleanup
rm -f ~/.claudeup/profiles/test-save-new.json
echo "=== All profile save tests passed ==="
```

### `scripts/test/profile-create-test.sh`

```bash
#!/bin/bash
# Tests profile create (copy/fork) functionality

set -e
echo "=== Profile Create Tests ==="

# Test 1: Create with --from
echo "[1] Testing: create with --from flag"
claudeup profile create test-fork-1 --from default
claudeup profile show test-fork-1
echo "✓ Pass"

# Test 2: Create with -y uses active profile
echo "[2] Testing: create with -y uses active profile"
claudeup profile use default -y
claudeup profile create test-fork-2 -y
claudeup profile show test-fork-2
echo "✓ Pass"

# Test 3: Create when name exists should error
echo "[3] Testing: create existing name errors"
if claudeup profile create test-fork-1 --from default 2>&1 | grep -q "already exists"; then
  echo "✓ Pass"
else
  echo "✗ Fail - should have errored"
  exit 1
fi

# Test 4: Interactive picker (tester selects manually)
echo "[4] Testing: interactive picker (select a profile when prompted)"
claudeup profile create test-fork-3
claudeup profile show test-fork-3
echo "✓ Pass"

# Cleanup
rm -f ~/.claudeup/profiles/test-fork-*.json
echo "=== All profile create tests passed ==="
```

### `scripts/test/profile-error-test.sh`

```bash
#!/bin/bash
# Tests error conditions

echo "=== Profile Error Tests ==="

# Test 1: create -y with no active profile
echo "[1] Testing: create -y with no active profile errors"
# Clear active profile first
claudeup config set preferences.activeProfile ""
if claudeup profile create test-err -y 2>&1 | grep -q "no active profile"; then
  echo "✓ Pass"
else
  echo "✗ Fail - should have errored about no active profile"
  exit 1
fi

# Test 2: create --from nonexistent profile
echo "[2] Testing: create --from nonexistent profile errors"
if claudeup profile create test-err --from nonexistent-xyz 2>&1 | grep -q "not found"; then
  echo "✓ Pass"
else
  echo "✗ Fail - should have errored about profile not found"
  exit 1
fi

echo "=== All error tests passed ==="
```

## Migration Notes

This is a breaking change to the `profile create` command semantics:

- **Old behavior:** `profile create <name>` snapshots current Claude Code state
- **New behavior:** `profile create <name>` copies an existing profile

Users who want the old behavior should use `profile save <name>` instead.

Consider adding a deprecation warning in the next minor release before removing the old behavior, or simply document the change in release notes since this is a young project.
