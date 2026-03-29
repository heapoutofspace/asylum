## 1. Kit Activation Tier

- [x] 1.1 Define `Tier` type (int enum: `TierAlwaysOn`, `TierDefault`, `TierOptIn`) in `internal/kit/kit.go`
- [x] 1.2 Replace `DefaultOn bool` with `Tier` field on Kit struct
- [x] 1.3 Update all kit registrations to use the new tier: `shell`/`node`/`title` → `TierAlwaysOn`, `docker`/`java`/`python`/`github`/`openspec` → `TierDefault`, `apt` → `TierOptIn`
- [x] 1.4 Update `Resolve` to check `k.Tier == TierAlwaysOn` instead of `k.DefaultOn`
- [x] 1.5 Update tests that reference `DefaultOn`

## 2. Kit ConfigNodes

- [x] 2.1 Add `ConfigNodes() []*yaml.Node` method or field to Kit struct (returns key + value node pairs for the kits mapping)
- [x] 2.2 Add `ConfigComment` field to Kit struct for opt-in/commented kits (raw comment text)
- [x] 2.3 Implement `ConfigNodes` for each kit: `docker`, `java`, `python`, `node`, `github`, `openspec`
- [x] 2.4 Implement `ConfigComment` for each opt-in kit: `apt`, `shell`, `title`
- [x] 2.5 Keep `ConfigSnippet`/`AssembleConfigSnippets` for initial config generation; `ConfigNodes` used for sync

## 3. State Tracking

- [x] 3.1 Create state file load/save functions (`LoadState`/`SaveState` for `~/.asylum/state.json`)
- [x] 3.2 Define state struct with `KnownKits []string` field
- [x] 3.3 Implement new kit detection: compare registered kit names against `KnownKits`
- [x] 3.4 Add tests for state load/save, new kit detection, and state-file-deleted scenario

## 4. Config Sync via yaml.Node

- [x] 4.1 Implement `findKitsNode` — parse config file as `yaml.Node`, walk to find the `kits` mapping node
- [x] 4.2 Implement `insertKitNode` — append a kit's `ConfigNodes` key-value pair to the kits mapping
- [x] 4.3 Implement `insertKitComment` — append a commented-out kit block as `FootComment` on the kits mapping
- [x] 4.4 Implement `kitExistsInConfig` — check if a kit key is already present in the kits mapping node
- [x] 4.5 Handle missing `kits` mapping — create it if absent
- [x] 4.6 Add tests: insert into existing config with comments, insert with no kits key, skip existing kit, round-trip preserves formatting

## 5. Activation Prompt

- [x] 5.1 Implement prompt function for `TierDefault` kits: "New kit: X — activate? [Y/n]"
- [x] 5.2 Implement info messages for `TierAlwaysOn` and `TierOptIn` kits
- [x] 5.3 Implement non-interactive fallback: `TierDefault` kits added as comments when not a terminal
- [x] 5.4 Skip sync for utility subcommands (`self-update`, `version`, `ssh-init`, `cleanup`) — handled by caller in main.go

## 6. Wire Into Startup

- [x] 6.1 Call kit sync in `main()` after project dir resolution, before `config.Load` — load state, detect new kits, run sync flow, save state
- [x] 6.2 Pass terminal detection (`term.IsTerminal()`) to the sync flow
- [x] 6.3 Update `WriteDefaults` and `DefaultConfig` — kept text-based for fresh config generation; yaml.Node used for incremental sync
- [x] 6.4 Update `migrateGlobalConfig` — kept text overlay for v1→v2 migration; new kit sync flow handles incremental additions

## 7. Specs

- [x] 7.1 Update `openspec/specs/kit-defaults/spec.md` with tier-based activation
- [x] 7.2 Update `openspec/specs/config-defaults/spec.md` with `ConfigNodes`-based assembly
