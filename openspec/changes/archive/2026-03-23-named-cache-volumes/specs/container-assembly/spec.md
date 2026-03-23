## MODIFIED Requirements

### Requirement: Common volume mounts
The container SHALL include all common mounts: project dir at real path, gitconfig, ssh, caches (as named Docker volumes), history, custom volumes, and direnv.

#### Scenario: All common mounts present
- **WHEN** gitconfig exists, ssh dir exists, and project has .envrc
- **THEN** all conditional mounts are included in the args

#### Scenario: Missing optional paths
- **WHEN** gitconfig and ssh dir do not exist
- **THEN** those mounts are omitted, all others remain

#### Scenario: Cache directories use named volumes
- **WHEN** the container is started
- **THEN** cache directories (npm, pip, maven, gradle) are mounted as named Docker volumes using `--mount type=volume,src=<container-name>-cache-<tool>,dst=<path>`

#### Scenario: No host cache directory created
- **WHEN** the container is started
- **THEN** no `~/.asylum/cache/` directory is created on the host
