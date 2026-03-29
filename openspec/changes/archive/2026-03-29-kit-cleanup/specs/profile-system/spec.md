## MODIFIED Requirements

### Requirement: Kit struct
A kit SHALL be a Go struct with fields: Name, Description, DockerSnippet, EntrypointSnippet, BannerLines, CacheDirs, OnboardingTasks, SubKits, Deps, and DefaultOn.

#### Scenario: Kit with dependencies
- **WHEN** a kit is defined with `Deps: ["node"]`
- **THEN** the dependency is validated during resolution

#### Scenario: Default-on kit
- **WHEN** a kit is defined with `DefaultOn: true`
- **THEN** it is included in the resolved set even when not explicitly listed in config
