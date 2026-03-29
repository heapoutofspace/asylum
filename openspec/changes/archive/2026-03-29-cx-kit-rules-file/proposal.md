## Why

The cx kit currently maintains a hand-written `RulesSnippet` in Go code that describes cx capabilities for agents. The cx tool itself provides a `cx skill` command that outputs a comprehensive, version-accurate description of its capabilities. By using `cx skill` output as a standalone rules file, the rules stay in sync with the installed cx version automatically and don't need manual maintenance in the Asylum codebase.

## What Changes

- During the Docker base image build, run `cx skill` and write its output to a temp location in the image (e.g., `/tmp/asylum-kit-rules/cx.md`)
- Add an `EntrypointSnippet` to the cx kit that copies the pre-generated rules file into the agent's rules directory (e.g., `~/.claude/rules/cx.md`) at container startup
- Remove the hardcoded `RulesSnippet` from the cx kit registration since the rules are now self-generated

## Capabilities

### New Capabilities
- `cx-kit-rules-gen`: cx kit generates its own rules file from `cx skill` during Docker build and mounts it into the agent's rules directory via the entrypoint

### Modified Capabilities
- `cx-kit`: Remove `RulesSnippet` field from cx kit, add `DockerSnippet` step to generate rules file, add `EntrypointSnippet` to place it at startup
- `sandbox-rules`: The cx kit's rules snippet will no longer appear in the assembled `asylum-sandbox.md` — it becomes a standalone rules file instead

## Impact

- `internal/kit/cx.go` — remove `RulesSnippet`, extend `DockerSnippet`, add `EntrypointSnippet`
- `assets/Dockerfile` — cx install step now also runs `cx skill` to generate the rules file
- `assets/entrypoint.sh` — cx entrypoint snippet copies rules file into place
- Existing tests for sandbox rules assembly may need updating since cx no longer contributes a `RulesSnippet`
