## Why

SDKMAN's shell initialization (`sdkman-init.sh`) is slow and adds noticeable latency to every container start. Replacing it with mise gives faster startup, a smaller image footprint, and a simpler entrypoint — while keeping the same Java version selection behavior.

## What Changes

- Replace SDKMAN with mise for managing Java (Temurin 17, 21, 25) and Gradle
- Simplify the entrypoint: `eval "$(mise activate bash)"` replaces the SDKMAN sourcing and manual `JAVA_HOME`/`PATH` manipulation
- `ASYLUM_JAVA_VERSION` env var continues to work, mapped to mise's version selection
- Remove SDKMAN entirely (install script, init sourcing, bashrc/zshrc setup)
- fnm (Node.js) and uv (Python) remain unchanged

## Capabilities

### New Capabilities

- `mise-java`: Java and Gradle version management via mise, replacing SDKMAN

### Modified Capabilities

- `container-image`: Entrypoint switches from SDKMAN sourcing to mise activation for Java version selection

## Impact

- `assets/Dockerfile`: Remove SDKMAN install, add mise install + Java/Gradle setup
- `assets/entrypoint.sh`: Replace SDKMAN block with mise activation, simplify Java version selection
- `integration/entrypoint_test.go`: Java tests may need adjustment if version strings change
- No config changes — `ASYLUM_JAVA_VERSION` and `versions.java` continue to work as before
