## MODIFIED Requirements

### Requirement: Orchestrator generates structured task plan via LLM
Orchestrator SHALL use a LangGraph tool-calling agent loop to generate an execution plan. The agent loop SHALL provide 5 tools (read_file, write_file, list_dir, run_skill, load_skill_resource) to the LLM. The LLM MAY invoke tools to gather context before producing a PlanOutput containing an overview string and a list of TaskDef objects (id, agent, title, content). The Planner prompt SHALL be augmented with: (1) L2 skill instructions from selected skills, (2) Pin Memory constraints, (3) recent Evolution experience, (4) tool usage instructions. The LLM SHALL iterate (call tools → receive results → continue thinking) until it outputs a PlanOutput JSON without tool_calls, or until max_iterations=10 is reached.

#### Scenario: Successful plan generation with tool usage
- **WHEN** LLM calls `run_skill("taskctl", "summary")` to check task state, receives result, then outputs PlanOutput JSON
- **THEN** plan_node SHALL return the parsed PlanOutput

#### Scenario: Successful plan generation with Pin constraints
- **WHEN** user has pinned "必须使用 Vue 3" and Planner builds prompt
- **THEN** Pin summary SHALL be injected into the Planner prompt under "## 必须遵守的约束（Pin）"

#### Scenario: Successful plan generation with Evolution experience
- **WHEN** 3 past orchestrations exist in evolution.yaml and Planner builds prompt
- **THEN** the 3 most recent entries SHALL be injected into the Planner prompt under "## 最近编排经验"

#### Scenario: Plan generation without Pin or Evolution
- **WHEN** no pins exist and evolution.yaml is empty
- **THEN** Planner SHALL use the base prompt without additional context sections

#### Scenario: Single agent request
- **WHEN** the user message only requires one agent
- **THEN** the LLM SHALL return a PlanOutput with exactly one TaskDef

#### Scenario: Plan generation fails (max iterations)
- **WHEN** LLM calls tools for 10 consecutive iterations without outputting PlanOutput
- **THEN** plan_node SHALL return `{"plan": None}` and the adapter SHALL yield an ERROR event

### Requirement: Orchestrator graph includes progressive disclosure nodes
The LangGraph graph SHALL include the following nodes in order: `discover` (L1) → `select` (LLM skill matching) → `load_l2` (skill instructions) → `plan` (tool-calling agent loop) → `write_shared` (file output). The graph SHALL be linear with no conditional edges between discover/select/load_l2. The `plan` node contains an internal agent loop (LLM ↔ tool_calls).

#### Scenario: Full graph execution
- **WHEN** a user message arrives
- **THEN** the graph SHALL execute discover → select → load_l2 → plan → write_shared in sequence

#### Scenario: Plan node enters agent loop
- **WHEN** plan_node starts and L2 content is available
- **THEN** LLM SHALL receive L2 instructions as part of system prompt and have access to 5 tools
