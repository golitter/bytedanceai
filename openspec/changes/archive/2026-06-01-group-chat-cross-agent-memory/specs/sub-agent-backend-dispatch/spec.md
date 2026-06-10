## MODIFIED Requirements

### Requirement: ExecutionEngine 通过后端 API 执行子 Agent
ExecutionEngine SHALL 不再直接调用 adapter，而是通过 HTTP 回调后端 `POST /api/tasks/:taskId/run` 发起子 Agent 执行。请求 SHALL 包含子 Agent 的 `session_id`、`message`（dispatch content）和 `agent_type`。ExecutionEngine MUST 删除 short-circuit 路径（`_get_adapter`、`_iter_adapter_with_timeout`）和 `adapter_registry` 参数，统一走 HTTP。

#### Scenario: 单个子 Agent 分派
- **WHEN** orchestrator 分派一个任务给 claude-code（session_id="sess-123"）
- **THEN** ExecutionEngine 调用 `POST /api/tasks/:taskId/run`，body 为 `{"session_id": "sess-123", "message": "<dispatch content>", "agent_type": "claude-code"}`，不经过 short-circuit 路径

#### Scenario: 多个子 Agent 并行分派
- **WHEN** orchestrator 分派 3 个任务（claude-code, opencode, claude-code）
- **THEN** ExecutionEngine 并发调用 3 次 `POST /api/tasks/:taskId/run`，每次使用对应的 session_id 和 agent_type，均走 HTTP 路径

#### Scenario: 无 adapter_registry 参数
- **WHEN** OrchestratorAdapter 创建 ExecutionEngine
- **THEN** ExecutionEngine 构造函数 SHALL 不接受 `adapter_registry` 参数
