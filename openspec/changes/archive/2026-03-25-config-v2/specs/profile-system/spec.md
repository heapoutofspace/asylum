## MODIFIED Requirements

### Requirement: Profile struct
A kit SHALL be a Go struct with fields: Name, Description, DockerSnippet, EntrypointSnippet, BannerLines, CacheDirs, Config, OnboardingTasks, and SubKits. The package SHALL be named `kit` and the type SHALL be named `Kit`.

#### Scenario: Kit with sub-kits
- **WHEN** a kit is defined with SubKits containing child kits
- **THEN** each child kit is a complete Kit struct accessible by name

### Requirement: Built-in profile registry
The system SHALL provide a registry of built-in kits: `java` (with sub-kits `maven`, `gradle`), `python` (with sub-kit `uv`), and `node` (with sub-kits `npm`, `pnpm`, `yarn`).

#### Scenario: Registry lookup by name
- **WHEN** a kit is looked up by name (e.g., "java")
- **THEN** the corresponding Kit struct is returned

#### Scenario: Unknown kit name
- **WHEN** a kit is looked up by a name not in the registry
- **THEN** an error is returned

### Requirement: Hierarchical activation
Activating a top-level kit SHALL activate it and all its sub-kits. Activating a sub-kit via `parent/child` syntax SHALL activate only the parent and that specific child.

#### Scenario: Activate top-level kit
- **WHEN** the kits map contains `java: {}`
- **THEN** the resolved list contains the java kit, maven sub-kit, and gradle sub-kit

#### Scenario: Activate specific sub-kit
- **WHEN** the kits map contains `java/maven: {}`
- **THEN** the resolved list contains the java kit and maven sub-kit, but not gradle

### Requirement: Default activation
When no `kits` key is specified in any config layer, the system SHALL default to activating all built-in top-level kits with all their sub-kits.

#### Scenario: No kits key in config
- **WHEN** no config layer specifies `kits`
- **THEN** all built-in kits (java, python, node) with all sub-kits are active

#### Scenario: Explicit empty kits
- **WHEN** any config layer specifies `kits: {}`
- **THEN** no language kits are active (core only)

### Requirement: Deduplication
A kit activated through multiple paths SHALL appear only once in the resolved list.

#### Scenario: Kit activated by name and by sub-kit
- **WHEN** the kits map contains both `java: {}` and `java/maven: {}`
- **THEN** the java kit and maven sub-kit each appear exactly once

## RENAMED Requirements

### Requirement: Profile struct
- **FROM**: Profile struct
- **TO**: Kit struct

### Requirement: Built-in profile registry
- **FROM**: Built-in profile registry
- **TO**: Built-in kit registry
