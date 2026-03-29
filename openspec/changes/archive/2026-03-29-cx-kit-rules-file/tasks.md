## 1. Update cx kit DockerSnippet

- [x] 1.1 Extend `DockerSnippet` in `internal/kit/cx.go` to run `cx skill > /tmp/asylum-kit-rules/cx.md` after cx installation (with `|| true` fallback)

## 2. Add cx kit EntrypointSnippet

- [x] 2.1 Add `EntrypointSnippet` to the cx kit in `internal/kit/cx.go` that bind-mounts `/tmp/asylum-kit-rules/cx.md` onto `~/.claude/rules/cx.md` if the file exists and is non-empty (using `sudo mount --bind`)

## 3. Remove RulesSnippet

- [x] 3.1 Remove the `RulesSnippet` field from the cx kit registration in `internal/kit/cx.go`

## 4. Tests

- [x] 4.1 Update any existing tests that assert cx contributes a `RulesSnippet` to sandbox rules assembly
- [x] 4.2 Verify `go test ./...` passes
