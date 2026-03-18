## Context

In a git worktree, the project directory's `.git` is a file containing:
```
gitdir: /path/to/main-repo/.git/worktrees/<worktree-name>
```

The worktree-specific dir (`/path/to/main-repo/.git/worktrees/<name>`) contains HEAD, index, and a few refs, but it also has a `commondir` file pointing back to the main `.git`:
```
/path/to/main-repo/.git
```

Git needs both:
1. The worktree-specific dir (pointed to by `.git` file) — for HEAD, index, worktree-specific refs
2. The main `.git` dir (pointed to by `commondir`) — for objects, packed-refs, config, hooks

Both directories are outside the mounted project directory, so both need to be mounted.

## Goals / Non-Goals

**Goals:**
- Git operations work inside the container when the project is a worktree
- No behavior change for regular (non-worktree) repos

**Non-Goals:**
- Supporting bare repositories
- Modifying the `.git` file content inside the container

## Decisions

### Detect worktree by checking if `.git` is a file

If `filepath.Join(projectDir, ".git")` is a regular file (not a directory), it's a worktree. Read the file, parse `gitdir: <path>`, resolve relative paths against the project dir.

### Mount both the worktree gitdir and the common gitdir

1. Parse `.git` file → get worktree-specific gitdir (e.g., `/repo/.git/worktrees/feature`)
2. Read `commondir` file inside that dir → get main `.git` (e.g., `/repo/.git`)
3. Mount both at their real host paths

This ensures git can resolve both the worktree-specific state and the shared object store.

### Mount read-write

Git inside the container needs to write (staging, committing, etc.), so both mounts are read-write. Use the same `z` SELinux option as the project directory mount.

## Risks / Trade-offs

- **Main repo's `.git` is shared state**: If the user modifies the main repo while the container is running, git inside the container sees the changes. This is the same as how the project directory mount works — expected and fine.
- **Symlinks inside `.git`**: Some git setups use symlinks within `.git`. Docker volume mounts handle symlinks that resolve within the mounted tree, so this should work as long as the target is within the mounted `.git` dir.
