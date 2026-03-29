## Context

`self-update` currently resolves either "latest stable" or "dev" via the GitHub Releases API. There is no way to target a specific version. The subcommand is parsed as `self-update` (hyphenated); `selfupdate` is rejected as an unknown argument.

## Goals / Non-Goals

**Goals:**
- Allow `asylum self-update 0.5.0` to install a specific tagged release
- Accept `selfupdate` as an alias for `self-update` in CLI dispatch
- Maintain backward compatibility — bare `self-update` still means "latest"

**Non-Goals:**
- Version range constraints or "upgrade to latest patch" semantics
- Downloading arbitrary URLs or non-GitHub assets
- Tab completion for available versions

## Decisions

### 1. Version argument is positional, not a flag

`asylum self-update 0.5.0` rather than `asylum self-update --version 0.5.0`. The version is the natural "what" of the command, not a modifier. This matches `go install pkg@v1.2.3` and `brew install pkg@1.2` conventions.

### 2. Fetch by tag via GitHub API

A version argument `X.Y.Z` (with or without `v` prefix) maps to `GET /repos/{owner}/{repo}/releases/tags/v{version}`. This reuses the existing `fetchRelease` machinery — just a third URL pattern alongside "latest" and "dev".

### 3. Version and `--dev` are mutually exclusive

Both specify what to install, so combining them is ambiguous. The CLI parser rejects `self-update --dev 0.5.0` with a clear error.

### 4. `selfupdate` alias handled in the arg parser

The `case` branch in `parseArgs` matches both `self-update` and `selfupdate`. No other code changes needed — by the time dispatch runs, the subcommand is always `"self-update"`.

## Risks / Trade-offs

- **Non-existent version tag** → GitHub returns 404. The existing `fetchRelease` error path already handles non-200 responses; the error message will say "HTTP 404" which is clear enough.
- **Version normalization** — Users may type `0.5.0` or `v0.5.0`. We normalize by ensuring a `v` prefix before querying the API, since all release tags use `v` prefix.
