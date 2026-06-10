## ADDED Requirements

### Requirement: Agent Profile API
系统 SHALL 提供 `GET /api/sessions/:sid/profile` 接口，返回 Hover 卡片所需的 Agent 摘要数据。

#### Scenario: 获取 profile 数据
- **WHEN** 前端调用 `GET /api/sessions/:sid/profile`
- **THEN** 返回 JSON 包含 `agent_name`、`agent_type`、`avatar_url`、`status`、`session_id`、`skills` 数组（每项含 `name`、`description`、`builtin`）

#### Scenario: Session 不存在
- **WHEN** 前端调用不存在的 session_id
- **THEN** 返回 404

### Requirement: Agent Detail API
系统 SHALL 提供 `GET /api/sessions/:sid/detail` 接口，返回详情页所需的完整 Agent 数据。

#### Scenario: 获取 detail 数据
- **WHEN** 前端调用 `GET /api/sessions/:sid/detail`
- **THEN** 返回 JSON 包含 profile 所有字段 + `task_id`、`workspace_path`（如有）、`created_at`、`message_count`

#### Scenario: Session 不存在
- **WHEN** 前端调用不存在的 session_id
- **THEN** 返回 404

### Requirement: Skills 数据 Mock
两个 API 中的 Skills 数据 SHALL 使用后端硬编码的 mock 数据，不从数据库读取。

#### Scenario: Mock Skills 返回
- **WHEN** 任意 Agent 的 profile/detail 被请求
- **THEN** Skills 数组返回硬编码的 taskctl 和 render 两条数据，`builtin` 字段为 `true`
