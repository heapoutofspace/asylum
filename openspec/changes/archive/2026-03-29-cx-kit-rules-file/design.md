## Context

The cx kit currently defines a `RulesSnippet` as a Go string literal in `internal/kit/cx.go`. This snippet is assembled into the monolithic `asylum-sandbox.md` rules file by `generateSandboxRules` in `internal/container/container.go`. The cx tool itself has a `cx skill` command that outputs a markdown description of its capabilities — this is the authoritative, version-matched source of truth for what cx can do.

Currently, all kit rules are embedded in the sandbox rules file which is generated on the host side and mounted as a single file. The cx kit's rules would instead be a standalone file generated during the Docker build (where cx is installed) and placed into the agent's rules directory at container startup.

## Goals / Non-Goals

**Goals:**
- Use `cx skill` output as the cx rules file so it stays in sync with the installed version
- Generate the rules file during Docker build (where cx is available) and store it in the image
- Place the rules file into the agent's rules directory at container startup via the entrypoint
- Remove the hardcoded `RulesSnippet` from the cx kit Go code

**Non-Goals:**
- Generalizing this pattern to all kits (other kits don't have a `skill` command)
- Changing how the main `asylum-sandbox.md` rules file is generated or mounted
- Modifying the `cx skill` command itself

## Decisions

### 1. Generate during Docker build, not at container startup

Run `cx skill > /tmp/asylum-kit-rules/cx.md` as part of the cx kit's `DockerSnippet`. This bakes the output into the image layer, so it's cached and doesn't add latency to every container start. The entrypoint rule against installing things doesn't apply — this is just capturing output from an already-installed tool.

**Alternative considered**: Running `cx skill` in the entrypoint. Rejected because it adds startup latency on every container start and violates the spirit of "entrypoint configures, Dockerfile installs."

### 2. Use `/tmp/asylum-kit-rules/` as the staging directory

A well-known directory in the image where kits can place pre-generated rules files. The cx kit creates `/tmp/asylum-kit-rules/cx.md` during build. This is a convention that other kits could adopt in the future if they gain similar self-description commands, but we're not designing for that now.

### 3. Bind mount in entrypoint, not copy

The cx kit's `EntrypointSnippet` uses `mount --bind` to overlay the pre-generated file onto `~/.claude/rules/cx.md`. This is critical because `~/.claude/rules/` is shared across all projects via the agent config mount — copying or writing there would leak cx rules into projects that don't have the cx kit active. A bind mount only exists for the container's lifetime and never modifies the underlying shared directory on the host.

### 4. Remove RulesSnippet, keep Tools

The `RulesSnippet` field is cleared since the standalone rules file replaces it. The `Tools` field (`[]string{"cx"}`) is kept so cx still appears in the "Kit Tools" aggregated list in `asylum-sandbox.md`.

## Risks / Trade-offs

- **[Risk] `cx skill` output format changes** → Low risk. The output is consumed as a rules file by the agent, not parsed programmatically. Any valid markdown works.
- **[Risk] `cx skill` fails during Docker build** → The `DockerSnippet` should handle this gracefully (e.g., `cx skill > /tmp/... || true`) so the image build doesn't fail if cx has an issue. The container just won't have the cx rules file.
- **[Trade-off] Standalone file vs. assembled snippet** → The cx rules become a separate file in `~/.claude/rules/` rather than part of `asylum-sandbox.md`. This means the agent loads it as a separate rules file, which is fine — Claude Code loads all files in `~/.claude/rules/`.
- **[Risk] Bind mount requires privileges** → The container needs `CAP_SYS_ADMIN` or privileged mode for `mount --bind`. Containers with the Docker kit already run privileged. For non-privileged containers, `sudo mount --bind` should work since the user has passwordless sudo.
