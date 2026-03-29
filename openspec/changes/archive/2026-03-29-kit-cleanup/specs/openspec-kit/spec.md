## ADDED Requirements

### Requirement: OpenSpec CLI kit
The system SHALL provide an `openspec` kit that installs the OpenSpec CLI (`@fission-ai/openspec`) via npm. The kit SHALL declare a dependency on the `node` kit and SHALL be default-on.

#### Scenario: OpenSpec kit active with node
- **WHEN** the openspec kit and node kit are both active
- **THEN** the `openspec` CLI is available in the container

#### Scenario: OpenSpec kit active without node
- **WHEN** the openspec kit is active but the node kit is not
- **THEN** a warning is emitted about the missing node dependency

#### Scenario: OpenSpec kit disabled
- **WHEN** the openspec kit is disabled
- **THEN** the OpenSpec CLI is not installed
