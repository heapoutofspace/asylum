## 1. Config Editing Functions

- [x] 1.1 Add `RemoveKitComment` function to `internal/config/sync.go` — finds and removes a commented-out kit block (the `# <name>:` line plus any subsequent deeper-indented comment lines) from the kits section
- [x] 1.2 Add `RemoveKitEntry` function to `internal/config/sync.go` — removes an active kit's YAML block (the kit key line plus all nested lines at deeper indentation)
- [x] 1.3 Add tests for `RemoveKitComment` and `RemoveKitEntry` covering: kit with commented options, kit without comments, kit with nested active config, kit not present

## 2. Tabbed TUI Component

- [x] 2.1 Add `tui.Tabs` bubbletea model in `internal/tui/tabs.go` — tab bar with left/right navigation, each tab contains either a selectModel or multiModel sub-model, Enter confirms all tabs, Escape cancels
- [x] 2.2 Add `Tab` struct and `RunTabs` public function that takes a slice of tab definitions (title, kind, options, defaults) and returns results per tab
- [x] 2.3 Add tests for tab navigation (left/right bounds, no wrap-around) and result collection

## 3. Config Command

- [x] 3.1 Add `runConfig` function in `cmd/asylum/config.go` — loads config, builds tab definitions from kit/credential/isolation state, runs tabbed TUI, applies changes
- [x] 3.2 Wire `"config"` subcommand in `cmd/asylum/main.go` dispatch switch (after `ssh-init`, before project dir resolution)
- [x] 3.3 Add `config` to the usage/help text in `printUsage`

## 4. Integration

- [x] 4.1 End-to-end manual test: run `asylum config`, verify kit toggle activates/deactivates correctly in config file, verify credential and isolation changes persist
- [x] 4.2 Add changelog entry under Unreleased/Added
