## ADDED Requirements

### Requirement: Message 数据模型
系统 SHALL 在后端新增 Message 模型，字段包括：id（uint 主键）、task_id（string，关联 task）、session_id（string，来自哪个 session）、role（string，user/agent）、content（text，消息内容）、agent_type（string，可选）、agent_name（string，可选）、created_at（time）。

#### Scenario: 数据库迁移
- **WHEN** 运行数据库迁移
- **THEN** 自动创建 messages 表，包含上述所有字段，task_id 建立索引

### Requirement: 保存 user message
系统 SHALL 在收到消息发送请求时，立即将 user message 保存到 Message 表，role 为 "user"，然后再转发给 agentend。

#### Scenario: user message 持久化
- **WHEN** 前端发送 `POST /api/tasks/:taskId/run`，body 含 message
- **THEN** 后端先创建 Message 记录（role=user, content=message），再转发请求到 agentend

#### Scenario: 转发失败不回滚消息
- **WHEN** user message 已保存但转发 agentend 失败
- **THEN** user message 仍保留在数据库中，前端可看到历史发送记录

### Requirement: 保存 agent message
系统 SHALL 在 SSE 流完成（收到 done 事件）后，将 agent 的完整响应内容保存到 Message 表，role 为 "agent"。

#### Scenario: agent message 持久化
- **WHEN** SSE 流收到 done 事件
- **THEN** 后端创建 Message 记录（role=agent, content=完整响应, agent_type, agent_name）

#### Scenario: SSE 流中断
- **WHEN** SSE 流中途断开未收到 done 事件
- **THEN** 已流式返回的部分内容仍保存为 agent message，标记为不完整

### Requirement: 加载 task 消息历史
系统 SHALL 提供 `GET /api/tasks/:taskId/messages` 端点，返回该 task 下所有消息，按 created_at 升序排列。

#### Scenario: 加载历史消息
- **WHEN** 前端发送 `GET /api/tasks/{taskId}/messages`
- **THEN** 返回 HTTP 200，data 为 Message 数组，按时间升序

#### Scenario: 空消息列表
- **WHEN** task 下没有任何消息
- **THEN** 返回 HTTP 200，data 为空数组

#### Scenario: task 不存在
- **WHEN** 查询不存在的 task_id
- **THEN** 返回 HTTP 404
