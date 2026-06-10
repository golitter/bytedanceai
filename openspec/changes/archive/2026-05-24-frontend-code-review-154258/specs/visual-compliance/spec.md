## MODIFIED Requirements

### Requirement: AgentAvatar SHALL 显示同色模糊光晕
Agent 头像 MUST 拥有 8px 同色 blur 的 `box-shadow` 光晕效果，颜色 MUST 通过 CSS 变量引用该 Agent 的标识色，不得硬编码 hex 值。

#### Scenario: 头像显示光晕
- **WHEN** 渲染 AgentAvatar
- **THEN** 元素具有 `box-shadow` 且 blur 值为 8px，颜色通过 `var(--agent-xxx)` 引用

#### Scenario: 无硬编码 hex
- **WHEN** 检查 AgentAvatar 源码
- **THEN** 不存在 `#DA7756`、`#10B981`、`#EAB308` 等硬编码色值

## ADDED Requirements

### Requirement: AgentAvatar 状态点尺寸 SHALL 为 4px
Agent 头像右下角状态灯 MUST 为 4px 圆点，不得为 8px 或其他尺寸。

#### Scenario: 状态点尺寸正确
- **WHEN** 渲染 AgentAvatar 的状态灯
- **THEN** 状态灯宽高为 4px

### Requirement: 组件 SHALL NOT 使用 rounded-full 在按钮或卡片上
按钮和卡片 MUST NOT 使用 `rounded-full`（胶囊形仅限于 Badge/标签），滚动按钮等 MUST 使用 ≤12px 的圆角。

#### Scenario: 滚动按钮圆角
- **WHEN** 渲染 "滚动到底部" 按钮
- **THEN** 圆角为 10px 或 12px，非 `rounded-full`
