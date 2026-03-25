## ADDED Requirements

### Requirement: Default config on first run
When `~/.asylum/config.yaml` does not exist, the system SHALL write a default config file with all default values and descriptive comments.

#### Scenario: First run creates config
- **WHEN** asylum starts and `~/.asylum/config.yaml` does not exist
- **THEN** the file is created with version 0.2, default agent (claude), default kits (java, python, node), and commented-out sections for additional options

#### Scenario: Config already exists
- **WHEN** asylum starts and `~/.asylum/config.yaml` already exists
- **THEN** no default config is written (existing config is loaded as-is or migrated)

### Requirement: Default config content
The default config SHALL include all currently-supported options with their default values, plus commented-out sections for optional features.

#### Scenario: Default kits present
- **WHEN** the default config is examined
- **THEN** it contains active kits for java (versions 17, 21, 25; default 21), python, and node with their default settings

#### Scenario: Optional sections commented out
- **WHEN** the default config is examined
- **THEN** optional agents (gemini, codex, opencode), optional kits (apt, shell), ports, volumes, and env sections are present but commented out with explanatory comments
