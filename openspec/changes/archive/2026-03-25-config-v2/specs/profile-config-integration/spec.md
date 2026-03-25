## MODIFIED Requirements

### Requirement: Profiles field in config
The Config struct SHALL include a `Kits` field that is a map of kit name to KitConfig. Kit presence in the map means activation. Nil map means all default kits; empty map means none.

#### Scenario: Kits in YAML config
- **WHEN** a config file contains `kits: {java: {default-version: "21"}, python: {}}`
- **THEN** the parsed Config has Kits map with java and python entries

#### Scenario: No kits key
- **WHEN** no config file specifies `kits`
- **THEN** the parsed Config has Kits as nil (interpreted as "all defaults" at resolution time)

#### Scenario: Empty kits map
- **WHEN** a config file contains `kits: {}`
- **THEN** the parsed Config has Kits as an empty non-nil map (no kits active)

### Requirement: Profiles last-wins across config layers
The `kits` field SHALL follow last-wins semantics: if a later config layer specifies `kits`, it replaces the value from earlier layers entirely.

#### Scenario: Project kits replace global kits
- **WHEN** global config has `kits: {java: {}, python: {}, node: {}}` and project config has `kits: {java: {}}`
- **THEN** the effective kits map contains only java

### Requirement: CLI flag for profiles
A `--kits` CLI flag SHALL allow overriding the kit list from the command line, following the same last-wins semantics.

#### Scenario: CLI overrides all config layers
- **WHEN** config has `kits: {java: {}}` and CLI passes `--kits python,node`
- **THEN** the effective kits are python and node (with default KitConfig)

### Requirement: Per-kit options
Each kit in the map MAY include kit-specific configuration: versions, default-version, packages, shadow-node-modules, onboarding, tab-title, allow-agent-terminal-title, build, start.

#### Scenario: Java kit with custom versions
- **WHEN** config has `kits: {java: {versions: [17, 21], default-version: "17"}}`
- **THEN** the java kit installs versions 17 and 21 with 17 as default

#### Scenario: Node kit with shadow-node-modules disabled
- **WHEN** config has `kits: {node: {shadow-node-modules: false}}`
- **THEN** shadow node_modules is disabled for this project
