## Context

PLAN.md sections 5.2–5.5 and 5.8 define container lifecycle, volume mounts, env vars, port forwarding, and shell modes. This package produces the `docker run` args; the CLI package does the `syscall.Exec`.

## Goals / Non-Goals

**Goals:**
- `RunArgs` struct/function that returns the complete `[]string` of Docker run arguments
- Container naming: `asylum-<sha256(project_dir)[:12]>`
- Hostname: `asylum-<project_name>`
- All common volume mounts per PLAN.md section 5.3
- Agent-specific mounts and env vars
- Custom volume parsing and mounting
- Port forwarding assembly
- Shell mode command selection
- Agent config directory seeding (first-run copy from native config)

**Non-Goals:**
- No `syscall.Exec` — that's in the CLI entry point
- No image building — that's in the image package

## Decisions

- **Single `Run` function**: Takes a `RunOpts` struct with all inputs (config, agent, image tag, project dir, mode, etc.) and returns `[]string` args plus the docker binary path.
- **Directory auto-creation**: Cache and history directories are created automatically if they don't exist.
- **Agent config seeding**: On first run for an agent, copies from the agent's native host config dir (e.g. `~/.claude`) to `~/.asylum/agents/<agent>/`.

## Risks / Trade-offs

- Many conditional mounts based on filesystem state. Each condition is a simple `os.Stat` check.
