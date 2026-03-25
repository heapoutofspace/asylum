## 1. Docker Kit Definition

- [x] 1.1 Create `internal/kit/docker.go` with docker kit: DockerSnippet (Docker engine install from current Dockerfile.core lines 47-56, plus `usermod -aG docker` from line 89), EntrypointSnippet (daemon startup from current entrypoint.core lines 7-26)
- [x] 1.2 Write test to verify docker kit is registered and has non-empty DockerSnippet and EntrypointSnippet

## 2. Dockerfile Core Cleanup

- [x] 2.1 Remove Docker engine install block from `assets/Dockerfile.core` (lines 47-56: docker-ce, buildx, compose install)
- [x] 2.2 Remove `usermod -aG docker ${USERNAME}` from user creation block in Dockerfile.core
- [x] 2.3 Keep Docker CLI in core: add `docker-ce-cli` to the apt-get install block (CLI only, no engine)
- [x] 2.4 Verify: core Dockerfile no longer contains `docker-ce ` (with space/end) but still has `docker-ce-cli`

## 3. Entrypoint Cleanup

- [x] 3.1 Remove Docker daemon startup block from `assets/entrypoint.core` (lines 7-26: the ASYLUM_DOCKER check and dockerd startup)
- [x] 3.2 Keep: the entrypoint.tail already handles the DOCKERD_PID cleanup — verify this still works when daemon is started from kit's EntrypointSnippet

## 4. Container Runtime Integration

- [x] 4.1 Update `container.go` RunArgs: make `--privileged` conditional on `opts.Config.KitActive("docker")`
- [x] 4.2 Update `container.go` appendEnvVars: set `ASYLUM_DOCKER=1` only when `opts.Config.KitActive("docker")`
- [x] 4.3 Update container tests for privileged/env var conditionality

## 5. Default Config and Wiring

- [x] 5.1 Add `docker:` kit to the default config template in `internal/config/defaults.go`
- [x] 5.2 Verify: all existing tests pass (docker kit active by default = same behavior)
- [x] 5.3 Add CHANGELOG entry under Unreleased

## 6. Testing

- [x] 6.1 Test: with docker kit active, assembled Dockerfile contains Docker engine install
- [x] 6.2 Test: without docker kit, assembled Dockerfile does NOT contain Docker engine install
- [x] 6.3 Test: with docker kit active, assembled entrypoint contains daemon startup
- [x] 6.4 Test: without docker kit, assembled entrypoint does NOT contain daemon startup
