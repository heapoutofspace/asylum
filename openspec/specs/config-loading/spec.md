## ADDED Requirements

### Requirement: Three-layer config loading
The config system SHALL load config from `~/.asylum/config.yaml`, `$project/.asylum`, and `$project/.asylum.local` in order, merging each layer on top of the previous. Before loading, each file SHALL be migrated from v1 format if necessary.

#### Scenario: All three files present
- **WHEN** all three config files exist with different values
- **THEN** values are merged according to merge semantics (scalars last-wins, lists concat, maps merge per-key with field-level merge within KitConfig)

#### Scenario: Missing files are skipped
- **WHEN** one or more config files do not exist
- **THEN** loading succeeds with values from the files that do exist

#### Scenario: Invalid YAML
- **WHEN** a config file contains invalid YAML
- **THEN** an error is returned

#### Scenario: Project kits supplement global kits
- **WHEN** global config has `kits: {node: {}, openspec: {}}` and project config has `kits: {shell: {}}`
- **THEN** the merged result has all three kits active

### Requirement: Scalar merge semantics
Scalar config values (agent, release-channel) SHALL use last-value-wins when merging layers.

#### Scenario: Agent override
- **WHEN** global config sets `agent: claude` and project config sets `agent: gemini`
- **THEN** the merged result has `agent: gemini`

### Requirement: List merge semantics
List config values (ports, volumes) SHALL be concatenated across layers.

#### Scenario: Ports concatenation
- **WHEN** global config has `ports: ["3000"]` and project config has `ports: ["8080"]`
- **THEN** the merged result has `ports: ["3000", "8080"]`

### Requirement: CLI flag overlay
CLI scalar flags SHALL override all config layers. CLI list flags SHALL be appended to merged config values.

#### Scenario: Agent flag overrides config
- **WHEN** config sets `agent: claude` and CLI flag sets `-a codex`
- **THEN** the final agent is `codex`

#### Scenario: Kits flag overrides config
- **WHEN** config has `kits: {java: {}, python: {}}` and CLI passes `--kits java`
- **THEN** the final kits map contains only java

### Requirement: Release channel config field
The config system SHALL support an optional `release-channel` scalar field with values `stable` or `dev`. It follows scalar merge semantics (last value wins across layers).

#### Scenario: Release channel set in global config
- **WHEN** `~/.asylum/config.yaml` contains `release-channel: dev`
- **THEN** the loaded config has `ReleaseChannel` set to `"dev"`

#### Scenario: Not set defaults to empty
- **WHEN** no config file sets `release-channel`
- **THEN** the loaded config has `ReleaseChannel` set to `""` (callers treat empty as stable)

### Requirement: Read Java version from .tool-versions
The config system SHALL read `.tool-versions` from the project directory and use the Java version as `kits.java.default-version` when not already set by asylum config or CLI flags.

#### Scenario: .tool-versions provides Java version
- **WHEN** `.tool-versions` contains `java 21.0.2` and no asylum config sets java's default-version
- **THEN** the loaded config has `kits.java.default-version` set to `"21.0.2"`

#### Scenario: Asylum config overrides .tool-versions
- **WHEN** `.tool-versions` contains `java 21.0.2` and config sets `kits: {java: {default-version: "17"}}`
- **THEN** the loaded config has java default-version set to `"17"`
