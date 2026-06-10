## MODIFIED Requirements

### Requirement: Orchestrator accepts user message and agent list
OrchestratorAdapter SHALL accept a user message string and a list of available agents (each with id, name, capabilities) via the standard `stream_chat` / `chat` interface. The agents list and shared directory path SHALL be passed through `kwargs` from `request.config`.

#### Scenario: Valid orchestrator request
- **WHEN** a request is sent to `/v1/agent/execute` with `agent_type: "orchestrator"` and `config` containing `agents`, `task_id`, `shared_dir`
- **THEN** OrchestratorAdapter SHALL receive `agents`, `task_id`, `shared_dir` in its kwargs and begin the full plan → dispatch → aggregate cycle

#### Scenario: Missing required config fields
- **WHEN** a request is sent with `agent_type: "orchestrator"` but `config` is missing `task_id` or `shared_dir`
- **THEN** the adapter SHALL raise a KeyError which the API layer surfaces as a 500 error

### Requirement: Orchestrator generates structured task plan via LLM
Orchestrator SHALL use a single LLM call to generate an execution plan containing an overview string and a list of TaskDef objects (id, agent, title, content). The Planner prompt SHALL be augmented with Pin Memory constraints and recent Evolution experience when available.

#### Scenario: Successful plan generation with Pin constraints
- **WHEN** user has pinned "必须使用 Vue 3" and Planner builds prompt
- **THEN** Pin summary SHALL be injected into the Planner prompt under "## 必须遵守的约束（Pin）"

#### Scenario: Successful plan generation with Evolution experience
- **WHEN** 3 past orchestrations exist in evolution.yaml and Planner builds prompt
- **THEN** the 3 most recent entries SHALL be injected into the Planner prompt under "## 最近编排经验"

#### Scenario: Plan generation without Pin or Evolution
- **WHEN** no pins exist and evolution.yaml is empty
- **THEN** Planner SHALL use the base PLAN_PROMPT without additional context sections

## ADDED Requirements

### Requirement: OrchestratorAdapter implements full closed-loop
OrchestratorAdapter.stream_chat SHALL sequentially execute: plan → dispatch (produce @agent instructions) → collect agent results → aggregate (LLM summary) → yield DONE event with aggregated report.

#### Scenario: Full closed-loop execution
- **WHEN** a 2-task plan is generated for claude-code and opencode
- **THEN** OrchestratorAdapter yields PLANNING events, then dispatch events with @mention instructions, then collects results, then aggregates via LLM, then yields DONE with the aggregated report

### Requirement: Orchestrator records experience after each cycle
After each orchestration cycle completes (success or failure), the system SHALL record an entry to EvolutionStore with timestamp, message, plan summary, results summary, success boolean, and per-agent performance.

#### Scenario: Successful cycle recorded
- **WHEN** plan → dispatch → aggregate completes without errors
- **THEN** evolution.yaml gains a new entry with `success: true`

#### Scenario: Failed cycle recorded
- **WHEN** an agent execution fails during dispatch
- **THEN** evolution.yaml gains a new entry with `success: false`
