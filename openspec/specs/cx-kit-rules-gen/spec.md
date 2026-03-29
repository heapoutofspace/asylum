# cx-kit-rules-gen Specification

## Purpose
cx kit generates its own rules file from `cx skill` during Docker build and bind-mounts it into the agent's rules directory via the entrypoint.

## Requirements

### Requirement: cx skill output generated during Docker build
The cx kit's `DockerSnippet` SHALL run `cx skill` after installing cx and write the output to `/tmp/asylum-kit-rules/cx.md`. The command SHALL NOT fail the image build if `cx skill` errors.

#### Scenario: cx skill succeeds during build
- **WHEN** the Docker image is built with the cx kit active
- **THEN** the image SHALL contain `/tmp/asylum-kit-rules/cx.md` with the output of `cx skill`

#### Scenario: cx skill fails during build
- **WHEN** `cx skill` returns a non-zero exit code during the Docker build
- **THEN** the build SHALL continue without error and `/tmp/asylum-kit-rules/cx.md` SHALL either not exist or be empty

### Requirement: cx rules file bind-mounted into agent rules directory at startup
The cx kit's `EntrypointSnippet` SHALL bind-mount `/tmp/asylum-kit-rules/cx.md` onto `~/.claude/rules/cx.md` if the source file exists and is non-empty. The bind mount ensures the shared agent rules directory on the host is not modified.

#### Scenario: Rules file mounted at container start
- **WHEN** the container starts and `/tmp/asylum-kit-rules/cx.md` exists and is non-empty
- **THEN** the file SHALL be bind-mounted to `~/.claude/rules/cx.md`
- **AND** the host filesystem SHALL NOT be modified

#### Scenario: Rules file missing
- **WHEN** the container starts and `/tmp/asylum-kit-rules/cx.md` does not exist
- **THEN** the entrypoint SHALL continue without error and no cx rules file SHALL be placed in the rules directory
