## 1. Docker package: session counting and container management

- [x] 1.1 Add `CountExecSessions(name string) int` — uses `docker top` to count user processes excluding the idle process
- [x] 1.2 Add `RemoveContainer(name string) error` — `docker rm -f`
- [x] 1.3 Add `RunDetached(args []string) error` — runs `docker run -d` and waits for it to return

## 2. Container package: detached run and agent exec

- [x] 2.1 Change `RunArgs` to produce detached args: replace `--rm -it` with `-d`, use `sleep infinity` as the command instead of the agent command
- [x] 2.2 Add `ExecOpts` struct with all fields needed for exec (container name, mode, agent, project dir, extra args, new-session flag, config)
- [x] 2.3 Refactor `ExecArgs` to accept `ExecOpts` and handle all modes including agent (builds agent command with resume/new-session/session-name logic)

## 3. Main: new lifecycle flow

- [x] 3.1 Refactor main to: check if running → if not, build images + start detached container → exec session → after exit, cleanup if last session
- [x] 3.2 Replace `syscall.Exec` with `exec.Command().Run()` for the exec path, with stdin/stdout/stderr connected and signal forwarding
- [x] 3.3 Add post-exit cleanup: count exec sessions, remove container if none remain

## 4. Entrypoint adjustments

- [x] 4.1 Ensure entrypoint works with `sleep infinity` as the command — setup runs, then idles

## 5. Tests

- [x] 5.1 Update `ExecArgs` tests for new `ExecOpts` struct and agent mode coverage
- [x] 5.2 Update `RunArgs` tests for detached mode (no `--rm`, no `-it`, command is `sleep infinity`)
- [x] 5.3 Manually test: start agent, open second agent in same project, exit first — second continues
