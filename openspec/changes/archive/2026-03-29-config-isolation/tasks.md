## 1. TUI Package

- [x] 1.1 Add `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss` dependencies to go.mod
- [x] 1.2 Create `internal/tui/option.go` with the `Option` struct (Label, Description)
- [x] 1.3 Create `internal/tui/select.go` with `Select(title string, options []Option, defaultIdx int) (int, error)` — Bubble Tea model with arrow key navigation, enter to confirm, esc to cancel; non-interactive fallback returns defaultIdx
- [x] 1.4 Create `internal/tui/multiselect.go` with `MultiSelect(title string, options []Option, defaultSelected []int) ([]int, error)` — Bubble Tea model with space to toggle, enter to confirm; non-interactive fallback returns defaultSelected
- [x] 1.5 Write tests for Select: returns default on non-interactive, error on cancel
- [x] 1.6 Write tests for MultiSelect: returns defaults on non-interactive

## 2. AgentConfig Isolation Field

- [x] 2.1 Add `Config string` field to `AgentConfig` struct with YAML tag `config`
- [x] 2.2 Add `AgentIsolation(agentName string) string` helper method on Config that returns the isolation level for the given agent (empty string if not set)
- [x] 2.3 Write test for AgentIsolation: returns value when set, empty when not set

## 3. Container Mount Logic

- [x] 3.1 Update `container.go` agent config mount: branch on isolation level — `shared` mounts native dir, `project` mounts per-project dir, `isolated`/default mounts asylum agents dir
- [x] 3.2 For `project` mode: create the per-project config dir at `~/.asylum/projects/<container>/<agent>-config/` if it doesn't exist
- [x] 3.3 For `shared` mode: skip `EnsureAgentConfig` (host dir used directly)
- [x] 3.4 For `project` mode: seed per-project dir from host config on first run (same as isolated but per-project)
- [x] 3.5 Write tests for mount path resolution: shared returns native dir, project returns project dir, isolated returns asylum dir

## 4. First-Run Prompt

- [x] 4.1 In `main.go`, after config loading and before container start: check if Claude agent's isolation is unconfigured
- [x] 4.2 If unconfigured and interactive: show TUI Select with three options (shared, isolated [default], project)
- [x] 4.3 If unconfigured and non-interactive: default to `isolated` with a log message
- [x] 4.4 After selection: write the value to `~/.asylum/config.yaml` using yaml.Node manipulation (preserve comments)
- [x] 4.5 Write a `config.SetAgentIsolation(path, agentName, level string) error` function for the yaml.Node write

## 5. Wiring

- [x] 5.1 Pass isolation level through RunOpts to container mount logic
- [x] 5.2 Update `EnsureAgentConfig` call to be skipped for shared mode
- [x] 5.3 Add CHANGELOG entry under Unreleased
- [x] 5.4 Verify all existing tests pass
