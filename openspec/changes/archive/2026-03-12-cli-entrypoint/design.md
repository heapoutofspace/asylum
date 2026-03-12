## Context

PLAN.md sections 4.5–4.6 define CLI flags and commands. The entry point must handle flag parsing, determine the mode, and ultimately exec into docker.

## Goals / Non-Goals

**Goals:**
- Parse all CLI flags per PLAN.md section 4.5
- Determine mode from positional args (agent, shell, admin shell, ssh-init, cleanup, arbitrary command)
- Load config, select agent, ensure images, build run args, exec
- Process replacement via `syscall.Exec` on Linux/macOS
- Unknown flags passed through to the agent

**Non-Goals:**
- No cobra/urfave dependency — manual flag parsing is cleaner for this passthrough-heavy CLI

## Decisions

- **Manual argument parsing**: Asylum's CLI has passthrough semantics — unknown flags go to the agent. Standard `flag` package doesn't support this well. A simple hand-rolled parser handles known flags and collects the rest.
- **syscall.Exec**: On Unix, replace the asylum process with `docker run`. This gives Docker the PID 1 signals and terminal control. The `docker` binary path is resolved via `exec.LookPath`.
- **Error flow**: All errors go through `log.Error` and `os.Exit(1)`. No error returns from main — it either execs or exits.

## Risks / Trade-offs

- Manual argument parsing is more code than using a framework, but the passthrough requirement makes frameworks awkward.
