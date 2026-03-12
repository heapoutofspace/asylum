## Why

The placeholder Dockerfile and entrypoint.sh need to be replaced with the real container assets. These are adapted from the existing AgentBox files, updated for Asylum (env var names, labels, all three agent CLIs, welcome banner).

## What Changes

- Replace `assets/Dockerfile` with full Asylum Dockerfile based on AgentBox, adding Gemini CLI and Codex installs
- Replace `assets/entrypoint.sh` with Asylum-adapted entrypoint (ASYLUM_ env vars, updated welcome banner)
- Rename all AGENTBOX_ references to ASYLUM_

## Capabilities

### New Capabilities
- `container-image`: Full Dockerfile and entrypoint.sh for the Asylum container with all three agent CLIs

### Modified Capabilities

None.

## Impact

- Replaces placeholder files in `assets/`
- After this change, `make build` produces a binary that can build a real Docker image
