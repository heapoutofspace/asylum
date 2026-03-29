## 1. Core Merge Logic

- [x] 1.1 Add `mergeKitConfig(base, overlay *KitConfig) *KitConfig` function in `internal/config/config.go` — handles field-level merge (scalars last-wins, Packages/Build concat, Versions replace)
- [x] 1.2 Replace the whole-map `Kits` assignment in `Merge()` with per-key deep merge using `mergeKitConfig`
- [x] 1.3 Replace the whole-map `Agents` assignment in `Merge()` with per-key merge

## 2. Tests

- [x] 2.1 Add unit tests for `mergeKitConfig` covering: scalar override, nil inputs, Packages concat, Build concat, Versions replace, Count zero-vs-nonzero
- [x] 2.2 Add/update `Merge` tests: project kits supplement global kits, project overrides single kit without losing others, overlay nil KitConfig preserves base
- [x] 2.3 Add/update `Merge` tests for agents per-key merge

## 3. Spec and Docs

- [x] 3.1 Update CHANGELOG.md with the breaking change under Unreleased
