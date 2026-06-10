## 1. Data Models

- [x] 1.1 Add `TaskResult` model to `orchestrator/models.py` (task_id, agent, success, content, duration)
- [x] 1.2 Add `DispatchResult` model to `orchestrator/models.py` (task_id, agent, mention, content, depends_on, workspace_path)

## 2. RuntimeState

- [x] 2.1 Create `orchestrator/state.py` with `TaskState` enum (PENDING, RUNNING, COMPLETED, FAILED)
- [x] 2.2 Implement `RuntimeState` class with tasks dict, results dict, running_agents dict

## 3. Pin Memory

- [x] 3.1 Create `orchestrator/pin_memory.py` with `PinMemory` class operating on `memory/common/`
- [x] 3.2 Implement `pin()` — write file to common/ + add _pins.yaml entry + AI summarize
- [x] 3.3 Implement `pin_existing()` — add _pins.yaml bookmark for existing common file + AI summarize
- [x] 3.4 Implement `unpin()` — remove _pins.yaml entry only, keep file
- [x] 3.5 Implement `get_context()` — return formatted summaries for prompt injection
- [x] 3.6 Implement `get_full_content()` — return full markdown content by filename
- [x] 3.7 Add Pin API endpoints in `api/v1/` — POST /v1/pin/add, POST /v1/pin/remove, GET /v1/pin/list

## 4. Self-Evolution

- [x] 4.1 Create `orchestrator/evolution.py` with `EvolutionStore` class
- [x] 4.2 Implement `record()` — append entry to evolution.yaml, cap at 20 entries
- [x] 4.3 Implement `get_recent_experience(n=5)` — return formatted string for prompt injection

## 5. Dispatcher

- [x] 5.1 Create `orchestrator/dispatcher.py` with `Dispatcher` class
- [x] 5.2 Implement `dispatch()` — convert PlanOutput + agents config to list[DispatchResult]
- [x] 5.3 Map agent IDs to workspace paths from agents config

## 6. Aggregator

- [x] 6.1 Create `orchestrator/aggregator.py` with `Aggregator` class
- [x] 6.2 Implement `aggregate()` — LLM call to summarize multiple agent results + overview

## 7. OrchestratorAdapter Upgrade

- [x] 7.1 Upgrade `adapters/orchestrator.py` stream_chat to implement plan → dispatch → collect → aggregate loop
- [x] 7.2 Yield PLANNING events during plan phase
- [x] 7.3 Yield dispatch events with DispatchResult JSON for each task
- [x] 7.4 Collect agent execution results (from caller/Go Backend via results callback)
- [x] 7.5 Call Aggregator and yield DONE event with aggregated report
- [x] 7.6 Record experience to EvolutionStore after cycle completes

## 8. Planner Prompt Upgrade

- [x] 8.1 Create `build_planner_prompt()` function that composes base prompt + Pin context + Evolution context
- [x] 8.2 Update `graph.py` plan_node to use `build_planner_prompt()` instead of raw `PLAN_PROMPT.format()`

## 9. Verification

- [x] 9.1 Test Pin Memory: pin new content, pin existing file, unpin, get_context, get_full_content
- [x] 9.2 Test EvolutionStore: record entries, cap at 20, get_recent_experience
- [x] 9.3 Test Dispatcher: PlanOutput → list[DispatchResult] with correct @mentions
- [x] 9.4 Test Aggregator: LLM aggregation of 2+ results
- [x] 9.5 Test full closed-loop via curl: orchestrator request → plan → dispatch events → aggregate → DONE
