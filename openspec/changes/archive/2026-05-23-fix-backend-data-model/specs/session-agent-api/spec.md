## ADDED Requirements

### Requirement: SSE 流式运行 Agent Session
系统 SHALL 提供 `POST /api/tasks/:taskId/run` 端点，接受 `session_id`（Agent 会话 ID）、`message` 和 `agent_type`，将请求转发给 AgentEnd 并以 SSE 流式返回响应。

#### Scenario: 成功运行并返回 SSE 流（新建 session）
- **WHEN** 发送 `POST /api/tasks/{taskId}/run` 并传入 `{"message": "hello", "agent_type": "claude-code", "session_id": "new-uuid"}`
- **THEN** Go 后端创建 Session 记录，转发请求到 AgentEnd，返回 HTTP 200，Content-Type 为 `text/event-stream`

#### Scenario: 成功运行并返回 SSE 流（resume session）
- **WHEN** 发送 `POST /api/tasks/{taskId}/run` 并传入已有的 `session_id`
- **THEN** Go 后端复用已有 Session 记录，AgentEnd 自动 resume cli_session_id

#### Scenario: SSE 流正常结束
- **WHEN** AgentEnd 发送 `done` 事件后关闭连接
- **THEN** Go 层将 session.status 更新为 `"completed"`

#### Scenario: SSE 流中途断开
- **WHEN** AgentEnd 连接中途断开或发生错误
- **THEN** Go 层将 session.status 更新为 `"failed"`

#### Scenario: Task 不存在
- **WHEN** 传入不存在的 task_id
- **THEN** 返回 HTTP 404，不创建 Session

#### Scenario: 缺少 message 参数
- **WHEN** 请求 body 中无 message 字段
- **THEN** 返回 HTTP 400，提示参数错误

### Requirement: Session 记录持久化
系统 SHALL 在运行时维护 Session 记录，包含 `session_id`（调用方传入）、`task_id`、`agent_type`、`status` 字段。

#### Scenario: 首次使用新 session_id
- **WHEN** 传入的 session_id 在数据库中不存在
- **THEN** 自动创建 Session 记录，status 为 `"running"`

#### Scenario: 复用已有 session_id
- **WHEN** 传入的 session_id 在数据库中已存在
- **THEN** 不重复创建，直接使用已有记录

#### Scenario: Session 状态更新
- **WHEN** SSE 流结束或中断
- **THEN** session.status 更新为 `"completed"` 或 `"failed"`
