## ADDED Requirements

### Requirement: Java versions managed by mise
The Dockerfile SHALL install mise and use it to install Java Temurin 17, 21, and 25, with 21 as the default.

#### Scenario: All Java versions available
- **WHEN** the container starts without `ASYLUM_JAVA_VERSION`
- **THEN** `java -version` reports a Java 21 Temurin build

#### Scenario: All versions pre-installed
- **WHEN** the image is built
- **THEN** `mise ls java` shows temurin-17, temurin-21, and temurin-25 installed

### Requirement: Gradle managed by mise
The Dockerfile SHALL install Gradle via mise.

#### Scenario: Gradle available
- **WHEN** the container starts
- **THEN** `gradle --version` succeeds

### Requirement: mise activation in entrypoint
The entrypoint SHALL activate mise to make managed tools available in PATH.

#### Scenario: mise activation
- **WHEN** the entrypoint runs
- **THEN** `mise` is on PATH and Java/Gradle are available without manual PATH setup

### Requirement: Java version selection via ASYLUM_JAVA_VERSION
The entrypoint SHALL select a Java version when `ASYLUM_JAVA_VERSION` is set, using mise.

#### Scenario: Select pre-installed version
- **WHEN** `ASYLUM_JAVA_VERSION` is set to a pre-installed version (17, 21, or 25)
- **THEN** `java -version` reports that Java version

#### Scenario: Select non-pre-installed version
- **WHEN** `ASYLUM_JAVA_VERSION` is set to a version not pre-installed (e.g., 11)
- **THEN** the project image Dockerfile installs it via `mise install java@temurin-<version>` and sets it as the global default

### Requirement: Non-pre-installed Java in project Dockerfile
When `versions.java` specifies a version not in the base image (17, 21, 25), the image package SHALL add a mise install command to the generated project Dockerfile.

#### Scenario: Custom Java version in project config
- **WHEN** `versions.java` is set to `11`
- **THEN** the project Dockerfile includes `mise install java@temurin-11` and `mise use --global java@temurin-11`
- **AND** the project image tag reflects this in its hash

#### Scenario: Pre-installed Java version in project config
- **WHEN** `versions.java` is set to `17`
- **THEN** no additional install is added to the project Dockerfile
- **AND** the entrypoint handles version selection at runtime
