## ADDED Requirements

### Requirement: SessionAgent model
系统 SHALL 提供 `session_agents` 数据表，包含以下字段：`id`（自增主键）、`session_id`（VARCHAR 128，UNIQUE KEY）、`agent_type`（VARCHAR 64）、`agent_name`（VARCHAR 128）、`avatar_url`（VARCHAR 512）、`created_at`、`updated_at`。与 Session 为 1:1 关系。

#### Scenario: AutoMigrate creates table
- **WHEN** 后端启动执行 AutoMigrate
- **THEN** SHALL 自动创建 `session_agents` 表，`session_id` 为 UNIQUE KEY

### Requirement: Session query joins session_agents
后端查询 Session 列表或详情时 SHALL LEFT JOIN `session_agents` 表，将 `agent_type`、`agent_name`、`avatar_url` 作为 Session 响应的一部分返回。API 响应 JSON 结构与拆表前保持一致。

#### Scenario: List sessions with agent info
- **WHEN** 调用 `GET /api/tasks/:taskId` 获取任务详情（含 sessions）
- **THEN** SHALL 返回每个 session 的 `agent_type`、`agent_name`、`avatar_url` 字段（来自 session_agents JOIN）

#### Scenario: Session without agent info
- **WHEN** 某个 session 在 session_agents 表中无对应记录
- **THEN** SHALL 返回 `agent_type`、`agent_name`、`avatar_url` 为 null/空值

### Requirement: Create session writes session_agents
创建 Session 时 SHALL 同时写入 `session_agents` 记录（agent_type、agent_name）。写入操作 MUST 在同一事务中完成。

#### Scenario: Create session with agent info
- **WHEN** 创建新 session 且指定 `agent_type: "claude-code"`、`agent_name: "Claude"`
- **THEN** SHALL 在 session_agents 表中创建对应记录

### Requirement: Update agent info writes session_agents
更新 agent 信息（agent_name、avatar_url）时 SHALL 写入 `session_agents` 表而非 Session 表。

#### Scenario: Update avatar via API
- **WHEN** 调用 `PUT /api/sessions/:sessionId` 更新 `avatar_url`
- **THEN** SHALL 更新 session_agents 表中对应记录的 `avatar_url`

#### Scenario: Update agent name via API
- **WHEN** 调用 `PUT /api/sessions/:sessionId` 更新 `agent_name`
- **THEN** SHALL 更新 session_agents 表中对应记录的 `agent_name`
