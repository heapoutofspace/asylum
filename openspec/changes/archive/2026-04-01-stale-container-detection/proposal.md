## Why

When a user changes kit packages (or other image-affecting config) in their project `.asylum` file and re-runs `asylum`, the changes are silently ignored if the container is still running. The entire image pipeline is gated behind `!docker.IsRunning(cname)` — so config changes never reach the hash check. Even users who manually stop the container may not realize a stale container was still running. This leads to confusion where packages appear missing despite correct configuration (see inventage-ai/asylum#16).

Beyond the running-container case, non-image config changes (volumes, env, ports) also require a restart but produce no warning at all.

## What Changes

- Move `EnsureBase` and `EnsureProject` calls before the running-container check so image freshness is always evaluated.
- When a running container's image doesn't match the expected image tag: kill and restart silently if no active exec sessions exist, or prompt the user if sessions are active.
- When a running container's image matches but non-image config (volumes, env, ports) has changed: warn the user that a restart is needed.
- Add docker helper functions: `ContainerImageID`, `ImageID`, `HasActiveSessions`.
- Add a config hash mechanism to detect non-image config changes against the running container.
- Add unit tests to confirm that kit packages from project config produce the correct Dockerfile and hash.

## Capabilities

### New Capabilities
- `stale-container-detection`: Detect when a running container's image or config is stale relative to current config, and handle restart automatically or via prompt.

### Modified Capabilities
- `image-build`: `EnsureBase`/`EnsureProject` are now called unconditionally (not gated behind container-not-running check), though their behavior is unchanged.

## Impact

- `cmd/asylum/main.go`: Restructure the main run flow to always check images and add stale-container logic.
- `internal/docker/docker.go`: New helper functions for container/image introspection and session detection.
- `internal/container/`: Store and compare a config hash to detect non-image config drift.
- Existing tests unaffected — new unit tests added for docker helpers and config hash comparison.
