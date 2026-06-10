## MODIFIED Requirements

### Requirement: Orchestrator accepts user message and agent list
OrchestratorAdapter SHALL accept a user message string and a list of available agents (each with id, name, capabilities) via the standard `stream_chat` / `chat` interface. The agents list and shared directory path SHALL be passed through `kwargs` from `request.config`。新增 `backend_client` 参数用于 ExecutionEngine 回调后端 API。

#### Scenario: Valid orchestrator request
- **WHEN** a request is sent to `/v1/agent/execute` with `agent_type: "orchestrator"` and `config` containing `agents`, `task_id`, `shared_dir`
- **THEN** OrchestratorAdapter SHALL receive `agents`, `task_id`, `shared_dir` in its kwargs and begin the full plan → dispatch → aggregate cycle

#### Scenario: Missing required config fields
- **WHEN** a request is sent with `agent_type: "orchestrator"` but `config` is missing `task_id` or `shared_dir`
- **THEN** the adapter SHALL raise a KeyError which the API layer surfaces as a 500 error
