## Why

`asylum cleanup` currently removes *all* asylum images, volumes, and cached data across every project. This is destructive for users with multiple projects — cleaning up one project nukes the base image and volumes for all the others, forcing expensive rebuilds. The common case is cleaning up the current project only.

## What Changes

- **Default behavior**: `asylum cleanup` (no flags) scopes cleanup to the current project — removes only that project's image, volumes, and cached project data. The base image and other projects are untouched.
- **`--all` flag**: `asylum cleanup --all` performs the existing global cleanup (all images, all volumes, all cached data) but first shows exactly what will be deleted and requires explicit user confirmation before proceeding.
- **`--cleanup` flag alias**: continues to work, now equivalent to scoped cleanup. `--cleanup --all` enables global cleanup.

## Capabilities

### Modified Capabilities
- `cleanup-command`: Scoped to current project by default, with `--all` for global cleanup with confirmation

## Impact

- **cmd/asylum/main.go**: `parseArgs` accepts `--all` flag on cleanup subcommand; `runCleanup` split into scoped (default) and global (`--all`) paths
- **internal/container**: `ContainerName` already exported — used to derive current project's container name and volume prefix
- **openspec/specs/cleanup-command/spec.md**: Updated scenarios for scoped vs global behavior
