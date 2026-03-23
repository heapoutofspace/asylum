## MODIFIED Requirements

### Requirement: Java installation in Dockerfile
The Dockerfile SHALL install Java versions using plain version numbers (`java@17`, `java@21`, `java@25`) instead of distribution-prefixed names.

#### Scenario: Java install commands
- **WHEN** the base image is built
- **THEN** the Dockerfile runs `mise install java@17 java@21 java@25`

#### Scenario: Default Java version
- **WHEN** the base image is built
- **THEN** the Dockerfile runs `mise use --global java@21`
