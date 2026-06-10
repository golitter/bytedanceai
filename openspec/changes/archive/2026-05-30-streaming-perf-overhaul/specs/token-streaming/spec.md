## MODIFIED Requirements

### Requirement: Adapter 逐 token yield TEXT 事件
Adapter SHALL 从 CLI stdout 逐行读取增量文本事件，每收到一个 token 级别的文本片段 SHALL 立即 yield 一个 `StreamEvent(type=TEXT, content={"text": "<incremental>"})` 事件，不等完整消息生成完毕。系统 SHALL 同时支持 `runtime_text` 事件类型，用于 Orchestrator 透传子 Agent 的增量文本。

#### Scenario: Claude CLI content_block_delta 事件
- **WHEN** Claude CLI 输出 `{"type": "stream_event", "event": {"type": "content_block_delta", "delta": {"type": "text_delta", "text": "你好"}}}`
- **THEN** adapter SHALL 立即 yield `StreamEvent(type=TEXT, content={"text": "你好"})`

#### Scenario: 多个 delta 事件连续到达
- **WHEN** CLI 连续输出三个 `content_block_delta` 事件，文本分别为 "你"、 "好"、 "世界"
- **THEN** adapter SHALL 依次 yield 三个独立的 TEXT 事件，每个包含对应的增量文本

#### Scenario: Orchestrator 透传子 Agent token
- **WHEN** ExecutionEngine 收到子 Agent 的 TEXT 事件
- **THEN** SHALL 包装为 `StreamEvent(type=RUNTIME_TEXT, content={"task_id": "...", "agent": "...", "text": "增量文本"})` 立即 yield
