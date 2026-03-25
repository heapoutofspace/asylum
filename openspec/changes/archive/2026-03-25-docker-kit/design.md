## Context

Docker-in-Docker (DinD) is currently baked into every asylum container: the Docker engine is installed in the base image, the entrypoint starts dockerd, and the container runs with `--privileged`. This adds ~200MB to the image, requires privileged mode (security concern), and slows down image builds — all for a feature that many projects don't use.

The kit system already handles optional language toolchains. Docker fits the same pattern: a Dockerfile snippet for installation, an entrypoint snippet for daemon startup, and container-level configuration (privileged mode, env var).

## Goals / Non-Goals

**Goals:**
- Docker-in-Docker as an optional kit, active by default
- When disabled: no docker-ce in image, no privileged flag, no ASYLUM_DOCKER env, no daemon startup
- When enabled: identical behavior to current (full DinD support)

**Non-Goals:**
- Changing how Docker socket mounting works (host Docker socket is independent of this kit)
- Removing Docker CLI from core (the `docker` CLI binary is useful even without the engine, e.g., for remote Docker hosts)

## Decisions

### 1. Kit definition in internal/kit/docker.go

```go
kit.Register(&Kit{
    Name:        "docker",
    Description: "Docker-in-Docker support",
    DockerSnippet: `# Install Docker engine...`,
    EntrypointSnippet: `# Start Docker daemon...`,
})
```

The docker kit has no sub-kits, no CacheDirs, no OnboardingTasks. It's a simple kit with a Dockerfile snippet and entrypoint snippet.

### 2. Container runtime effects via config check

The `--privileged` flag and `ASYLUM_DOCKER=1` env var in container.go are conditional:

```go
if cfg.KitActive("docker") {
    args = append(args, "--privileged")
}
```

This is different from language kits which only affect the image build. The docker kit also affects container runtime arguments.

### 3. Docker CLI stays in core

The Docker CLI (`docker` binary) remains in the core Dockerfile because:
- It's needed for host Docker socket mounting (which works without DinD)
- Agent CLIs and onboarding may use `docker exec` via the host socket
- It's tiny compared to the engine

Only the Docker *engine* (docker-ce, buildx, compose) moves to the kit.

### 4. Entrypoint socket detection stays in core

The entrypoint currently checks for a Docker socket before starting the daemon. The socket check stays in core (it's useful to detect host-mounted sockets). The daemon startup logic moves to the docker kit's EntrypointSnippet.

### 5. Default config includes docker kit

The default config generated on first run includes `docker:` under kits, making DinD active by default. Users disable it by removing or commenting out the `docker:` line.

## Risks / Trade-offs

**Existing users who rely on Docker-in-Docker** → No impact. Docker kit is active by default and migration preserves existing behavior.

**`--privileged` tied to a kit** → Slightly unusual since it's a container runtime flag, not a build-time concern. But it's the right trade-off: privileged mode is only needed for DinD, so it should be tied to the docker kit.
