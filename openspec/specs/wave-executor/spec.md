## Requirements

### Requirement: Execute node runs as a subgraph with wave-based task execution

The EXECUTE node SHALL be implemented as a LangGraph subgraph. Tasks SHALL be grouped into execution waves based on dependency topological sort. Tasks within the same wave execute in parallel; waves execute sequentially.

#### Scenario: All tasks have no dependencies
- **WHEN** all DispatchResults have empty depends_on
- **THEN** all tasks are in Wave 0 and execute in parallel

#### Scenario: Tasks have chain dependencies
- **WHEN** task-B depends on task-A
- **THEN** execution_waves = [[A], [B]], A executes first, B waits for A to complete

#### Scenario: Diamond dependency
- **WHEN** task-C and task-D depend on task-A, task-E depends on task-C and task-D
- **THEN** execution_waves = [[A], [C, D], [E]], each wave waits for previous wave

### Requirement: Each task execution invokes the agent adapter

Within each wave, the wave executor SHALL invoke the appropriate agent adapter (via ExecutionEngine or backend_client) for each task. Results SHALL be collected as TaskResult objects.

#### Scenario: Task executes via ExecutionEngine
- **WHEN** a task with backend_client available is executed
- **THEN** ExecutionEngine is used to call the agent, and TaskResult is collected

#### Scenario: Task execution yields progress events
- **WHEN** a task is running
- **THEN** task.started, agent.delta, and task.completed/failed events are yielded through the stream

### Requirement: Partial failure within a wave does not abort remaining waves

If one task in a wave fails, other tasks in the same wave SHALL continue. Failed tasks SHALL be recorded with success=false. The REVIEW node decides whether to replan.

#### Scenario: One task fails in a wave of two
- **WHEN** task-A succeeds and task-B fails in Wave 1
- **THEN** both results are collected, Wave 2 still executes (if it doesn't depend on task-B), REVIEW decides replan
