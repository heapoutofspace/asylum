## MODIFIED Requirements

### Requirement: Command dispatch
The CLI SHALL dispatch to version display, agent mode (default), shell mode, ssh-init, cleanup, self-update, or arbitrary command based on flags and positional args. The CLI SHALL accept `selfupdate` as an alias for `self-update`.

#### Scenario: Self-update
- **WHEN** `asylum self-update` is run
- **THEN** the self-update logic executes and the process exits before any container setup

#### Scenario: Self-update with dev flag
- **WHEN** `asylum self-update --dev` is run
- **THEN** the self-update targets the dev channel

#### Scenario: Self-update with version argument
- **WHEN** `asylum self-update 0.4.0` is run
- **THEN** the self-update targets the specified version

#### Scenario: Selfupdate alias
- **WHEN** `asylum selfupdate` is run
- **THEN** the self-update logic executes identically to `asylum self-update`

#### Scenario: Selfupdate alias with arguments
- **WHEN** `asylum selfupdate --dev` is run
- **THEN** the self-update targets the dev channel, same as `asylum self-update --dev`
