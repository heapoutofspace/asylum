## Why

Running `asylum` in a project where a container is already running fails with a Docker name conflict. Users want to open multiple agent sessions in the same project — e.g., one for coding and one for reviewing. The current workaround is manual `docker exec`, which defeats the purpose.

Additionally, the current lifecycle ties the container to the first session: when it exits, all other sessions die. Each session should be able to exit independently.

## What Changes

- Start containers in detached mode with an idle process instead of exec'ing the agent directly
- All sessions (first and subsequent, all modes) use `docker exec` into the container
- After each session exits, check if other sessions remain — remove the container when the last one exits
- Replace `syscall.Exec` with `cmd.Run()` for exec'd sessions so asylum can do cleanup afterward

## Capabilities

### New Capabilities

### Modified Capabilities

- `container-exec`: All modes (including agent) exec into running containers; detached container lifecycle with automatic cleanup
- `container-assembly`: Container starts detached with idle process; `--rm` replaced by manual cleanup

## Impact

- `cmd/asylum/main.go`: Major refactor — detached container start, exec for all modes, post-exit cleanup, signal forwarding
- `internal/container/container.go`: `RunArgs` produces detached run args; `ExecArgs` handles all modes including agent; new cleanup functions
- `internal/docker/docker.go`: Add functions for counting exec sessions and removing containers
- `assets/entrypoint.sh`: Adjust for detached mode (idle process instead of agent as `$@`)
