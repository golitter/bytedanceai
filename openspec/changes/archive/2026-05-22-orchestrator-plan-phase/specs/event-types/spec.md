## ADDED Requirements

### Requirement: PLANNING event type for orchestrator progress
EventType enum SHALL include a `PLANNING` value (string `"planning"`) used to signal orchestrator progress events to the frontend. OrchestratorAdapter SHALL emit PLANNING events when planning starts and when each graph node completes.

#### Scenario: Planning start event
- **WHEN** OrchestratorAdapter.stream_chat begins
- **THEN** it SHALL yield a StreamEvent with type "planning" and content containing `status: "started"`

#### Scenario: Node completion event
- **WHEN** the LangGraph "plan" or "write_shared" node completes
- **THEN** OrchestratorAdapter SHALL yield a StreamEvent with type "planning" and content containing the completed `node` name
