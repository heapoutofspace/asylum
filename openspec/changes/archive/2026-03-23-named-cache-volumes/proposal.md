## Why

Cache directories (npm, pip, maven, gradle) are currently bind-mounted from `~/.asylum/cache/<container>/` on the host. On macOS with Docker Desktop, bind mounts go through a filesystem sharing layer (VirtioFS/gRPC-FUSE) that adds significant overhead to IO-heavy operations like package installation and builds. Named Docker volumes bypass this layer entirely since they live inside the Linux VM, giving near-native IO performance.

## What Changes

- **Cache mounts switch from bind to named volumes**: The four cache directories (`~/.npm`, `.cache/pip`, `.m2`, `.gradle`) use `--mount type=volume` with named volumes instead of `-v host:container` bind mounts.
- **Volume naming**: Named volumes follow the existing pattern: `<container-name>-cache-<tool>` (e.g., `asylum-a1b2c3d4e5f6-cache-npm`).
- **Host cache directory removed**: `~/.asylum/cache/` is no longer created or used. The `--cleanup` command removes cache volumes alongside shadow volumes.
- **Shell history stays as bind mount**: History must survive container removal and is tiny — no IO benefit from switching.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `container-assembly`: Cache mounts change from bind volumes to named Docker volumes
- `cleanup-command`: Cleanup removes cache volumes (named `*-cache-*`) in addition to shadow volumes and images

## Impact

- **`internal/container/container.go`**: `appendVolumes` replaces bind mount loop with named volume mounts for caches.
- **`cmd/asylum/main.go`**: `runCleanup` removes cache volumes (the existing `ListVolumes("asylum-")` already catches them by prefix).
- **`~/.asylum/cache/`**: No longer created. Existing host cache directories become orphaned (harmless, cleaned up by `--cleanup`).
