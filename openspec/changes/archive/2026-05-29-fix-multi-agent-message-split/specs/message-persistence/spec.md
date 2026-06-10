## MODIFIED Requirements

### Requirement: 保存 agent message（增量持久化）
系统 SHALL 在 Agent 流开始时创建 Message 记录（status=streaming, content=""），在流式输出过程中通过 StreamWriter 增量 UPDATE content。当 SSE TEXT 事件中 agent_type 变化时，SHALL finalize 当前 Message（status=completed）并在同 session 下创建新 Message。流结束时设置最后一条 Message 和原始 Message 的 status 为 completed/failed。

#### Scenario: agent message 创建
- **WHEN** 后端启动 goroutine 消费 agentend 流
- **THEN** 创建 Message 记录（role=agent, status=streaming, content="", message_id=UUID, agent_type, agent_name）

#### Scenario: agent message 增量更新
- **WHEN** StreamWriter 触发 flush
- **THEN** UPDATE messages SET content=累积内容, last_seq=Redis Stream ID WHERE message_id=?

#### Scenario: agent message Agent 切换拆分
- **WHEN** StreamWriter 检测到 TEXT 事件中 agent_type 与当前 Message 不同
- **THEN** finalize 当前 Message（status=completed），创建新 Message（agent_type=新值, status=streaming），后续内容写入新 Message

#### Scenario: agent message 完成
- **WHEN** goroutine 消费完 agentend 流（收到 Done event）
- **THEN** final flush 当前 Message，设置当前 Message 和原始 Message 的 status=completed

#### Scenario: agent message 失败
- **WHEN** goroutine 消费过程中遇到错误或超时
- **THEN** flush 已有内容，设置当前 Message 和原始 Message 的 status=failed

### Requirement: 加载 task 消息历史
系统 SHALL 提供 `GET /api/tasks/:taskId/messages` 端点，支持 cursor 分页参数 `limit`（默认 20）和 `before`（消息 uint ID），返回消息数组和 `has_more` 标识。不传 `before` 时返回最新的 limit 条消息。消息按 created_at 升序排列。SHALL 支持可选 `session_id` query param，传入时只返回匹配 session_id 的消息。

#### Scenario: 无分页参数加载
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages`（不带 limit 和 before）
- **THEN** 返回该 task 下所有消息，has_more 为 false

#### Scenario: 带分页参数加载最新消息
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages?limit=20`
- **THEN** 返回最新的 20 条消息（按时间正序），如果总消息数 > 20 则 has_more 为 true

#### Scenario: 按 session_id 过滤
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages?session_id=xxx&limit=20`
- **THEN** 只返回 session_id=xxx 的消息，分页逻辑不变

#### Scenario: 不传 session_id 向后兼容
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages?limit=20`（无 session_id）
- **THEN** 返回 task 下所有 session 的消息（当前行为不变）

#### Scenario: cursor 向前翻页
- **WHEN** 前端发送 `GET /api/tasks/:taskId/messages?limit=20&before=42`
- **THEN** 返回 ID < 42 的最新 20 条消息（按时间正序），如果前面还有更早的消息则 has_more 为 true

#### Scenario: 空消息列表
- **WHEN** task 下没有任何消息
- **THEN** 返回 HTTP 200，messages 为空数组，has_more 为 false
