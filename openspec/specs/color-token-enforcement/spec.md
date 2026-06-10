## ADDED Requirements

### Requirement: 消除 opacity 调色模式
组件 SHALL NOT 使用 `text-foreground/N`、`text-primary/N`、`bg-foreground/N` 等 opacity 调色语法。文字和背景色 MUST 使用定义的语义 token（`text-secondary`、`text-tertiary`、`text-muted-foreground`）或精确 CSS 变量。

#### Scenario: AgentMeta 描述文字
- **WHEN** AgentMeta 组件渲染描述文字
- **THEN** 文字颜色使用 `text-secondary`（#8B91A0）或 `text-tertiary`（#5A6070），不使用 `text-foreground/50`

#### Scenario: SkillCard 描述文字
- **WHEN** SkillCard 渲染技能描述
- **THEN** 文字颜色使用 `text-secondary`，不使用 `text-foreground/75`

#### Scenario: SkillCard 元信息
- **WHEN** SkillCard 渲染来源信息
- **THEN** 文字颜色使用 `text-tertiary`，不使用 `text-foreground/55`

### Requirement: 品牌色仅用于合规场景
品牌色（`#6366F1` / `--color-brand` / `text-primary`）SHALL 仅出现在以下场景：选中态、发送按钮、当前会话标记、焦点环。其他场景（streaming 光标、状态文字、装饰元素）MUST 使用语义色。

#### Scenario: Streaming 光标
- **WHEN** MessageBubble 渲染 streaming 中的闪烁光标
- **THEN** 光标颜色使用 `text-foreground` 或 `text-secondary`，不使用品牌色（`text-primary`）

#### Scenario: Streaming 状态文字
- **WHEN** ChatArea 显示 streaming 状态指示
- **THEN** 状态文字使用 `text-tertiary`，不使用品牌色（`text-primary`）

### Requirement: AskAgentCard 使用中性色
AskAgentCard 的边框和背景 SHALL 使用中性色 token（`border`、`bg-card`、`bg-hover`），不使用品牌色 tint。

#### Scenario: AskAgentCard 默认态
- **WHEN** AskAgentCard 渲染默认状态
- **THEN** 边框使用 `border`（rgba(255,255,255,0.06)），背景使用 `bg-card`，不使用 `border-primary/20` 或 `bg-primary/[0.03]`

#### Scenario: AskAgentCard hover 态
- **WHEN** AskAgentCard 被 hover
- **THEN** 背景变为 `bg-hover`（#22262F），不使用 `bg-primary/[0.06]`

### Requirement: 状态 Badge 使用语义色
Agent 状态 Badge SHALL 使用语义色 token（`text-success`、`bg-success/10`），不使用原始 Tailwind 色（`text-emerald-400`、`bg-emerald-500/10`）。

#### Scenario: Agent 就绪状态
- **WHEN** AskAgentCard 显示就绪状态 Badge
- **THEN** 使用 `text-success` 和 `bg-success/10`，不使用 `text-emerald-400` 和 `bg-emerald-500/10`

### Requirement: 字号对齐规范体系
所有文字元素 SHALL 使用 visual-style-guide.md 定义的字号档位（11/12/13/14/20px），不使用未定义的字号（9px、10px）。

#### Scenario: 侧边栏导航标签
- **WHEN** IconSidebar 的 NavItem label 渲染
- **THEN** 字号至少为 11px（微标注），不使用 9px

#### Scenario: 状态 Badge 文字
- **WHEN** AskAgentCard 状态 Badge 渲染
- **THEN** 字号为 11px（`text-[11px]`），不使用 10px（`text-[10px]`）
