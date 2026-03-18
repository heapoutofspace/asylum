## MODIFIED Requirements

### Requirement: ASYLUM_ environment variables
The entrypoint SHALL use ASYLUM_ prefixed environment variables instead of AGENTBOX_.

#### Scenario: Docker flag check
- **WHEN** `ASYLUM_DOCKER=1` is set
- **THEN** the entrypoint starts the Docker daemon

#### Scenario: Java version
- **WHEN** `ASYLUM_JAVA_VERSION` is set
- **THEN** the entrypoint selects that Java version via mise

## REMOVED Requirements

### Requirement: SDKMAN initialization
**Reason**: Replaced by mise for Java/Gradle version management. mise activates faster and has simpler version selection.
**Migration**: No user-facing changes. `ASYLUM_JAVA_VERSION` continues to work. SDKMAN is no longer available inside the container.
