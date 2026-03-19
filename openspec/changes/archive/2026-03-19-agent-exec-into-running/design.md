## Context

Currently the container lifecycle is tied to a single process: `docker run --rm --init` starts the container, the entrypoint execs into the agent, and when the agent exits, PID 1 exits, Docker removes the container. This means exec'd sessions die when the first session exits.

The user wants any session to be able to exit independently — the container should only be removed when all sessions are gone.

## Goals / Non-Goals

**Goals:**
- Any `asylum` invocation can exit without affecting other running sessions
- Container is automatically cleaned up when the last session exits
- First invocation still handles image build and container creation
- Subsequent invocations skip image build and exec into the running container

**Non-Goals:**
- Hot-reloading config changes into a running container
- Keeping the container alive with zero sessions (there's no daemon)

## Decisions

### Detached container with idle entrypoint

Start the container in detached mode (`docker run -d`) with `tail -f /dev/null` as the entrypoint command. The entrypoint script still runs (git config, mise, SSH permissions, etc.) but instead of exec'ing into an agent, it exec's into an idle process. All sessions — including the first — use `docker exec`.

The flow becomes:
1. **No container running**: `docker run -d --init --name <name> ... tail -f /dev/null` → container starts and idles → `docker exec -it <name> <agent-command>`
2. **Container running**: `docker exec -it <name> <agent-command>`
3. **After exec exits**: count remaining exec sessions → if none, `docker rm -f <name>`

### Use cmd.Run() instead of syscall.Exec for exec

We need control back after `docker exec` exits to do cleanup. Replace `syscall.Exec("docker", "exec", ...)` with `exec.Command("docker", "exec", ...).Run()`. The asylum process stays alive as a thin wrapper, forwarding signals and exit codes.

For `docker run` (container creation), we can still use `cmd.Run()` since the container is detached — it returns immediately.

Keep `syscall.Exec` only for `docker run` in agent mode when no exec-into-running behavior is needed (but actually, since we always use the detached pattern now, this goes away entirely).

### Remove `--rm` flag

Since we manage container removal ourselves (after last session exits), `--rm` is no longer used. Container removal happens in the cleanup step.

### File-based session counter

Each asylum process increments a counter file at `~/.asylum/projects/<container>/sessions` before exec'ing, and decrements it after. When it hits 0, the container is removed. This is simple, reliable, and avoids platform-specific differences in Docker process introspection.

Alternative considered: `docker top` to count processes. Rejected because Docker Desktop on macOS shows different process trees than Linux, and the entrypoint's dockerd/containerd processes made counting unreliable.

### Signal forwarding

Since `asylum` stays alive as a wrapper around `docker exec`, it needs to forward SIGINT/SIGTERM to the docker process so Ctrl+C works correctly. Use `cmd.Process.Signal()` in a signal handler goroutine.

### Entrypoint changes

The entrypoint no longer receives the agent command as `$@`. Instead it always gets `tail -f /dev/null`. The welcome banner should still print on the first exec (not on container start, since that's detached). Move the banner to a wrapper script or skip it — the agent has its own startup output.

Actually, simpler: keep the entrypoint as-is but make the detached command `sleep infinity`. The entrypoint runs all setup (git, mise, SSH, etc.) and then execs into sleep. The welcome banner prints to the detached container's stdout (invisible), which is fine.

## Risks / Trade-offs

- **Stale containers**: If asylum crashes without cleanup, the container stays running. Mitigated by `asylum --cleanup` and `docker rm -f`.
- **No `--rm`**: We manage removal ourselves. If the cleanup check fails, containers accumulate. Could add a startup check that removes stale containers from previous sessions.
- **Signal forwarding complexity**: Using `cmd.Run()` instead of `syscall.Exec` means we need to handle signals, stdin/stdout/stderr, and terminal modes. `docker exec -it` handles most of this, but we need to ensure SIGINT propagates.
- **Welcome banner**: Currently prints during entrypoint. With detached start, it's invisible. Could print on first exec, but that adds complexity. The agent's own startup output is sufficient.
