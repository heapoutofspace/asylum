## ADDED Requirements

### Requirement: Tab state preservation across tab switches
The tabbed TUI SHALL preserve user selections when switching between tabs. When the user navigates away from a tab and later returns, the tab SHALL display the selections as the user left them, not the original defaults.

#### Scenario: Multiselect selections preserved after round-trip
- **WHEN** the user toggles selections on a multiselect tab, switches to another tab, then switches back
- **THEN** the multiselect tab SHALL display the user's modified selections, not the original defaults

#### Scenario: Select cursor preserved after round-trip
- **WHEN** the user moves the cursor on a select tab, switches to another tab, then switches back
- **THEN** the select tab SHALL display the cursor at the user's last position, not the default

#### Scenario: First visit uses defaults
- **WHEN** the user navigates to a tab for the first time
- **THEN** the tab SHALL display the default selections (`DefaultSel` for multiselect, `DefaultIdx` for select)

#### Scenario: Empty selection preserved
- **WHEN** the user deselects all options on a multiselect tab, switches away, then switches back
- **THEN** the tab SHALL show no options selected (not revert to defaults)
