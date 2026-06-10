## Requirements

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

### Requirement: ExecutionEngine 通过 SSE 收集子 Agent 结果
ExecutionEngine SHALL 在调用 RunTask 获得返回的 `message_id` 后，通过 SSE 连接 `GET /api/tasks/:taskId/stream?message_id=<id>` 订阅子 Agent 的流式输出，收集文本内容构建 `TaskResult`。

#### Scenario: 子 Agent 成功完成
- **WHEN** RunTask 返回 `202 + {"message_id": "msg-456"}`
- **THEN** ExecutionEngine 连接 SSE 流，收集所有 TEXT 事件中的文本，收到 DONE 事件后构建 `TaskResult(success=True, content=累积文本)`

#### Scenario: 子 Agent 执行失败
- **WHEN** SSE 流中出现 ERROR 事件或连接超时
- **THEN** ExecutionEngine 构建 `TaskResult(success=False, content=错误信息)`

#### Scenario: 子 Agent 超时
- **WHEN** 子 Agent 执行超过 timeout（默认 300s）
- **THEN** ExecutionEngine 断开 SSE 连接，构建 `TaskResult(success=False, content="[Timeout] ...")`

### Requirement: agentend 新增 BackendClient
agentend SHALL 新增 `BackendClient` 类，封装对后端 API 的 HTTP 调用。SHALL 支持调用 `POST /api/tasks/:taskId/run` 和 `GET /api/tasks/:taskId/stream`。

#### Scenario: BackendClient 初始化
- **WHEN** agentend 启动
- **THEN** BackendClient 使用配置中的后端地址（默认 `http://localhost:8080`）初始化 HTTP 客户端

#### Scenario: 调用 RunTask
- **WHEN** ExecutionEngine 调用 `backend_client.run_task(task_id, session_id, message, agent_type)`
- **THEN** 发送 `POST /api/tasks/:taskId/run`，返回 response 中提取 `message_id`

### Requirement: orchestrator SSE 流保留子 Agent 状态事件
OrchestratorAdapter 的 SSE 流 SHALL 继续推送 `RUNTIME_EXECUTING` 和 `RUNTIME_COMPLETED` 事件，保持前端对子 Agent 状态的实时可见性。子 Agent 的 TEXT 事件由后端 Redis Stream 独立推送。

#### Scenario: 子 Agent 开始执行
- **WHEN** ExecutionEngine 开始等待子 Agent 完成
- **THEN** orchestrator 流推送 `RUNTIME_EXECUTING(task_id, agent, status="running")`

#### Scenario: 子 Agent 完成执行
- **WHEN** ExecutionEngine 通过 SSE 收到子 Agent 的 DONE 事件
- **THEN** orchestrator 流推送 `RUNTIME_COMPLETED(task_id, agent, success, duration)`
