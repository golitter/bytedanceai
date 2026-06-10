## ADDED Requirements

### Requirement: Group chat hides empty agent messages
The system SHALL NOT render agent or system messages that have empty textual content and no structured blocks.

#### Scenario: Empty orchestrator message in history
- **WHEN** history contains an agent message with empty content and no parsed blocks
- **THEN** the message list omits that message and does not reserve visual space for it

#### Scenario: Empty message with structured blocks
- **WHEN** history contains an agent message with empty textual content but parsed structured blocks
- **THEN** the message list renders the structured blocks

### Requirement: Delegation cards summarize assignment state
The system SHALL render ask-agent delegation as a structured card showing source agent, target agent, question summary, and answer status. The card SHALL NOT duplicate the full child-agent answer when that answer is also present as a child-agent message.

#### Scenario: Answered delegation with child message
- **WHEN** an ask-agent card is answered and a separate child-agent message exists for the same question
- **THEN** the card shows the question summary and answered status, while the full answer appears only in the child-agent message

#### Scenario: Answered delegation without child message
- **WHEN** an ask-agent card is answered and no separate child-agent message exists
- **THEN** the card may show a concise answer summary but MUST remain collapsible

#### Scenario: Duplicate start and done markers
- **WHEN** persisted content contains both pending and answered markers for the same `question_id`
- **THEN** the renderer coalesces them into one card with final answered or failed status

### Requirement: Runtime progress is transient
The system SHALL treat runtime execution text as transient progress and SHALL NOT duplicate it with durable child-agent output in the final visible conversation.

#### Scenario: Runtime progress before child output
- **WHEN** a child agent is running and only `runtime_text` has arrived
- **THEN** the UI displays progress inside the corresponding runtime/task area

#### Scenario: Durable child output arrives
- **WHEN** normal child-agent text arrives for a task that already displayed runtime progress
- **THEN** the UI clears the runtime transcript and renders the normal child-agent message as the durable answer

#### Scenario: Reload after completion
- **WHEN** the user reloads after a group-chat task completes
- **THEN** the visible conversation does not show both the runtime transcript and the same child-agent answer as separate final content

### Requirement: Task failures are structured
The system SHALL render task timeouts and errors as structured failure indicators, not as raw `[Timeout]` or `[Error]` suffixes in normal prose.

#### Scenario: Runtime timeout marker
- **WHEN** a task result includes a timeout marker such as `[Timeout] Task task-004 exceeded 300.0s`
- **THEN** the UI renders a failure block/card identifying `task-004` and the timeout reason

#### Scenario: Error marker appended to prose
- **WHEN** persisted message content contains normal prose followed by a raw `[Error]` marker
- **THEN** the prose and the error indicator are rendered as separate blocks

#### Scenario: Final summary references failures
- **WHEN** the orchestrator final result includes failed tasks
- **THEN** the final summary shows the overall status as partial or failed and lists failed task IDs without relying on raw marker strings

### Requirement: Final orchestrator summary is concise by default
The system SHALL present group-chat final reports as summary-first cards with expandable details.

#### Scenario: Partial completion report
- **WHEN** an orchestrator final report includes completed and failed tasks
- **THEN** the default view shows overall status, completed count, failed count, and next recommended action

#### Scenario: Detailed report expansion
- **WHEN** the user expands the final report
- **THEN** the UI shows detailed per-task tables, logs, and evidence

#### Scenario: Long Markdown final report
- **WHEN** final report content is long Markdown
- **THEN** the default message remains scannable and the long content is available through scroll or expansion

### Requirement: Agent labels use display names
The system SHALL display user-defined agent names when available and use agent type labels only as fallback.

#### Scenario: Agent type fallback in persisted row
- **WHEN** a persisted message has `agent_name` equal to `Orchestrator` but the active group session name is `ddd`
- **THEN** the UI displays `ddd`

#### Scenario: Unknown agent name
- **WHEN** no session metadata can resolve the agent
- **THEN** the UI displays the best available fallback type label
