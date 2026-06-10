## ADDED Requirements

### Requirement: Dispatcher converts PlanOutput to @agent dispatch instructions
Dispatcher SHALL accept a `PlanOutput` and agent config list, and produce a `list[DispatchResult]` where each item contains: task_id, agent name, @mention string, task content, depends_on list, and workspace_path.

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
