## 1. Dockerfile: Replace SDKMAN with mise

- [x] 1.1 Remove the SDKMAN install block (curl get.sdkman.io, sdk install java/gradle, bashrc/zshrc sourcing)
- [x] 1.2 Add mise install via `curl https://mise.run | sh`
- [x] 1.3 Install Java versions: `mise install java@temurin-17 java@temurin-21 java@temurin-25`
- [x] 1.4 Set default Java 21: `mise use --global java@temurin-21`
- [x] 1.5 Install Gradle: `mise install gradle@latest` and `mise use --global gradle@latest`
- [x] 1.6 Add mise activation to bashrc/zshrc: `eval "$(mise activate bash)"` / `eval "$(mise activate zsh)"`

## 2. Entrypoint: Replace SDKMAN with mise activation

- [x] 2.1 Remove the SDKMAN sourcing block (`source sdkman-init.sh` and the ASYLUM_JAVA_VERSION case statement)
- [x] 2.2 Add `eval "$(mise activate bash)"` to the entrypoint
- [x] 2.3 Implement ASYLUM_JAVA_VERSION handling for pre-installed versions (17, 21, 25): `mise use --global java@temurin-$version`, warning for unrecognized values

## 3. Project Dockerfile: Non-pre-installed Java versions

- [x] 3.1 Add Java version handling to `generateProjectDockerfile` — when `versions.java` is set to a version not in {17, 21, 25}, emit `mise install java@temurin-<version>` and `mise use --global java@temurin-<version>` in the project Dockerfile
- [x] 3.2 Thread `versions.java` through to `EnsureProject` so it can be included in the project Dockerfile generation

## 4. Verify and test

- [x] 4.1 Build the image and verify `java -version`, `gradle --version`, and `mise ls` work
- [x] 4.2 Update integration tests if Java version output format changed
- [x] 4.3 Run integration tests (`make test-integration`)
