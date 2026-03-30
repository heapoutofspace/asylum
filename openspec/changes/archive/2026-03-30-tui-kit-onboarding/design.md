## Context

`SyncNewKits` in `internal/config/kitsync.go` detects kits not yet tracked in `state.json` and prompts for each `TierDefault` kit individually via `fmt.Scanln`. The rest of the onboarding UX (config isolation, credentials) uses the bubbletea-based `tui.MultiSelect` and `tui.Wizard`. The TUI multiselect component already supports labels, descriptions, and pre-selected defaults — everything needed to replace the plain-text prompt.

## Goals / Non-Goals

**Goals:**
- Replace per-kit Y/n prompts with a single `tui.MultiSelect` call that shows all new `TierDefault` kits at once, pre-selected by default
- Show kit descriptions in the multiselect so users can make informed choices
- Preserve existing behavior for `TierAlwaysOn` (log only) and `TierOptIn` (comment only) kits
- Preserve non-interactive fallback (all defaults activated without prompting)

**Non-Goals:**
- Integrating this into the wizard flow (it runs at a different lifecycle point — startup sync vs first-run)
- Changing kit tiers or adding new tiers
- Modifying the first-run wizard's kit selection

## Decisions

### Use `tui.MultiSelect` directly, not the wizard

The kit sync runs on every startup when new kits are detected, not just during first-run. The wizard is designed for multi-step onboarding flows. A single multiselect call is simpler and appropriate — there's only one question to ask.

**Alternative considered:** Wrapping in a `tui.Wizard` with one step. Rejected because it adds visual overhead (tab bar, progress) for a single prompt, and the wizard pattern is reserved for the first-run flow.

### Collect all new `TierDefault` kits before prompting

Currently the loop processes kits one at a time, prompting and writing config for each. The new approach collects all new `TierDefault` kits into a slice first, presents them in one multiselect, then processes the results. This requires restructuring the loop into two passes: classify, then act.

**Alternative considered:** Keeping the per-kit loop and replacing `promptActivateKit` with `tui.Confirm` per kit. Rejected because it doesn't solve the batch UX problem — users still answer N separate prompts.

### Pass `tui.MultiSelect` as a function parameter

`SyncNewKits` lives in `internal/config` which should not import `internal/tui` directly (config is a lower-level package). Instead, accept an optional prompt function parameter that `main.go` provides. When nil (or in non-interactive mode), the default behavior activates all default-on kits.

Signature: `func SyncNewKits(asylumDir string, interactive bool, promptFn func([]kit.Kit) []string) (bool, error)`

The `promptFn` receives the list of new default kits and returns the names the user selected. `main.go` wires this to a closure that calls `tui.MultiSelect`.

**Alternative considered:** Having `kitsync.go` import `internal/tui` directly. Rejected to keep the config package free of TUI dependencies.

### Pre-select all default-on kits

All `TierDefault` kits are selected by default in the multiselect. Users deselect what they don't want. This matches the current behavior where the default answer is "Y" — the change is visual, not behavioral.

## Risks / Trade-offs

- **[Visual noise for single kit]** If only one new kit is detected, a multiselect for a single item is slightly heavier than a Y/n prompt. → Acceptable trade-off for consistency. A single-item multiselect still works fine and the user can just press Enter.
- **[Import boundary]** Passing a prompt function adds a parameter to `SyncNewKits`. → Clean separation is worth the slightly longer signature. The function is only called from one place (`main.go`).
