## MODIFIED Requirements

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
