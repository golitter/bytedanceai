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
