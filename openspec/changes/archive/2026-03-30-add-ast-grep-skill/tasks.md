# Tasks: add-ast-grep-skill

Add the upstream ast-grep Claude Code skill (`ast-grep/agent-skill`) to the ast-grep kit, following the same pattern as the agent-browser kit: generate the skill at Docker build time via `npx skills add`, stage it at a known path, and mount it into `~/.claude/skills/` at runtime via the entrypoint.

## Tasks

- [x] Update `internal/kit/astgrep.go`: add DockerSnippet step to run `npx skills add ast-grep/agent-skill --skill ast-grep --yes --copy` and stage the result; add EntrypointSnippet to mount the skill directory; set `NeedsMount: true`
- [x] Update `docs/kits/ast-grep.md` to mention the auto-mounted Claude Code skill
- [x] Verify build compiles and tests pass
- [x] Add CHANGELOG entry under Unreleased
