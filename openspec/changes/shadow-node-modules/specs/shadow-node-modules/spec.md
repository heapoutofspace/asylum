## ADDED Requirements

### Requirement: Detect node_modules directories
The system SHALL walk the project directory to find `node_modules` directories, skipping nested `node_modules` inside `node_modules` and irrelevant directories.

#### Scenario: Top-level node_modules found
- **WHEN** the project has a `package.json` and a `node_modules` directory
- **THEN** the `node_modules` path is returned

#### Scenario: Monorepo with multiple node_modules
- **WHEN** the project has `node_modules` at the root and under `packages/app/node_modules`
- **THEN** both paths are returned

#### Scenario: Nested node_modules inside node_modules skipped
- **WHEN** `node_modules/some-pkg/node_modules` exists
- **THEN** only the outer `node_modules` is returned

#### Scenario: No package.json
- **WHEN** the project has no `package.json` at the root
- **THEN** no walk is performed and no paths are returned

#### Scenario: Heavy directories skipped
- **WHEN** `node_modules` exists inside `.venv`, `.git`, `vendor`, `target`, `build`, or `dist`
- **THEN** those directories are not walked and the inner `node_modules` is not found

### Requirement: Shadow node_modules with named volumes
Each detected `node_modules` directory SHALL be shadowed with a named Docker volume using `--mount type=volume,src=<name>,dst=<path>`.

#### Scenario: Volume naming
- **WHEN** a `node_modules` at relative path `node_modules` is detected for container `asylum-a1b2c3d4e5f6`
- **THEN** the volume is named `asylum-a1b2c3d4e5f6-npm-<hash>` where `<hash>` is the first 11 hex chars of SHA-256 of the relative path

#### Scenario: Volume persists across restarts
- **WHEN** dependencies are installed inside the container and the container is restarted
- **THEN** the named volume retains the installed dependencies

### Requirement: Feature can be disabled
The shadow feature SHALL be disabled when `features: { shadow-node-modules: false }` is set in config.

#### Scenario: Feature disabled
- **WHEN** config has `features: { shadow-node-modules: false }`
- **THEN** no `--mount` flags for `node_modules` are added

#### Scenario: Feature enabled by default
- **WHEN** config does not mention `shadow-node-modules`
- **THEN** the shadow behavior is active

### Requirement: Cleanup removes shadow volumes
The `--cleanup` command SHALL remove all Docker volumes with the `asylum-` prefix alongside image removal.

#### Scenario: Volumes removed on cleanup
- **WHEN** `asylum --cleanup` is run and asylum-prefixed volumes exist
- **THEN** the volumes are removed

#### Scenario: No volumes to remove
- **WHEN** `asylum --cleanup` is run and no asylum-prefixed volumes exist
- **THEN** cleanup proceeds without error
