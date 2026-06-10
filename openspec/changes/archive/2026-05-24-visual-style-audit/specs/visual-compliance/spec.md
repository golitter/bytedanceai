## ADDED Requirements

### Requirement: Agent 标识色 SHALL 匹配规范
所有 Agent 标识色 MUST 与 visual-style-guide.md 定义一致：Claude `#DA7756`、OpenCode `#10B981`、Orchestrator `#EAB308`。

#### Scenario: CSS 变量值正确
- **WHEN** 渲染 `.dark` 主题
- **THEN** `--agent-claude` 为 `#DA7756`、`--agent-opencode` 为 `#10B981`、`--agent-orchestrator` 为 `#EAB308`

#### Scenario: AgentAvatar 使用正确的标识色
- **WHEN** 渲染 AgentAvatar 且 agentType 为 `claude-code`
- **THEN** 头像背景色为 `#DA7756`

### Requirement: AgentAvatar SHALL 显示同色模糊光晕
Agent 头像 MUST 拥有 8px 同色 blur 的 `box-shadow` 光晕效果，颜色为该 Agent 的标识色。

#### Scenario: 头像显示光晕
- **WHEN** 渲染 AgentAvatar
- **THEN** 元素具有 `box-shadow` 且 blur 值为 8px，颜色为该 Agent 标识色

### Requirement: 状态灯动效 SHALL 匹配规范
Agent 头像右下角状态灯的动效 MUST 符合规范：`ready` 状态为 opacity 脉冲（2s 循环，0.6→1→0.6），`running` 状态为旋转（1.5s 循环）。

#### Scenario: ready 状态灯脉冲
- **WHEN** AgentAvatar status 为 `ready`
- **THEN** 状态灯具有 opacity 脉冲动画，周期 2s

#### Scenario: running 状态灯旋转
- **WHEN** AgentAvatar status 为 `running`
- **THEN** 状态灯具有旋转动画，周期 1.5s

### Requirement: UI 文字 SHALL 使用 -0.01em letter-spacing
所有 UI 文字 MUST 应用 `-0.01em` 的 letter-spacing，代码文字保持 `0`。

#### Scenario: 正文 letter-spacing
- **WHEN** 渲染页面文字
- **THEN** body 应用了 `letter-spacing: -0.01em`

#### Scenario: 代码块 letter-spacing 为 0
- **WHEN** 渲染代码块
- **THEN** 代码块文字 letter-spacing 为 `0`

### Requirement: 输入框背景 SHALL 使用 card 色
MessageInput 的 textarea 背景 MUST 为 `var(--bg-card)` (`#1A1D24`)，而非 `var(--bg-hover)`。

#### Scenario: 输入框背景色
- **WHEN** 渲染 MessageInput
- **THEN** textarea 背景色为 `var(--bg-card)`

### Requirement: hover 过渡时长 SHALL 为 120ms
所有交互元素的 hover 背景色过渡 MUST 为 120ms ease-out。

#### Scenario: 列表项 hover 过渡
- **WHEN** 用户 hover 侧栏列表项
- **THEN** 背景色过渡时长为 120ms，缓动函数为 ease-out

### Requirement: 分割线颜色 SHALL 比边框更淡
分割线 MUST 使用 `rgba(255, 255, 255, 0.04)`，区别于边框的 `rgba(255, 255, 255, 0.06)`。

#### Scenario: header 底部分割线
- **WHEN** 渲染侧栏 header 底部分割线
- **THEN** 边框颜色为 `rgba(255, 255, 255, 0.04)`

### Requirement: 图标 strokeWidth SHALL 接近 Phosphor Light 风格
在使用 Lucide 图标的场景中，`strokeWidth` MUST 设为 `1.25` 以模拟 Phosphor Light 的细线视觉。

#### Scenario: 图标线宽
- **WHEN** 渲染任何 Lucide 图标
- **THEN** strokeWidth 为 `1.25`
