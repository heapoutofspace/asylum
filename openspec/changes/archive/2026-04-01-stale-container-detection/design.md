## Context

The main run flow in `cmd/asylum/main.go` gates the entire image pipeline (`EnsureBase`, `EnsureProject`, `docker run`) behind `!docker.IsRunning(cname)`. Since asylum containers run detached with `sleep infinity` as PID 1, a running container persists across agent sessions. Any config change (kit packages, volumes, env vars) is silently ignored on subsequent `asylum` invocations because the code jumps straight to `docker exec`.

The image build functions themselves (`EnsureBase`, `EnsureProject`) are cheap when images are up to date — they compute a hash, do a `docker inspect`, and return. The expensive `docker build` only happens when the hash changes.

Key files:
- `cmd/asylum/main.go:214` — the `if !docker.IsRunning(cname)` gate
- `internal/image/image.go:197` — `EnsureProject` hash-based rebuild detection
- `internal/docker/docker.go:124` — `IsRunning` check

## Goals / Non-Goals

**Goals:**
- Detect stale running containers when image-affecting config changes (kit packages, profile snippets, java version, base image changes).
- Silently restart when no active exec sessions exist; prompt when sessions are active.
- Detect non-image config drift (volumes, env, ports) and warn the user.
- Confirm the fresh-start rebuild path works correctly for kit packages from project config (test coverage).

**Non-Goals:**
- Hot-reload config changes into a running container without restart.
- Track which specific config keys changed (just detect "changed" vs "unchanged").
- Detect changes to files mounted into the container (e.g., agent config edits).

## Decisions

### 1. Always run `EnsureBase` / `EnsureProject`

Move both calls before the `docker.IsRunning` check. They return quickly when images are up to date (hash check + `docker inspect`). This gives us the expected `imageTag` for comparison regardless of container state.

`runOnboarding` and `EnsureAgentConfig` remain gated behind "about to start a new container" — they prompt the user and seed config, which should only happen on fresh starts.

**Alternative considered:** Compute the expected hash without calling `Ensure*` (to avoid any side effects). Rejected because `Ensure*` already handle the "up to date" case efficiently and we'd be duplicating hash logic.

### 2. Compare image IDs, not tags

Tags like `asylum:latest` can point to different images after a base rebuild. Compare the SHA256 image ID of the running container against the image ID of the expected tag:

```
docker inspect --format '{{.Image}}' <container>   → container's image ID
docker inspect --format '{{.Id}}' <tag>             → current ID for that tag
```

**Alternative considered:** Compare tags directly. Rejected because `EnsureProject` returns `baseTag` ("asylum:latest") when no project image is needed, and a base rebuild changes the image behind that tag.

### 3. Detect active sessions via `docker top`

Count processes in the container. The base container runs only `sleep infinity` (PID 1). Any additional processes indicate active exec sessions (agents, shells, background tasks).

```go
func HasActiveSessions(cname string) bool {
    // docker top <cname> -o pid → count lines minus header
    // > 1 process means active sessions
}
```

**Alternative considered:** Check for specific agent processes (e.g., `pgrep claude`). Rejected because it's fragile — new agents, shell sessions, or user-started processes wouldn't be detected. Counting any extra process is simpler and more robust.

### 4. Config hash for non-image drift detection

Store a hash of the runtime-relevant config (volumes, env, ports, agent isolation) as a label on the container at creation time. On subsequent runs, recompute the hash from the current config and compare.

```go
// At container creation:
docker run --label asylum.config.hash=<hash> ...

// On subsequent runs:
docker inspect --format '{{index .Config.Labels "asylum.config.hash"}}' <container>
```

The hash input includes serialized volumes, env, and ports from the merged config — the same values that feed into `container.RunArgs`. This catches any drift without tracking individual keys.

**Alternative considered:** Skip non-image drift detection entirely. Rejected because the user explicitly requested it, and volume/env changes are a common source of "why isn't this working" confusion.

### 5. Restructured main flow

```
EnsureBase(globalKits, ...)
imageTag := EnsureProject(projectKits, collectPackages(cfg), ...)
configHash := computeConfigHash(cfg)

if docker.IsRunning(cname) {
    imageStale := docker.ContainerImageID(cname) != docker.ImageID(imageTag)
    configStale := docker.ContainerLabel(cname, "asylum.config.hash") != configHash

    if imageStale {
        if !docker.HasActiveSessions(cname) {
            log.Info("config changed, restarting container...")
            docker.RemoveContainer(cname)
        } else {
            confirmed := tui.Confirm("Image has changed. Restart container?", true)
            if confirmed {
                docker.RemoveContainer(cname)
            }
        }
    } else if configStale {
        log.Warn("config changed (volumes/env/ports) — restart with --rebuild to apply")
    }
}

if !docker.IsRunning(cname) {
    runOnboarding(...)
    ensureAgentConfig(...)
    runArgs := container.RunArgs(...)  // includes --label asylum.config.hash=<hash>
    docker.RemoveContainer(cname)
    docker.RunDetached(runArgs)
    ...
}
```

Image staleness triggers automatic or prompted restart. Config-only staleness warns but doesn't kill — because a warning is less disruptive and `--rebuild` already exists as the explicit mechanism.

## Risks / Trade-offs

- **`EnsureBase` now runs even when container is running.** This adds a `docker inspect` call (~50ms) to every `asylum` invocation. Acceptable for correctness. Risk: if `docker inspect` fails (daemon restarting), asylum now errors instead of exec'ing into the running container. Mitigation: treat inspect failures as "assume up to date" and fall through.

- **Killing a container with no active sessions could interrupt background work.** A user might have started a long-running process (build, test) and detached. Mitigation: `docker top` would show that process, so `HasActiveSessions` would return true and trigger a prompt instead.

- **Config hash false positives.** If the hash computation isn't stable (e.g., map iteration order), it could trigger spurious warnings. Mitigation: sort map keys before hashing, use a deterministic serialization.

- **`EnsureProject` returns `baseTag` when no project image is needed.** If the user removes packages from config, `EnsureProject` returns `baseTag` but the container is running with a project image. The image ID comparison correctly detects this as stale.
