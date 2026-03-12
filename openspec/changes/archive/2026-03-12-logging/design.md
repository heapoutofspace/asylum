## Context

PLAN.md specifies five log levels with colored prefixes. This is a small, focused package — hand-rolled ANSI codes are simpler than pulling in a dependency.

## Goals / Non-Goals

**Goals:**
- Five log functions matching PLAN.md spec: Info, Success, Warn, Error, Build
- Each prints a colored prefix followed by the message
- Error writes to stderr; all others to stdout
- Printf-style formatting (format string + args)

**Non-Goals:**
- No log levels, filtering, or verbosity control
- No file logging
- No structured logging

## Decisions

- **No external dependency**: Five ANSI color codes and `fmt.Fprintf` are sufficient. Adding `fatih/color` would be overkill for this.
- **Function signatures**: `Info(format, args...)`, etc. — matching `fmt.Printf` convention. Simple and familiar.
- **No interface**: Concrete functions, not a logger object. There's only one implementation.

## Risks / Trade-offs

- ANSI codes won't render in non-terminal contexts (e.g., piped output). Acceptable — this is a CLI tool meant for interactive use.
