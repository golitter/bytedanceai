## Requirements

### Requirement: GroupChatRule 内置规则
系统 MUST 新增 `GroupChatRule` 内置规则，继承 `BaseRule`，字段为 `name="group_chat"`、`description="Injects group chat context from other agents"`、`phase="pre"`、`priority=6`。

#### Scenario: GroupChatRule 注册到 Registry
- **WHEN** 系统启动并加载内置 Rule
- **THEN** RuleRegistry 中 SHALL 包含 GroupChatRule 实例

### Requirement: GroupChatRule check 条件
GroupChatRule 的 `check` 方法 SHALL 在 `context["group_chat_messages"]` 非空列表时返回 `True`，否则返回 `False`。

#### Scenario: 有群聊消息
- **WHEN** context 包含 `group_chat_messages: [{role: "agent", agent_name: "claude", content: "..."}]`
- **THEN** `check` SHALL 返回 `True`

#### Scenario: 无群聊消息（单聊或首次执行）
- **WHEN** context 包含 `group_chat_messages: []` 或不包含该 key
- **THEN** `check` SHALL 返回 `False`

### Requirement: GroupChatRule enforce 输出
GroupChatRule 的 `enforce` 方法 SHALL 调用 `build_group_chat_context` 格式化消息列表，返回 `{"system_prompt_append": <格式化文本>}`。空结果 SHALL 返回 `{}`。

#### Scenario: 格式化多条消息
- **WHEN** group_chat_messages 包含 2 条 agent 消息
- **THEN** `enforce` SHALL 返回 `{"system_prompt_append": "## 群聊上下文\n..."}` 包含格式化后的消息

#### Scenario: 格式化后为空（无有效内容）
- **WHEN** group_chat_messages 的所有消息 role 既非 "user" 也非 "agent"
- **THEN** `enforce` SHALL 返回 `{}`

### Requirement: 群聊 Prompt 模板
系统 MUST 提供 `build_group_chat_context` 函数，接收 `cross_round_messages` 列表，返回格式化字符串。每条消息 SHALL 标注来源 Agent 名称（agent 消息用 `🤖 {name}:`，user 消息用 `👤 用户:`）。

#### Scenario: 混合 user 和 agent 消息
- **WHEN** 消息列表包含 `[{role: "user", content: "帮我改"}, {role: "agent", agent_name: "claude", content: "已改完"}]`
- **THEN** 输出 SHALL 包含 `👤 用户:\n帮我改` 和 `🤖 claude:\n已改完`

#### Scenario: 空消息列表
- **WHEN** `cross_round_messages` 为 `None` 或空列表
- **THEN** 函数 SHALL 返回空字符串 `""`

### Requirement: API 层解析 group_chat_messages
AgentEnd 的 `/v1/agent/stream` endpoint SHALL 从请求体中提取 `group_chat_messages` 字段，放入 rule context 中。

#### Scenario: 请求体包含 group_chat_messages
- **WHEN** 请求体包含 `group_chat_messages: [{...}]`
- **THEN** rule context SHALL 包含 `group_chat_messages` 字段，GroupChatRule 可正常 check/enforce

#### Scenario: 请求体不包含 group_chat_messages
- **WHEN** 请求体不包含该字段（单聊场景）
- **THEN** rule context 中 `group_chat_messages` SHALL 默认为 `[]`
