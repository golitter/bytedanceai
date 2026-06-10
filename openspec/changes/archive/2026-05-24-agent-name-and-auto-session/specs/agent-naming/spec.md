## ADDED Requirements

### Requirement: Session SHALL support agent_name
Session 模型 MUST 包含 `agent_name` 字段，用于存储用户自定义的 agent 显示名称。

#### Scenario: 创建 session 时指定 agent_name
- **WHEN** 前端发送 `POST /api/tasks { agents: [{ type: "claude-code", name: "代码审查助手" }] }`
- **THEN** 创建的 session 的 `agent_name` 为 `"代码审查助手"`

#### Scenario: agent_name 未指定时 fallback
- **WHEN** 前端发送 `POST /api/tasks { agents: [{ type: "claude-code" }] }` 不含 name
- **THEN** 创建的 session 的 `agent_name` 为空，前端 fallback 显示 `agent_type`

### Requirement: 前端 SHALL 显示 agent_name
对话列表和聊天区 MUST 优先显示 `agent_name`，当 `agent_name` 为空时 fallback 到 `agent_type`。

#### Scenario: ConversationItem 显示自定义名称
- **WHEN** 渲染 ConversationItem 且 session 有 `agent_name`
- **THEN** 显示 `agent_name` 而非 `agent_type`

#### Scenario: ChatArea header 显示自定义名称
- **WHEN** 渲染 ChatArea header 且 session 有 `agent_name`
- **THEN** 显示 `agent_name`

### Requirement: agent_name 不穿透到 agentend
后端调用 agentend 时 MUST 只传 `agent_type`，不传 `agent_name`。

#### Scenario: RunTask 不传 agent_name
- **WHEN** 后端调用 agentend 的 stream 接口
- **THEN** 请求中只有 `agent_type` 字段，没有 `agent_name`
