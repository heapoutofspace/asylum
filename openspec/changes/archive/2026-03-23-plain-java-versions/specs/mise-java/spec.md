## MODIFIED Requirements

### Requirement: Java versions managed by mise
The Dockerfile SHALL install mise and use it to install Java 17, 21, and 25 using plain version numbers, with 21 as the default.

#### Scenario: All Java versions available
- **WHEN** the container starts without `ASYLUM_JAVA_VERSION`
- **THEN** `java -version` reports a Java 21 build

#### Scenario: All versions pre-installed
- **WHEN** the image is built
- **THEN** `mise ls java` shows java 17, 21, and 25 installed

#### Scenario: .tool-versions compatibility
- **WHEN** the project has `.tool-versions` with `java 25`
- **THEN** mise resolves it to the pre-installed version without warnings

### Requirement: Java version selection via ASYLUM_JAVA_VERSION
The entrypoint SHALL select a Java version when `ASYLUM_JAVA_VERSION` is set, using `mise use --global java@<version>` with plain version numbers.

#### Scenario: Select pre-installed version
- **WHEN** `ASYLUM_JAVA_VERSION` is set to `25`
- **THEN** `mise use --global java@25` is run and `java -version` reports Java 25

#### Scenario: Select non-pre-installed version
- **WHEN** `ASYLUM_JAVA_VERSION` is set to a version not pre-installed (e.g., `11`)
- **THEN** the project image Dockerfile installs it via `mise install java@11` and sets it as the global default
