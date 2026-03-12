## Why

The CLI entry point ties everything together: argument parsing, config loading, agent selection, image building, container assembly, and process replacement via `syscall.Exec`.

## What Changes

- Rewrite `cmd/asylum/main.go` with full argument parsing
- Flag handling: `-a/--agent`, `-p`, `-v`, `--java`, `-n/--new`, `--cleanup`, `--rebuild`, `-h/--help`
- Positional args: `shell`, `ssh-init`, or passthrough to agent
- Config loading, agent lookup, image ensure, container run args, exec

## Capabilities

### New Capabilities
- `cli-dispatch`: CLI argument parsing, config loading, and command dispatch with process replacement

### Modified Capabilities

None.

## Impact

- Rewrites `cmd/asylum/main.go` (currently just prints version)
- This is the only user-facing entry point
