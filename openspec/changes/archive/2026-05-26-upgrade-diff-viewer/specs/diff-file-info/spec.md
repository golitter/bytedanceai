## ADDED Requirements

### Requirement: Tab displays relative file path
DiffFileTabs SHALL display the file's relative path (e.g. `src/hooks/use-chat.ts`) instead of only the filename.

#### Scenario: Tab with relative path
- **WHEN** a diff contains file `src/hooks/use-chat-stream.ts`
- **THEN** the tab SHALL display `src/hooks/use-chat-stream.ts` (truncated with ellipsis if exceeding max width)

#### Scenario: Hover shows full path
- **WHEN** user hovers over a truncated tab
- **THEN** a tooltip SHALL display the full file path

### Requirement: Tab displays change type label
Each file tab SHALL display a single-letter change type label (M/A/D/R/C).

#### Scenario: Modified file tab
- **WHEN** a file has `type: "modify"`
- **THEN** the tab SHALL display an "M" label with blue styling

#### Scenario: Added file tab
- **WHEN** a file has `type: "add"`
- **THEN** the tab SHALL display an "A" label with green styling

#### Scenario: Deleted file tab
- **WHEN** a file has `type: "delete"`
- **THEN** the tab SHALL display a "D" label with red styling

### Requirement: File info bar above diff content
DiffCard SHALL render a file info bar between tabs and diff content showing full path, change type, and add/delete counts.

#### Scenario: File info bar content
- **WHEN** a file `src/api.ts` with `type: "modify"`, `+12`, `-5` is active
- **THEN** the info bar SHALL display: `src/api.ts` + `[M]` + `+12 -5`

### Requirement: Change type label mapping
The system SHALL map diff types to labels and colors consistently: add→A/green, delete→D/red, modify→M/blue, rename→R/purple, copy→C/gray.

#### Scenario: Rename file label
- **WHEN** a file has `type: "rename"`
- **THEN** the label SHALL display "R" with purple styling
