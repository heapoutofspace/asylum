## MODIFIED Requirements

### Requirement: Default config content
The default config SHALL be assembled from a header template, kit-provided `ConfigNodes` (for active kits) and comments (for opt-in kits), and a footer template. Kits provide their own structured config via `ConfigNodes` rather than text snippets.

#### Scenario: Default kits present
- **WHEN** the default config is examined
- **THEN** it contains active kits for java (versions 17, 21, 25; default 21), python, and node with their default settings, assembled from each kit's `ConfigNodes` output

#### Scenario: Optional sections commented out
- **WHEN** the default config is examined
- **THEN** optional agents (gemini, codex, opencode), optional kits (apt, shell), ports, volumes, and env sections are present but commented out with explanatory comments
