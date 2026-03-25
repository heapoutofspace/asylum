## MODIFIED Requirements

### Requirement: Three-layer config loading
The config system SHALL load config from `~/.asylum/config.yaml`, `$project/.asylum`, and `$project/.asylum.local` in order, merging each layer on top of the previous. Before loading, each file SHALL be migrated from v1 format if necessary.

#### Scenario: All three files present
- **WHEN** all three config files exist with different values
- **THEN** values are merged according to merge semantics (scalars last-wins, lists concat, maps merge per-key)

#### Scenario: Missing files are skipped
- **WHEN** one or more config files do not exist
- **THEN** loading succeeds with values from the files that do exist

#### Scenario: Invalid YAML
- **WHEN** a config file contains invalid YAML
- **THEN** an error is returned

### Requirement: CLI flag overlay
CLI scalar flags SHALL override all config layers. CLI list flags SHALL be appended to merged config values.

#### Scenario: Agent flag overrides config
- **WHEN** config sets `agent: claude` and CLI flag sets `-a codex`
- **THEN** the final agent is `codex`

#### Scenario: Kits flag overrides config
- **WHEN** config has `kits: {java: {}, python: {}}` and CLI passes `--kits java`
- **THEN** the final kits map contains only java

### Requirement: Read Java version from .tool-versions
The config system SHALL read `.tool-versions` from the project directory and use the Java version as `kits.java.default-version` when not already set by asylum config or CLI flags.

#### Scenario: .tool-versions provides Java version
- **WHEN** `.tool-versions` contains `java 21.0.2` and no asylum config sets java's default-version
- **THEN** the loaded config has `kits.java.default-version` set to `"21.0.2"`

#### Scenario: Asylum config overrides .tool-versions
- **WHEN** `.tool-versions` contains `java 21.0.2` and config sets `kits: {java: {default-version: "17"}}`
- **THEN** the loaded config has java default-version set to `"17"`

## REMOVED Requirements

### Requirement: Map-of-lists merge semantics
**Reason**: Top-level `packages` map no longer exists; packages are per-kit fields.
**Migration**: Use `kits.apt.packages`, `kits.node.packages`, `kits.python.packages` instead.

### Requirement: FeatureOff method for default-on features
**Reason**: Top-level `features` map no longer exists; feature flags move into their respective kits as typed fields.
**Migration**: Use kit-specific fields (e.g., `kits.node.shadow-node-modules: false`).
