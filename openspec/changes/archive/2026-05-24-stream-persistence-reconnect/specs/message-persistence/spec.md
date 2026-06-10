## MODIFIED Requirements

### Requirement: Message 数据模型
系统 SHALL 在后端新增 Message 模型，字段包括：id（uint 主键）、message_id（string UUID，对外唯一标识）、task_id（string，关联 task）、session_id（string，来自哪个 session）、role（string，user/agent）、content（text，消息内容）、status（string，streaming/completed/failed，默认 completed）、agent_type（string，可选）、agent_name（string，可选）、last_seq（string，Redis Stream 最后已 flush 的 event ID，默认空）、created_at（time）。

#### Scenario: 数据库迁移
- **WHEN** 运行数据库迁移
- **THEN** 自动创建 messages 表，包含上述所有字段，task_id 建立索引，message_id 建立 UNIQUE 索引

#### Scenario: 旧数据兼容
- **WHEN** 已有 Message 记录没有 message_id、status、last_seq 字段
- **THEN** message_id 默认为空，status 默认 `completed`，last_seq 默认空字符串

### Requirement: 保存 user message
系统 SHALL 在收到消息发送请求时，立即将 user message 保存到 Message 表，role 为 "user"，然后再转发给 agentend。

#### Scenario: user message 持久化
- **WHEN** 前端发送 `POST /api/tasks/:taskId/run`，body 含 message
- **THEN** 后端先创建 Message 记录（role=user, content=message），再转发请求到 agentend

#### Scenario: 转发失败不回滚消息
- **WHEN** user message 已保存但转发 agentend 失败
- **THEN** user message 仍保留在数据库中，前端可看到历史发送记录

### Requirement: 保存 agent message（增量持久化）
系统 SHALL 在 Agent 流开始时创建 Message 记录（status=streaming, content=""），在流式输出过程中通过 BatchWriter 增量 UPDATE content，流结束时设置 status=completed/failed。

#### Scenario: agent message 创建
- **WHEN** 后端启动 goroutine 消费 agentend 流
- **THEN** 创建 Message 记录（role=agent, status=streaming, content="", message_id=UUID, agent_type, agent_name）

#### Scenario: agent message 增量更新
- **WHEN** BatchWriter 触发 flush
- **THEN** UPDATE messages SET content=累积内容, last_seq=Redis Stream ID WHERE message_id=?

#### Scenario: agent message 完成
- **WHEN** goroutine 消费完 agentend 流（收到 Done event）
- **THEN** final flush 全部内容，UPDATE status=completed

#### Scenario: agent message 失败
- **WHEN** goroutine 消费过程中遇到错误或超时
- **THEN** flush 已有内容，UPDATE status=failed

#### Scenario: SSE 流中断（保留旧行为）
- **WHEN** SSE 流中途断开未收到 done 事件
- **THEN** goroutine 继续运行至完成或超时，已流式返回的部分内容通过 BatchWriter 持久化

### Requirement: 加载 task 消息历史
系统 SHALL 提供 `GET /api/tasks/:taskId/messages` 端点，返回该 task 下所有消息，按 created_at 升序排列。

#### Scenario: 加载历史消息
- **WHEN** 前端发送 `GET /api/tasks/{taskId}/messages`
- **THEN** 返回 HTTP 200，data 为 Message 数组，按时间升序，包含 status 字段

#### Scenario: 空消息列表
- **WHEN** task 下没有任何消息
- **THEN** 返回 HTTP 200，data 为空数组

#### Scenario: task 不存在
- **WHEN** 查询不存在的 task_id
- **THEN** 返回 HTTP 404
