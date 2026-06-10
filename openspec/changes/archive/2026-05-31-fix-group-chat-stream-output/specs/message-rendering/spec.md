## ADDED Requirements

### Requirement: Structured chat blocks coalesce duplicates
The system SHALL coalesce duplicate structured blocks that refer to the same logical object before rendering or freezing a message into history.

#### Scenario: Duplicate ask-agent blocks
- **WHEN** a message contains multiple ask-agent blocks with the same `question_id`
- **THEN** the renderer shows a single ask-agent card with the latest status and merged metadata

#### Scenario: Runtime blocks with same task ID
- **WHEN** a message contains multiple runtime status blocks with the same `task_id`
- **THEN** the renderer shows a single runtime status block with the latest status and accumulated progress only while it is transient

### Requirement: Long message cards remain scrollable when collapsed
The system SHALL render long messages at a fixed preview height by default, with internal vertical scrolling and a click or control to open an expanded detail view.

#### Scenario: Collapsed long message scroll
- **WHEN** a long message exceeds the preview height
- **THEN** the user can scroll within the collapsed card to inspect nearby content

#### Scenario: Expanded long message
- **WHEN** the user clicks the long message preview or expansion control
- **THEN** the UI opens an expanded view with more room for reading the full message

#### Scenario: Structured card label
- **WHEN** a message contains a self-identifying structured card such as ask-agent
- **THEN** the outer message bubble does not add a redundant `@agent` label

### Requirement: Legacy timeout and error markers render as failure blocks
The system SHALL parse legacy timeout/error marker text into structured failure blocks where possible.

#### Scenario: Timeout marker in Markdown text
- **WHEN** message text contains `[Timeout] Task task-003 exceeded 300.0s`
- **THEN** the message renderer separates it from surrounding Markdown and renders a task failure block

#### Scenario: Error marker in Markdown text
- **WHEN** message text contains `[Error]` followed by an error description
- **THEN** the message renderer separates it from surrounding Markdown and renders an error block

#### Scenario: Unknown marker format
- **WHEN** an error-like string cannot be parsed safely
- **THEN** the renderer leaves it as plain text rather than dropping information
