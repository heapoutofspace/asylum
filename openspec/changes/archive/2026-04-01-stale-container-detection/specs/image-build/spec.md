## MODIFIED Requirements

### Requirement: Base image auto-rebuild
The image package SHALL detect when the embedded Dockerfile or entrypoint.sh has changed and rebuild the base image automatically. `EnsureBase` SHALL be called on every asylum invocation regardless of container state. When a running container exists and `docker inspect` fails, asylum SHALL treat images as up to date rather than erroring out.

#### Scenario: First build
- **WHEN** no `asylum:latest` image exists
- **THEN** the base image is built with hash and version labels

#### Scenario: Hash matches
- **WHEN** the `asylum.hash` label on `asylum:latest` matches the current asset hash
- **THEN** no rebuild occurs

#### Scenario: Hash differs
- **WHEN** the `asylum.hash` label differs from the current asset hash
- **THEN** the base image is rebuilt and dangling images are pruned

#### Scenario: Called with running container
- **WHEN** a container is already running
- **THEN** `EnsureBase` SHALL still be called and return the expected tag for comparison

### Requirement: Project image generation
The image package SHALL generate a project-specific Dockerfile from the packages config and build it when packages are configured OR when project kits have `EntrypointSnippet`s or `BannerLines`. `EnsureProject` SHALL be called on every asylum invocation regardless of container state.

#### Scenario: No packages configured
- **WHEN** packages config is empty and no project kits have entrypoint snippets or banner lines
- **THEN** `asylum:latest` is returned as the image tag

#### Scenario: Packages configured
- **WHEN** packages config has apt, npm, pip, or run entries
- **THEN** a project image `asylum:proj-<hash>` is built from a generated Dockerfile

#### Scenario: Project image up to date
- **WHEN** `asylum:proj-<hash>` already exists with matching packages hash
- **THEN** no rebuild occurs

#### Scenario: Project kits with entrypoint snippets only
- **WHEN** packages config is empty but project kits have `EntrypointSnippet`s
- **THEN** a project image SHALL be built containing the project entrypoint script

#### Scenario: Called with running container
- **WHEN** a container is already running
- **THEN** `EnsureProject` SHALL still be called and return the expected tag for comparison
