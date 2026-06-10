## Requirements

### Requirement: Dispatcher converts PlanOutput to @agent dispatch instructions
Dispatcher SHALL accept a `PlanOutput` and agent config list, and produce a `list[DispatchResult]` where each item contains: task_id, agent name, @mention string, task content, depends_on list, and workspace_path。DispatchResult 的消费方从 agentend 内部 ExecutionEngine 变为后端 RunTask handler。

#### Scenario: Two-agent dispatch
- **WHEN** Planner outputs 2 tasks (one for claude-code, one for opencode) with agents config containing both agents
- **THEN** Dispatcher produces 2 DispatchResult items with `@claude-code` and `@opencode` mentions respectively

#### Scenario: Agent not found in config
- **WHEN** Planner assigns a task to agent "codex" but agents config only has "claude-code" and "opencode"
- **THEN** DispatchResult workspace_path SHALL be empty string

### Requirement: DispatchResult is structured JSON for frontend consumption
Each DispatchResult SHALL be serializable to JSON with fields: task_id, agent, mention, content, depends_on, workspace_path.

#### Scenario: SSE stream includes dispatch events
- **WHEN** OrchestratorAdapter streams dispatch results
- **THEN** each DispatchResult is yielded as a StreamEvent with type "dispatch" and payload containing the full DispatchResult JSON

### Requirement: Dispatcher produces execution waves via topological sort

The Dispatcher SHALL topologically sort DispatchResults based on their depends_on field, producing execution_waves — a list of lists where each inner list contains tasks that can execute in parallel.

#### Scenario: No dependencies — single wave
- **WHEN** all tasks have empty depends_on
- **THEN** execution_waves = [[task-001, task-002, task-003]]

#### Scenario: Chain dependencies — sequential waves
- **WHEN** task-002 depends on task-001, task-003 depends on task-002
- **THEN** execution_waves = [[task-001], [task-002], [task-003]]

#### Scenario: Mixed dependencies
- **WHEN** task-002 and task-003 depend on task-001, no other dependencies
- **THEN** execution_waves = [[task-001], [task-002, task-003]]
