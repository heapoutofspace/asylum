## Context

Cache directories (npm, pip, maven, gradle) are bind-mounted from `~/.asylum/cache/<container-name>/` into the container. On macOS with Docker Desktop, all bind mounts go through a filesystem sharing layer (VirtioFS or gRPC-FUSE) that adds latency to every IO operation. For IO-heavy workloads like `npm install` or `gradle build`, this causes measurable slowdowns compared to native Linux Docker.

Named Docker volumes live inside the Linux VM and bypass the sharing layer entirely, giving near-native IO performance. Shadow node_modules volumes already use this pattern successfully.

## Goals / Non-Goals

**Goals:**
- Replace bind-mounted cache directories with named Docker volumes for better IO on macOS.
- Ensure cleanup still works (cache volumes removed by `--cleanup`).
- Keep volume naming consistent with shadow node_modules pattern.

**Non-Goals:**
- Shell history is not switched — it's tiny and must be accessible across container recreations.

## Decisions

### 1. Named volume pattern

Cache volumes use `<container-name>-cache-<tool>` naming:
- `asylum-a1b2c3d4e5f6-cache-npm`
- `asylum-a1b2c3d4e5f6-cache-pip`
- `asylum-a1b2c3d4e5f6-cache-maven`
- `asylum-a1b2c3d4e5f6-cache-gradle`

This is consistent with shadow volumes (`<container-name>-npm-<hash>`) and ensures all asylum volumes share the `asylum-` prefix, which `--cleanup` already uses to find them via `ListVolumes("asylum-")`.

### 2. Mount syntax

Use `--mount type=volume,src=<name>,dst=<path>` (same as shadow volumes) rather than `-v name:path`. The `--mount` form is explicit and consistent with the existing shadow volume code.

### 3. No host directory creation

The `os.MkdirAll` calls that create `~/.asylum/cache/<container>/` directories are removed. Named volumes are created automatically by Docker on first use.

### 4. Cleanup already works

`runCleanup` in `main.go` calls `docker.ListVolumes("asylum-")` which matches all asylum volumes by prefix. Cache volumes named `asylum-*-cache-*` are already caught by this. No cleanup code changes needed.

### 5. Temporary migration from bind mounts

On first container start, if old host cache directories exist at `~/.asylum/cache/<cname>/<tool>/`, their contents are copied into the new named volumes via `docker cp`. The old `~/.asylum/cache/<cname>/` directory is then removed. This is non-fatal — errors are logged but don't prevent the container from starting. The migration code can be removed in a future release once users have had time to upgrade.

## Risks / Trade-offs

- **Cache not inspectable from host**: Bind mounts let users browse cache contents on the host. Named volumes require `docker volume inspect` or exec into a container. Acceptable — users rarely inspect caches.
- **Orphaned host directories**: Existing `~/.asylum/cache/` directories on user machines become unused. `--cleanup` already offers to remove `~/.asylum/cache/`, so this is handled.
- **Cache lost on volume removal**: Named volumes are removed by `--cleanup` and `docker volume prune`. This matches the current behavior (bind-mounted caches are removed by `--cleanup` too).
