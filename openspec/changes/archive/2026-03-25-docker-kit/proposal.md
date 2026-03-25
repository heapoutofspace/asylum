## Why

Docker-in-Docker support (docker-ce, buildx, compose, daemon startup, privileged mode) is currently hardcoded into the core Dockerfile, entrypoint, and container setup. Not every project needs Docker inside the container, but every user pays the image size (~200MB) and gets a privileged container. Making Docker a kit lets users disable it for projects that don't need it, resulting in smaller images and non-privileged containers.

## What Changes

- New `docker` kit: DockerSnippet (docker-ce/buildx/compose install, docker group membership), EntrypointSnippet (daemon startup), plus container-level effects (privileged flag, ASYLUM_DOCKER env var)
- Remove Docker engine install from `Dockerfile.core` (keep Docker socket check in entrypoint core for when host socket is mounted)
- Remove hardcoded `--privileged` and `ASYLUM_DOCKER=1` from container.go — make them conditional on docker kit being active
- Docker kit is active by default (present in default config)

## Capabilities

### New Capabilities
- `docker-kit`: Docker-in-Docker as an optional kit with Dockerfile snippet, entrypoint snippet, and container runtime effects

### Modified Capabilities
- `profile-image-build`: Core Dockerfile no longer includes Docker engine installation
- `profile-container-setup`: Container `--privileged` flag and `ASYLUM_DOCKER` env var are conditional on docker kit

## Impact

- **internal/kit/docker.go** (new): Docker kit definition with DockerSnippet, EntrypointSnippet
- **assets/Dockerfile.core**: Remove Docker engine install block and `usermod -aG docker`
- **assets/entrypoint.core**: Remove Docker daemon startup (moves to kit EntrypointSnippet); keep socket detection for host-mounted sockets
- **internal/container/container.go**: `--privileged` and `ASYLUM_DOCKER=1` conditional on docker kit active in config
- **internal/config/defaults.go**: Docker kit included in default config
