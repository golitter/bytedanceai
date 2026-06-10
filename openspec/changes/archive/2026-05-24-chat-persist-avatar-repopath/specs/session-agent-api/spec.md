## MODIFIED Requirements

### Requirement: SSE 流式运行 Agent Session
系统 SHALL 提供 `POST /api/tasks/:taskId/run` 端点，接受 `session_id`（Agent 会话 ID）、`message` 和 `agent_type`。在转发请求前，系统 SHALL 先将 user message 保存到 Message 表。SSE 流完成后，系统 SHALL 将 agent 响应保存到 Message 表。

#### Scenario: 成功运行并返回 SSE 流（新建 session）
- **WHEN** 发送 `POST /api/tasks/{taskId}/run` 并传入 `{"message": "hello", "agent_type": "claude-code", "session_id": "new-uuid"}`
- **THEN** Go 后端先保存 user message 到 Message 表，再创建 Session 记录，转发请求到 AgentEnd，返回 HTTP 200，Content-Type 为 `text/event-stream`

#### Scenario: 成功运行并返回 SSE 流（resume session）
- **WHEN** 发送 `POST /api/tasks/{taskId}/run` 并传入已有的 `session_id`
- **THEN** Go 后端先保存 user message，再复用已有 Session 记录，AgentEnd 自动 resume cli_session_id

#### Scenario: SSE 流正常结束并保存 agent message
- **WHEN** AgentEnd 发送 `done` 事件后关闭连接
- **THEN** Go 层将完整 agent 响应保存到 Message 表，将 session.status 更新为 `"completed"`

#### Scenario: SSE 流中途断开并保存部分内容
- **WHEN** AgentEnd 连接中途断开或发生错误
- **THEN** Go 层将已接收的部分内容保存到 Message 表，session.status 更新为 `"failed"`

#### Scenario: Task 不存在
- **WHEN** 传入不存在的 task_id
- **THEN** 返回 HTTP 404，不创建 Message 和 Session

#### Scenario: 缺少 message 参数
- **WHEN** 请求 body 中无 message 字段
- **THEN** 返回 HTTP 400，提示参数错误
