## ADDED Requirements

### Requirement: AgentRequest data model
系统 SHALL 定义 `AgentRequest` Pydantic 模型，包含字段：`task_id`（str）、`conversation_id`（str）、`session_id`（可选 str）、`message`（str）、`agent_type`（str，默认 "claude-code"）、`stream`（bool，默认 True）、`system_prompt`（可选 str）、`rules`（list[str]，默认空）、`workspace_path`（可选 str）、`config`（可选 dict）。

#### Scenario: Parse valid request JSON
- **WHEN** 接收包含 `task_id`、`conversation_id`、`message` 字段的 JSON
- **THEN** SHALL 成功解析为 `AgentRequest` 对象，`agent_type` 默认为 `"claude-code"`，`stream` 默认为 `True`

#### Scenario: Reject request missing required fields
- **WHEN** 接收缺少 `task_id` 或 `message` 字段的 JSON
- **THEN** SHALL 抛出 Pydantic `ValidationError`

### Requirement: AgentResponse data model
系统 SHALL 定义 `AgentResponse` Pydantic 模型，包含字段：`session_id`（str）、`content`（str）、`artifacts`（list[dict]，默认空）、`usage`（dict，默认空）。

#### Scenario: Create response from adapter output
- **WHEN** Adapter 完成一次同步执行
- **THEN** SHALL 将结果封装为 `AgentResponse`，包含 session_id、文本内容、产物列表和 token 使用量

### Requirement: StreamEvent data model
系统 SHALL 定义 `StreamEvent` Pydantic 模型，包含字段：`type`（str，枚举：text / tool_call / tool_result / artifact / done / error）、`content`（dict）、`timestamp`（float）。

#### Scenario: Create text stream event
- **WHEN** Adapter 解析出文本内容
- **THEN** SHALL 创建 `type="text"` 的 StreamEvent，`content` 包含 `{"text": "..."}`

#### Scenario: Create tool_call stream event
- **WHEN** Adapter 解析出工具调用
- **THEN** SHALL 创建 `type="tool_call"` 的 StreamEvent，`content` 包含 `{"tool": "...", "args": {...}}`

#### Scenario: Create done stream event
- **WHEN** Adapter 检测到 CLI 输出结束
- **THEN** SHALL 创建 `type="done"` 的 StreamEvent，`content` 包含 `{"usage": {...}}`

### Requirement: SSE event format
每个 `StreamEvent` SHALL 以 SSE 格式输出：`event: <type>\ndata: <json>\n\n`。`data` 字段 MUST 为 `StreamEvent` 的 JSON 序列化字符串。

#### Scenario: SSE format output
- **WHEN** 生成一个 `type="text"` 的 StreamEvent
- **THEN** SSE 输出 SHALL 为 `event: text\ndata: {"type":"text","content":{"text":"..."},"timestamp":1234.5}\n\n`

### Requirement: StreamEvent type mapping from Claude Code output
系统 SHALL 定义 Claude Code CLI stream-json 输出类型到 StreamEvent 类型的映射规则。未识别的类型 SHALL 映射为 `type="text"`。

#### Scenario: Map known Claude Code output types
- **WHEN** CLI 输出 `{"type": "assistant"}` → 映射为 `StreamEvent.type="text"`
- **WHEN** CLI 输出 `{"type": "tool_use"}` → 映射为 `StreamEvent.type="tool_call"`
- **WHEN** CLI 输出 `{"type": "tool_result"}` → 映射为 `StreamEvent.type="tool_result"`
- **WHEN** CLI 输出 `{"type": "result"}` → 映射为 `StreamEvent.type="done"`

#### Scenario: Map unknown output type
- **WHEN** CLI 输出无法识别的 type 字段
- **THEN** SHALL 映射为 `StreamEvent.type="text"`，完整 content 保留在 `content.raw` 中
