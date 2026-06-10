## ADDED Requirements

### Requirement: Light theme CSS variables
系统 SHALL 在 `:root` 块中定义完整的 Light 主题 CSS 变量集，覆盖所有 `.dark` 块中存在的自定义语义变量，包括背景色阶（5 级）、文字色阶（3 级）、品牌色、Agent 标识色、Diff 对比色、Prose 增强色、滚动条色。

#### Scenario: Light mode renders all custom variables
- **WHEN** `<html>` 元素没有 `dark` class
- **THEN** 所有自定义 CSS 变量（`--bg-canvas`、`--text-primary`、`--color-brand`、`--agent-claude`、`--diff-insert-bg`、`--prose-link` 等）SHALL 有明确的 Light 主题色值
- **AND** 这些值 SHALL 与 `frontend/docs/reference/color-palette.md` L-八 中建议的值一致

#### Scenario: Primary color is Indigo in light mode
- **WHEN** 处于 Light 主题
- **THEN** `--primary` SHALL 解析为 `#6366F1`（Indigo），而非 shadcn 默认的 oklch 无彩黑

### Requirement: useTheme hook
系统 SHALL 提供 `useTheme` hook（`src/hooks/use-theme.ts`），返回当前主题偏好、实际生效主题、以及设置函数。

#### Scenario: Default theme follows system preference
- **WHEN** 用户首次访问（localStorage 中无主题偏好）
- **THEN** hook SHALL 返回 `theme: 'system'`，`resolvedTheme` 跟随 `prefers-color-scheme` 媒体查询结果

#### Scenario: User sets theme to light
- **WHEN** 用户调用 `setTheme('light')`
- **THEN** `<html>` 元素 SHALL 移除 `dark` class
- **AND** `localStorage` 中 SHALL 存储 `'light'`
- **AND** hook SHALL 返回 `resolvedTheme: 'light'`

#### Scenario: User sets theme to dark
- **WHEN** 用户调用 `setTheme('dark')`
- **THEN** `<html>` 元素 SHALL 添加 `dark` class
- **AND** `localStorage` 中 SHALL 存储 `'dark'`
- **AND** hook SHALL 返回 `resolvedTheme: 'dark'`

#### Scenario: User sets theme to system
- **WHEN** 用户调用 `setTheme('system')`
- **THEN** `<html>` 的 `dark` class SHALL 跟随 `prefers-color-scheme` 实时变化
- **AND** `localStorage` 中 SHALL 存储 `'system'`

#### Scenario: Theme persists across page reload
- **WHEN** 用户上次设置了主题并刷新页面
- **THEN** 页面加载后 SHALL 立即应用上次保存的主题，不出现闪烁

### Requirement: Settings panel in IconSidebar
系统 SHALL 在 `IconSidebar` 底部的 Settings 图标上启用交互（移除 `disabled`），点击后弹出 Settings 面板，面板包含主题切换功能。

#### Scenario: Settings icon opens popover
- **WHEN** 用户点击 Settings 图标
- **THEN** SHALL 弹出 Popover 面板，位于图标右侧
- **AND** 面板 SHALL 包含"外观"标题和三段式主题切换（深色 / 浅色 / 跟随系统）
- **AND** 当前生效的主题选项 SHALL 有视觉高亮

#### Scenario: Settings icon in dark mode
- **WHEN** 当前为 Dark 主题
- **THEN** "深色"选项 SHALL 显示为选中状态（bg-primary-soft + text-primary）

#### Scenario: Settings icon in light mode
- **WHEN** 当前为 Light 主题
- **THEN** "浅色"选项 SHALL 显示为选中状态

### Requirement: FOUC prevention
系统 SHALL 在 React 渲染之前同步设置 `<html>` 的 `dark` class，避免主题闪烁。

#### Scenario: First paint uses correct theme
- **WHEN** 页面首次加载（HTML 解析阶段）
- **THEN** `<html>` 元素的 `dark` class SHALL 在任何 CSS 渲染前由内联脚本确定
- **AND** 用户 SHALL 看到正确的主题色，无白屏闪烁

### Requirement: Inline dark-only styles adaptation
系统中使用硬编码 `rgba(0,0,0,...)` 做 boxShadow 的内联样式 SHALL 改为 CSS 变量引用，以适配 Light 模式。

#### Scenario: Popup shadow adapts to theme
- **WHEN** Light 主题激活
- **THEN** 弹出菜单的 boxShadow SHALL 使用适合浅色背景的阴影值（如 `0 4px 24px rgba(0,0,0,0.1)`）
- **WHEN** Dark 主题激活
- **THEN** boxShadow SHALL 使用暗色背景下的阴影值（如 `0 4px 24px rgba(0,0,0,0.4)`）
