## ADDED Requirements

### Requirement: Agent 常量 SHALL 集中管理
`AGENT_NAMES` 和 `AGENT_DESCRIPTIONS` MUST 定义在 `lib/constants.ts` 中，所有消费组件从此模块导入。

#### Scenario: ChatArea 使用集中常量
- **WHEN** ChatArea 需要显示 agent 名称
- **THEN** 从 `lib/constants.ts` 导入 `AGENT_NAMES`，不内联定义

#### Scenario: ConversationItem 使用集中常量
- **WHEN** ConversationItem 需要显示 agent 名称
- **THEN** 从 `lib/constants.ts` 导入 `AGENT_NAMES`，不内联定义

#### Scenario: NewChatDialog 使用集中常量
- **WHEN** NewChatDialog 需要 agent 描述文案
- **THEN** 从 `lib/constants.ts` 导入 `AGENT_DESCRIPTIONS`，不内联定义

### Requirement: createConversation SHALL 防御空 sessions
`createConversation` MUST 在访问 `detail.sessions[0]` 前检查数组非空，空时抛出明确错误。

#### Scenario: 后端未返回 session
- **WHEN** createTask + agents 调用成功但 sessions 为空
- **THEN** 抛出 Error('Backend failed to create session')，不抛 TypeError

### Requirement: Session.agent_name SHALL 为可选字段
`Session` 接口的 `agent_name` MUST 标记为可选（`agent_name?: string`），兼容旧 session 无此字段的场景。

#### Scenario: 旧 session 无 agent_name
- **WHEN** fetchConversations 获取到旧 session（无 agent_name 字段）
- **THEN** `agent_name` 为 undefined，映射到 `agentName: ''`

### Requirement: SSE 事件类型 SHALL 使用契约生成常量
`use-chat-stream.ts` 中的事件类型字符串比较（如 `event === 'token'`）MUST 使用从 `generated/` 导出的常量或 `lib/constants.ts` 中定义的枚举值，不得使用魔法字符串。

#### Scenario: token 事件使用常量
- **WHEN** use-chat-stream 处理 SSE 事件
- **THEN** 事件类型比较使用常量（如 `EventTypeValues.Text`），非字符串字面量 `'text'`

#### Scenario: 事件类型集中定义
- **WHEN** 新增一种 SSE 事件类型
- **THEN** 只需在常量文件中添加一个值，所有消费方自动同步

### Requirement: MessageList 非虚拟化路径 SHALL 传递完整 props
`MessageList` 中非虚拟化渲染路径的 `MessageRenderer` MUST 接收与虚拟化路径相同的 props（包括 `avatarUrl`、`agentName`），确保两种路径行为一致。

#### Scenario: 非虚拟化路径有 agent 信息
- **WHEN** 消息数量 ≤ 虚拟化阈值，走非虚拟化渲染
- **THEN** MessageRenderer 仍接收 `avatarUrl` 和 `agentName` props
