## ADDED Requirements

### Requirement: Config subcommand launches tabbed TUI
The system SHALL provide an `asylum config` subcommand that launches an interactive tabbed TUI with three tabs: Kits, Credentials, and Isolation.

#### Scenario: Launch config command
- **WHEN** the user runs `asylum config`
- **THEN** a tabbed TUI is displayed with tabs "Kits", "Credentials", and "Isolation", with the "Kits" tab active by default

#### Scenario: Non-TTY environment
- **WHEN** `asylum config` is run without a TTY (e.g. piped)
- **THEN** the command SHALL exit with an error message indicating a terminal is required

### Requirement: Tab navigation with arrow keys
The user SHALL be able to switch between tabs using the left and right arrow keys (or h/l).

#### Scenario: Switch to next tab
- **WHEN** the user presses the right arrow key while on the "Kits" tab
- **THEN** the "Credentials" tab becomes active and its content is displayed

#### Scenario: Switch to previous tab
- **WHEN** the user presses the left arrow key while on the "Credentials" tab
- **THEN** the "Kits" tab becomes active and its content is displayed

#### Scenario: Wrap around at edges
- **WHEN** the user presses the right arrow key while on the "Isolation" tab (last tab)
- **THEN** the active tab SHALL NOT change (no wrap-around)

### Requirement: Kits tab shows multiselect of all kits
The Kits tab SHALL display a multiselect list of all registered kits (excluding always-on kits). Kits that are currently active in the config SHALL be pre-selected.

#### Scenario: Kit list population
- **WHEN** the Kits tab is displayed
- **THEN** all registered kits with tier Default or OptIn are listed, with currently active kits checked

#### Scenario: Toggle kit selection
- **WHEN** the user presses space on a kit entry
- **THEN** the kit's selection state is toggled (checked/unchecked)

### Requirement: Credentials tab shows multiselect of credential-capable kits
The Credentials tab SHALL display a multiselect list of all kits that have credential providers. Kits with `credentials: auto` SHALL be pre-selected.

#### Scenario: Credentials list population
- **WHEN** the Credentials tab is displayed
- **THEN** all credential-capable kits are listed with their credential labels, pre-selected if currently configured as `auto`

### Requirement: Isolation tab shows single-select for config isolation
The Isolation tab SHALL display a single-select list with options: Shared, Isolated, and Project-isolated. The current isolation level SHALL be pre-selected.

#### Scenario: Isolation options display
- **WHEN** the Isolation tab is displayed
- **THEN** three options are shown: "Shared with host", "Isolated (recommended)", "Project-isolated", with the current setting highlighted

### Requirement: Confirm applies all changes
When the user presses Enter, all changes across all tabs SHALL be applied to `~/.asylum/config.yaml`.

#### Scenario: Apply kit activation
- **WHEN** the user selects a previously inactive kit and presses Enter
- **THEN** any commented-out config for that kit (including commented options) is removed and the kit's active ConfigSnippet is inserted into the kits section

#### Scenario: Apply kit deactivation
- **WHEN** the user deselects a previously active kit and presses Enter
- **THEN** the kit's active config entry is removed from the kits section and a commented version is inserted

#### Scenario: Apply credential change
- **WHEN** the user toggles a credential kit and presses Enter
- **THEN** the kit's `credentials` value is updated to `auto` (if selected) or `false` (if deselected)

#### Scenario: Apply isolation change
- **WHEN** the user selects a different isolation level and presses Enter
- **THEN** the agent's config isolation is updated in the config file

#### Scenario: Cancel discards changes
- **WHEN** the user presses Escape
- **THEN** no changes are written and the command exits

### Requirement: Kit comment removal on activation
When activating a kit, the system SHALL detect and remove any existing commented-out config block for that kit before inserting the active config snippet.

#### Scenario: Commented kit with options
- **WHEN** a kit has a commented-out block like `# python:` followed by commented option lines (e.g. `#   versions:`, `#     - 3.14`)
- **THEN** all commented lines belonging to that kit block are removed before the active snippet is inserted

#### Scenario: No commented version exists
- **WHEN** a kit has no commented-out config in the file
- **THEN** the active snippet is inserted normally without any removal step

### Requirement: Kit entry removal on deactivation
When deactivating a kit, the system SHALL remove the active config entry (the kit key and all its nested YAML lines) before inserting the commented version.

#### Scenario: Active kit with nested config
- **WHEN** a kit `java:` has nested lines like `versions:`, `default-version:`
- **THEN** the entire block (key + all nested lines) is removed before the commented version is inserted
