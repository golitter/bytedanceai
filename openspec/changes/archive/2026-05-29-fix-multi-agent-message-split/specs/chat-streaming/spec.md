## MODIFIED Requirements

### Requirement: SSE 连接改为 GET 请求
系统 SHALL 将 SSE 连接从 `POST /api/tasks/:taskId/run` 改为 `GET /api/tasks/:taskId/stream?message_id=`。消息提交和流订阅分离为两步：先 POST 获取 message_id，再 GET 订阅 SSE。

#### Scenario: 提交消息并订阅流
- **WHEN** 用户发送消息
- **THEN** 先调用 `POST /run` 获取 `{ message_id }`，再连接 `GET /stream?message_id=xxx` 接收 SSE events

#### Scenario: POST /run 返回 202
- **WHEN** 前端发送 `POST /api/tasks/:taskId/run`
- **THEN** 后端返回 HTTP 202，body 为 `{ message_id: "uuid", status: "streaming" }`

### Requirement: SSE 重放携带 agent 元数据
`ServeStream` 在发送 MySQL 历史内容（Phase 1 和 serveCompleted）时，SHALL 在 SSE 事件中包含消息自身的 `agent_type` 和 `agent_name`。SSE 事件格式 SHALL 为 `{type:"text", content:{text, agent, agent_type}}`。

#### Scenario: serveCompleted 带 agent 元数据
- **WHEN** 前端连接到一个已完成消息的 SSE 端点
- **THEN** 每个 text 事件包含 `agent` 和 `agent_type` 字段，与 Message 记录中的值一致

#### Scenario: serveStreaming Phase 1 带 agent 元数据
- **WHEN** 前端连接到一个正在流式传输的消息，MySQL 中已有历史内容
- **THEN** Phase 1 发送的历史内容事件包含 `agent` 和 `agent_type` 字段

#### Scenario: 无 agent_type 的消息向后兼容
- **WHEN** 消息的 agent_type 为空（旧数据或单 Agent 场景）
- **THEN** SSE 事件中不包含 agent 和 agent_type 字段（与当前行为一致）
