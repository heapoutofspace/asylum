## ADDED Requirements

### Requirement: Mount git worktree directories
When the project directory is a git worktree, the volume assembly SHALL mount both the worktree-specific gitdir and the main repo's `.git` directory into the container.

#### Scenario: Project is a git worktree
- **WHEN** the project directory's `.git` is a file containing `gitdir: /repo/.git/worktrees/feature`
- **THEN** both `/repo/.git/worktrees/feature` and `/repo/.git` are mounted at their real host paths

#### Scenario: Project is a regular repo
- **WHEN** the project directory's `.git` is a directory
- **THEN** no additional git-related volumes are added (`.git` is already inside the mounted project dir)

#### Scenario: Project has no .git
- **WHEN** the project directory has no `.git` file or directory
- **THEN** no additional git-related volumes are added
