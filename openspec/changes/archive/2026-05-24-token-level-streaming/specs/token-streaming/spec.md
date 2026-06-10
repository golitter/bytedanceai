## ADDED Requirements

### Requirement: Adapter 逐 token yield TEXT 事件
Adapter SHALL 从 CLI stdout 逐行读取增量文本事件，每收到一个 token 级别的文本片段 SHALL 立即 yield 一个 `StreamEvent(type=TEXT, content={"text": "<incremental>"})` 事件，不等完整消息生成完毕。

#### Scenario: Claude CLI content_block_delta 事件
- **WHEN** Claude CLI 输出 `{"type": "stream_event", "event": {"type": "content_block_delta", "delta": {"type": "text_delta", "text": "你好"}}}`
- **THEN** adapter SHALL 立即 yield `StreamEvent(type=TEXT, content={"text": "你好"})`

#### Scenario: 多个 delta 事件连续到达
- **WHEN** CLI 连续输出三个 `content_block_delta` 事件，文本分别为 "你"、 "好"、 "世界"
- **THEN** adapter SHALL 依次 yield 三个独立的 TEXT 事件，每个包含对应的增量文本

### Requirement: 忽略 stream_event 元事件
Adapter SHALL 忽略以下 `stream_event` 子类型：`message_start`、`content_block_start`、`content_block_stop`、`message_delta`、`message_stop`。这些元事件不产生 StreamEvent yield。

#### Scenario: message_start 事件
- **WHEN** CLI 输出 `{"type": "stream_event", "event": {"type": "message_start", ...}}`
- **THEN** adapter SHALL 不 yield 任何事件

#### Scenario: content_block_stop 事件
- **WHEN** CLI 输出 `{"type": "stream_event", "event": {"type": "content_block_stop", ...}}`
- **THEN** adapter SHALL 不 yield 任何事件
