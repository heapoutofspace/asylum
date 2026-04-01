## Why

The tabbed TUI (`tui.RunTabs`) does not preserve user selections when switching between tabs. If a user selects kits in the Kits tab, navigates to Credentials, then back, all kit selections are reset to defaults. This makes the config TUI unusable for multi-tab workflows since users must configure everything in a single tab visit.

## What Changes

- Fix `initTab()` in `internal/tui/tabs.go` to restore previously saved state from `results[]` instead of always rebuilding from `DefaultSel`/`DefaultIdx`
- Fix the existing test `TestTabsModelRestoresStateOnSwitch` which currently only verifies the results array was saved but doesn't verify the visual model reflects saved state after switching back

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `tui-prompts`: Tab switching must preserve user selections — `initTab()` should restore from saved results when available

## Impact

- `internal/tui/tabs.go` — `initTab()` method
- `internal/tui/tabs_test.go` — strengthen restore test
- Affects `asylum config` command (primary consumer of `RunTabs`)
