## Why

OpenSpec is used for structured change management in projects that use Asylum containers. The original AgentBox Dockerfile included it, but it was removed during the Asylum adaptation. It should be available inside the container so agents can use the `/opsx:*` workflow.

## What Changes

- Add `@fission-ai/openspec@latest` to the global npm packages installed in the Asylum Dockerfile

## Capabilities

### New Capabilities
- `openspec-in-container`: OpenSpec CLI available inside Asylum containers

### Modified Capabilities

None.

## Impact

- Modifies `assets/Dockerfile` — adds one package to the existing npm install step
- Triggers a base image rebuild on next run (asset hash changes)
