## 1. Package Rename: profile → kit

- [x] 1.1 Rename `internal/profile/` directory to `internal/kit/`
- [x] 1.2 Rename all types: `Profile` → `Kit`, `SubProfiles` → `SubKits`, update all field names and method names
- [x] 1.3 Update all imports across the codebase: `internal/profile` → `internal/kit`, `profile.Profile` → `kit.Kit`, etc.
- [x] 1.4 Rename CLI flag `--profiles` → `--kits` in `cmd/asylum/main.go` parseArgs and printUsage
- [x] 1.5 Update all tests to use new names
- [x] 1.6 Verify `go build ./...` and `go test ./...` pass after rename

## 2. Config Struct Restructure

- [x] 2.1 Define `KitConfig` struct with fields: Versions, DefaultVersion, Packages, ShadowNodeModules, Onboarding, TabTitle, AllowAgentTermTitle, Build, Start
- [x] 2.2 Define `AgentConfig` struct (empty for now)
- [x] 2.3 Rewrite `Config` struct: replace `Profiles *[]string` with `Kits map[string]*KitConfig`, replace `Agents *[]string` with `Agents map[string]*AgentConfig`, remove `TabTitle`, `Versions`, `Packages`, `Features`, `Onboarding` fields, add `Version` field
- [x] 2.4 Update `CLIFlags`: replace `Profiles *[]string` with `Kits *[]string` (CLI still takes comma-separated names), replace `Agents *[]string` with `Agents *[]string` (same), remove `Java` flag (now via kits.java.default-version)
- [x] 2.5 Rewrite `Merge` function for new struct: kits map last-wins, agents map last-wins, env/ports/volumes as before
- [x] 2.6 Rewrite `applyFlags` for new struct: `--kits` converts list to map with nil values, `--agents` converts list to map with nil values
- [x] 2.7 Add helper methods on Config: `KitActive(name) bool`, `KitOption(name) *KitConfig`, `AgentActive(name) bool`
- [x] 2.8 Update `.tool-versions` Java reading to set `kits.java.default-version`
- [x] 2.9 Write tests for new merge logic, KitActive, AgentActive

## 3. Config Migration

- [x] 3.1 Create `internal/config/migrate.go` with `MigrateV1ToV2(path string) error` — reads file as `yaml.Node`, detects version, transforms node tree, writes back with backup
- [x] 3.2 Implement field mapping: profiles → kits map, versions.java → kits.java.default-version, packages → per-kit packages, features → per-kit fields, onboarding → per-kit onboarding, tab-title → kits.title.tab-title, agents list → agents map
- [x] 3.3 Implement removal of old keys after mapping (features, packages, versions, onboarding, tab-title, profiles)
- [x] 3.4 Add version stamp: set `version: 0.2` after migration
- [x] 3.5 Implement project config detection: migrate if `features` key is present (no version field needed)
- [x] 3.6 Wire migration into `config.Load`: call `MigrateV1ToV2` for each config path before loading
- [x] 3.7 Write tests for migration: v1 global → v2, v1 project → v2, already-v2 skipped, backup created, various field mappings

## 4. Default Config Generation

- [x] 4.1 Create `internal/config/defaults.go` with `WriteDefaults(path string) error` — writes a well-commented default config YAML string to the given path
- [x] 4.2 Define the default config template string matching the config.yaml template: version 0.2, agent claude, agents map (claude active, others commented), kits (java/python/node with defaults, apt/title/shell commented), ports/volumes/env commented
- [x] 4.3 Wire into startup: in `main.go`, before `config.Load`, check if `~/.asylum/config.yaml` exists; if not, call `WriteDefaults`
- [x] 4.4 Write tests: default file content matches expected structure, existing file not overwritten

## 5. Kit Resolution Update

- [x] 5.1 Update `kit.Resolve` to accept `map[string]*KitConfig` instead of `*[]string` — nil map = all defaults, empty map = none, map keys = active kits
- [x] 5.2 Extract kit names from map keys for resolution, pass KitConfig options to resolved kits where applicable
- [x] 5.3 Update all callers of kit.Resolve (main.go, resolveProfileTiers → resolveKitTiers)
- [x] 5.4 Update tests for new Resolve signature

## 6. Agent Resolution Update

- [x] 6.1 Update `agent.ResolveInstalls` to accept `map[string]*AgentConfig` instead of `*[]string` — nil map = claude-only, empty map = none, map keys = active agents
- [x] 6.2 Update all callers (main.go)
- [x] 6.3 Update tests for new resolution signature

## 7. Wiring and Integration

- [x] 7.1 Update `cmd/asylum/main.go`: use new Config fields throughout (kits map for shadow-node-modules check, onboarding check, java version, tab-title, etc.)
- [x] 7.2 Replace all `cfg.Feature(...)` and `cfg.FeatureOff(...)` calls with kit-specific field access
- [x] 7.3 Replace `cfg.Packages` usage with kit-specific packages
- [x] 7.4 Replace `cfg.Versions["java"]` with java kit's DefaultVersion
- [x] 7.5 Replace `cfg.TabTitle` with title kit's TabTitle
- [x] 7.6 Update `container.go`: replace `DefaultCacheDirs` usage, update shadow-node-modules feature check
- [x] 7.7 Update `image.go`: EnsureProject to use kit packages instead of top-level packages
- [x] 7.8 Update CHANGELOG entry under Unreleased
- [x] 7.9 Verify `go build ./...` and `go test ./...` pass

## 8. Cleanup

- [x] 8.1 Remove old `internal/profile/` directory if not already handled by rename
- [x] 8.2 Remove `config.yaml` template from repo root (it was a working document, not a shipped file)
- [x] 8.3 Final pass: grep for any remaining references to "profile" that should be "kit"
