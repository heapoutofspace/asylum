## Context

Mise's Java plugin resolves bare version numbers (e.g., `java@25`) to an appropriate distribution automatically. The base image currently forces the `temurin-` prefix, which creates a mismatch when projects use `.tool-versions` with bare numbers (the standard form). The entrypoint has a `case` statement that maps bare numbers back to temurin-prefixed names — this indirection is unnecessary if the base image uses bare versions from the start.

## Goals / Non-Goals

**Goals:**
- `.tool-versions` with `java 25` works without "missing" warnings
- Uniform version handling: base image, entrypoint, and project image all use bare version numbers
- Simpler entrypoint code (no temurin prefix mapping)

**Non-Goals:**
- Pinning a specific JDK distribution (mise picks the default, which is fine)
- Changing which Java versions are pre-installed (still 17, 21, 25)

## Decisions

### 1. Use bare versions in Dockerfile

Change `mise install java@temurin-17 java@temurin-21 java@temurin-25` to `mise install java@17 java@21 java@25`. Mise resolves these to the appropriate distribution. The exact distribution may differ from temurin but is functionally equivalent for development use.

### 2. Simplify entrypoint case statement

The current `case` maps `17|21|25` to `temurin-` prefixed names and everything else to bare names. With bare versions in the base image, all versions use the same form: `mise use --global java@<version>`. The entire `case` block becomes a single line.

### 3. No project image changes

`generateProjectDockerfile` in `image.go` already uses `java@<version>` (bare). No change needed.

## Risks / Trade-offs

- **Distribution may change**: Mise's default Java distribution for bare versions could change between mise updates. For development containers this is acceptable — the JDK version matters more than the distribution.
- **Image rebuild required**: Existing base images have temurin-prefixed installs. Users need a rebuild to pick up the change (happens automatically on version updates).
