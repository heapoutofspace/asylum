## ADDED Requirements

### Requirement: Default-on kits
Kits with `DefaultOn: true` SHALL be included in the resolved set when the user specifies explicit kits but does not mention the default-on kit, unless explicitly disabled.

#### Scenario: Explicit kits plus default-on
- **WHEN** config has `kits: {java: {}}` and the shell kit has `DefaultOn: true`
- **THEN** both java and shell are active

#### Scenario: Default-on kit explicitly disabled
- **WHEN** config has `kits: {java: {}, shell: {disabled: true}}`
- **THEN** java is active but shell is not

#### Scenario: No kits key (nil)
- **WHEN** no config layer specifies kits
- **THEN** all kits are active (including default-on kits) — unchanged behavior

#### Scenario: Empty kits map
- **WHEN** config has `kits: {}`
- **THEN** no kits are active (default-on kits are NOT added to explicit empty)

### Requirement: Kit disabling
A kit SHALL be disableable by setting `disabled: true` in its KitConfig. This overrides default-on behavior and can disable globally-configured kits at project level.

#### Scenario: Disable global kit at project level
- **WHEN** global config has `kits: {java: {}, github: {}}` and project config has `kits: {github: {disabled: true}}`
- **THEN** java is active but github is not

#### Scenario: Disabled kit not resolved
- **WHEN** a kit has `disabled: true` in its KitConfig
- **THEN** it is excluded from the resolved kit list and its DockerSnippet is not included
