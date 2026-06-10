## MODIFIED Requirements

### Requirement: 加载 task 消息历史
系统 SHALL 提供 `GET /api/tasks/:taskId/messages` 端点，支持 cursor 分页参数 `limit`（默认 20）和 `before`（消息 uint ID），返回消息数组和 `has_more` 标识。不传 `before` 时返回最新的 limit 条消息。消息按 created_at 升序排列。

#### Scenario: 无分页参数加载
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages`（不带 limit 和 before）
- **THEN** 返回该 task 下所有消息，has_more 为 false

#### Scenario: 带分页参数加载最新消息
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages?limit=20`
- **THEN** 返回最新的 20 条消息（按时间正序），如果总消息数 > 20 则 has_more 为 true

#### Scenario: cursor 向前翻页
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages?limit=20&before=42`
- **THEN** 返回 ID < 42 的最新 20 条消息（按时间正序），如果前面还有更早的消息则 has_more 为 true

#### Scenario: 已加载全部历史消息
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages?limit=20&before=5`，且 ID < 5 的消息不足 20 条
- **THEN** 返回所有更早的消息，has_more 为 false

#### Scenario: 空消息列表
- **WHEN** task 下没有任何消息
- **THEN** 返回 HTTP 200，messages 为空数组，has_more 为 false

#### Scenario: task 不存在
- **WHEN** 查询不存在的 task_id
- **THEN** 返回 HTTP 404
