## MODIFIED Requirements

### Requirement: Dynamic cache directories
Cache directory volume mounts SHALL be aggregated from active kits' CacheDirs fields instead of a hardcoded constant. The container's `--privileged` flag and `ASYLUM_DOCKER` environment variable SHALL be set only when the docker kit is active.

#### Scenario: Docker kit active
- **WHEN** the docker kit is present in the kits config
- **THEN** the container runs with `--privileged` and `ASYLUM_DOCKER=1` is set

#### Scenario: Docker kit inactive
- **WHEN** the docker kit is not present in the kits config
- **THEN** the container runs without `--privileged` and `ASYLUM_DOCKER` is not set
