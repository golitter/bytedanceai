## MODIFIED Requirements

### Requirement: 返回可用 Agent 类型列表
系统 SHALL 提供 `GET /api/agent-types` 端点，返回当前支持的 Agent 类型列表。

#### Scenario: 获取 Agent 类型列表
- **WHEN** 发送 `GET /api/agent-types`
- **THEN** 返回 HTTP 200，data 为 `["claude-code", "opencode", "orchestrator"]`

#### Scenario: 列表内容与契约一致
- **WHEN** 返回 Agent 类型列表
- **THEN** 列表中的值 SHALL 与 `contracts/schemas/agent-request.yaml` 中定义的 `AgentType` 枚举值一致
