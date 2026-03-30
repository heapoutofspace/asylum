## 1. Refactor SyncNewKits signature

- [x] 1.1 Add `promptFn func([]kit.Kit) []string` parameter to `SyncNewKits` in `internal/config/kitsync.go`
- [x] 1.2 Update the call site in `cmd/asylum/main.go` to pass `nil` for now (preserving current behavior when `promptFn` is nil)

## 2. Restructure kit processing loop

- [x] 2.1 Split the loop in `SyncNewKits` into two passes: first collect new kits by tier (always-on, default, opt-in), then process each group
- [x] 2.2 Process `TierAlwaysOn` and `TierOptIn` kits as before (log messages, config comments)
- [x] 2.3 For `TierDefault` kits: when interactive and `promptFn` is non-nil, call `promptFn` with the collected default kits and use the returned names to determine which to activate vs comment out
- [x] 2.4 When `promptFn` is nil or non-interactive, activate all `TierDefault` kits automatically (existing fallback)

## 3. Wire up TUI multiselect in main.go

- [x] 3.1 Create a closure in `main.go` that builds `tui.Option` slice from `[]kit.Kit` (label=name, description=kit.Description), pre-selects all, calls `tui.MultiSelect`, and returns selected kit names
- [x] 3.2 Handle `tui.ErrCancelled` by returning an empty slice (all kits declined)
- [x] 3.3 Pass the closure as `promptFn` to `SyncNewKits`

## 4. Remove old prompt function

- [x] 4.1 Delete `promptActivateKit` from `kitsync.go`
- [x] 4.2 Remove the `fmt` import if no longer needed

## 5. Test and verify

- [x] 5.1 Run `go build ./...` and `go vet ./...` to verify compilation
- [x] 5.2 Run `go test ./internal/config/...` to verify existing tests pass
- [x] 5.3 Add changelog entry under Unreleased/Changed
