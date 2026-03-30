## Why

There is no way to reconfigure asylum after the initial onboarding — users must manually edit YAML config files. A `asylum config` command with an interactive TUI would let users toggle kits, manage credentials, and change isolation settings using the same familiar prompts from onboarding.

## What Changes

- New `asylum config` CLI subcommand that launches a tabbed TUI
- Three tabs: **Kits**, **Credentials**, **Isolation** — switchable via arrow keys
- Each tab reuses the same interaction model as the onboarding wizard (multiselect for kits/credentials, single-select for isolation)
- Kit activation detects and removes any existing commented-out config (including commented options) before inserting the active config snippet
- Kit deactivation removes the active config entry and optionally inserts a commented version

## Capabilities

### New Capabilities
- `config-command`: CLI subcommand routing, tabbed TUI model, and config editing logic for toggling kits (with comment cleanup), credentials, and isolation settings

### Modified Capabilities

_(none — the existing config editing functions in `config/sync.go` and `config/edit.go` will be extended, but no spec-level requirement changes to existing capabilities)_

## Impact

- **New files**: `internal/tui/tabs.go` (tabbed TUI component), config command handler
- **Modified files**: `cmd/asylum/main.go` (subcommand dispatch), `internal/config/sync.go` (comment removal on activation)
- **Dependencies**: No new dependencies — uses existing bubbletea/lipgloss
