## Why

When new kits are added in a release, existing users never see them — their `config.yaml` already exists so `WriteDefaults` is a no-op, and migration only runs for version-gated schema changes. There is no mechanism to introduce new kits to existing installations, prompt for activation, or even inform the user. This also means kit activation tiers (always-on, default-on, opt-in) have no runtime enforcement beyond the initial config generation.

## What Changes

- **Kit state tracking** — a new `~/.asylum/state.json` file (machine-managed, not user-edited) records which kits the installation has seen. On each run, the registered kit set is compared against this state to detect newly added kits.
- **Config sync via `yaml.Node`** — when new kits are detected, their config entries are inserted into the existing `config.yaml` using `yaml.Node` tree manipulation (preserving comments, ordering, and user edits). Active kits are added as real YAML nodes; opt-in kits are added as comments.
- **Three-tier activation model** — kits declare their tier: `always-on` (active regardless of config, like `shell`), `default-on` (added uncommented to config, like `docker`), or `opt-in` (added as commented YAML, like `apt`). The tier determines both config insertion behavior and whether a user prompt is shown.
- **User consent for default-on kits** — when a new default-on kit is detected in an interactive terminal, the user is prompted before it is activated. In non-interactive mode, new default-on kits are added as commented-out entries (visible but inactive) to avoid changing the sandbox without consent.
- **Kit `ConfigNode` method** — each kit provides structured `yaml.Node` trees for config insertion instead of (or alongside) the existing text-based `ConfigSnippet`. This enables robust YAML manipulation that survives user edits to the config file.

## Capabilities

### New Capabilities
- `kit-state-tracking`: Persistent tracking of known kits in `~/.asylum/state.json`, detection of newly added kits on startup
- `kit-config-sync`: Inserting new kit entries into existing `config.yaml` using `yaml.Node` tree manipulation with comment preservation
- `kit-activation-prompt`: Interactive consent flow for activating new default-on kits, with non-interactive fallback

### Modified Capabilities
- `config-defaults`: `WriteDefaults` and `DefaultConfig` will use `yaml.Node`-based assembly instead of text concatenation
- `kit-defaults`: Kit struct gains a `ConfigNode` method/field and a formal activation tier replacing the boolean `DefaultOn`

## Impact

- **`internal/kit/`**: Kit struct changes — `DefaultOn bool` replaced by an activation tier enum, `ConfigSnippet string` replaced or supplemented by `ConfigNode` returning `[]*yaml.Node`
- **`internal/config/defaults.go`**: Rewritten to assemble config via `yaml.Node` instead of string concatenation
- **`internal/config/migrate.go`**: `migrateGlobalConfig` updated to use `yaml.Node` insertion for kit sync
- **New `internal/state/` package** (or `internal/config/state.go`): Manages `~/.asylum/state.json` load/save
- **`cmd/asylum/main.go`**: Wires kit sync check into startup, before config load. Passes terminal detection for interactive prompts.
- **All kit files** (`internal/kit/*.go`): Each kit updated to provide `ConfigNode` and declare its activation tier
