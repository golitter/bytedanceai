## Requirements

### Requirement: Orchestrator lifecycle is a LangGraph StateGraph

The Orchestrator lifecycle SHALL be implemented as a LangGraph StateGraph with nodes: REASON, DISPATCH, EXECUTE, REVIEW, EVOLVE, SAVE_MEM. The graph SHALL use conditional edges for routing decisions.

#### Scenario: Text output path
- **WHEN** REASON outputs output_type="text"
- **THEN** graph routes REASON → SAVE_MEM → END

#### Scenario: Plan output path success
- **WHEN** REASON outputs output_type="plan" and REVIEW determines needs_replan=false
- **THEN** graph routes REASON → DISPATCH → EXECUTE → REVIEW → EVOLVE → SAVE_MEM → END

#### Scenario: Plan output path with replan
- **WHEN** REVIEW determines needs_replan=true and iteration < max_iterations
- **THEN** graph routes REVIEW → REASON with replan_reason injected into state

### Requirement: REVIEW node decides replan with iteration limit

The REVIEW node SHALL check task results for failures. If failures exist and iteration < max_iterations (default 3), it SHALL set needs_replan=true with a formatted failure context. Otherwise needs_replan=false.

#### Scenario: All tasks succeed
- **WHEN** all TaskResult.success is true
- **THEN** REVIEW returns needs_replan=false

#### Scenario: Some tasks fail, iteration under limit
- **WHEN** some TaskResult.success is false and state["iteration"] < state["max_iterations"]
- **THEN** REVIEW returns needs_replan=true, replan_reason contains failure details

#### Scenario: Max iterations reached
- **WHEN** some TaskResult.success is false and state["iteration"] >= state["max_iterations"]
- **THEN** REVIEW returns needs_replan=false (accept partial results)

### Requirement: Graph uses MemorySaver as checkpointer

The compiled graph SHALL use `MemorySaver()` as its checkpointer. The thread_id SHALL be the session_id, ensuring per-session state isolation.

#### Scenario: Two concurrent sessions
- **WHEN** session "a" and session "b" run concurrently
- **THEN** each has independent graph state and memory
