## ADDED Requirements

### Requirement: SSE 流式运行 Agent Task
系统 SHALL 提供 `POST /api/sessions/:sid/tasks/run` 端点，创建 Task 记录并启动 Agent 运行，以 SSE 流式返回 AgentEnd 的响应。

#### Scenario: 成功运行并返回 SSE 流
- **WHEN** 发送 `POST /api/sessions/{sid}/tasks/run` 并传入 `{"message": "hello", "agent_type": "claude-code"}`
- **THEN** 返回 HTTP 200，Content-Type 为 `text/event-stream`，逐行转发 AgentEnd 的 SSE 数据

#### Scenario: SSE 流正常结束
- **WHEN** AgentEnd 发送 `data: {"type": "done", ...}` 后关闭连接
- **THEN** Go 层将 task.status 更新为 `"completed"` 并关闭 SSE 流

#### Scenario: SSE 流中途断开
- **WHEN** AgentEnd 连接中途断开或发生错误
- **THEN** Go 层将 task.status 更新为 `"failed"` 并关闭连接

#### Scenario: Session 不存在
- **WHEN** 传入不存在的 session_id
- **THEN** 返回 HTTP 404，不创建 Task

#### Scenario: 缺少 message 参数
- **WHEN** 请求 body 中无 message 字段
- **THEN** 返回 HTTP 400，提示参数错误

### Requirement: Task 记录持久化
系统 SHALL 在运行 Task 时创建数据库记录，包含 `task_id`（UUID）、`session_id`、`agent_type`、`status`、`message` 字段。

#### Scenario: 创建 Task 记录
- **WHEN** 开始运行 Agent Task
- **THEN** 数据库中写入一条 Task 记录，status 为 `"running"`，task_id 为 UUID

#### Scenario: Task 完成后更新状态
- **WHEN** SSE 流正常结束或中途断开
- **THEN** task.status 更新为 `"completed"` 或 `"failed"`，更新 `updated_at` 时间戳
