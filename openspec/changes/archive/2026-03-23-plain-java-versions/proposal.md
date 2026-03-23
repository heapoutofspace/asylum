## Why

The base image pre-installs Java using `java@temurin-17`, `java@temurin-21`, `java@temurin-25`. When a project has `.tool-versions` with `java 25` (the common form), mise treats it as a different version (`java@25` vs `java@temurin-25`) and reports it as missing. Using plain version numbers (`java@17`, `java@21`, `java@25`) in the base image avoids this mismatch — mise resolves these to the appropriate distribution automatically, and `.tool-versions` references work out of the box.

## What Changes

- **Dockerfile**: Install Java with `mise install java@17 java@21 java@25` instead of `java@temurin-17 java@temurin-21 java@temurin-25`. Default set with `mise use --global java@21`.
- **Entrypoint**: Remove the `temurin-` prefix mapping in the `case` statement. All versions (pre-installed and custom) use `mise use --global java@<version>` uniformly.
- **Project image**: Already uses bare versions (`java@<version>`) — no change needed.
- **Pre-installed version check**: `preinstalledJava` map in `image.go` stays the same (keys are bare numbers).

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `mise-java`: Pre-installed Java versions use plain version numbers instead of temurin-prefixed names
- `container-image`: Dockerfile Java install commands change from temurin-prefixed to plain versions

## Impact

- **`assets/Dockerfile`**: Three `mise install` arguments change, one `mise use --global` argument changes.
- **`assets/entrypoint.sh`**: The `case` statement simplifies — all versions use the same `mise use --global java@<version>` form.
- **No Go code changes**: `preinstalledJava` map and project Dockerfile generation already use bare versions.
- **Image rebuild required**: Users will get the new base image on next `--rebuild` or version update.
