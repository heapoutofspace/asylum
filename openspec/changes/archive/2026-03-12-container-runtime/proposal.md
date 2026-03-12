## Why

This is the main orchestration layer. It assembles all the Docker run arguments — volumes, environment variables, ports, working directory, container name — by pulling together config, agent, and image information.

## What Changes

- Create `internal/container` package that builds the full `docker run` argument list
- Common mounts (project dir, git, ssh, caches, history, direnv)
- Agent-specific mounts and env vars
- Custom volumes from config
- Port forwarding
- Container naming from project directory hash
- Shell mode selection (agent, shell, admin, arbitrary command)

## Capabilities

### New Capabilities
- `container-assembly`: Full Docker run argument assembly combining config, agent, and image into a runnable command

### Modified Capabilities

None.

## Impact

- Adds `internal/container/container.go`
- Ties together config, agent, and docker packages
- Called by the CLI entry point to get the final Docker command
