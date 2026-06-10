## ADDED Requirements

### Requirement: Pin files are stored in common/ directory with _pins.yaml bookmark
Pin files SHALL reside in `shared/.agent/memory/common/` alongside existing shared memory. A `_pins.yaml` file in the same directory tracks pinned entries with filename, title, source, pinned_at timestamp, and AI-generated summary.

#### Scenario: Pin new content
- **WHEN** user pins content with title "API 规范"
- **THEN** a file `api-spec.md` is written to `common/`, an entry is added to `_pins.yaml` with AI-generated 1-3 sentence summary

#### Scenario: Pin existing common file
- **WHEN** user pins a file `coding-standards.md` that already exists in `common/`
- **THEN** no new file is written, only `_pins.yaml` entry is added with AI summary

### Requirement: Unpin removes bookmark only, file stays in common
Unpinning SHALL only remove the entry from `_pins.yaml`. The underlying markdown file SHALL remain in `common/` as regular shared memory accessible via `taskctl common-memory`.

#### Scenario: Unpin a pinned file
- **WHEN** user unpins `api-spec.md`
- **THEN** `_pins.yaml` entry is removed, `common/api-spec.md` file still exists

#### Scenario: Unpin non-pinned file
- **WHEN** user unpins a filename not in `_pins.yaml`
- **THEN** operation returns False, no error raised

### Requirement: Pin summaries are injected into Planner prompt
`PinMemory.get_context()` SHALL return a formatted string of all pinned entries' summaries with file path references, suitable for injecting into the Planner LLM prompt.

#### Scenario: Planner receives pinned constraints
- **WHEN** Planner builds prompt and calls `pin_memory.get_context()`
- **THEN** returned string includes all pinned summaries under "## 必须遵守的约束（Pin）" heading, each with `> 完整内容: common/<filename>` reference

### Requirement: Agents can read full pin content on demand
`PinMemory.get_full_content(filename)` SHALL return the complete markdown content of a pinned file for Agent on-demand reading.

#### Scenario: Agent needs full API spec detail
- **WHEN** Agent calls `get_full_content("api-spec.md")`
- **THEN** full markdown content of `common/api-spec.md` is returned

### Requirement: Pin API endpoints
System SHALL expose REST endpoints for pin management: POST `/v1/pin/add` (content + title), POST `/v1/pin/remove` (filename), GET `/v1/pin/list`.

#### Scenario: Add pin via API
- **WHEN** POST `/v1/pin/add` with `{content: "...", title: "API 规范"}`
- **THEN** file is written, `_pins.yaml` updated with AI summary, response returns filename

#### Scenario: Remove pin via API
- **WHEN** POST `/v1/pin/remove` with `{filename: "api-spec.md"}`
- **THEN** `_pins.yaml` entry removed, response returns success
