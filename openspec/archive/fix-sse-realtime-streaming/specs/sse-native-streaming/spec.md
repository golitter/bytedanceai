## ADDED Requirements

### Requirement: SSE client uses native EventSource API
The system SHALL use the browser-native `EventSource` API for SSE connections instead of `fetch` + `ReadableStream`. The `connectSSE` function in `frontend/src/lib/sse.ts` SHALL create an `EventSource` instance and dispatch parsed events through the `onEvent` callback.

#### Scenario: SSE connection established
- **WHEN** `connectSSE` is called with a valid URL and params
- **THEN** an `EventSource` instance is created targeting `{url}?{params}` and events are forwarded to `onEvent`

#### Scenario: SSE event received and parsed
- **WHEN** the backend sends `data: {"type":"text","content":{"text":"hello"}}`
- **THEN** the `onEvent` callback is called with `{ type: "text", content: { text: "hello" } }`

#### Scenario: SSE connection aborted
- **WHEN** the AbortController signal fires
- **THEN** the `EventSource` is closed via `es.close()`

#### Scenario: SSE connection error without reconnect
- **WHEN** `reconnect` is false and the connection errors
- **THEN** `onError` is called with an Error and the EventSource is closed

#### Scenario: SSE connection error with reconnect enabled
- **WHEN** `reconnect` is true and the connection drops
- **THEN** `EventSource` automatically reconnects without calling `onError`

### Requirement: User message displayed immediately via optimistic update
The `sendMessage` function in `use-chat-stream.ts` SHALL add the user message to the Zustand store BEFORE awaiting the `submitMessage` API call. If the API call fails, the store SHALL be updated to error state.

#### Scenario: Message sent successfully
- **WHEN** user sends a message
- **THEN** the message appears in the chat immediately (before API response)
- **AND** after API response, SSE stream begins delivering agent reply

#### Scenario: Message send fails
- **WHEN** the `submitMessage` API call fails
- **THEN** `store.streamError` is called with the error, showing error state in the UI

### Requirement: Backend serveStreaming handles goroutine registration delay
The `serveStreaming` function in `backend/internal/handler/stream.go` SHALL NOT immediately fail when `stream.IsActive()` returns false and the message status is still "streaming". It SHALL continue the XRead blocking loop to wait for the goroutine to register or for the status to change to "completed"/"failed".

#### Scenario: Goroutine not yet registered
- **WHEN** `stream.IsActive()` returns false and message status is "streaming"
- **THEN** the handler SHALL fall through to the XRead block instead of sending an error
- **AND** the handler SHALL wait up to 5 seconds before re-checking

#### Scenario: Goroutine completed
- **WHEN** `stream.IsActive()` returns false and message status is "completed"
- **THEN** the handler SHALL send remaining content diff and a `done` event

#### Scenario: Goroutine failed
- **WHEN** `stream.IsActive()` returns false and message status is "failed"
- **THEN** the handler SHALL send an error event
