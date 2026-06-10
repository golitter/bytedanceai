## ADDED Requirements

### Requirement: 统一 CSS 变量到 shadcn 语义 token
系统 SHALL 在 `.dark` 块中仅保留 shadcn 语义 token（`--background`、`--foreground`、`--primary`、`--destructive`、`--border` 等）和少量无 shadcn 对应的自定义 token（`--text-tertiary`、`--bg-active`、`--code-bg`、`--color-danger-bg`、`--agent-claude`、`--agent-opencode`、`--agent-orchestrator`）。所有与 shadcn token 重复的自定义 app token（`--bg-canvas`、`--bg-sidebar`、`--bg-hover`、`--text-primary`、`--text-secondary`、`--color-brand`、`--color-error`、`--divider`）SHALL 被移除。

#### Scenario: 重复变量移除
- **WHEN** 检查 `index.css` 的 `.dark` 块
- **THEN** 不存在 `--bg-canvas`、`--bg-sidebar`、`--bg-hover`、`--text-primary`、`--text-secondary`、`--color-brand`、`--color-error`、`--divider` 变量定义
- **AND** 所有 shadcn 语义 token 值与移除前的自定义 token 值一致（如 `--background` = `#0A0B0E`，`--foreground` = `#E8EBF0`，`--primary` = `#6366F1`）

#### Scenario: 自定义 token 保留
- **WHEN** 检查 `index.css` 的 `.dark` 块
- **THEN** `--text-tertiary`、`--bg-active`、`--code-bg`、`--color-danger-bg`、`--agent-claude`、`--agent-opencode`、`--agent-orchestrator` 仍存在
- **AND** 这些变量已添加到 `@theme inline` 块中，可通过 Tailwind 工具类访问

### Requirement: 内联样式迁移到 Tailwind 工具类
所有应用组件中的静态内联 `style={{}}` SHALL 替换为等价的 Tailwind 工具类。动态条件样式可保留内联形式，但 SHALL 使用 CSS 变量而非硬编码颜色值。

#### Scenario: 静态颜色内联样式替换
- **WHEN** 组件使用 `style={{ color: 'var(--text-primary)' }}`
- **THEN** 替换为 Tailwind 类 `text-foreground`

#### Scenario: 静态背景内联样式替换
- **WHEN** 组件使用 `style={{ backgroundColor: 'var(--bg-canvas)' }}`
- **THEN** 替换为 Tailwind 类 `bg-background`

#### Scenario: 边框内联样式替换
- **WHEN** 组件使用 `style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}`
- **THEN** 替换为 Tailwind 类 `border-b border-border`

#### Scenario: 动态样式保留内联
- **WHEN** 组件的内联样式包含动态表达式（如基于 state 的条件颜色）
- **THEN** 保留内联样式形式，但其中的硬编码颜色值 SHALL 替换为 CSS 变量引用

### Requirement: 硬编码颜色值消除
组件中 SHALL 不存在硬编码的颜色字面量（hex、rgb、rgba、hsl）。所有颜色 SHALL 通过 CSS 变量引用。

#### Scenario: 硬编码边框色替换
- **WHEN** 检查所有组件文件
- **THEN** 不存在硬编码的 `rgba(255,255,255,0.06)`，全部使用 `border-border` 或 `var(--border)`

#### Scenario: 硬编码状态色替换
- **WHEN** 检查 `AgentAvatar.tsx`、`NewChatDialog.tsx` 等文件
- **THEN** `#22C55E` 替换为 `var(--color-success)`，`#EF4444` 替换为 `var(--destructive)`，`#F59E0B` 替换为 `var(--color-warning)`，`#5A6070` 替换为 `var(--text-tertiary)`

#### Scenario: 品牌色边框替换
- **WHEN** 检查 `MessageBubble.tsx`
- **THEN** `rgba(99,102,241,0.08)` 和 `rgba(99,102,241,0.15)` 替换为基于 `--primary` 的 CSS 变量

### Requirement: visual-style-guide.md 规范符合
所有样式 SHALL 符合 `visual-style-guide.md` 中的 Dark Utilitarian 设计规范。

#### Scenario: 无违禁阴影
- **WHEN** 检查所有组件的 CSS/Tailwind 类
- **THEN** 不存在 `box-shadow`、`shadow-*`（弹出菜单除外）、`drop-shadow-*`

#### Scenario: 无违禁渐变
- **WHEN** 检查所有组件的样式
- **THEN** 不存在 `bg-gradient-*` 或 `linear-gradient`/`radial-gradient`

#### Scenario: 圆角规范
- **WHEN** 检查所有组件的 `border-radius` / `rounded-*`
- **THEN** 按钮 ≤ 6px，输入框 ≤ 8px，卡片 ≤ 10px，面板 ≤ 12px

#### Scenario: 动画规范
- **WHEN** 检查所有组件的 `transition` / `animate-*`
- **THEN** 持续时间在 120-300ms，仅使用 `transform` 和 `opacity`，timing function 为 `ease-out`
