## Context

Asylum's config was built incrementally: flat fields for agent, ports, volumes, env, then profiles/agents as lists, features/onboarding as bool maps, packages as typed maps, versions as a string map. The kit system (née profiles) now provides a natural grouping — per-kit options replace scattered top-level fields.

The target format is defined by the `config.yaml` template in the repo root. Key structural changes: kits is a map (presence = activated), agents is a map, and per-kit fields (versions, packages, shadow-node-modules, onboarding) replace top-level equivalents.

## Goals / Non-Goals

**Goals:**
- Single config restructure: kits map as the primary organizing concept
- Rename profile → kit everywhere (code, config, CLI, docs)
- Well-commented default config written on first run
- Automatic migration of existing v1 configs (global + project)
- Preserve user values during migration
- Use `yaml.Node` API for comment-preserving config generation

**Non-Goals:**
- Changing kit behavior (installation logic, snippets, resolution — those stay the same)
- Adding new kits beyond renaming existing profiles + new groupings (apt, title, shell)
- Changing how agents work at runtime (command generation, session detection)
- Supporting rolling back from v2 to v1 format

## Decisions

### 1. New Config struct

```go
type Config struct {
    Version        string                    `yaml:"version,omitempty"`
    ReleaseChannel string                    `yaml:"release-channel,omitempty"`
    Agent          string                    `yaml:"agent,omitempty"`
    Agents         map[string]*AgentConfig   `yaml:"agents,omitempty"`
    Kits           map[string]*KitConfig     `yaml:"kits,omitempty"`
    Ports          []string                  `yaml:"ports,omitempty"`
    Volumes        []string                  `yaml:"volumes,omitempty"`
    Env            map[string]string         `yaml:"env,omitempty"`
}

type KitConfig struct {
    Versions              []string `yaml:"versions,omitempty"`
    DefaultVersion        string   `yaml:"default-version,omitempty"`
    Packages              []string `yaml:"packages,omitempty"`
    ShadowNodeModules     *bool    `yaml:"shadow-node-modules,omitempty"`
    Onboarding            *bool    `yaml:"onboarding,omitempty"`
    TabTitle              string   `yaml:"tab-title,omitempty"`
    AllowAgentTermTitle   *bool    `yaml:"allow-agent-terminal-title,omitempty"`
    Build                 []string `yaml:"build,omitempty"`
    Start                 []string `yaml:"start,omitempty"`
}

type AgentConfig struct {
    // Empty for now — placeholder for future per-agent config
}
```

Kit activation is by presence: if `kits.java` exists in the map (even with empty/nil value), Java is active. No kits key at all means default kits (java, python, node — backwards compatible). Empty `kits: {}` means no kits.

**Alternative considered**: Separate `KitConfig` types per kit (JavaKitConfig, NodeKitConfig). Rejected — a single struct with optional fields is simpler and covers all current kits. Fields that don't apply to a kit are simply not set.

### 2. Agents as a map

`agents` changes from `*[]string` to `map[string]*AgentConfig`. Keys are agent names, values are (currently empty) per-agent config. Presence = installed.

- `nil` map (not specified): defaults to `{"claude": nil}` (same as current nil-means-claude)
- Empty map `agents: {}`: no agents
- `agents: {claude:, gemini:}`: both installed

### 3. Kit activation semantics

Kits map `nil` (key absent from config) means default kits — same backwards-compatible behavior as old `profiles: nil`. Explicitly present `kits: {}` (empty map) means no kits.

For language kits (java, python, node), activation triggers the same profile resolution as before — the kit system internally maps kit names to the existing profile/kit definitions.

Non-language kits (apt, title, shell) are configuration-only groupings — they don't have Dockerfile snippets but their options affect the build or runtime.

### 4. Migration strategy

Two migration paths:

**Global config** (`~/.asylum/config.yaml`): Detected by missing or old `version` field. Read as `yaml.Node`, transform the tree, write back. The old flat fields are reorganized into the kits map.

**Project configs** (`.asylum`, `.asylum.local`): No version field. Detected by presence of `features` key (which doesn't exist in v2). Same node-level transformation.

Migration mapping:
- `profiles: [java, node]` → `kits: {java: {}, node: {}}`
- `versions: {java: "21"}` → `kits: {java: {default-version: "21"}}`
- `packages: {apt: [...]}` → `kits: {apt: {packages: [...]}}`
- `packages: {npm: [...]}` → `kits: {node: {packages: [...]}}`
- `packages: {pip: [...]}` → `kits: {python: {packages: [...]}}`
- `packages: {run: [...]}` → `kits: {shell: {build: [...]}}`
- `features: {shadow-node-modules: true}` → `kits: {node: {shadow-node-modules: true}}`
- `features: {onboarding: false}` → removed (onboarding is now per-kit)
- `features: {allow-agent-terminal-title: false}` → `kits: {title: {allow-agent-terminal-title: false}}`
- `onboarding: {npm: false}` → `kits: {node: {onboarding: false}}`
- `tab-title: "..."` → `kits: {title: {tab-title: "..."}}`
- `agents: [claude, gemini]` → `agents: {claude: {}, gemini: {}}`

After migration, old keys are removed from the YAML node tree.

### 5. Default config generation

On first run (no `~/.asylum/config.yaml` exists), write a default config using a Go template string that includes comments. This is NOT parsed from a struct — it's a literal YAML string with inline comments, written directly to disk.

### 6. Profile → Kit rename

The `internal/profile` package is renamed to `internal/kit`. All types: `Profile` → `Kit`, `SubProfiles` → `SubKits`. The `profile.Resolve` function becomes `kit.Resolve`. CLI flag `--profiles` → `--kits`. Config field `profiles` → `kits`.

The old `profiles` field is still recognized during migration but not used in v2 configs.

## Risks / Trade-offs

**Breaking change for existing configs** → Mitigated by automatic migration. Users don't need to do anything manually.

**Migration could corrupt configs** → Mitigated by writing a `.asylum.backup` before migrating. If migration fails, the original is preserved.

**Kit config struct is generic (not per-kit typed)** → Acceptable trade-off. A `shadow-node-modules` field on a Java kit config is harmless (ignored). Per-kit types would add significant complexity for minimal benefit.

**Non-language kits (apt, title, shell) blur the kit concept** → They're convenient groupings that keep the config organized. Users won't be confused because the config comments explain each section.
