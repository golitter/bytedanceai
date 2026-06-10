## ADDED Requirements

### Requirement: 次要文字必须使用 text-text-secondary
所有表示"次要文字"语义的 Tailwind 类名 SHALL 使用 `text-text-secondary`，禁止使用 `text-secondary`（后者映射到 shadcn secondary 色阶，暗色背景下不可见）。

#### Scenario: MarkdownRenderer em 标签文字可见
- **WHEN** MarkdownRenderer 渲染 `<em>` 标签
- **THEN** 使用 `text-text-secondary` 类名，文字在暗色背景下清晰可见

#### Scenario: MarkdownRenderer 表头文字可见
- **WHEN** MarkdownRenderer 渲染 `<th>` 表头
- **THEN** 使用 `text-text-secondary` 类名，文字在暗色背景下清晰可见

### Requirement: 背景色必须使用语义 token
组件背景色 SHALL 使用 `bg-bg-card`、`bg-bg-hover`、`bg-bg-active` 等 CSS 变量 token，禁止使用 `bg-secondary` 表示"卡片背景"语义。

#### Scenario: ToolCard 背景色正确
- **WHEN** ToolCard 渲染工具调用卡片
- **THEN** 背景使用 `bg-bg-card` 或 `bg-bg-hover`，与 style guide 色阶一致

#### Scenario: GitGraphPanel 当前分支标识背景正确
- **WHEN** GitGraphPanel 渲染分支标签
- **THEN** 当前分支使用 `bg-primary`，非当前分支使用 `bg-bg-hover`

### Requirement: 禁止纯白文字
CSS 变量和 Tailwind 类名中 SHALL NOT 使用 `#FFFFFF` 纯白色。主要文字 SHALL 使用 `#E8EBF0`（通过 `text-text-primary` 或 `--text-primary` 变量）。

#### Scenario: h1 标题不使用纯白
- **WHEN** Markdown 渲染 h1 标题
- **THEN** 文字颜色为 `--prose-heading`（`#F0F2F7`），不使用 `#FFFFFF`

### Requirement: 边框必须使用半透明白色
暗色模式下边框 SHALL 使用 `rgba(255, 255, 255, 0.06)`（通过 `border-border` 类名或 `--border` 变量），禁止使用实色边框。

#### Scenario: 所有卡片和面板边框统一
- **WHEN** 组件渲染带边框的容器
- **THEN** 边框颜色通过 `border-border` 或 `border-white/6` 类名引用，不硬编码色值
