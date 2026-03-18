## Context

Java versions are currently managed by SDKMAN, which requires sourcing `sdkman-init.sh` on every container start. This script is noticeably slow (~3-5s) because it sets up shell functions, checks for updates, and configures candidates. The entrypoint then manually overrides `JAVA_HOME` and `PATH` when `ASYLUM_JAVA_VERSION` is set.

mise is a polyglot version manager written in Rust. It activates in ~10ms and can manage Java (via adoptium/temurin builds), Gradle, and many other tools. It uses `~/.tool-versions` or `~/.config/mise/config.toml` for configuration.

## Goals / Non-Goals

**Goals:**
- Replace SDKMAN with mise for Java and Gradle management
- Preserve `ASYLUM_JAVA_VERSION` behavior (17, 21, 25 selection at runtime)
- Support arbitrary Java versions via project Dockerfile when not pre-installed
- Faster entrypoint startup
- Simpler entrypoint logic for Java version selection

**Non-Goals:**
- Replacing fnm (Node.js) or uv (Python) — these stay as-is
- Changing the user-facing config (`versions.java`)

## Decisions

### Use mise's global config for default Java version

Set the default Java 21 via `~/.config/mise/config.toml` at image build time. The entrypoint uses `mise use --global java@temurin-<version>` when `ASYLUM_JAVA_VERSION` is set, which is simpler than the current manual `JAVA_HOME`/`PATH` manipulation.

Alternative: Use environment variables (`MISE_JAVA_VERSION`). Rejected because `mise use` integrates better with mise's shim/activation model and keeps PATH management in one place.

### Install mise via official install script

`curl https://mise.run | sh` — consistent with how we install uv and fnm. The binary lands in `~/.local/bin` which is already on PATH.

### Keep all three Java versions pre-installed in the image

`mise install java@temurin-17 java@temurin-21 java@temurin-25` at build time. This matches the current behavior where all versions are available without network access at runtime.

### Install Gradle via mise instead of SDKMAN

`mise install gradle@latest` replaces `sdk install gradle`. The system Gradle package from apt is not installed — SDKMAN was the only provider.

### Non-pre-installed Java versions via project Dockerfile

The base image ships with Java 17, 21, and 25. If `versions.java` is set to something else (e.g., `11`), `generateProjectDockerfile` adds `mise install java@temurin-<version>` and `mise use --global java@temurin-<version>` to the project Dockerfile. This is consistent with how apt/npm/pip packages work — custom versions get baked into the project image.

Pre-installed versions (17, 21, 25) are handled at runtime in the entrypoint, so they don't trigger a project image build just for Java version selection.

The set of pre-installed versions is defined as a constant in the image package, used by both `generateProjectDockerfile` (to decide whether to add an install step) and referenced by the entrypoint (which only does runtime switching for versions it knows are present).

## Risks / Trade-offs

- **mise Java version identifiers may differ from SDKMAN** → The entrypoint maps `ASYLUM_JAVA_VERSION=17` to `java@temurin-17`, which mise resolves to the latest patch. Test that version strings still contain "17", "21", "25" as expected by integration tests.
- **mise is a newer tool than SDKMAN** → mise is actively maintained and widely adopted. If it breaks, we can pin a version in the install command.
- **Image rebuild required** → This changes the Dockerfile significantly, so all users will get a full base image rebuild on next run. This is expected and acceptable.
