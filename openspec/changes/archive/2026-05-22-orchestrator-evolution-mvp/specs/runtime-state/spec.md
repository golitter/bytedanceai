## ADDED Requirements

### Requirement: RuntimeState tracks task status in memory
RuntimeState SHALL maintain an in-memory dict mapping task_id to TaskState (PENDING / RUNNING / COMPLETED / FAILED), a results dict mapping task_id to result summary string, and a running_agents dict mapping agent_id to task_id.

#### Scenario: Task lifecycle
- **WHEN** a task is created, started, and completed
- **THEN** state transitions from PENDING → RUNNING → COMPLETED, result is stored

#### Scenario: Task failure
- **WHEN** an agent fails to execute a task
- **THEN** state transitions to FAILED, no result stored

### Requirement: State transitions are direct mutations
State changes SHALL be simple dict assignments. No Event Sourcing, Reducer, or event log required.

#### Scenario: Mark task as running
- **WHEN** `state.set_running("task-001")` is called
- **THEN** `state.tasks["task-001"]` equals TaskState.RUNNING
