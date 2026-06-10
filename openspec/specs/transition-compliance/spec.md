## ADDED Requirements

### Requirement: 过渡动画仅限 transform 和 opacity
所有 CSS transition SHALL 仅应用于 `transform` 和 `opacity` 属性。组件 SHALL NOT 使用 `transition-colors`、`transition-all` 或其他动画非合规属性的过渡类。

#### Scenario: 按钮 hover 过渡
- **WHEN** 任何按钮的 hover 状态变化
- **THEN** 使用 `transition-[transform,opacity]` 或无过渡，不使用 `transition-colors`

#### Scenario: 会话列表项 hover 过渡
- **WHEN** ConversationItem 或 ConversationList 的 hover 状态变化
- **THEN** 不使用 `transition-colors duration-120`，如需过渡仅限 `transition-[transform,opacity]`

#### Scenario: 弹出浮层入场动画
- **WHEN** IconSidebar 的 hover card 入场
- **THEN** 过渡仅限 transform/opacity，不使用 `transition-all`

### Requirement: 过渡时长符合规范
过渡时长 SHALL 限制在 120-300ms 范围内，缓动函数使用 `ease-out`。

#### Scenario: hover 微交互
- **WHEN** 组件的 hover 状态过渡
- **THEN** 时长为 120ms ease-out

#### Scenario: 面板出入
- **WHEN** 面板或浮层出入动画
- **THEN** 时长为 250ms ease-out

### Requirement: 弹出菜单阴影统一为规范值
弹出菜单和下拉框的阴影 SHALL 使用 `0 4px 24px rgba(0, 0, 0, 0.4)`，不使用其他阴影值。

#### Scenario: Hover Card 阴影
- **WHEN** IconSidebar 的 hover card 渲染
- **THEN** box-shadow 为 `0 4px 24px rgba(0, 0, 0, 0.4)`，不使用 `0 8px 32px rgba(0,0,0,0.4)`
