## Why

When the project directory is a git worktree, `.git` is a file pointing to the main repo's `.git/worktrees/<name>` directory, which lives outside the project dir. Since asylum only mounts the project directory, git operations inside the container fail — the git object store, refs, and config aren't accessible.

## What Changes

- Detect when the project directory is a git worktree (`.git` is a file, not a directory)
- Resolve the main repo's `.git` directory from the worktree's `.git` file
- Mount the main repo's `.git` directory into the container at the same path

## Capabilities

### New Capabilities

### Modified Capabilities

- `container-assembly`: Volume assembly detects git worktrees and mounts the main repo's `.git` directory

## Impact

- `internal/container/container.go`: `appendVolumes` gains worktree detection and an additional volume mount
