## 1. Worktree detection and volume mounting

- [x] 1.1 Add `resolveGitWorktree(projectDir string) (worktreeDir, commonDir string)` to `internal/container/container.go` — reads `.git` file, parses `gitdir:`, reads `commondir`, returns both paths (empty strings if not a worktree)
- [x] 1.2 In `appendVolumes`, call `resolveGitWorktree` and mount both directories if non-empty

## 2. Tests

- [x] 2.1 Add unit test for `resolveGitWorktree` with a simulated worktree structure (`.git` file, worktree dir with `commondir` file)
- [x] 2.2 Add unit test for the non-worktree case (`.git` is a directory or missing)
