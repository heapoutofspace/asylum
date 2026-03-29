## MODIFIED Requirements

### Requirement: Profiles field in config
The Config struct SHALL include a `Kits` field that is a map of kit name to KitConfig. Kit presence in the map means activation. KitConfig SHALL include a `Disabled` field that when true excludes the kit from resolution.

#### Scenario: Kit disabled in config
- **WHEN** config has `kits: {shell: {disabled: true}}`
- **THEN** `KitActive("shell")` returns false

#### Scenario: Kit disabled at project level overrides global
- **WHEN** global config has `kits: {github: {}}` and project config has `kits: {github: {disabled: true}}`
- **THEN** the github kit is not active in the merged config
