## Context

Users can only configure asylum during the initial onboarding wizard (isolation + credentials) or by manually editing `~/.asylum/config.yaml`. There is no interactive way to change settings after the first run. The existing TUI wizard (`internal/tui/wizard.go`) already supports multi-step flows with select/multiselect, and the config editing layer (`internal/config/sync.go`, `internal/config/isolation.go`) already handles text-based YAML manipulation.

## Goals / Non-Goals

**Goals:**
- Add `asylum config` subcommand that opens a tabbed TUI for managing kits, credentials, and isolation
- Reuse existing TUI building blocks (select, multiselect) within a new tabbed container model
- Handle kit activation cleanly: when activating a kit that was previously commented out, remove the commented version before inserting the active config snippet
- Handle kit deactivation: remove the active config entry and insert the commented version

**Non-Goals:**
- Per-project config editing (this command edits `~/.asylum/config.yaml` only)
- Adding new config options or kit features
- Changing the onboarding wizard flow

## Decisions

### 1. Tabbed TUI model as a new bubbletea component

The config command needs a tabbed interface where left/right arrows switch tabs and each tab contains its own select/multiselect content. This is a new `tui.Tabs` component that wraps the existing `selectModel` and `multiModel` sub-models.

**Tab layout:** Tab bar at the top (`[ Kits | Credentials | Isolation ]`), content below. Active tab is highlighted. Left/right arrows switch tabs. Up/down and space/enter work within the active tab's content. The component is standalone — it is not a wizard step but a top-level bubbletea program.

**Alternative considered:** Reuse the wizard with three steps. Rejected because the wizard is sequential (you move through steps linearly), while tabs allow random access to any section.

### 2. Config command in its own file

Add `cmd/asylum/config.go` with a `runConfig()` function, dispatched from `main.go` via the existing subcommand switch. This keeps the config logic separate from the main flow.

### 3. Kit comment removal before activation

When a kit is toggled from inactive to active, the existing commented-out config block must be found and removed before inserting the active `ConfigSnippet`. This requires a new function `RemoveKitComment` in `internal/config/sync.go` that:

1. Finds the `kits:` section
2. Scans for a commented line matching `# <kitName>:` at the kit entry indent level (2 spaces)
3. Removes that line plus any subsequent commented lines at deeper indentation (the kit's commented options)
4. Removes any trailing blank line left behind

This runs before `SyncKitToConfig` so the insertion logic doesn't need to change.

### 4. Kit deactivation: remove active entry, insert comment

When a kit is toggled from active to inactive, a new function `RemoveKitEntry` removes the active YAML block (the kit key and all its nested lines), then `SyncKitCommentToConfig` inserts the commented version.

### 5. State derivation from config file

The current state of each tab is derived by reading the config:
- **Kits tab**: All registered kits shown. Checked = kit key exists uncommented in `kits:` section. Uses `kit.All()` for the full list and parses the config to determine which are active.
- **Credentials tab**: All credential-capable kits. Checked = `credentials: auto` under the kit's config. Same as onboarding.
- **Isolation tab**: Current agent isolation level. Single-select from shared/isolated/project. Read from `cfg.AgentIsolation(agentName)`.

### 6. Apply changes on confirm, not on toggle

Changes are collected as the user toggles items across all tabs, then applied all at once when the user presses Enter on any tab (confirming the whole configuration). This avoids partial writes if the user cancels.

## Risks / Trade-offs

- **Comment format fragility**: Comment removal relies on the format produced by `SyncKitCommentToConfig`. If a user manually edits comments into a different format, removal may not find them. Mitigation: match on the kit name pattern (`# <name>:`) which is unlikely to be altered.
- **No undo**: Once confirmed, changes are written to the config file. Mitigation: the file is version-controllable and changes are idempotent (re-running config restores the desired state).
