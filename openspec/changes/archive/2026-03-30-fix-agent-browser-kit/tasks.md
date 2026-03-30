# Tasks: fix-agent-browser-kit

Replace the browser kit's Playwright-based implementation with agent-browser. Rename the kit from "browser" to "agent-browser", install agent-browser instead of Playwright, generate the Claude Code skill at build time and mount it at runtime, and update the rules snippet to document the agent-browser workflow.

## Tasks

- [x] Rename and rewrite `internal/kit/browser.go` to `internal/kit/agent_browser.go` with agent-browser implementation (new name, description, DockerSnippet, EntrypointSnippet, RulesSnippet, Tools, BannerLines; drop CacheDirs and Playwright)
- [x] Add kit alias handling in `kit.Resolve` so `"browser"` in existing configs silently resolves to `"agent-browser"`
- [x] Update `docs/kits/browser.md` for agent-browser
- [x] Update `openspec/specs/browser-kit/spec.md` for agent-browser
- [x] Update any tests that reference the old "browser" kit name
- [x] Add CHANGELOG entry under Unreleased
