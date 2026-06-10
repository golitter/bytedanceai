## ADDED Requirements

### Requirement: SSE 事件类型 SHALL 使用契约生成常量
`use-chat-stream.ts` 中的事件类型字符串比较（如 `event === 'token'`）MUST 使用从 `generated/` 导出的常量或 `lib/constants.ts` 中定义的枚举值，不得使用魔法字符串。

#### Scenario: token 事件使用常量
- **WHEN** use-chat-stream 处理 SSE 事件
- **THEN** 事件类型比较使用常量（如 `EventType.TOKEN`），非字符串字面量 `'token'`

#### Scenario: 事件类型集中定义
- **WHEN** 新增一种 SSE 事件类型
- **THEN** 只需在常量文件中添加一个值，所有消费方自动同步

### Requirement: MessageList 非虚拟化路径 SHALL 传递完整 props
`MessageList` 中非虚拟化渲染路径的 `MessageRenderer` MUST 接收与虚拟化路径相同的 props（包括 `avatarUrl`、`agentName`），确保两种路径行为一致。

#### Scenario: 非虚拟化路径有 agent 信息
- **WHEN** 消息数量 ≤ 虚拟化阈值，走非虚拟化渲染
- **THEN** MessageRenderer 仍接收 `avatarUrl` 和 `agentName` props
