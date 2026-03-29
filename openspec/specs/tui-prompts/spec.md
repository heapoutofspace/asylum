## ADDED Requirements

### Requirement: Single-choice prompt
The TUI package SHALL provide a `Select` function that displays a list of options and returns the selected index.

#### Scenario: User selects an option
- **WHEN** `Select` is called with 3 options and defaultIdx 1
- **THEN** the prompt displays all options with option 1 pre-selected, and returns the user's final selection

#### Scenario: User cancels
- **WHEN** the user presses Escape or Ctrl+C during a Select prompt
- **THEN** the function returns -1 and an error

#### Scenario: Non-interactive mode
- **WHEN** `Select` is called but stdin is not a TTY
- **THEN** the function returns the default index with no prompt

### Requirement: Multi-choice prompt
The TUI package SHALL provide a `MultiSelect` function that displays a list of options with checkboxes and returns the selected indices.

#### Scenario: User selects multiple options
- **WHEN** `MultiSelect` is called with 4 options and 2 pre-selected
- **THEN** the prompt displays all options with the pre-selected ones checked, and returns the final selection

#### Scenario: User cancels multi-select
- **WHEN** the user presses Escape during a MultiSelect prompt
- **THEN** the function returns nil and an error

### Requirement: Option descriptions
Each option SHALL support a `Label` and an optional `Description` displayed below the label.

#### Scenario: Option with description
- **WHEN** an option has both Label and Description set
- **THEN** the label is displayed prominently and the description in dimmer text below it
