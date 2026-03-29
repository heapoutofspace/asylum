## ADDED Requirements

### Requirement: GitHub CLI kit
The system SHALL provide a `github` kit that installs the GitHub CLI (`gh`) via the official apt repository. The kit SHALL be default-on.

#### Scenario: GitHub kit active
- **WHEN** the github kit is active
- **THEN** the `gh` CLI is available in the container

#### Scenario: GitHub kit disabled
- **WHEN** the github kit is disabled
- **THEN** the `gh` CLI is not installed
