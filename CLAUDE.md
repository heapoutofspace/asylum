# Asylum

Agent-agnostic Docker sandbox for AI coding agents (Claude Code, Gemini CLI, Codex). Single Go binary, cross-compiled for ARM and x86. See `PLAN.md` for the full specification.

## Change Management

This project uses [OpenSpec](https://openspec.dev) for structured change management. Use the `/opsx:propose` skill to start a new change, `/opsx:apply` to implement, and `/opsx:archive` to archive completed changes. See `openspec/` for specs and change history.

## Architecture

- **Go** (latest stable) — single binary, no runtime dependencies beyond Docker
- Cross-compiled for `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`
- Shells out to Docker CLI via `os/exec` and `syscall.Exec`
- Layered YAML config: `~/.asylum/config.yaml` → `$project/.asylum` → `$project/.asylum.local` → CLI flags
- Embedded assets (Dockerfile, entrypoint.sh) via `go:embed`
- One external dependency: `gopkg.in/yaml.v3`

## Code Style

### General

- **Less code is better.** Every line must earn its place. Avoid defensive boilerplate, speculative abstractions, and "just in case" code paths.
- Use modern Go: generics where they reduce duplication, errors as values, `slices`/`maps` packages.
- No unnecessary interfaces — don't create an interface until there are two implementations. A concrete type is fine.
- Keep functions short: one concern per function, early returns for error cases.
- Use `if err != nil { return err }` — don't wrap errors unless the wrapper adds information the caller doesn't already have.
- **Do not add fields, config options, or functionality without consulting the user.** If something seems needed but isn't explicitly requested, ask first.

### Comments

Code comments are used sparingly. Comprehensible and expressive code (consistent, logical naming) is preferred.

Comments are added when they contribute to much faster, better understanding in two cases:
- To explain **why** something was done, when it is not apparent from the context.
- To explain **what** is being done, if the code is necessarily difficult to understand.

If a log line explains what is happening, any comment above that line which essentially says the same thing is redundant and should not be added.

### Naming

- Package names: short, lowercase, no underscores. Avoid stutter (`config.Config` is fine, `config.ConfigConfig` is not).
- Functions/methods: verb-noun (`buildImage`, `loadConfig`). Getters drop the `Get` prefix (`Name()`, not `GetName()`).
- Variables: short-lived vars can be short (`f`, `err`, `cmd`). Longer-lived vars get descriptive names.
- Constants: `CamelCase`, not `SCREAMING_SNAKE`.

### Error Handling

- Return `error` from functions that can fail. Don't panic except for programmer errors.
- Wrap errors with `fmt.Errorf("context: %w", err)` only when the wrapper adds value.
- Log errors at the point of handling, not at the point of returning.
- Use the project's `log` package for user-facing output, not `fmt.Println` or the standard `log` package.

### Testing

- Use Go's built-in `testing` package. No test frameworks.
- Table-driven tests for functions with multiple input/output cases.
- Test files live next to the code they test (`config_test.go` next to `config.go`).
- Test the important logic: config merging, volume shorthand parsing, session detection, command generation, hash computation. Don't test trivial getters.
- Use `testdata/` directories for fixture files.

### Project Structure

Follow the layout defined in PLAN.md section 8. `cmd/asylum/` for the entry point, `internal/` for all packages, `assets/` for embedded files.

## Dependencies

Use libraries freely when they save meaningful effort. Prefer well-maintained, focused libraries over rolling your own. Some natural fits:

- `gopkg.in/yaml.v3` — YAML parsing
- `github.com/fatih/color` or similar — colored terminal output (if simpler than hand-rolling ANSI)
- `github.com/spf13/cobra` or `github.com/urfave/cli/v2` — CLI framework (if `flag` feels limiting)

Avoid large dependency trees that pull in the world (e.g., Docker SDK when shelling out to `docker` CLI works fine).

## What NOT to Do

- Do not add Docker SDK. Shell out to the `docker` CLI — it's simpler and avoids a huge dependency tree.
- Do not create unnecessary abstractions, utility packages, or helper functions for one-off operations.
- Do not add config options, features, or agent support beyond what PLAN.md specifies.
- Do not attempt to fix git corruption (broken packfiles, bad objects, etc.) yourself. Always prompt the user to resolve it.
