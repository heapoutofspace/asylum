## ADDED Requirements

### Requirement: First-run detection
The system SHALL detect a first-run condition by checking whether `~/.asylum/agents/` directory exists. If it does not exist, the system SHALL trigger the first-run onboarding flow before loading config. The `agents/` directory is created by `EnsureAgentConfig` on the first actual run, making it a reliable signal that distinguishes fresh installs from existing users (since the installer only creates `~/.asylum/bin/`).

#### Scenario: First run — agents directory does not exist
- **WHEN** the user runs `asylum` and `~/.asylum/agents/` does not exist
- **THEN** the system SHALL run the first-run onboarding flow before proceeding

#### Scenario: Subsequent run — agents directory exists
- **WHEN** the user runs `asylum` and `~/.asylum/agents/` already exists
- **THEN** the system SHALL skip first-run onboarding and proceed normally

### Requirement: Credential file detection
The system SHALL check for the existence of the following files on the host:
- `~/.m2/settings.xml` (Maven)

Only files that exist SHALL be offered for mounting.

#### Scenario: Maven settings exist
- **WHEN** `~/.m2/settings.xml` exists on the host
- **THEN** the system SHALL list it in the onboarding prompt

#### Scenario: No credential files exist
- **WHEN** no credential files exist on the host
- **THEN** the system SHALL skip the prompt entirely

### Requirement: Interactive credential mount prompt
When credential files are detected, the system SHALL display the found files and ask the user whether to make them available in the sandbox, using a `[Y/n]` prompt.

#### Scenario: User accepts
- **WHEN** the user responds with empty input, "y", or "Y"
- **THEN** the system SHALL generate `~/.asylum/config.yaml` with read-only volume mounts for the detected credential files

#### Scenario: User declines
- **WHEN** the user responds with "n" or "N"
- **THEN** the system SHALL NOT generate a config file

### Requirement: Config file generation
When the user accepts, the system SHALL write `~/.asylum/config.yaml` with volume entries mapping each detected credential file into the container at the same path, with read-only (`:ro`) option.

#### Scenario: Generated config for Maven
- **WHEN** Maven settings exist and user accepts
- **THEN** `~/.asylum/config.yaml` SHALL contain volume entry `~/.m2/settings.xml:ro`
