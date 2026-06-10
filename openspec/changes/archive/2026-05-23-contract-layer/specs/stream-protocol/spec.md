## MODIFIED Requirements

### Requirement: SSE event format
每个 `StreamEvent` SHALL 以 SSE 格式输出：`event: <type>\ndata: <json>\n\n`。`data` 字段 MUST 为 `StreamEvent` 的 JSON 序列化字符串。`<type>` SHALL 使用 `contracts/schemas/event-types.json` 中定义的枚举值，SHALL NOT 使用硬编码字符串。

#### Scenario: SSE format output
- **WHEN** 生成一个 EventType.TEXT 的 StreamEvent
- **THEN** SSE 输出 SHALL 为 `event: text\ndata: {"type":"text","content":{"text":"..."},"timestamp":1234.5}\n\n`，其中 "text" 来源于生成的 EventType 枚举

### Requirement: StreamEvent type mapping from Claude Code output
系统 SHALL 定义 Claude Code CLI stream-json 输出类型到 StreamEvent 类型的映射规则。映射的目标类型 SHALL 从 `contracts/schemas/event-types.json` 生成的 EventType 枚举中取值。未识别的类型 SHALL 映射为 EventType.TEXT。

#### Scenario: Map known Claude Code output types
- **WHEN** CLI 输出 `{"type": "assistant"}` → 映射为生成的 EventType.TEXT
- **WHEN** CLI 输出 `{"type": "tool_use"}` → 映射为生成的 EventType.TOOL_CALL
- **WHEN** CLI 输出 `{"type": "tool_result"}` → 映射为生成的 EventType.TOOL_RESULT
- **WHEN** CLI 输出 `{"type": "result"}` → 映射为生成的 EventType.DONE

#### Scenario: Map unknown output type
- **WHEN** CLI 输出无法识别的 type 字段
- **THEN** SHALL 映射为生成的 EventType.TEXT，完整 content 保留在 `content.raw` 中
