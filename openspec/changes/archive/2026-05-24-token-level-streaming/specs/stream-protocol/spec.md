## MODIFIED Requirements

### Requirement: StreamEvent type mapping from Claude Code output
系统 SHALL 定义 Claude Code CLI stream-json 输出类型到 StreamEvent 类型的映射规则。映射的目标类型 SHALL 从 `contracts/schemas/event-types.json` 生成的 EventType 枚举中取值。未识别的类型 SHALL 映射为空事件（不 yield）。

映射规则：
- `{"type": "system"}` → EventType.INIT
- `{"type": "stream_event", "event": {"type": "content_block_delta", "delta": {"type": "text_delta"}}}` → EventType.TEXT（取 `event.delta.text`）
- `{"type": "stream_event", "event": {"type": "message_start" | "content_block_start" | "content_block_stop" | "message_delta" | "message_stop"}}` → 忽略（不 yield）
- `{"type": "assistant"}` → EventType.TEXT（兼容无 `--include-partial-messages` 的旧模式）
- `{"type": "tool_use"}` → EventType.TOOL_CALL
- `{"type": "tool_result"}` → EventType.TOOL_RESULT
- `{"type": "result"}` → EventType.DONE
- 其他未知 type → 忽略（不 yield）

#### Scenario: Map stream_event content_block_delta to TEXT
- **WHEN** CLI 输出 `{"type": "stream_event", "event": {"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": "你好"}}}`
- **THEN** SHALL 映射为 EventType.TEXT，content 为 `{"text": "你好"}`

#### Scenario: Map stream_event meta events to ignore
- **WHEN** CLI 输出 `{"type": "stream_event", "event": {"type": "message_start", ...}}`
- **THEN** SHALL 不 yield 任何事件

#### Scenario: Map known Claude Code output types (unchanged)
- **WHEN** CLI 输出 `{"type": "system"}` → 映射为 EventType.INIT
- **WHEN** CLI 输出 `{"type": "tool_use"}` → 映射为 EventType.TOOL_CALL
- **WHEN** CLI 输出 `{"type": "tool_result"}` → 映射为 EventType.TOOL_RESULT
- **WHEN** CLI 输出 `{"type": "result"}` → 映射为 EventType.DONE

#### Scenario: Map unknown output type to ignore
- **WHEN** CLI 输出无法识别的 type 字段
- **THEN** SHALL 不 yield 任何事件
