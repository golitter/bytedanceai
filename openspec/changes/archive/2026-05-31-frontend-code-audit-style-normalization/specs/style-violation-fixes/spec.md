## ADDED Requirements

### Requirement: 消除硬编码颜色值
组件中 SHALL NOT 直接使用 `#fff`、`#ffffff`、`#000`、`#000000` 等硬编码色值。所有颜色 SHALL 通过 CSS 变量（`var(--xxx)`)或 Tailwind 语义类（`text-primary-foreground`）引用。

#### Scenario: AdminPasswordDialog 按钮文字颜色
- **WHEN** AdminPasswordDialog 中的品牌色按钮渲染时
- **THEN** 按钮文字颜色 SHALL 使用 `var(--primary-foreground)` 而非硬编码 `#fff`

#### Scenario: SessionCleanupPage 危险按钮文字颜色
- **WHEN** SessionCleanupPage 中的错误色按钮渲染时
- **THEN** 按钮文字颜色 SHALL 使用 `var(--primary-foreground)` 而非硬编码 `#fff`

#### Scenario: AgentOverviewPage 按钮文字颜色
- **WHEN** AgentOverviewPage 中的品牌色按钮渲染时
- **THEN** 按钮文字颜色 SHALL 使用 `var(--primary-foreground)` 而非硬编码 `#fff`

#### Scenario: 图标 text-white 替换
- **WHEN** AgentProfilePage 或 UserManagementPage 中的 Camera 图标渲染时
- **THEN** 图标颜色 SHALL 使用 `text-primary-foreground` 而非 `text-white`

### Requirement: Lucide 图标 strokeWidth 统一为 1.25
所有 Lucide React 图标组件 SHALL 使用 `strokeWidth={1.25}`，不得使用 1.5、1.6、1.7 等其他值。

#### Scenario: AskAgentCard 图标 strokeWidth
- **WHEN** AskAgentCard 渲染任何 Lucide 图标时
- **THEN** 所有图标 SHALL 使用 `strokeWidth={1.25}`

#### Scenario: MessageBubble 图标 strokeWidth
- **WHEN** MessageBubble 渲染 Copy 图标时
- **THEN** 图标 SHALL 使用 `strokeWidth={1.25}`

#### Scenario: FinalSummaryCard 图标 strokeWidth
- **WHEN** FinalSummaryCard 渲染 CheckCircle2 或 XCircle 图标时
- **THEN** 图标 SHALL 使用 `strokeWidth={1.25}`

#### Scenario: TaskFailureCard 图标 strokeWidth
- **WHEN** TaskFailureCard 渲染 AlertTriangle 图标时
- **THEN** 图标 SHALL 使用 `strokeWidth={1.25}`

### Requirement: 圆角值符合设计规范
组件圆角 SHALL 遵循 visual-style-guide 的圆角体系：按钮 6px、输入框 8px、卡片 10px、面板 12px、头像 8px（圆角方形）、Badge 9999px（胶囊）。不得使用 `rounded-xl`（12px）以上圆角用于非面板元素。

#### Scenario: AgentProfilePage 头像区域圆角
- **WHEN** AgentProfilePage 渲染头像容器时
- **THEN** 容器 SHALL 使用 `rounded-lg`（8px）而非 `rounded-xl`

#### Scenario: UserManagementPage 头像区域圆角
- **WHEN** UserManagementPage 渲染头像容器时
- **THEN** 容器 SHALL 使用 `rounded-lg`（8px）而非 `rounded-xl`

#### Scenario: 头像 overlay 圆角跟随父容器
- **WHEN** 头像上方的 hover overlay 渲染时
- **THEN** overlay 圆角 SHALL 与父容器保持一致（`rounded-lg`）

### Requirement: 阴影使用受限
`shadow-lg`、`shadow-xl`、`shadow-2xl` SHALL 仅用于弹出菜单、模态框、下拉框等明确脱离页面的浮动元素。卡片、列表项、侧边栏等固定元素 SHALL NOT 使用阴影。

#### Scenario: Dialog 组件阴影保留
- **WHEN** shadcn Dialog 组件渲染模态框时
- **THEN** `shadow-lg` SHALL 保留，因为模态框属于弹出元素

#### Scenario: IconSidebar tooltip 阴影保留
- **WHEN** IconSidebar 的 tooltip 渲染时
- **THEN** `shadow-lg` SHALL 保留，因为 tooltip 属于弹出元素
