## Context

The `Merge()` function in `internal/config/config.go` currently replaces the entire `Kits` map when the overlay has a non-nil `Kits` field. This means any project `.asylum` file that sets even one kit option wipes all globally-configured kits. The same applies to `Agents`.

The config-loading spec already describes "maps merge per-key" but the implementation never matched.

## Goals / Non-Goals

**Goals:**
- Deep-merge `Kits` maps: overlay keys add to or override base keys, non-overlapping base keys are preserved
- Deep-merge `Agents` maps with the same per-key semantics
- Within `KitConfig`, merge individual fields with appropriate semantics (scalars replace, accumulating lists concat)
- Maintain existing behavior for everything else (scalars, ports, volumes, env)

**Non-Goals:**
- Changing CLI flag behavior (`--kits` still does whole-map replacement — it's an explicit override)
- Adding a way to remove a global kit from a project config (existing `disabled: true` already handles this)
- Changing any config file format or adding new fields

## Decisions

### Per-key map merge for Kits and Agents
Overlay keys are merged into the base map. If a key exists in both, the overlay's value is used (with field-level merge for KitConfig). Keys only in the base are preserved.

**Alternative**: Keep whole-map replacement, require users to re-declare global kits in project configs. Rejected because it defeats the purpose of layered config.

### Field-level merge within KitConfig
A new `mergeKitConfig(base, overlay *KitConfig) *KitConfig` function handles field-by-field merging:
- **Scalars replace**: `Disabled`, `DefaultVersion`, `ShadowNodeModules`, `Onboarding`, `TabTitle`, `AllowAgentTermTitle`, `Count` — overlay wins when non-zero/non-nil
- **Lists concatenate**: `Packages`, `Build` — project adds to global
- **Lists replace**: `Versions` — overlay replaces entirely (version lists are a complete declaration, not additive)

**Rationale**: Packages and build commands are naturally additive (project needs extra packages on top of global ones). Versions are a complete set declaration (project wants exactly these versions, not global + local).

### AgentConfig stays simple
`AgentConfig` is currently an empty struct. Merge is trivial: overlay presence means the agent is active. No field-level merge needed until AgentConfig gains fields.

## Risks / Trade-offs

- **[Breaking for replacement users]** → Anyone relying on project config to suppress global kits must switch to `disabled: true`. Mitigated by: this is the correct way to disable a kit, and the previous behavior was a bug relative to the spec.
- **[Nil vs empty KitConfig]** → A kit entry with `nil` KitConfig (e.g., `openspec:` with no sub-fields) vs a non-nil empty KitConfig must both work correctly in merge. The `mergeKitConfig` function handles nil inputs explicitly.
