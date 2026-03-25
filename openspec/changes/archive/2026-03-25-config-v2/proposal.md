## Why

The config format grew organically across multiple changes — profiles, agents, features, packages, onboarding, and versions are all separate top-level fields with no clear grouping. The new kit system provides an opportunity to restructure config around kits as the primary organizing concept. Options that previously lived in flat maps (`features`, `packages`, `versions`, `onboarding`) now live inside the kit they belong to, making it obvious what affects what. Additionally, first-run users get a well-commented default config, and existing configs are migrated automatically.

## What Changes

- **Rename profiles → kits** everywhere: config field, CLI flag, internal package name, all references
- **Restructure Config struct** to match new format: `kits` is a map of kit-name → kit-config (presence = activated), `agents` is a map of agent-name → agent-config
- **New kit types beyond languages**: `apt` (system packages), `title` (terminal title config), `shell` (custom build/start commands)
- **Per-kit options**: `versions`, `default-version`, `packages`, `shadow-node-modules`, `onboarding` move from top-level into their respective kits
- **Remove top-level fields**: `features`, `packages`, `versions`, `onboarding`, `tab-title`, `profiles` — all absorbed into kits
- **Config version field**: `version: 0.2` in `~/.asylum/config.yaml` (not in project configs)
- **Default config generation**: on first run, write `~/.asylum/config.yaml` with defaults and comments using `yaml.Node` to preserve formatting
- **Config migration**: detect old format and rewrite to new format, preserving user values
  - Global config: detected by `version` field (absent or < 0.2)
  - Project configs (`.asylum`, `.asylum.local`): detected by presence of `features` key
- **CLI flag rename**: `--profiles` → `--kits`

## Capabilities

### New Capabilities
- `config-migration`: Automatic migration from v1 to v2 config format, applied to global and project configs
- `config-defaults`: First-run default config generation with comments describing every option

### Modified Capabilities
- `config-loading`: Config struct restructured around kits map, new parsing logic, version field
- `profile-system`: Renamed to kit system — package name, types, references all change from profile → kit
- `profile-config-integration`: Kit activation by presence in config map, per-kit options replace top-level fields
- `agent-install`: Agents config field changes from `*[]string` to `map[string]AgentConfig`

## Impact

- **internal/profile/** → **internal/kit/**: Package rename, Profile → Kit
- **internal/config/config.go**: Complete Config struct rewrite, new KitConfig/AgentConfig types, migration logic
- **internal/config/defaults.go** (new): Default config template with comments
- **internal/config/migrate.go** (new): v1 → v2 migration logic
- **internal/image/image.go**: Update imports/references from profile → kit
- **internal/container/container.go**: Update references
- **internal/agent/install.go**: AgentInstall resolution changes for map-based agents
- **cmd/asylum/main.go**: Wire new config structure, rename --profiles → --kits, trigger migration
- **assets/entrypoint.tail**: No change (already dynamic)
- Every file that imports `internal/profile` or references `profile.` types
