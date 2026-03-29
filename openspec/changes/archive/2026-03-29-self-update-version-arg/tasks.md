## 1. CLI Parsing

- [x] 1.1 Add `selfupdate` as alias in the `parseArgs` case branch (match both `self-update` and `selfupdate`)
- [x] 1.2 Parse optional positional version argument after `self-update` (stored in a new `cliFlags.TargetVersion` field)
- [x] 1.3 Reject `--dev` combined with a version argument with a clear error message
- [x] 1.4 Add test cases for: `selfupdate` alias, version argument, `v`-prefixed version, `--dev` + version conflict

## 2. Selfupdate Package

- [x] 2.1 Add `fetchRelease` support for a specific tag via `GET /repos/{owner}/{repo}/releases/tags/v{version}`
- [x] 2.2 Update `Run` to accept an optional target version parameter; skip "already up to date" when version matches current
- [x] 2.3 Normalize version input (strip or add `v` prefix as needed for API and comparison)

## 3. Dispatch

- [x] 3.1 Pass the version argument from `cliFlags` through to `selfupdate.Run` in the `self-update` dispatch block
- [x] 3.2 Update `printUsage()` help text to show `asylum self-update [version] [--dev]`

## 4. Specs

- [x] 4.1 Update `openspec/specs/self-update/spec.md` with version-targeted scenarios
- [x] 4.2 Update `openspec/specs/cli-dispatch/spec.md` with `selfupdate` alias scenarios
