## MODIFIED Requirements

### Requirement: StreamEvent parsing
系统 SHALL 解析每个 SSE `data:` 行为符合 StreamEvent 契约的 JSON 对象（`type`, `content`, `timestamp`）。支持的 event 类型：`init`, `text`, `runtime_text`, `tool_call`, `tool_result`, `artifact`, `planning`, `done`, `error`。前端 SHALL 识别 `runtime_text` 事件并将其归属到对应子 Agent 的消息块中。

#### Scenario: Text event received
- **WHEN** 收到 `{ "type": "text", "content": { "text": "..." } }` event
- **THEN** 追加文本内容到当前 streaming agent message

#### Scenario: Runtime text event received
- **WHEN** 收到 `{ "type": "runtime_text", "content": { "task_id": "task-001", "agent": "claude-code", "text": "增量" } }` event
- **THEN** 追加文本到对应 task_id + agent 标识的子 Agent 消息块，若该消息块不存在 SHALL 自动创建

#### Scenario: Done event received
- **WHEN** 收到 `{ "type": "done" }` event
- **THEN** 完成 streaming agent message，状态转为 `done`

#### Scenario: Error event received
- **WHEN** 收到 `{ "type": "error", "content": { "message": "..." } }` event
- **THEN** 显示错误信息，状态转为 `error`

#### Scenario: Malformed SSE data
- **WHEN** SSE `data:` 行包含无效 JSON
- **THEN** 记录解析错误并跳过该 event，不中断流

#### Scenario: 未知 event type
- **WHEN** 收到未知 type 的 event
- **THEN** 忽略该 event，不中断流，不报错
