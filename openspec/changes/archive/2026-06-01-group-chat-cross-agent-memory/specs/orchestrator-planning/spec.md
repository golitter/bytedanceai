## MODIFIED Requirements

### Requirement: Orchestrator accepts user message and agent list
OrchestratorAdapter SHALL accept a user message string and a list of available agents (each with id, name, capabilities) via the standard `stream_chat` / `chat` interface. The agents list and shared directory path SHALL be passed through `kwargs` from `request.config`。新增 `backend_client` 参数用于 ExecutionEngine 回调后端 API。`stream_chat` SHALL 在构建 initial_state 前查询 Orchestrator 自身的跨 Agent 窗口消息，注入 `orchestrator_context` 到 state 中。

#### Scenario: Valid orchestrator request
- **WHEN** a request is sent to `/v1/agent/execute` with `agent_type: "orchestrator"` and `config` containing `agents`, `task_id`, `shared_dir`
- **THEN** OrchestratorAdapter SHALL receive `agents`, `task_id`, `shared_dir` in its kwargs and begin the full plan → dispatch → aggregate cycle

#### Scenario: Orchestrator 查询自身窗口
- **WHEN** OrchestratorAdapter.stream_chat 被调用且 backend_client 可用
- **THEN** SHALL 调用 `backend_client.get_agent_window_messages(task_id, session_id)` 获取窗口消息，通过 `build_group_chat_context` 格式化后注入 initial_state 的 `orchestrator_context`

### Requirement: System prompt supports chitchat and orchestration modes

The system prompt SHALL instruct the LLM that it is a conversational agent capable of both answering questions directly and orchestrating multi-agent tasks. When the user's request requires multi-agent collaboration, the LLM SHALL call the plan_and_dispatch tool. When the request can be answered directly, the LLM SHALL respond with text. REASON prompt SHALL 新增 `{orchestrator_context}` 占位，用于注入 Orchestrator 的跨 Agent 窗口上下文。

#### Scenario: Prompt enables dual-mode behavior
- **WHEN** REASON node constructs the system prompt
- **THEN** the prompt includes instructions for both direct reply and orchestration, available tools list, agents description, skills content, and `orchestrator_context` section if available

#### Scenario: Orchestrator context 为空
- **WHEN** 无跨 Agent 窗口消息（首次执行或单 Agent）
- **THEN** `{orchestrator_context}` 占位 SHALL 渲染为空字符串，不影响 prompt 其余部分

## ADDED Requirements

### Requirement: BackendClient 新增窗口查询方法
BackendClient SHALL 新增 `get_agent_window_messages(task_id, session_id)` 方法，调用 `GET /api/tasks/:taskId/messages/window?session_id=xxx`，返回窗口消息列表。

#### Scenario: 成功查询
- **WHEN** BackendClient 调用 `get_agent_window_messages("task-1", "sess-orch")`
- **THEN** SHALL 发送 GET 请求到 Backend，返回解析后的消息列表

#### Scenario: 查询失败降级
- **WHEN** Backend 返回错误或超时
- **THEN** SHALL 返回空列表，不阻塞 Orchestrator 执行
