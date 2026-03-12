## ADDED Requirements

### Requirement: Flag parsing
The CLI SHALL parse `-a/--agent`, `-p`, `-v`, `--java`, `-n/--new`, `--cleanup`, `--rebuild`, and `-h/--help` flags. Unknown flags SHALL be passed through to the agent.

#### Scenario: Known flags consumed
- **WHEN** `asylum -a gemini -p 3000` is run
- **THEN** agent is set to gemini, port 3000 is forwarded, no passthrough args

#### Scenario: Unknown flags passed through
- **WHEN** `asylum -a gemini -p "fix the bug"` is run
- **THEN** `-p "fix the bug"` is passed to the agent as extra args

### Requirement: Command dispatch
The CLI SHALL dispatch to agent mode (default), shell mode, ssh-init, cleanup, or arbitrary command based on positional args.

#### Scenario: Default invocation
- **WHEN** `asylum` is run with no positional args
- **THEN** the selected agent starts in YOLO mode

#### Scenario: Shell mode
- **WHEN** `asylum shell` is run
- **THEN** an interactive zsh shell starts

#### Scenario: Arbitrary command
- **WHEN** `asylum ls -la` is run
- **THEN** `ls -la` runs in the container

### Requirement: Process replacement
The CLI SHALL use `syscall.Exec` to replace itself with the docker process.

#### Scenario: Exec into docker
- **WHEN** the docker run args are assembled
- **THEN** `syscall.Exec` is called with the docker binary path and args
