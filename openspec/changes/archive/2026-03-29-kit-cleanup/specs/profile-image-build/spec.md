## MODIFIED Requirements

### Requirement: Dockerfile decomposition
The core Dockerfile SHALL contain only OS-level system packages, Docker CLI, language version managers (fnm, mise, uv), and user creation. Language-specific apt packages (maven, python3-dev), tool CLIs (gh, openspec), and shell configuration (oh-my-zsh, tmux, direnv hooks) SHALL be provided by their respective kits.

#### Scenario: Core fragment content
- **WHEN** the core Dockerfile is examined
- **THEN** it contains OS packages, Docker CLI, user creation, and version managers — but no maven, no python3-dev, no gh, no oh-my-zsh, and no tmux config
