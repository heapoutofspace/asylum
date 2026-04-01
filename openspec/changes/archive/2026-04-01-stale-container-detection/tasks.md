## 1. Docker Helper Functions

- [x] 1.1 Add `ContainerImageID(cname string) (string, error)` to `internal/docker/docker.go` — runs `docker inspect --format '{{.Image}}' <cname>` and returns the image digest
- [x] 1.2 Add `ImageID(tag string) (string, error)` to `internal/docker/docker.go` — runs `docker inspect --format '{{.Id}}' <tag>` and returns the image digest
- [x] 1.3 Add `ContainerLabel(cname, label string) (string, error)` to `internal/docker/docker.go` — runs `docker inspect --format '{{index .Config.Labels "<label>"}}' <cname>` (reuses the pattern from `InspectLabel` but for containers)
- [x] 1.4 Add `HasActiveSessions(cname string) bool` to `internal/docker/docker.go` — runs `docker top <cname> -o pid`, counts lines minus header, returns true if more than 1 process. Returns true on error (safe default).

## 2. Config Hash

- [x] 2.1 Add `ConfigHash(cfg Config) string` function to `internal/config/config.go` — computes a SHA256 hash from deterministic serialization of volumes (sorted), env (sorted keys), and ports (sorted). Returns hex string.
- [x] 2.2 Add unit tests for `ConfigHash` — verify determinism across different map iteration orders, verify different configs produce different hashes, verify empty config produces a stable hash.

## 3. Container Label Plumbing

- [x] 3.1 Add `asylum.config.hash` label to container creation in `internal/container/container.go` `RunArgs` — accept the config hash and include it as a `--label` in the docker run args.

## 4. Restructure Main Flow

- [x] 4.1 Move `EnsureBase` and `EnsureProject` calls before the `docker.IsRunning(cname)` check in `cmd/asylum/main.go`. Wrap in a helper that returns `(imageTag, baseRebuilt, error)` and falls through gracefully on inspect errors when a container is running.
- [x] 4.2 Add stale-image detection block: when `docker.IsRunning(cname)` is true, compare `docker.ContainerImageID(cname)` with `docker.ImageID(imageTag)`. If stale and no active sessions, log and kill. If stale and active sessions, prompt via `tui.Confirm`.
- [x] 4.3 Add config-drift warning: when image is up to date but `docker.ContainerLabel(cname, "asylum.config.hash")` differs from `config.ConfigHash(cfg)`, log a warning suggesting `--rebuild`.
- [x] 4.4 Pass config hash to `container.RunArgs` so it's stored as a label on new containers.

## 5. Tests

- [x] 5.1 Add unit test in `internal/config/` confirming that kit packages from project config merge correctly and produce non-empty output from `KitPackages` — the fresh-start regression test.
- [x] 5.2 Add unit test in `internal/image/` confirming that `generateProjectDockerfile` with npm packages produces a Dockerfile containing `npm install -g` and that different packages produce different hashes.
- [x] 5.3 Add unit test for `HasActiveSessions` parsing logic (mock `docker top` output with 1 process vs multiple).

## 6. Changelog

- [x] 6.1 Add changelog entry under Unreleased/Fixed: "Kit packages from project config now trigger project image rebuild when container is restarted (inventage-ai/asylum#16)". Add under Added: "Asylum now detects stale containers and restarts automatically when no active sessions exist, or prompts when sessions are active."
