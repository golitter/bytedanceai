## MODIFIED Requirements

### Requirement: StreamEvent data model
系统 SHALL 从 `contracts/schemas/` 生成的类型中导入 StreamEvent 定义，SHALL NOT 在 `agentend/src/schemas/events.py` 中手写定义。StreamEvent 的类型枚举 SHALL 为 `text / tool_call / tool_result / artifact / planning / done / error / init`，与 `contracts/schemas/event-types.json` 保持一致。AgentEnd 代码 SHALL 从 `src/generated/events.py` 导入 EventType 和 StreamEvent。

#### Scenario: Create text stream event
- **WHEN** Adapter 解析出文本内容
- **THEN** SHALL 使用生成的 EventType.TEXT 创建 StreamEvent，`content` 包含 `{"text": "..."}`

#### Scenario: Create tool_call stream event
- **WHEN** Adapter 解析出工具调用
- **THEN** SHALL 使用生成的 EventType.TOOL_CALL 创建 StreamEvent，`content` 包含 `{"tool": "...", "args": {...}}`

#### Scenario: Create done stream event
- **WHEN** Adapter 检测到 CLI 输出结束
- **THEN** SHALL 使用生成的 EventType.DONE 创建 StreamEvent，`content` 包含 `{"usage": {...}}`

#### Scenario: Import from generated
- **WHEN** AgentEnd 代码需要使用 EventType 或 StreamEvent
- **THEN** SHALL 从 `src/generated/events.py` 导入，SHALL NOT 从 `src/schemas/events.py` 导入
