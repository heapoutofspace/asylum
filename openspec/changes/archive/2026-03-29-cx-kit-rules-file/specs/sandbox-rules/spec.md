## MODIFIED Requirements

### Requirement: Kit-provided rules snippets
Each kit that installs tools or provides capabilities MUST populate `RulesSnippet` with a concise markdown description, EXCEPT for kits that generate their own standalone rules files (e.g., cx). At minimum, the docker, java, python, and node kits SHALL have rules snippets.

#### Scenario: Docker kit snippet
- **WHEN** the docker kit is active
- **THEN** its rules snippet SHALL mention that full Docker engine is available (not just CLI) and the container runs in privileged mode

#### Scenario: cx kit no longer contributes snippet
- **WHEN** the cx kit is active
- **THEN** its `RulesSnippet` SHALL be empty and the assembled `asylum-sandbox.md` SHALL NOT contain a cx-specific section
- **AND** the cx tool SHALL still appear in the "Kit Tools" list via the `Tools` field
