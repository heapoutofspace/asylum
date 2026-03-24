## Context

Asylum stores global config at `~/.asylum/config.yaml` and agent data under `~/.asylum/agents/`. The `~/.asylum` directory is created implicitly when agent config is seeded (`container.EnsureAgentConfig`), but no global config file is generated. Users who want Maven or Docker credentials in the sandbox must manually create the config file with volume entries.

The existing project onboarding framework (`internal/onboarding/`) handles per-project setup tasks (e.g., `npm install`) after container start. First-run onboarding is different: it runs once globally, before any container or config loading, and produces a config file rather than executing container commands.

## Goals / Non-Goals

**Goals:**
- Detect first-run (no `~/.asylum` directory) and prompt for package manager credential mounts
- Generate `~/.asylum/config.yaml` with appropriate volume entries
- Only offer to mount files that actually exist on the host
- Ensure the prompt never appears again after first run

**Non-Goals:**
- General-purpose first-run wizard (other settings, agent selection, etc.) — keep it minimal for now
- Modifying existing config loading or merging logic
- Supporting non-interactive mode (CI) — asylum is inherently interactive

## Decisions

**New `internal/firstrun` package** rather than extending `internal/onboarding`.
The project onboarding framework is designed around detecting project-level workloads, running commands in containers, and tracking state per container. First-run is a one-shot global concern that produces a config file before the container even exists. Mixing these into one package would conflate two different lifecycles.

**Check for `~/.asylum/agents/` directory existence** as the first-run signal.
The `agents/` directory is created by `EnsureAgentConfig` (which runs later in the flow), not by the installer. This reliably distinguishes fresh installs from existing users. No explicit marker file or directory creation is needed — `EnsureAgentConfig` handles it as part of the normal flow.

**Only mount files that exist on the host.** If `~/.m2/settings.xml` doesn't exist, don't include it in the config — it would cause Docker bind-mount errors. Show the user which files were found and will be mounted.

**Mount as read-only.** Credentials should not be writable from inside the sandbox. Use `:ro` option.

**Use `~` prefix in volume paths.** The config file stores `~/...` paths which `ParseVolume` + `ExpandTilde` already handle. This keeps the config portable and readable.

## Risks / Trade-offs

**User skips prompt but later wants credentials** → They can manually add volume entries to `~/.asylum/config.yaml`. This is documented and straightforward.

**New credential files appear after first run** → Not auto-detected. User adds them manually. Acceptable for v1; we could add an `asylum setup` command later.
