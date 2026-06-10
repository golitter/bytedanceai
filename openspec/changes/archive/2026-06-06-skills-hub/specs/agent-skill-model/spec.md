## MODIFIED Requirements

### Requirement: Agent Profile API
系统 SHALL 提供 `GET /api/sessions/:sid/profile` 接口，返回 Hover 卡片所需的 Agent 摘要数据。

#### Scenario: 获取 profile 数据
- **WHEN** 前端调用 `GET /api/sessions/:sid/profile`
- **THEN** 返回 JSON 包含 `agent_name`、`agent_type`、`avatar_url`、`status`、`session_id`、`skills` 数组（每项含 `name`、`description`、`builtin`、`source`）。`source` 字段为 `"builtin"` 或 `"hub"`

#### Scenario: Session 不存在
- **WHEN** 前端调用不存在的 session_id
- **THEN** 返回 404

### Requirement: Agent Detail API
系统 SHALL 提供 `GET /api/sessions/:sid/detail` 接口，返回详情页所需的完整 Agent 数据。Skills 数据 SHALL 从 Agentend 实时读取（`GET /api/v1/skills/:agent_type?workspace_path=xxx`），不再使用硬编码 mock。

#### Scenario: 获取 detail 数据
- **WHEN** 前端调用 `GET /api/sessions/:sid/detail`
- **THEN** 返回 JSON 包含 profile 所有字段 + `task_id`、`workspace_path`（如有）、`created_at`、`message_count`。`skills` 数组从 Agentend 实时获取，每项含 `name`、`description`、`builtin`、`source`

#### Scenario: Agentend unreachable fallback
- **WHEN** Agentend 不可达
- **THEN** Backend 返回 profile 数据但 `skills` 为空数组，不阻塞页面渲染

#### Scenario: Session 不存在
- **WHEN** 前端调用不存在的 session_id
- **THEN** 返回 404

### Requirement: Skills 数据从 Agentend 实时读取
Skills 数据 SHALL 从 Agentend `GET /api/v1/skills/:agent_type` 实时读取，区分 builtin/external，增加 `source` 字段。不再使用硬编码 mock 数据。

#### Scenario: Builtin skill data
- **WHEN** Agentend 返回 `{name: "render", builtin: true, source: "builtin"}`
- **THEN** 前端展示绿色 builtin 标签

#### Scenario: External skill data
- **WHEN** Agentend 返回 `{name: "my-skill", builtin: false, source: "hub"}`
- **THEN** 前端展示 indigo external 标签

#### Scenario: No skills
- **WHEN** Agentend 返回空数组
- **THEN** 前端显示空状态提示
