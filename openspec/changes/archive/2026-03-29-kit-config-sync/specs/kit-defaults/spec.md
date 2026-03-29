## MODIFIED Requirements

### Requirement: Default-on kits
Kits with tier `TierAlwaysOn` SHALL be included in the resolved set when the user specifies explicit kits but does not mention the always-on kit, unless explicitly disabled. Kits with tier `TierDefault` are active only when present (uncommented) in the config.

#### Scenario: Explicit kits plus always-on
- **WHEN** config has `kits: {java: {}}` and the shell kit has tier `TierAlwaysOn`
- **THEN** both java and shell are active

#### Scenario: Always-on kit explicitly disabled
- **WHEN** config has `kits: {java: {}, shell: {disabled: true}}`
- **THEN** java is active but shell is not

#### Scenario: No kits key (nil)
- **WHEN** no config layer specifies kits
- **THEN** all kits are active (including always-on kits) — unchanged behavior

#### Scenario: Empty kits map
- **WHEN** config has `kits: {}`
- **THEN** no kits are active (always-on kits are NOT added to explicit empty)

### Requirement: Kit activation tier
Each kit SHALL declare an activation tier: `TierAlwaysOn` (active even without config), `TierDefault` (active when present in config, added by default), or `TierOptIn` (only active if user explicitly enables). This replaces the `DefaultOn bool` field.

#### Scenario: Always-on tier
- **WHEN** a kit has tier `TierAlwaysOn`
- **THEN** it is active even if not listed in the config's `kits` map

#### Scenario: Default tier
- **WHEN** a kit has tier `TierDefault`
- **THEN** it is active only if its key is present and uncommented in the config's `kits` map

#### Scenario: Opt-in tier
- **WHEN** a kit has tier `TierOptIn`
- **THEN** it is active only if the user explicitly adds it to their config
