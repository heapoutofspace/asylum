## Why

Users cannot pin or rollback to a specific version — `self-update` always fetches the latest release. When a new version introduces a regression, the only workaround is `--safe` (dev channel) or manually downloading from GitHub. Adding an optional version argument makes upgrades and rollbacks trivial. The `selfupdate` alias removes a common typo friction point.

## What Changes

- `asylum self-update` accepts an optional version argument (e.g., `asylum self-update 0.4.0` or `asylum self-update v0.4.0`).
- When a version is specified, that exact release is fetched from GitHub instead of "latest".
- The `--dev` flag and version argument are mutually exclusive (error if both given).
- `asylum selfupdate` is accepted as an alias for `self-update` in CLI dispatch.

## Capabilities

### New Capabilities
- `version-targeted-update`: Accept an optional version argument in `self-update` to install a specific GitHub release

### Modified Capabilities
- `self-update`: Add version argument support and version/dev mutual exclusivity
- `cli-dispatch`: Accept `selfupdate` alias and route version argument to self-update

## Impact

- **CLI parsing** (`cmd/asylum/main.go`): New positional arg after `self-update`, `selfupdate` alias in subcommand dispatch
- **selfupdate package** (`internal/selfupdate/`): New `fetchRelease` path for tagged versions, `Run` accepts optional target version
- **Specs**: `self-update` and `cli-dispatch` specs updated with new scenarios
