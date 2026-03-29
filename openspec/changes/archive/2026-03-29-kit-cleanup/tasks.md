## 1. Kit System Extensions

- [x] 1.1 Add `Deps []string` and `DefaultOn bool` fields to Kit struct in `internal/kit/kit.go`
- [x] 1.2 Add `Disabled *bool` field to KitConfig in `internal/config/config.go`; update `KitActive` to return false when disabled
- [x] 1.3 Update `kit.Resolve`: after resolving requested kits, add default-on kits that are not already present and not disabled; then validate dependencies and warn for missing ones
- [x] 1.4 Update config merge: when merging kits maps, a kit with `disabled: true` in overlay replaces the base entry (preserving the disable)
- [x] 1.5 Write tests: default-on kit included when not listed, default-on kit excluded when disabled, default-on kit excluded with empty map, dependency warning emitted, no warning when dep satisfied

## 2. New Kits

- [x] 2.1 Create `internal/kit/github.go`: GitHub CLI kit with DockerSnippet (apt repo setup + gh install), `DefaultOn: true`
- [x] 2.2 Create `internal/kit/openspec.go`: OpenSpec CLI kit with DockerSnippet (npm install), `Deps: ["node"]`, `DefaultOn: true`
- [x] 2.3 Create `internal/kit/shell.go`: Shell kit with DockerSnippet (oh-my-zsh install, theme, PATH re-add for fnm/mise in zshrc, direnv hooks for bash/zsh, terminal size handling, tmux config), `DefaultOn: true`

## 3. Move Existing Hardcoded Items

- [x] 3.1 Move maven from core apt-get to `java/maven` sub-kit DockerSnippet (root apt-get install)
- [x] 3.2 Move python build deps (`python3-dev`, `python3-pip`, `python3-venv`, `libssl-dev`, `libffi-dev`) from core apt-get into python kit's DockerSnippet (root apt-get prepended before uv tool installs)
- [x] 3.3 Remove `@fission-ai/openspec@latest` from node kit's npm packages in `internal/kit/node.go`

## 4. Dockerfile and Tail Cleanup

- [x] 4.1 Remove `maven` from core apt-get block in `assets/Dockerfile.core`
- [x] 4.2 Remove `python3-dev python3-pip python3-venv libssl-dev libffi-dev` from core apt-get block
- [x] 4.3 Remove GitHub CLI install block from `assets/Dockerfile.core` (lines 58-67)
- [x] 4.4 Remove GitLab CLI install block from `assets/Dockerfile.core` (lines 69-78) — move to a `gitlab` kit or remove entirely if not needed
- [x] 4.5 Remove oh-my-zsh install, direnv hooks, terminal size handling, and tmux config from `assets/Dockerfile.tail` (moves to shell kit)
- [x] 4.6 Verify: Dockerfile.tail only contains git config, workspace dir, root switch, entrypoint COPY, WORKDIR, USER, ENV after cleanup

## 5. Default Config Update

- [x] 5.1 Update `internal/config/defaults.go`: add github, openspec, and shell to default kits (active by default, commented-out `disabled: true` example)
- [x] 5.2 Update migration: add default-on kits (github, openspec, shell) during v1→v2 migration

## 6. Wiring

- [x] 6.1 Update `cmd/asylum/main.go`: pass disabled kit info through to resolve; update `--kits` CLI flag to support `--kits -shell` or similar syntax for disabling
- [x] 6.2 Remove `defaultCacheDirs` fallback from `internal/container/container.go` (dead code)
- [x] 6.3 Remove unused `/home/claude/workspace` mkdir from Dockerfile.tail
- [x] 6.4 Add CHANGELOG entry under Unreleased

## 7. Testing

- [x] 7.1 Unit tests for new kits: github, openspec, shell registered with correct fields
- [x] 7.2 Unit tests for dependency validation: warning on missing dep, no warning when satisfied
- [x] 7.3 Unit tests for default-on: included when not listed, excluded when disabled, excluded with empty map
- [x] 7.4 Unit tests for kit disabling at project level overriding global config
- [x] 7.5 Verify all existing tests pass after Dockerfile cleanup
