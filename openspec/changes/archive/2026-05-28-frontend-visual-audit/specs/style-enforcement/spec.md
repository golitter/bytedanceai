## ADDED Requirements

### Requirement: No pure white text
所有组件 SHALL NOT 使用 `#FFFFFF`、`#fff`、`text-white` 或任何纯白色文字。文字颜色 MUST 使用 CSS 变量 `var(--text-primary)`（对应 Tailwind 的 `text-primary-foreground` 或 `text-foreground`）。

#### Scenario: AgentAvatar initials text color
- **WHEN** AgentAvatar 渲染 Agent 名称首字母
- **THEN** 文字颜色为 `text-foreground`，不是 `text-white`

#### Scenario: GroupAvatar initials text color
- **WHEN** GroupAvatar 渲染多 Agent 头像中的首字母
- **THEN** 文字颜色为 `text-foreground`，不是 `text-white`

### Requirement: No shadows outside popups
`box-shadow` SHALL ONLY 用于弹出菜单、下拉框、Popover 等脱离页面的浮层组件。头像、卡片、列表项、按钮 MUST NOT 使用 box-shadow。

#### Scenario: AgentAvatar no glow effect
- **WHEN** AgentAvatar 渲染 Agent 头像
- **THEN** 无 `box-shadow` 发光效果；如需视觉层次，MUST 通过背景色差实现（如同色低透明度背景圈）

### Requirement: Agent identity colors scope limited
Agent 标识色（`--agent-claude`、`--agent-opencode`、`--agent-orchestrator`、`--agent-codex`）SHALL ONLY 用于：
1. Agent 头像背景
2. 消息左侧 3px 色条

其他 UI 元素（badge、标签、面板背景）MUST NOT 直接使用 Agent 标识色。

#### Scenario: PlanCard agent badge neutral colors
- **WHEN** PlanCard 渲染 task agent badge
- **THEN** badge 使用 `text-secondary` 文字色 + `bg-card` 或 `bg-accent` 背景，不使用 AGENT_COLORS

### Requirement: Animation timing uses ease-out
所有 CSS 动画和过渡的缓动函数 MUST 为 `ease-out`。SHALL NOT 使用 `ease-in-out`、`ease-in`、`linear`、弹簧物理或 bounce 缓动。

#### Scenario: AgentAvatar status pulse ease-out
- **WHEN** AgentAvatar ready 状态灯脉冲动画播放
- **THEN** 动画缓动为 `ease-out`，不是 `ease-in-out`

### Requirement: Font sizes within defined scale
组件字号 MUST 限定在以下档位：11px（微标注）、12px（辅助文字）、13px（代码）、14px（正文/区块标题）、20px（页面标题）。SHALL NOT 使用 8px、9px、10px 等未定义字号。

#### Scenario: GroupAvatar initials font size
- **WHEN** GroupAvatar 渲染首字母
- **THEN** 字号为 `text-[11px]`（微标注档位），不使用 `text-[9px]`

#### Scenario: GroupAvatar count badge font size
- **WHEN** GroupAvatar 渲染成员数量角标
- **THEN** 字号为 `text-[11px]`（微标注档位），不使用 `text-[8px]`

### Requirement: Prefer Tailwind classes over inline styles
当 Tailwind 工具类能表达样式时，MUST 使用 Tailwind 类而非内联 `style` 属性。仅在 Tailwind 无法覆盖的动态值（如运行时计算的颜色）时允许内联样式。

#### Scenario: AgentAvatar size via Tailwind
- **WHEN** AgentAvatar 渲染不同尺寸的头像
- **THEN** 宽高通过 Tailwind 尺寸类（如 `h-8 w-8`、`h-12 w-12`）控制，不使用内联 `width`/`height`

### Requirement: Border-radius consistent with spec
圆角 MUST 遵循以下规则：
- 按钮 6px
- 输入框 8px
- 卡片/消息气泡 10px
- 面板 12px
- Agent 头像 8px
- 用户头像 50%
- Badge/胶囊标签 9999px
- SHALL NOT 使用 >12px 圆角在非面板元素上
