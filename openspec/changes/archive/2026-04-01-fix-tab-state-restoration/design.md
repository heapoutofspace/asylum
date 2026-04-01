## Context

The `tabsModel` in `internal/tui/tabs.go` manages tab switching via two methods:
- `saveTab()` — captures the current tab's selection state into `results[m.active]`
- `initTab(idx)` — initializes the active tab's visual model (`selModel` or `multiModel`)

`saveTab()` is called before every tab switch (lines 122, 130) and correctly persists state. However, `initTab()` always rebuilds the model from `DefaultSel`/`DefaultIdx` on the `Tab` struct, ignoring any previously saved state in `results[]`. This means user selections are visually lost on every tab switch, even though they exist in the results array.

## Goals / Non-Goals

**Goals:**
- `initTab()` restores from saved results when available, falling back to defaults for first visit
- Existing tests are strengthened to verify the visual model (not just the results array)

**Non-Goals:**
- Changing the save/restore architecture (the `results[]` approach is sound)
- Adding undo/reset-to-defaults within tabs

## Decisions

**Check results before defaults in `initTab()`**: When `results[idx]` has been populated (non-nil `MultiIdx` or a previously saved `SelectIdx`), use that to rebuild the model. Otherwise fall back to `DefaultSel`/`DefaultIdx`. This is the minimal fix — one conditional branch added to `initTab()`.

**Use a `visited` flag or check `MultiIdx != nil`**: For multiselect tabs, a nil `MultiIdx` means "never saved" while an empty `[]int{}` means "user deselected everything." These must be distinguishable. The simplest approach: `saveTab()` always sets a non-nil slice (empty `[]int{}` if nothing selected). Then `initTab()` checks `results[idx].MultiIdx != nil` to detect prior visits. For select tabs, we can track this with a boolean on `TabResult` (e.g., `Visited bool`) since `SelectIdx: 0` is ambiguous with the zero value.

## Risks / Trade-offs

- **Zero-value ambiguity for StepSelect**: `SelectIdx: 0` is both "never visited" and "user chose first option." Adding a `Visited` field to `TabResult` resolves this cleanly without breaking the public API (it's an additive field). Alternatively, we could use a pointer (`*int`), but that's less idiomatic for this codebase.
