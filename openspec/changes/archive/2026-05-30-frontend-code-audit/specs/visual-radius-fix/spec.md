## ADDED Requirements

### Requirement: 按钮 border-radius 统一为 6px
所有按钮元素（包括 icon button、send button、sidebar nav item）SHALL 使用 `rounded-[6px]`，不使用 `rounded-md`、`rounded-lg` 或其他值。

#### Scenario: 发送按钮圆角
- **WHEN** MessageInput 的发送按钮渲染
- **THEN** border-radius 为 6px（`rounded-[6px]`）

#### Scenario: 侧边栏导航按钮圆角
- **WHEN** IconSidebar 的 NavItem 渲染
- **THEN** border-radius 为 6px（`rounded-[6px]`），不使用 `rounded-[10px]`

#### Scenario: 新建会话按钮圆角
- **WHEN** ConversationList 的新建会话按钮渲染
- **THEN** border-radius 为 6px（`rounded-[6px]`），不使用 `rounded-md`

### Requirement: 输入框 border-radius 统一为 8px
所有输入框和 textarea 元素 SHALL 使用 `rounded-[8px]`。

#### Scenario: 消息输入框圆角
- **WHEN** MessageInput 的 textarea 渲染
- **THEN** border-radius 为 8px（`rounded-[8px]`），不使用 `rounded-lg`

#### Scenario: 搜索输入框圆角
- **WHEN** ConversationList 的搜索输入框渲染
- **THEN** border-radius 为 8px（`rounded-[8px]`），不使用 `rounded-lg`

### Requirement: 面板 border-radius 统一为 12px
所有面板和弹出浮层 SHALL 使用 `rounded-[12px]`，不使用 `rounded-xl`（14px）或更大值。

#### Scenario: Hover Card 圆角
- **WHEN** IconSidebar 的 hover card 渲染
- **THEN** border-radius 为 12px（`rounded-[12px]`），不使用 `rounded-xl`

### Requirement: Badge/标签 border-radius 统一为 9999px
所有 Badge 和状态标签 SHALL 使用 `rounded-full`（即 `rounded-[9999px]`）实现胶囊形状。

#### Scenario: Agent 状态 Badge
- **WHEN** AskAgentCard 中的状态 Badge 渲染
- **THEN** border-radius 为 9999px（`rounded-full`），不使用 `rounded-md`

### Requirement: 用户头像 border-radius 统一为 50%
所有用户头像 SHALL 使用 `rounded-full`（即 50% 圆形）。

#### Scenario: 消息中用户头像
- **WHEN** MessageBubble 中的用户头像渲染
- **THEN** border-radius 为 50%（`rounded-full`），不使用 `rounded-lg`

#### Scenario: 侧边栏用户头像
- **WHEN** IconSidebar 底部的用户头像渲染
- **THEN** border-radius 为 50%（`rounded-full`），不使用 `rounded-[10px]`

### Requirement: Agent 头像保持 8px 方形圆角
Agent 头像 SHALL 使用 `rounded-[8px]` 或等效的 8px 方形圆角。

#### Scenario: AgentAvatar 圆角
- **WHEN** AgentAvatar 组件渲染
- **THEN** border-radius 为 8px，形状为圆角方形
