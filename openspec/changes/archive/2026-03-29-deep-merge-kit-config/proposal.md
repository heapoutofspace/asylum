## Why

Project `.asylum` configs that define any `kits:` key completely replace the global `~/.asylum/config.yaml` kits map. This means globally-enabled kits (like `openspec`) silently disappear when a project adds even one kit option (like `shell: {build: ...}`). The existing config-loading spec already says "maps merge per-key" but the implementation does whole-map replacement.

## What Changes

- **BREAKING**: `Merge()` deep-merges the `Kits` map per-key instead of replacing it entirely. A project config that defines `shell:` now adds/overrides only the `shell` kit entry, preserving all other global kits.
- **BREAKING**: `Merge()` deep-merges the `Agents` map per-key (same treatment for consistency).
- Within each `KitConfig`, fields merge with field-appropriate semantics:
  - Scalars (Disabled, DefaultVersion, ShadowNodeModules, Onboarding, TabTitle, AllowAgentTermTitle, Count): last-wins
  - Lists that accumulate (Packages, Build): concatenated
  - Lists that replace (Versions): last-wins

## Capabilities

### New Capabilities
- `deep-merge-kit-config`: Per-key deep merge of kits and agents maps across config layers, with field-level merge semantics within KitConfig

### Modified Capabilities
- `config-loading`: The kits/agents merge semantics change from whole-map replacement to per-key deep merge

## Impact

- `internal/config/config.go`: `Merge()` function rewritten for kits/agents
- `internal/config/config_test.go`: New/updated tests for deep merge behavior
- Any project `.asylum` files that relied on kits replacement to disable global kits will need to use `disabled: true` instead
