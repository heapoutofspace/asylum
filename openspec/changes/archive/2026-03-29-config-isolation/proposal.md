## Why

Users are confused that their Claude config inside the sandbox differs from their host config. The current behavior (copy host config to `~/.asylum/agents/claude/` on first run, then diverge) is one of three valid strategies. Users should choose their preferred isolation level. This is also the first feature requiring an interactive prompt beyond a simple `y/n`, making it the right time to introduce a proper TUI framework for future prompts (kit activation, onboarding improvements, etc.).

## What Changes

- **New `internal/tui` package**: Wraps Bubble Tea to provide reusable `Select` (single-choice) and `MultiSelect` functions that return the user's choice. Clean API: `tui.Select(title, options, defaultIdx) → selectedIdx, error`
- **New `config` field on AgentConfig**: `Config string` with values `shared`, `isolated`, `project` controlling how the agent's config directory is mounted
- **First-run prompt**: When `agents.claude.config` is not set, prompt the user with a single-select TUI to choose their isolation level before the first container start
- **Three isolation modes**:
  - `shared`: Mount host `~/.claude` directly into the container (symlinks work, changes propagate both ways)
  - `isolated`: Mount from `~/.asylum/agents/claude/` (current behavior — shared across projects but separate from host)
  - `project`: Mount from `~/.asylum/projects/<container>/claude-config/` (per-project isolation, no state shared between projects)
- **Config written after prompt**: The chosen value is written back to `~/.asylum/config.yaml` so the user isn't asked again
- **New dependency**: `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss` for TUI rendering

## Capabilities

### New Capabilities
- `tui-prompts`: Reusable TUI prompt framework using Bubble Tea (Select, MultiSelect)
- `config-isolation`: Per-agent config isolation level (shared/isolated/project) with first-run prompt

### Modified Capabilities
- `profile-container-setup`: Agent config volume mount respects the isolation level
- `agent-install`: AgentConfig struct gains a Config field

## Impact

- **internal/tui/** (new package): `select.go` (single-choice), `multiselect.go` (multi-choice), using Bubble Tea
- **internal/config/config.go**: `AgentConfig` gains `Config string` field
- **internal/container/container.go**: Agent config mount logic branches on isolation level
- **cmd/asylum/main.go**: First-run prompt before container start if isolation not configured
- **go.mod**: New dependency on `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`
