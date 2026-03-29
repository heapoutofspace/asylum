## MODIFIED Requirements

### Requirement: cx rules snippet
The cx kit SHALL NOT populate the `RulesSnippet` field. Instead, the cx kit SHALL generate a standalone rules file via `cx skill` during the Docker build and place it in the agent's rules directory via the entrypoint. The kit SHALL retain its `Tools` field so that `cx` continues to appear in the aggregated "Kit Tools" list.

#### Scenario: Rules file contains cx section
- **WHEN** sandbox rules are assembled with cx kit active
- **THEN** the assembled `asylum-sandbox.md` SHALL NOT contain a cx-specific rules snippet
- **AND** the cx tool SHALL still appear in the "Kit Tools" section

#### Scenario: Standalone cx rules file present
- **WHEN** the container starts with cx kit active and `cx skill` succeeded during build
- **THEN** a standalone `cx.md` rules file SHALL be bind-mounted into the agent's rules directory

### Requirement: cx installation via install script
The kit SHALL provide a DockerSnippet that installs the cx CLI by downloading and running the install script from the cx repository. After installation, the snippet SHALL also run `cx skill` to generate the rules file at `/tmp/asylum-kit-rules/cx.md`. The installed binary SHALL be available on PATH.

#### Scenario: cx installed and rules generated in image
- **WHEN** the cx kit is active and the Docker image is built
- **THEN** the `cx` command is available on PATH inside the container
- **AND** `/tmp/asylum-kit-rules/cx.md` contains the output of `cx skill`
