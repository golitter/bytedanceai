## ADDED Requirements

### Requirement: 组件文件不超过 300 行
所有组件文件 SHALL 控制在 300 行以内。超过 200 行的文件 SHALL 被标记为需要拆分审查，超过 300 行的文件 MUST 拆分。

#### Scenario: 新组件文件行数检查
- **WHEN** 审查所有 `components/` 和 `pages/` 下的 .tsx 文件
- **THEN** 无文件超过 300 行

### Requirement: NewChatDialog 拆分
NewChatDialog 组件（当前 331 行）SHALL 拆分为：
1. `AgentSelectList` 组件 — 负责 Agent 列表渲染和选择逻辑
2. `NewChatDialog` 组件 — 保留对话框结构和布局，引用 `AgentSelectList`

#### Scenario: NewChatDialog 拆分后行数
- **WHEN** 检查拆分后的 NewChatDialog 文件
- **THEN** 行数不超过 200 行

#### Scenario: AgentSelectList 独立导出
- **WHEN** 其他组件需要 Agent 选择列表
- **THEN** 可以直接从 `components/im/AgentSelectList` 导入使用

### Requirement: MessageList 拆分
MessageList 组件（当前 282 行）SHALL 将消息渲染逻辑抽为独立组件，虚拟列表配置和消息处理逻辑保留在 MessageList 中。

#### Scenario: MessageList 拆分后行数
- **WHEN** 检查拆分后的 MessageList 文件
- **THEN** 行数不超过 200 行
