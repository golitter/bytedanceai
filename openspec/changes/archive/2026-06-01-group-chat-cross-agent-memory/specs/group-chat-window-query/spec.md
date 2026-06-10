## ADDED Requirements

### Requirement: Backend RunTask 窗口消息查询
Go Backend 的 RunTask handler 在调用 AgentEnd 之前，SHALL 查询当前 task 下其他 session 的消息窗口。查询条件 SHALL 为：`task_id` 匹配、`session_id` 不等于当前 session、`status IN ("completed", "streaming")`、`created_at` 大于当前 session 最后一条 agent 消息的时间。结果 SHALL 按 `created_at ASC, id ASC` 排序。

#### Scenario: 首次发言的 Agent（无历史消息）
- **WHEN** 当前 session 从未有过 agent 消息
- **THEN** SHALL 查询该 task 下所有其他 session 的 completed/streaming 消息（不设时间下界）

#### Scenario: 有历史消息的 Agent
- **WHEN** 当前 session 最后一条 agent 消息时间为 T
- **THEN** SHALL 只查询 `created_at > T` 且 session_id 不同的 completed/streaming 消息

#### Scenario: 单聊场景（只有 1 个 session）
- **WHEN** task 下只有当前 session
- **THEN** 查询结果 SHALL 为空列表 `[]`，RunTask 正常执行不注入任何内容

### Requirement: 窗口消息截断
每条窗口消息的 `content` 字段 SHALL 被截断为最多 2000 个 rune 字符。超过部分 SHALL 替换为 `"\n...[截断]"` 后缀。

#### Scenario: 短消息不截断
- **WHEN** 消息内容为 500 字符
- **THEN** content 字段 SHALL 保持原样

#### Scenario: 长消息截断
- **WHEN** 消息内容为 5000 字符
- **THEN** content 字段 SHALL 为前 2000 字符 + `"\n...[截断]"`

### Requirement: 窗口消息注入请求体
RunTask handler SHALL 将查询结果作为 `group_chat_messages` 字段注入到发送给 AgentEnd 的请求体中。每条消息 SHALL 包含 `role`、`agent_name`、`content`（已截断）三个字段。

#### Scenario: 有窗口消息时注入
- **WHEN** 查询到 3 条窗口消息
- **THEN** 请求体 SHALL 包含 `group_chat_messages: [{role, agent_name, content}, ...]`

#### Scenario: 无窗口消息时空列表
- **WHEN** 查询结果为空
- **THEN** 请求体 SHALL 包含 `group_chat_messages: []`

### Requirement: 窗口消息 API 路由
Go Backend SHALL 暴露 `GET /api/tasks/:taskId/messages/window?session_id=xxx` API，供 Orchestrator 自身查询窗口消息。返回格式和查询逻辑与 RunTask 内部查询一致。

#### Scenario: Orchestrator 查询自身窗口
- **WHEN** OrchestratorAdapter 调用 `get_agent_window_messages(task_id, session_id)`
- **THEN** SHALL 返回该 session 的窗口消息列表，包含 `role`、`agent_name`、`content`（已截断）

### Requirement: 窗口查询失败降级
当窗口查询发生数据库错误时，RunTask SHALL 降级为空窗口（`group_chat_messages: []`）继续执行，不阻塞任务。

#### Scenario: 数据库查询失败
- **WHEN** MySQL 查询返回错误
- **THEN** RunTask SHALL 记录 warning 日志，使用空列表继续执行
