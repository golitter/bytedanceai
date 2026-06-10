## ADDED Requirements

### Requirement: SSE connection via fetch ReadableStream
The system SHALL establish SSE connections using native `fetch` + `ReadableStream` + `TextDecoder`. The connection MUST support POST method with JSON body containing `message`, `agent_type`, and `session_id`.

#### Scenario: Successful streaming connection
- **WHEN** user sends a message via `POST /api/tasks/:taskId/run`
- **THEN** system opens ReadableStream, reads chunks via TextDecoder, and emits parsed StreamEvent objects

#### Scenario: Cross-chunk SSE event boundary
- **WHEN** a `\n\n` SSE event boundary spans two ReadableStream chunks
- **THEN** system buffers incomplete data across chunks and only emits events when a complete `\n\n`-delimited event is received

#### Scenario: Connection abort
- **WHEN** user navigates away or cancels during streaming
- **THEN** system calls `AbortController.abort()` to terminate the connection and clean up resources

### Requirement: StreamEvent parsing
The system SHALL parse each SSE `data:` line as a JSON object conforming to the StreamEvent contract (`type`, `content`, `timestamp`). Supported event types: `init`, `text`, `tool_call`, `tool_result`, `artifact`, `planning`, `done`, `error`.

#### Scenario: Text event received
- **WHEN** a `{ "type": "text", "content": { "text": "..." } }` event arrives
- **THEN** system appends the text content to the current streaming agent message

#### Scenario: Done event received
- **WHEN** a `{ "type": "done" }` event arrives
- **THEN** system finalizes the streaming agent message and transitions to `done` status

#### Scenario: Error event received
- **WHEN** a `{ "type": "error", "content": { "message": "..." } }` event arrives
- **THEN** system displays the error message and transitions to `error` status

#### Scenario: Malformed SSE data
- **WHEN** an SSE `data:` line contains invalid JSON
- **THEN** system logs the parse error and skips that event without crashing the stream

### Requirement: Chat state machine
The system SHALL manage chat lifecycle using a discriminated union state machine with states: `idle`, `loading`, `streaming`, `tool_running`, `done`, `error`. The implementation MUST use `useReducer`, NOT Zustand, for streaming state management.

#### Scenario: State transitions on happy path
- **WHEN** user sends a message
- **THEN** state transitions: `idle` → `loading` → `streaming` (on first text event) → `done` (on done event)

#### Scenario: State transition on error
- **WHEN** an error event is received or network fails
- **THEN** state transitions to `{ status: 'error', error: Error }` and displays error UI

#### Scenario: State transition on tool_call
- **WHEN** a `tool_call` event is received during streaming
- **THEN** state transitions to `{ status: 'tool_running', toolName: string }` while maintaining the abort controller

### Requirement: Message accumulation during streaming
The system SHALL accumulate streaming text into a single agent message object, NOT create a new message per chunk. The final message is only added to the message list when streaming completes. Only the last message's DOM SHALL be updated on each chunk, not the entire list.

#### Scenario: Progressive text accumulation
- **WHEN** multiple `text` events arrive in sequence
- **THEN** their content is concatenated into a single streaming message

#### Scenario: Streaming message finalization
- **WHEN** a `done` event arrives
- **THEN** the accumulated streaming message is finalized and appended to the message history as a complete agent message
