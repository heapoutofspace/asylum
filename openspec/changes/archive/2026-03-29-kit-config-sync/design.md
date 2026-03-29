## Context

Today, the default config is written once on first run via `WriteDefaults` (O_EXCL — no-op if file exists). When new kits are added in a release, existing users never see them. The migration system only fires for version-gated schema changes, not for "new kit available." There is no mechanism to introduce kits to existing installations.

The kit system has an informal two-tier model: `DefaultOn: true` (always active) and everything else (active if in config). But there's actually a third tier needed: kits that should be in the config by default but can be removed (like `docker`, `java`), vs kits that are active even without config (like `shell`). The `DefaultOn` bool conflates these.

Each kit currently provides a `ConfigSnippet` string for text-based config assembly. This works for generating fresh configs but is fragile for modifying existing configs that the user may have hand-edited.

## Goals / Non-Goals

**Goals:**
- Detect newly added kits on startup and insert them into existing configs
- Preserve user edits (comments, ordering, custom values) when modifying config
- Prompt users before activating new default-on kits in their sandbox
- Formalize three kit activation tiers with clear semantics
- Make adding a new kit to the registry automatically surface it to all users

**Non-Goals:**
- Removing kits from config when they're unregistered (deleted kits just become inert keys)
- Per-project kit state tracking (kit awareness is global — image-level concern)
- Migrating existing `ConfigSnippet` text to `yaml.Node` in this change (both can coexist during transition)

## Decisions

### 1. State file at `~/.asylum/state.json`

A separate machine-managed JSON file tracks which kits the installation has seen:

```json
{
  "known_kits": ["apt", "docker", "github", "java", "node", "openspec", "python", "shell", "title"]
}
```

**Why not config.yaml?** The config is the user's domain — adding machine metadata pollutes it and risks merge conflicts. JSON (not YAML) reinforces "don't hand-edit this."

**Why not version-gated?** Feature branches adding kits would need coordinated version bumps. With explicit tracking, a branch just registers a kit; when it runs, the new kit is detected and handled. No ordering conflicts across branches.

**Why not infer from config keys?** Commented-out kits aren't YAML keys — they're invisible to the parser. We'd need text scanning, which is fragile. Explicit state is unambiguous.

### 2. Three-tier activation enum

Replace `DefaultOn bool` with a `Tier` field:

| Tier | Behavior | Config presence | Prompt |
|------|----------|-----------------|--------|
| `TierAlwaysOn` | Active even if absent from config | Not required | Info only: "Kit X is now available" |
| `TierDefault` | Added uncommented to config | Required for activation | "Activate kit X? [Y/n]" (interactive) or added commented (non-interactive) |
| `TierOptIn` | Added commented to config | Required for activation | Info only: "Kit X available — uncomment in config to enable" |

Current mapping: `shell`, `node`, `title` → `TierAlwaysOn`. `docker`, `java`, `python`, `github`, `openspec` → `TierDefault`. `apt` → `TierOptIn`.

### 3. Config modification via `yaml.Node`

Use `yaml.Unmarshal` into a `yaml.Node` (not a struct) to parse existing config. This preserves comments, key ordering, and formatting. Walk the node tree to find the `kits` mapping, then append new kit entries as child nodes.

For active kits: insert real `yaml.Node` key-value pairs (the kit provides these via a new `ConfigNodes() []*yaml.Node` method).

For opt-in/commented kits: insert as `FootComment` on the last existing kit node, since commented-out YAML isn't representable as real nodes.

**Why not text splicing?** Users reformat, reorder, add blank lines. Text anchors like "# Port forwarding" are fragile. `yaml.Node` is the parser's own representation — it handles any valid YAML.

### 4. Kit sync runs before config load

The sync check (compare registry vs state, prompt, modify config) runs early in `main()`, after `state.json` is loaded but before `config.Load`. This ensures the config is up to date before it's parsed into the typed `Config` struct.

Sequence:
1. Load `state.json` (or create empty)
2. Compare registered kit names vs `known_kits`
3. If new kits: run sync (prompt if interactive, modify config, update state)
4. Proceed with `config.Load` as normal

### 5. Non-interactive fallback

When stdin is not a terminal (CI, piped input, scripts), new default-on kits are added as commented-out entries instead of uncommented. This avoids silently changing the sandbox. The user sees them on next interactive run or when they open their config.

## Risks / Trade-offs

- **`yaml.Node` complexity** — Node tree manipulation is more code than text splicing. But it's correct for all edge cases (reordered keys, inline comments, flow style). Mitigated by keeping the node-walking logic in a single well-tested function.
- **`state.json` gets out of sync** — If the user deletes it, the next run re-detects all kits as "new" and prompts again. This is acceptable — it's a one-time re-prompt, not destructive. The state file can document this: "Delete to re-run kit setup."
- **Comment insertion ordering** — Commented-out kits appended as `FootComment` may not look perfectly formatted. Mitigated by testing the output carefully and adjusting whitespace.
- **Always-on kits with no config presence** — Users can't see or disable them from the config. This is intentional for `shell`/`node`/`title` but could confuse users who want to understand what's active. Mitigated by the sandbox rules file which lists all active kits.
