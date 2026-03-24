## 1. First-run package

- [x] 1.1 Create `internal/firstrun/firstrun.go` with `Run(homeDir string) error` that detects first-run (`~/.asylum` not existing), checks for credential files, prompts the user, and writes `~/.asylum/config.yaml`
- [x] 1.2 Create `internal/firstrun/firstrun_test.go` with tests for credential detection and config generation (use temp dirs, no prompting in unit tests)

## 2. CLI integration

- [x] 2.1 Call `firstrun.Run(home)` in `cmd/asylum/main.go` before `config.Load`, gated on non-utility subcommands (skip for `--version`, `--help`, `ssh-init`, `self-update`, `--cleanup`)

## 3. Changelog

- [x] 3.1 Add entry under Unreleased > Added in CHANGELOG.md
