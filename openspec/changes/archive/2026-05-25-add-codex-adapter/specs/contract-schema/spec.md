## MODIFIED Requirements

### Requirement: JSON Schema 契约格式
每个跨端协议定义 SHALL 以 JSON Schema 文件存储在 `contracts/schemas/` 目录中。文件名 SHALL 使用 kebab-case（如 `agent-request.yaml`）。`agent-request.yaml` 中的 `AgentType` 枚举 SHALL 包含 `claude-code`、`opencode`、`orchestrator`、`codex` 四个值。

#### Scenario: AgentType 包含 codex 枚举值
- **WHEN** 查看 `contracts/schemas/agent-request.yaml` 的 AgentType 枚举定义
- **THEN** SHALL 包含枚举值 `codex`，描述为 "Codex CLI"

#### Scenario: 三端生成类型包含 codex
- **WHEN** 运行 `make generate`
- **THEN** 三端的 `generated/request.*` 文件中 AgentType 枚举 SHALL 包含 `codex` 值
