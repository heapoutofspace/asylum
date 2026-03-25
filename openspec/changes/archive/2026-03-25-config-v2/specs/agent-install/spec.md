## MODIFIED Requirements

### Requirement: Agents config field
The Config struct SHALL include an `Agents` field that is a map of agent name to AgentConfig. Agent presence in the map means installation. Nil map defaults to `{"claude": {}}`. Empty map means no agents.

#### Scenario: Agents in YAML config
- **WHEN** a config file contains `agents: {claude: {}, gemini: {}}`
- **THEN** the parsed Config has Agents map with claude and gemini entries

#### Scenario: No agents key
- **WHEN** no config file specifies `agents`
- **THEN** the parsed Config has Agents as nil (defaults to claude-only at resolution time)

#### Scenario: Empty agents map
- **WHEN** a config file contains `agents: {}`
- **THEN** no agent CLIs are installed

#### Scenario: CLI overrides agents
- **WHEN** config has `agents: {claude: {}}` and CLI passes `--agents gemini,codex`
- **THEN** the effective agents are gemini and codex

### Requirement: Agent install resolution
The system SHALL resolve agent installs from the agents map. Map keys are agent names; nil map defaults to claude-only.

#### Scenario: Map-based resolution
- **WHEN** agents map is `{claude: {}, gemini: {}}`
- **THEN** claude and gemini install definitions are returned
