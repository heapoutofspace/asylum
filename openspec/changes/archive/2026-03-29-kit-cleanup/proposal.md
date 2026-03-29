## Why

Several tools remain hardcoded in `Dockerfile.core` that should be kits: GitHub CLI, Maven, OpenSpec, and Python build dependencies. The shell configuration (oh-my-zsh, tmux, direnv hooks) is also hardcoded in the tail. Moving these into kits makes them optional and configurable. This also requires two new kit system features: kit dependencies (OpenSpec depends on node) and default-on kits (shell should be active even with no config, but explicitly disableable).

## What Changes

- **New kits**: `github` (gh CLI), `openspec` (OpenSpec CLI, depends on node), `shell` (oh-my-zsh, tmux, direnv hooks, terminal size handling)
- **Move maven** from core apt-get into `java/maven` sub-kit's DockerSnippet
- **Move python build deps** (`python3-dev`, `python3-pip`, `python3-venv`, `libssl-dev`, `libffi-dev`) from core apt-get into `python` kit's DockerSnippet
- **Move shell config** (oh-my-zsh install, theme, direnv hooks, terminal size handling, tmux config) from `Dockerfile.tail` into `shell` kit
- **Remove `@fission-ai/openspec@latest`** from node kit's npm packages — it becomes its own kit
- **Kit dependencies**: Kit struct gains a `Deps []string` field. Resolution validates deps are active; missing deps emit a warning (same as agent deps)
- **Default-on kits**: Kit struct gains a `DefaultOn bool` field. Default-on kits are included when config is nil (no kits key) AND when kits are explicitly listed but don't mention the default-on kit. They can be explicitly disabled with `kit-name: false` in config.
- **Kit disabling**: `KitConfig` supports a `false` YAML value (or `disabled: true` field) to explicitly exclude a kit, useful for disabling globally-configured kits at project level

## Capabilities

### New Capabilities
- `kit-dependencies`: Kits can declare dependencies on other kits, validated during resolution
- `kit-defaults`: Default-on kits that are active unless explicitly disabled
- `github-kit`: GitHub CLI as optional kit
- `openspec-kit`: OpenSpec CLI as optional kit with node dependency
- `shell-kit`: Shell configuration (oh-my-zsh, tmux, direnv) as default-on kit

### Modified Capabilities
- `profile-system`: Kit struct gains Deps and DefaultOn fields, resolution handles dependencies and defaults
- `profile-image-build`: Core Dockerfile shrinks (maven, python deps, shell config removed); tail shrinks (shell config moves to kit)
- `profile-config-integration`: Kit disabling via `kit-name: false` or `disabled: true` in KitConfig

## Impact

- **internal/kit/kit.go**: Add `Deps`, `DefaultOn` fields; update `Resolve` for dependency validation and default-on logic
- **internal/kit/github.go** (new): GitHub CLI kit
- **internal/kit/openspec.go** (new): OpenSpec CLI kit with node dependency
- **internal/kit/shell.go** (new): Shell config kit (oh-my-zsh, tmux, direnv hooks, terminal size)
- **internal/kit/java.go**: Move maven apt package into java/maven sub-kit DockerSnippet
- **internal/kit/python.go**: Add python build deps to python kit DockerSnippet
- **internal/kit/node.go**: Remove `@fission-ai/openspec@latest` from npm packages
- **internal/config/config.go**: KitConfig gains `Disabled` field; `KitActive` respects it
- **internal/config/defaults.go**: Update default config with new kits
- **assets/Dockerfile.core**: Remove maven, python build deps from apt-get block
- **assets/Dockerfile.tail**: Remove oh-my-zsh, direnv hooks, terminal size, tmux config (move to shell kit)
