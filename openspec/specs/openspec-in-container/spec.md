## ADDED Requirements

### Requirement: OpenSpec CLI available in container
The Asylum container image SHALL have the `openspec` CLI installed globally via npm.

#### Scenario: OpenSpec available
- **WHEN** the container starts
- **THEN** `openspec --version` runs successfully
