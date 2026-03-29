## 1. Flag Parsing

- [x] 1.1 Add `All bool` field to `cliFlags` struct
- [x] 1.2 Update `cleanup` subcommand parsing in `parseArgs` to accept `--all` flag (currently rejects any args after cleanup)
- [x] 1.3 Update `--cleanup` flag alias path to also check for `--all`
- [x] 1.4 Pass `All` flag through to `runCleanup`

## 2. Scoped Cleanup (Default)

- [x] 2.1 In `runCleanup`, when `--all` is not set: resolve `projectDir` from cwd, compute container name via `container.ContainerName(projectDir)`
- [x] 2.2 Remove container: `docker rm -f <container-name>` (ignore error if not exists)
- [x] 2.3 List and remove volumes prefixed with `<container-name>-` (shadow + cache volumes)
- [x] 2.4 Remove project data dir: `~/.asylum/projects/<container-name>/` (with port release, respecting active sessions)
- [x] 2.5 If not in a valid project dir (or cwd resolution fails), warn and suggest `asylum cleanup --all`

## 3. Global Cleanup (`--all`)

- [x] 3.1 Enumerate resources: list images (`asylum:latest` + `asylum:proj-*`), volumes (`asylum-*`), check cache/projects dirs
- [x] 3.2 Print enumerated resources to terminal with clear labels
- [x] 3.3 Prompt for confirmation (y/N, default no). If not a terminal, warn and exit
- [x] 3.4 On confirmation, execute existing cleanup logic (remove images, volumes, prompt for cache removal)

## 4. Spec Update

- [x] 4.1 Update `openspec/specs/cleanup-command/spec.md` with new scenarios for scoped cleanup, `--all` flag, and confirmation prompt

## 5. Testing

- [x] 5.1 Unit test: `parseArgs` accepts `cleanup --all`, rejects `cleanup --unknown`
- [x] 5.2 Unit test: `parseArgs` accepts `--cleanup --all` flag alias
- [x] 5.3 Verify existing tests still pass

## 6. Changelog

- [x] 6.1 Add entry under Unreleased: Changed — `cleanup` now scoped to current project by default; Added — `cleanup --all` for global cleanup with confirmation
