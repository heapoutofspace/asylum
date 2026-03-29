## ADDED Requirements

### Requirement: Interactive prompt for default-on kits
When new `TierDefault` kits are detected in an interactive terminal session, the system SHALL prompt the user for each kit before activating it.

#### Scenario: User accepts kit
- **WHEN** a new `TierDefault` kit `rust` is detected and the user responds `Y` (or presses Enter for default yes)
- **THEN** the kit is inserted as an active entry in `config.yaml`

#### Scenario: User declines kit
- **WHEN** a new `TierDefault` kit `rust` is detected and the user responds `n`
- **THEN** the kit is inserted as a commented-out entry in `config.yaml`

#### Scenario: Multiple new kits
- **WHEN** kits `rust` and `zig` are both newly detected as `TierDefault`
- **THEN** the user is prompted for each one separately

### Requirement: Info message for always-on kits
When new `TierAlwaysOn` kits are detected, the system SHALL display an informational message but not prompt, since these kits activate regardless of config.

#### Scenario: Always-on kit announced
- **WHEN** a new `TierAlwaysOn` kit `shell` is detected
- **THEN** an info message is displayed (e.g., "New kit: shell (always active)")
- **AND** no prompt is shown and no config modification is made

### Requirement: Info message for opt-in kits
When new `TierOptIn` kits are detected, the system SHALL display an informational message directing the user to their config.

#### Scenario: Opt-in kit announced
- **WHEN** a new `TierOptIn` kit `apt` is detected
- **THEN** an info message is displayed (e.g., "New kit available: apt — uncomment in config.yaml to enable")
- **AND** the kit is added as a commented-out entry in the config

### Requirement: Skip prompts for utility subcommands
Kit sync prompts SHALL NOT run for utility subcommands that don't involve container setup.

#### Scenario: Self-update skips sync
- **WHEN** asylum is invoked with `self-update`
- **THEN** kit sync does not run

#### Scenario: Version skips sync
- **WHEN** asylum is invoked with `version`
- **THEN** kit sync does not run
