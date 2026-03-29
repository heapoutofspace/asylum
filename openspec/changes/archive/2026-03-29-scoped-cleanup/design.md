## Context

Currently `runCleanup()` removes `asylum:latest`, all `asylum:proj-*` images, all `asylum-*` volumes, and optionally `~/.asylum/cache/` and `~/.asylum/projects/`. This is a "nuke everything" operation. The change scopes cleanup to the current project by default and adds `--all` for the global case with a confirmation step.

## Goals / Non-Goals

**Goals:**
- Default cleanup scoped to current project: remove project image, project volumes, project cache dir
- `--all` flag for global cleanup with a preview of what will be deleted and y/N confirmation
- Base image preserved during scoped cleanup (it's shared across projects)

**Non-Goals:**
- Selective cleanup of specific resources (e.g., only volumes, only images)
- Dry-run flag (the `--all` preview serves this purpose)

## Decisions

### 1. Scoped cleanup (default)

`asylum cleanup` resolves the current working directory to a container name via `container.ContainerName(projectDir)`, then removes only resources belonging to that project:

- **Project image**: `asylum:proj-<hash>` — need to find which project image is associated. Since the project image tag is derived from config content (not the project dir hash), we list `asylum:proj-*` images and check labels, OR we simply look up what image the container was last run with. Simpler approach: the project image tag is stored nowhere persistently, but we can inspect the container (if it exists) or just list project images. However, the most reliable approach is to look at `~/.asylum/projects/<container-name>/` for state. Actually, the simplest approach: skip project image removal in scoped mode. The project image is small (layer on top of base) and shared if config is identical across projects. Focus on what's project-specific: **container, volumes, and cached project data**.

Revised scoped cleanup removes:
1. **Container**: `docker rm -f <container-name>` (if running/stopped)
2. **Volumes**: all volumes prefixed with `<container-name>-` (shadow npm volumes + cache volumes)
3. **Project data**: `~/.asylum/projects/<container-name>/` (port allocations, session counters)

This is clean because container name is deterministic from project dir, and all volume names are prefixed with the container name.

### 2. Global cleanup (`--all`)

`asylum cleanup --all` performs the existing global cleanup but adds a confirmation step:

1. Enumerate what will be deleted:
   - Images: `asylum:latest` + all `asylum:proj-*`
   - Volumes: all `asylum-*` prefixed volumes
   - Optionally: `~/.asylum/cache/` and `~/.asylum/projects/`
2. Print the list to the terminal
3. Prompt: "Proceed? (y/N)"
4. Only delete if user confirms

Non-terminal: warn and exit (same as current cache prompt behavior).

### 3. Flag parsing

The `cleanup` subcommand accepts `--all`:
```
asylum cleanup          → scoped to current project
asylum cleanup --all    → global cleanup with confirmation
asylum --cleanup        → scoped (same as `asylum cleanup`)
asylum --cleanup --all  → global with confirmation
```

`parseArgs` allows `--all` after `cleanup` subcommand (currently it rejects any args after cleanup).

### 4. No project dir requirement for `--all`

Scoped cleanup needs to resolve the project directory (current working directory). If we're not in a project directory, scoped cleanup should warn and suggest `--all`.

Global cleanup (`--all`) does not need a project directory — it works the same as the current behavior.

## Risks / Trade-offs

**Scoped cleanup doesn't remove the project image** → Acceptable. Project images are thin layers and are rebuilt cheaply. They're also content-addressed, so they may be shared across projects with identical config. `--all` still removes them.

**Running `cleanup` outside a project dir fails** → We warn with a clear message suggesting `--all`. This is a reasonable UX since cleanup outside a project context is inherently ambiguous.
