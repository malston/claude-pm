# Hobson Profile Design

A community profile for claudeup that provides full access to the wshobson/agents plugin marketplace with an interactive category-based setup wizard.

## Overview

The hobson profile enables users to install the entire [wshobson/agents](https://github.com/wshobson/agents) marketplace (65+ plugins) and select which categories to enable through an interactive wizard on first run.

## Goals

- Provide full marketplace access so users can pick and choose plugins
- Guide users through category selection rather than overwhelming with 65+ individual choices
- Run setup wizard automatically on first install, skip on re-apply
- Support power users who want to customize individual plugins within categories

## Profile Structure

The profile is stored at `~/.claudeup/profiles/hobson.json` (or embedded as built-in):

```json
{
  "name": "hobson",
  "description": "Full access to wshobson/agents plugin marketplace with interactive category selection",
  "marketplaces": [
    {"source": "github", "repo": "wshobson/agents"}
  ],
  "plugins": [],
  "postApply": {
    "script": "hobson-setup.sh",
    "condition": "first-run"
  }
}
```

Key decisions:
- **Empty plugins array** - No plugins enabled by default; the post-apply hook handles selection
- **Marketplace only** - Adds wshobson/agents so all 67 plugins become available
- **New postApply field** - Triggers interactive wizard after profile is applied
- **Condition: first-run** - Only runs if no plugins from this marketplace are enabled

## Post-Apply Hook Mechanism

claudeup requires a new feature to support profile hooks.

### Profile Struct Addition

```go
type Profile struct {
    // ... existing fields ...
    PostApply *PostApplyHook `json:"postApply,omitempty"`
}

type PostApplyHook struct {
    Script    string `json:"script,omitempty"`    // Script path relative to profile
    Command   string `json:"command,omitempty"`   // Or direct command
    Condition string `json:"condition,omitempty"` // "always" | "first-run"
}
```

### Apply Flow Change

After existing apply logic completes:

1. If `Condition == "first-run"`, check if any plugins from the profile's marketplaces are already enabled. If yes, skip.
2. Resolve script path (relative to profile file for file-based profiles, embedded for built-ins).
3. Execute script interactively (stdin/stdout connected to terminal).

### Script Location

For built-in profiles, store alongside the profile in `internal/profiles/builtin/hobson/` with both `profile.json` and `hobson-setup.sh`.

## Interactive Wizard

The `hobson-setup.sh` script provides category-based selection:

```
┌─────────────────────────────────────────────────────┐
│  Welcome to the Hobson Profile Setup                │
│                                                     │
│  Select development categories (space to toggle):  │
│                                                     │
│  [x] Core Development (workflows, debugging, docs) │
│  [ ] Quality & Testing                             │
│  [ ] AI & Machine Learning                         │
│  [x] Infrastructure & DevOps                       │
│  [ ] Security & Compliance                         │
│  [ ] Data & Databases                              │
│  [x] Languages                                     │
│  [ ] Business & Specialty                          │
│                                                     │
│  [Enter] Continue   [c] Customize   [q] Quit       │
└─────────────────────────────────────────────────────┘
```

### Implementation

- Written in bash using `gum` (charmbracelet/gum) for interactive UI
- Falls back to simple numbered prompts if `gum` isn't installed
- Each theme maps to predefined plugins
- Pressing `c` shows individual plugins within selected themes
- Runs `claudeup enable <plugin>@wshobson-agents` for each selection

### Output

```
Enabling plugins...
  ✓ git-pr-workflows@wshobson-agents
  ✓ debugging-toolkit@wshobson-agents
  ✓ kubernetes-operations@wshobson-agents
  ...
Setup complete! Run 'claudeup status' to see your configuration.
```

## Category-to-Plugin Mapping

| Theme | Plugins |
|-------|---------|
| **Core Development** | code-documentation, debugging-toolkit, git-pr-workflows, backend-development, frontend-mobile-development, full-stack-orchestration, code-refactoring, dependency-management, error-debugging, team-collaboration, documentation-generation, c4-architecture, multi-platform-apps, developer-essentials |
| **Quality & Testing** | unit-testing, tdd-workflows, code-review-ai, comprehensive-review, performance-testing-review, framework-migration, codebase-cleanup |
| **AI & Machine Learning** | llm-application-dev, agent-orchestration, context-management, machine-learning-ops |
| **Infrastructure & DevOps** | deployment-strategies, deployment-validation, kubernetes-operations, cloud-infrastructure, cicd-automation, incident-response, error-diagnostics, distributed-debugging, observability-monitoring |
| **Security & Compliance** | security-scanning, security-compliance, backend-api-security, frontend-mobile-security |
| **Data & Databases** | data-engineering, data-validation-suite, database-design, database-migrations, application-performance, database-cloud-optimization |
| **Languages** | python-development, javascript-typescript, systems-programming, jvm-languages, web-scripting, functional-programming, julia-development, arm-cortex-microcontrollers, shell-scripting |
| **Business & Specialty** | api-scaffolding, api-testing-observability, seo-content-creation, seo-technical-optimization, seo-analysis-monitoring, business-analytics, hr-legal-compliance, customer-sales-automation, content-marketing, blockchain-web3, quantitative-trading, payment-processing, game-development, accessibility-compliance |

## Implementation Tasks

### 1. claudeup Changes (Prerequisite)

- Add `PostApplyHook` struct to `internal/profile/profile.go`
- Add hook execution logic to `internal/profile/apply.go`
- Add first-run detection (check for existing marketplace plugins)
- Add `--no-interactive` flag to skip hooks for CI/scripting
- Tests for new hook behavior

### 2. Hobson Profile Assets

- `internal/profiles/builtin/hobson/profile.json` - Profile definition
- `internal/profiles/builtin/hobson/hobson-setup.sh` - Interactive wizard
- Embed using Go's `embed` package
- Documentation in `docs/profiles.md`

### Dependencies

- `gum` CLI (optional, degrades gracefully)

## User Experience

```bash
# First time
$ claudeup setup --profile hobson
Installing Hobson profile...
  ✓ Added marketplace: wshobson/agents

[Interactive category selection wizard runs]

# Re-apply later
$ claudeup profile use hobson
Applying profile: hobson
  (skipping setup wizard - already configured)
  ✓ Marketplace wshobson/agents already present
```
