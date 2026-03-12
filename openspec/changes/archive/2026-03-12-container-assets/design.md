## Context

The AgentBox Dockerfile and entrypoint.sh at `/agentbox/` are the starting point. They need adaptation for Asylum: installing all three agent CLIs, renaming env vars from AGENTBOX_ to ASYLUM_, and updating the welcome banner.

## Goals / Non-Goals

**Goals:**
- Full Dockerfile per PLAN.md section 7: all toolchains, all three agent CLIs, proper labels
- Entrypoint per PLAN.md section 6: Docker-in-Docker, NVM, SDKMAN, Python venv, SSH, direnv, git, welcome banner
- Rename AGENTBOX_ env vars to ASYLUM_
- Remove openspec npm package (not needed for Asylum)

**Non-Goals:**
- No functional changes to the build toolchain — same base image, same tool versions

## Decisions

- **Keep OpenSpec removal**: The AgentBox Dockerfile installs `@fission-ai/openspec@latest` — this is not needed in Asylum and is removed.
- **Add Gemini + Codex**: Two additional `npm install -g` commands for the agent CLIs.
- **ASYLUM_ prefix**: All environment variables use ASYLUM_ prefix instead of AGENTBOX_.

## Risks / Trade-offs

- Large Dockerfile — this is inherent to a multi-language dev environment. Docker layer caching mitigates rebuild time.
