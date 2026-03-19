## MODIFIED Requirements

### Requirement: Detect running container
The docker package SHALL provide a function to check if a container with a given name is currently running.

#### Scenario: Container is running
- **WHEN** a container named `asylum-<hash>` is running
- **THEN** `IsRunning("asylum-<hash>")` returns `true`

#### Scenario: Container is not running
- **WHEN** no container with that name exists
- **THEN** `IsRunning("asylum-<hash>")` returns `false`

#### Scenario: Container exists but is stopped
- **WHEN** a container with that name exists but is in exited/dead state
- **THEN** `IsRunning("asylum-<hash>")` returns `false`

### Requirement: Exec into running container for shell mode
When a container is already running for the current project and the user runs `asylum shell`, asylum SHALL exec into the running container instead of starting a new one.

#### Scenario: Shell with running container
- **WHEN** the user runs `asylum shell` and a container is running for the project
- **THEN** asylum runs `docker exec -it <container-name> /bin/zsh`

#### Scenario: Admin shell with running container
- **WHEN** the user runs `asylum shell --admin` and a container is running for the project
- **THEN** asylum runs `docker exec -it -u root <container-name> /bin/zsh`

### Requirement: Exec into running container for run mode
When a container is already running and the user runs `asylum run <cmd>`, asylum SHALL exec the command in the running container.

#### Scenario: Run command with running container
- **WHEN** the user runs `asylum run echo hello` and a container is running
- **THEN** asylum runs `docker exec -it <container-name> echo hello`

### Requirement: Skip image build when exec-ing
When asylum detects it will exec into a running container, it SHALL skip the image build step.

#### Scenario: No image build on exec
- **WHEN** a container is running and any mode is used
- **THEN** `EnsureBase` and `EnsureProject` are not called

## REMOVED Requirements

### Requirement: Agent mode does not exec
**Reason**: Agent mode now execs into running containers, same as shell/run modes.
**Migration**: No changes needed. Running `asylum` when a container is running now starts the agent inside the existing container instead of failing.

## ADDED Requirements

### Requirement: Detached container lifecycle
When no container is running, asylum SHALL start the container in detached mode with an idle process, then exec the session into it.

#### Scenario: First invocation starts detached container
- **WHEN** no container is running and the user runs `asylum`
- **THEN** the container is started detached with an idle process, then the agent is exec'd into it

#### Scenario: First invocation still builds images
- **WHEN** no container is running
- **THEN** `EnsureBase` and `EnsureProject` are called before starting the container

### Requirement: Exec agent into running container
When a container is already running for the current project and the user runs `asylum` (agent mode), asylum SHALL exec the agent into the running container.

#### Scenario: Agent exec with running container
- **WHEN** the user runs `asylum` and a container is running for the project
- **THEN** asylum execs the agent command into the running container via `docker exec -it`

#### Scenario: Agent exec respects resume
- **WHEN** the user runs `asylum` with a running container and a previous session exists
- **THEN** the exec'd agent uses `--continue` for Claude

#### Scenario: Agent exec respects new session flag
- **WHEN** the user runs `asylum -n` with a running container
- **THEN** the exec'd agent starts a new session (no `--continue`)

### Requirement: Container cleanup after last session
After any exec'd session exits, asylum SHALL check if other sessions remain in the container and remove the container if none do.

#### Scenario: Last session exits
- **WHEN** the last exec'd session in a container exits
- **THEN** the container is stopped and removed

#### Scenario: Other sessions still running
- **WHEN** an exec'd session exits but other sessions are still running
- **THEN** the container continues running

### Requirement: Independent session exit
Each asylum session SHALL be able to exit independently without affecting other running sessions in the same container.

#### Scenario: First session exits, others continue
- **WHEN** the first `asylum` session exits and a second session is still running
- **THEN** the second session continues running, the container stays alive
