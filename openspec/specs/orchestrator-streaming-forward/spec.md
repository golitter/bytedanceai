## ADDED Requirements

### Requirement: Forwarded runtime text is not durable answer content
The system SHALL distinguish child-agent runtime progress from durable child-agent answer text when forwarding orchestrator streams.

#### Scenario: Runtime text during execution
- **WHEN** ExecutionEngine receives a child-agent text token while a task is running
- **THEN** it emits runtime progress suitable for live display

#### Scenario: Task result content emitted later
- **WHEN** the child-agent task completes and its collected result content is emitted as normal agent text
- **THEN** downstream renderers can identify that this is durable answer content and avoid freezing prior runtime progress as duplicate final content

#### Scenario: Reloaded history
- **WHEN** the user reloads after an orchestrated run
- **THEN** durable child-agent messages remain visible and runtime progress does not appear as a second copy of the same answer

### Requirement: Task failures are emitted separately from result prose
The system SHALL expose task timeouts and execution errors as structured failure information instead of only appending marker strings to `TaskResult.content`.

#### Scenario: Child task timeout
- **WHEN** a child task exceeds its timeout
- **THEN** orchestrator streaming emits failure metadata including the task ID, agent, and timeout reason

#### Scenario: Child task error
- **WHEN** a child task returns an error event
- **THEN** orchestrator streaming emits failure metadata including the task ID, agent, and error message

#### Scenario: Aggregation with failed tasks
- **WHEN** the aggregator receives completed and failed task results
- **THEN** it generates a summary that separates completed work from failed work and does not bury raw error markers inside normal prose

### Requirement: Final aggregation is summary-first
The system SHALL generate orchestrator final output that is concise by default and structured enough for the frontend to render details separately.

#### Scenario: Successful multi-agent run
- **WHEN** all child tasks complete successfully
- **THEN** the final output includes a short success summary and references detailed child-agent outputs without duplicating all logs inline

#### Scenario: Partial multi-agent run
- **WHEN** one or more child tasks fail or time out
- **THEN** the final output identifies the run as partial, lists failed tasks, and gives the next recommended action

#### Scenario: Internal execution logs
- **WHEN** child agents produce step-by-step execution logs
- **THEN** those logs are not copied wholesale into the default final summary
