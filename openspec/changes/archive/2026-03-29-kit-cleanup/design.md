## Context

The kit system currently handles language toolchains (java, python, node) and Docker. But several tools remain hardcoded in the core Dockerfile and tail: GitHub CLI, Maven (in the core apt block), Python build deps, OpenSpec (bundled into node's npm packages), and the entire shell configuration (oh-my-zsh, tmux, direnv). Moving these into kits requires two new features: kit dependencies and default-on kits.

## Goals / Non-Goals

**Goals:**
- Extract remaining hardcoded tools into kits
- Kit dependencies: OpenSpec depends on node (needs npm), validated at resolve time
- Default-on kits: shell kit active unless explicitly disabled, even when user lists specific kits
- Kit disabling: `kit-name: false` in config to exclude a kit at any layer

**Non-Goals:**
- Transitive dependency resolution (if A depends on B which depends on C — not needed yet)
- Auto-activating missing dependencies (warn only, same as agent deps)
- Making core system packages (git, curl, zsh) into kits

## Decisions

### 1. Kit dependencies via `Deps` field

```go
type Kit struct {
    // ...existing fields...
    Deps      []string // kit names this kit depends on
    DefaultOn bool     // active unless explicitly disabled
}
```

During `Resolve`, after building the active kit list, validate that each kit's `Deps` are satisfied. Missing deps emit a warning via `log.Warn` but don't block — the Docker build will fail with a clear error.

**Alternative considered**: Auto-activate missing deps. Rejected because it creates surprising behavior — the user explicitly chose their kits.

### 2. Default-on kits via `DefaultOn` field

A kit with `DefaultOn: true` is included in the resolved set unless the user explicitly disables it. The behavior:

- **Config nil** (no kits key): all kits active (including default-on) — unchanged behavior
- **Config has explicit kits** (e.g., `kits: {java: {}}`): listed kits + default-on kits
- **Kit explicitly disabled** (e.g., `kits: {shell: false}`): default-on kit excluded

This means `shell` is always present unless you write `shell: false`. It doesn't change behavior for users who list `kits: {}` (empty = none, still means none — default-on kits are NOT added to an explicit empty map).

### 3. Kit disabling via `false` or `Disabled` field

In YAML, `shell: false` parses as a boolean, not a map. To support both `kit-name: {}` (enabled with options) and `kit-name: false` (disabled), `KitConfig` needs a `Disabled` field or we use a custom unmarshaler.

Simplest approach: add `Disabled *bool` to KitConfig. In YAML:
```yaml
kits:
  java:                    # active, no options
  shell:
    disabled: true         # explicitly disabled
  python:
    packages: [ansible]    # active with options
```

`KitActive` returns false when `Disabled` is true. `Resolve` skips disabled kits.

**Alternative considered**: Custom YAML unmarshaler for `false` literal. Rejected — adds complexity and `disabled: true` is more explicit and grep-friendly.

### 4. New kits

**github** — installs `gh` CLI via the official apt repository:
```go
Kit{
    Name: "github",
    DockerSnippet: `# Install GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/...`,
    DefaultOn: true,
}
```

**openspec** — installs OpenSpec CLI via npm:
```go
Kit{
    Name: "openspec",
    DockerSnippet: `# Install OpenSpec CLI
RUN bash -c '... npm install -g @fission-ai/openspec@latest'`,
    Deps: []string{"node"},
    DefaultOn: true,
}
```

**shell** — oh-my-zsh, tmux config, direnv hooks, terminal size handling:
```go
Kit{
    Name: "shell",
    DockerSnippet: `# Install oh-my-zsh and configure shell...`,
    DefaultOn: true,
}
```

The shell kit's DockerSnippet contains everything currently in `Dockerfile.tail` lines 1-25 (oh-my-zsh, theme, PATH re-add, direnv hooks, terminal size, tmux). After the shell kit snippet, the tail only has: git config, workspace dir, entrypoint COPY, USER/WORKDIR, ENV.

### 5. Maven moves to java/maven sub-kit

The `maven` apt package is currently in the core apt-get block. It moves to the java/maven sub-kit's DockerSnippet as a root-user apt install:

```go
Kit{
    Name: "java/maven",
    DockerSnippet: `USER root
RUN apt-get update && apt-get install -y --no-install-recommends maven && rm -rf /var/lib/apt/lists/*
USER claude`,
}
```

### 6. Python build deps move to python kit

`python3-dev`, `python3-pip`, `python3-venv`, `libssl-dev`, `libffi-dev` move from the core apt-get into the python kit's DockerSnippet as a root-user apt install, prepended before the uv tool installs.

### 7. Snippet insertion order

Kit DockerSnippets run as USER claude (set by core). Kits that need root (maven apt, python apt) must switch to USER root and back. Order:

```
Dockerfile.core (ends with USER claude)
→ kit snippets (java, java/maven with USER root/claude, python with USER root/claude, node, ...)
→ agent snippets
→ Dockerfile.tail (starts with shell kit content if active, then git config, entrypoint COPY)
```

Wait — shell kit content needs to go into the Dockerfile snippet, not the tail. The tail shrinks to just: git config, workspace dir, root switch, entrypoint COPY, WORKDIR, USER, ENV.

## Risks / Trade-offs

**Shell kit disabled breaks interactive experience** → Acceptable. If you disable shell, you get bare zsh with no oh-my-zsh, no tmux config, no direnv hooks. The container still works, just less polished.

**Kit dependency warnings may be noisy** → Only warns when a dep is missing, which is a user configuration error. The warning helps them fix it.

**Apt installs in kit snippets less efficient than single apt-get** → Each kit that needs apt runs its own `apt-get update && install`. This adds a few seconds per kit but keeps kits self-contained. Could optimize with a shared apt cache mount.
