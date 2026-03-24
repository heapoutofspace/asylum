## Why

When a user runs `asylum` for the very first time (`~/.asylum` doesn't exist yet), there's no opportunity to configure global settings before the first session launches. Package manager credentials (Maven `settings.xml`, Docker `config.json`) are common mounts that most Java/Docker users would want, but currently require manually creating `~/.asylum/config.yaml` before the first run. A brief interactive onboarding makes this zero-friction.

## What Changes

- Detect first-run condition (no `~/.asylum` directory) in the main CLI flow, before config loading.
- Prompt the user: "Make package manager credentials available in the sandbox? (Maven settings.xml, Docker config.json)" with a Y/n prompt.
- If yes, generate `~/.asylum/config.yaml` with volume mounts for `~/.m2/settings.xml` and `~/.docker/config.json` (only for files that actually exist on the host).
- If no (or no relevant files exist), create `~/.asylum` directory but skip config generation, so subsequent runs don't re-prompt.

## Capabilities

### New Capabilities
- `first-run-onboarding`: Interactive first-run detection and config generation for package manager credentials

### Modified Capabilities

## Impact

- `cmd/asylum/main.go`: New first-run check before `config.Load`, calling into a new onboarding function.
- New file in `internal/` (e.g., `internal/firstrun/`) or extending the existing onboarding package.
- No new dependencies. Uses `fmt.Scanln` for prompts (same pattern as `--cleanup`), `os` for file detection, `gopkg.in/yaml.v3` for YAML generation.
- No changes to existing config loading or merging logic — this runs before config load and produces a standard config file.
