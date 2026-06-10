## ADDED Requirements

### Requirement: hover 交互 SHALL 通过 useHoverStyle hook 实现
当多个组件使用相同的 mouseenter/leave backgroundColor 切换模式时，MUST 提取为 `useHoverStyle` hook，不在组件内手写。

#### Scenario: 列表项 hover
- **WHEN** ConversationItem 或 ConversationList 需要实现 hover 背景变色
- **THEN** 使用 `useHoverStyle()` 返回的 `onMouseEnter`/`onMouseLeave`

#### Scenario: NewChatDialog 按钮 hover
- **WHEN** NewChatDialog 中 agent 选择按钮需要 hover 效果
- **THEN** 使用 `useHoverStyle()` 返回的事件处理器

### Requirement: useHoverStyle SHALL 支持自定义色值
`useHoverStyle` MUST 接受可选的 hover 和 normal 背景色参数，默认为 `var(--bg-hover)` 和 `transparent`。

#### Scenario: 自定义 hover 色
- **WHEN** 调用 `useHoverStyle('var(--card)')`
- **THEN** onMouseEnter 设置 `backgroundColor: var(--card)`
