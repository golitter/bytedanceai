## ADDED Requirements

### Requirement: Agent list display
Agent 概览页面 SHALL 展示所有已注册 Agent 的只读信息卡片。

#### Scenario: Display agent cards
- **WHEN** 用户打开 Agent 概览页面
- **THEN** 显示 4 个 Agent 卡片（Claude Code、OpenCode、Codex、Orchestrator），每个包含名称、类型标识、描述、配置文件路径

### Requirement: CLI config file display
用户 SHALL 能展开查看每个 Agent 的 CLI 配置文件内容。

#### Scenario: Toggle config display
- **WHEN** 用户点击某个 Agent 卡片的"查看配置"按钮
- **THEN** 展开显示该 Agent 配置文件的完整内容（代码块样式），按钮文本变为"收起配置"

#### Scenario: Config content is sanitized
- **WHEN** 配置文件内容包含敏感信息（API Key、Token、Secret）
- **THEN** 敏感字段值 SHALL 被替换为 `***`，不暴露真实密钥

### Requirement: Agent config API
后端 SHALL 提供 `GET /api/admin/agents` 接口，返回 Agent 列表及配置文件内容。

#### Scenario: Fetch agents with config
- **WHEN** 前端请求 `GET /api/admin/agents`
- **THEN** 返回 Agent 列表，每项包含 type、name、description、configDir、configContent（已脱敏）
