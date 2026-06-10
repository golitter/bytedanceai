## ADDED Requirements

### Requirement: Split view as default rendering mode
DiffFileView SHALL use `viewType="split"` (left-right panes) as the default diff rendering mode.

#### Scenario: Diff loaded with default view
- **WHEN** a diff card is loaded without prior user preference
- **THEN** the diff SHALL render in split view (left = old content, right = new content)

### Requirement: Unified/split toggle button
DiffCard header SHALL provide a toggle button to switch between split and unified view modes.

#### Scenario: User switches to unified view
- **WHEN** user clicks the unified toggle button
- **THEN** the diff SHALL re-render in unified (top-bottom) mode

#### Scenario: User switches back to split view
- **WHEN** user clicks the split toggle button while in unified mode
- **THEN** the diff SHALL re-render in split (left-right) mode

### Requirement: Toggle button visual state
The active view mode button SHALL be visually distinct from the inactive one.

#### Scenario: Split mode active
- **WHEN** split view is active
- **THEN** the split button SHALL have an active visual style and the unified button SHALL have an inactive style
