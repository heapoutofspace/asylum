## Why

The kit sync flow (`SyncNewKits`) prompts for new `TierDefault` kits one at a time using plain `fmt.Scanln`, which is inconsistent with the rest of the onboarding UX that uses our bubbletea-based TUI components. When multiple new kits land at once (e.g., after an update), the user must answer separate Y/n prompts for each one — no batch selection, no descriptions, no visual consistency with the wizard.

## What Changes

- **Replace plain-text kit prompts with a single TUI multiselect.** When new `TierDefault` kits are detected, collect them all and present one `tui.MultiSelect` prompt with kit descriptions. All default-on kits are pre-selected; the user can deselect any they don't want.
- **Remove per-kit Y/n prompting.** The sequential `promptActivateKit()` function and its `fmt.Scanln` call are replaced entirely.
- **Preserve non-interactive behavior.** When stdin is not a TTY, default-on kits are still activated automatically (multiselect returns defaults without prompting).

## Capabilities

### New Capabilities

### Modified Capabilities

- `kit-activation-prompt`: The activation prompt changes from per-kit Y/n text input to a single batched TUI multiselect with all new default-on kits pre-selected.

## Impact

- **Code:** `internal/config/kitsync.go` — `SyncNewKits` and `promptActivateKit` are reworked; new dependency on `internal/tui` package.
- **UX:** Users see one multiselect instead of N sequential prompts. Descriptions are visible. Muscle memory for Y/n responses changes.
- **No API/config/dependency changes.** The TUI package and multiselect component already exist.
