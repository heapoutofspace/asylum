## Why

All subsequent packages need colored terminal output for user-facing messages. PLAN.md section 5.9 defines five log levels with specific colors and prefixes. This must exist before any feature code is written.

## What Changes

- Create `internal/log` package with functions for each log level: Info (blue `i`), Success (green `ok`), Warn (yellow `!`), Error (red `x` to stderr), Build (cyan `#`)
- Use ANSI escape codes directly — no external dependency needed for this simple use case

## Capabilities

### New Capabilities
- `colored-logging`: Colored terminal output with five log levels per PLAN.md section 5.9

### Modified Capabilities

None.

## Impact

- Adds `internal/log/log.go` — all subsequent packages import this instead of `fmt.Println` or stdlib `log`
- No external dependencies added
