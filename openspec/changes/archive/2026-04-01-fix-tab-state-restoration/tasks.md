## 1. Fix state restoration

- [x] 1.1 Add `Visited bool` field to `TabResult` struct
- [x] 1.2 Update `saveTab()` to set `Visited = true` and ensure `MultiIdx` is always non-nil (empty slice, not nil)
- [x] 1.3 Update `initTab()` to check `results[idx].Visited` and restore from saved state instead of defaults

## 2. Fix tests

- [x] 2.1 Update `TestTabsModelRestoresStateOnSwitch` to verify the visual model (`multiModel.selected`) reflects saved state after switching back, not just the results array
- [x] 2.2 Add test for empty-selection preservation (deselect all, switch away, switch back — should remain empty, not revert to defaults)
- [x] 2.3 Add test for select tab cursor preservation across tab switches
