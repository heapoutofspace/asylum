## Context

Agent config is currently always mounted from `~/.asylum/agents/<agent>/` — a copy seeded from the host on first run. Users who expect their host Claude config to work identically inside the sandbox are confused. Three isolation levels make sense: share with host, share across projects (current), or isolate per project.

Asylum currently uses `fmt.Print` + `fmt.Scanln` for all prompts (cleanup, onboarding). This works for y/n but is inadequate for multi-option selection. Bubble Tea provides a mature TUI framework that will be reused for kit activation prompts, onboarding improvements, and future interactive features.

## Goals / Non-Goals

**Goals:**
- Reusable TUI prompt abstractions (`tui.Select`, `tui.MultiSelect`)
- Agent config isolation level as a per-agent config field
- First-run prompt for Claude isolation when not configured
- Persist the choice to config.yaml so the user isn't asked again
- All three isolation modes working: shared, isolated, project

**Non-Goals:**
- Migrating existing prompts (cleanup, onboarding) to Bubble Tea — separate follow-up
- Isolation for non-Claude agents — they can use the same mechanism later
- Symlink resolution in shared mode — already handled by the host-user-alignment change

## Decisions

### 1. TUI package API

```go
package tui

// Option represents a choice in a select prompt.
type Option struct {
    Label       string
    Description string // shown below the label in dimmer text
}

// Select shows a single-choice prompt and returns the selected index.
// defaultIdx is pre-selected. Returns -1 and error if cancelled (Ctrl+C/Esc).
func Select(title string, options []Option, defaultIdx int) (int, error)

// MultiSelect shows a multi-choice prompt and returns selected indices.
// defaultSelected contains initially checked indices.
func MultiSelect(title string, options []Option, defaultSelected []int) ([]int, error)
```

Both functions block until the user makes a choice. They use Bubble Tea internally but the caller never touches Bubble Tea types. Non-interactive mode (no TTY) falls back to the default selection with a log message.

### 2. AgentConfig gains Config field

```go
type AgentConfig struct {
    Config string `yaml:"config,omitempty"` // shared, isolated, project
}
```

Values:
- `shared`: mount `~/.claude` → `~/.claude` (host config used directly)
- `isolated`: mount `~/.asylum/agents/claude/` → `~/.claude` (current behavior, default)
- `project`: mount `~/.asylum/projects/<container>/claude-config/` → `~/.claude` (per-project)

Empty string means "not configured yet" — triggers the prompt.

### 3. First-run prompt flow

```
┌─ Claude Configuration ───────────────────────────────┐
│                                                       │
│  How should Claude's config be managed?               │
│                                                       │
│    ○ Shared with host                                │
│        Use your host ~/.claude directly.              │
│        Changes sync both ways.                        │
│                                                       │
│  ● Isolated (recommended)                            │
│        Separate from host, shared across projects.    │
│        This is the current default.                   │
│                                                       │
│    ○ Project-isolated                                │
│        Separate config per project.                   │
│        No state shared between projects.              │
│                                                       │
│  ↑/↓ navigate  •  enter select  •  esc cancel        │
└───────────────────────────────────────────────────────┘
```

The prompt appears after config loading but before container start, only when:
1. The agent is Claude
2. `agents.claude.config` is empty/unset
3. Stdin is a TTY

After selection, the value is written to `~/.asylum/config.yaml` using `yaml.Node` manipulation to preserve comments.

### 4. Volume mount logic

In `container.go`, the agent config mount branches on the isolation level:

```go
switch isolationLevel {
case "shared":
    // Mount host native config dir directly
    vol(nativeDir, containerDir, "")
case "project":
    // Mount from per-project directory
    projConfigDir := filepath.Join(home, ".asylum", "projects", containerName, agentName+"-config")
    os.MkdirAll(projConfigDir, 0755)
    vol(projConfigDir, containerDir, "")
default: // "isolated" or empty (backwards compatible)
    // Mount from asylum agents dir (current behavior)
    vol(agentDir, containerDir, "")
}
```

### 5. EnsureAgentConfig changes

For `shared` mode, `EnsureAgentConfig` is skipped entirely — the host dir is used as-is.
For `project` mode, `EnsureAgentConfig` seeds from the host config into the per-project dir on first run (same as current isolated behavior, but per-project).
For `isolated` mode, behavior is unchanged.

### 6. Config persistence

After the TUI prompt, the selected value is written back to `~/.asylum/config.yaml`. This uses `yaml.Node` tree manipulation (same approach as kit config sync) to insert the `config:` field under the `claude:` agent entry without destroying comments or formatting.

## Risks / Trade-offs

**New dependency (Bubble Tea)** → Significant addition to a single-dependency project. But it's well-maintained, widely used, and the alternative (hand-rolling TUI) would be worse long-term. The dependency is justified by the roadmap of interactive features.

**Shared mode exposes host config to container writes** → By design. Users who choose shared explicitly want this. The prompt description makes the trade-off clear.

**Config write after prompt** → Could fail (permissions, disk full). In that case the prompt would re-appear next run. Acceptable UX — the user can also manually set it in config.
