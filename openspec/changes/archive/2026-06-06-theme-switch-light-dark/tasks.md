## 1. Light 主题 CSS 变量补全

- [x] 1.1 在 `src/index.css` 的 `:root` 块中补全背景色阶变量：`--bg-canvas: #FFFFFF`、`--bg-sidebar: #F8F9FA`、`--bg-card: #FFFFFF`、`--bg-hover: #F1F3F5`、`--bg-active: #E9ECEF`
- [x] 1.2 补全文字色阶变量：`--text-primary: #1A1A1A`、`--text-secondary: #6B7280`、`--text-tertiary: #9CA3AF`
- [x] 1.3 覆盖品牌色变量：`--primary: #6366F1`、`--primary-foreground: #FAFAFA`、`--ring: #6366F1`、`--color-brand: #6366F1`
- [x] 1.4 补全品牌色变体：`--primary-soft: rgba(99,102,241,0.06)`、`--primary-border: rgba(99,102,241,0.12)`
- [x] 1.5 补全语义色：`--color-success: #22C55E`、`--color-warning: #F59E0B`、`--color-error: #EF4444`、`--color-danger-bg: rgba(239,68,68,0.06)`
- [x] 1.6 补全 Agent 标识色：`--agent-claude: #DA7756`、`--agent-opencode: #10B981`、`--agent-orchestrator: #EAB308`、`--agent-codex: #6366F1`
- [x] 1.7 补全 Diff 对比色（6 个变量：insert/delete 的 color、bg、bg-strong）
- [x] 1.8 补全 Prose 增强色（9 个变量：heading、link、link-hover、bold、bq-border、bq-bg、code-bg、code-text、li-marker、hr）
- [x] 1.9 补全特殊变量：`--code-bg: #F6F8FA`、`--shadow-popup`（新增，Light = `0 4px 24px rgba(0,0,0,0.1)`，Dark = `0 4px 24px rgba(0,0,0,0.4)`）

## 2. useTheme Hook

- [x] 2.1 创建 `src/hooks/use-theme.ts`，导出 `useTheme()` 返回 `{ theme, resolvedTheme, setTheme }`
- [x] 2.2 实现主题偏好读写：`localStorage` key 为 `theme`，合法值为 `'dark' | 'light' | 'system'`，默认 `'system'`
- [x] 2.3 实现 `resolvedTheme` 逻辑：`system` 时通过 `window.matchMedia('(prefers-color-scheme: dark)')` 解析，`dark`/`light` 直接返回
- [x] 2.4 实现 `setTheme` 时同步更新 `<html>` 的 `dark` class 和 `localStorage`
- [x] 2.5 监听 `matchMedia` 变化事件，当 `theme === 'system'` 时实时响应系统主题切换

## 3. FOUC 防闪烁

- [x] 3.1 在 `frontend/index.html` 的 `<head>` 中添加内联 `<script>`，同步读取 `localStorage.getItem('theme')` 并在 `<html>` 上设置 `dark` class，确保首帧渲染即使用正确主题

## 4. Settings 面板 UI

- [x] 4.1 修改 `IconSidebar.tsx`：移除 Settings NavItem 的 `disabled` 属性，将其从 `NavItem` 改为可点击按钮（不再作为 tab 导航，而是触发 Popover）
- [x] 4.2 创建 `src/components/layout/SettingsPanel.tsx`，包含"外观"标题 + 三段式切换按钮（深色 / 浅色 / 跟随系统），使用 `useTheme` hook
- [x] 4.3 在 SettingsPanel 中实现三段式切换 UI：当前选项显示 `bg-primary-soft text-primary`，其余为 `bg-muted text-muted-foreground`，hover 时 `bg-bg-hover`
- [x] 4.4 在 `IconSidebar` 中用 `Popover`（复用 `components/ui/popover.tsx`）包裹 Settings 按钮，弹出 `SettingsPanel`

## 5. 内联暗色样式适配

- [x] 5.1 排查并修复 `IconSidebar.tsx` 中 hover 卡片的内联 `boxShadow: '0 4px 24px rgba(0,0,0,0.4)'`，改用 CSS 变量 `var(--shadow-popup)`
- [x] 5.2 排查并修复 `AgentAvatar.tsx` 和 `AgentProfilePage.tsx` 中的内联 boxShadow（agent glow），确认在 Light 模式下视觉效果合理
- [x] 5.3 排查并修复 `SidebarActions.tsx` 中的 toast 内联样式，确认 Light 模式下文字可读

## 6. 验证

- [x] 6.1 启动前端，在 Light 模式下逐页检查所有页面（聊天、通讯录、SkillsHub、Admin 面板），确认无颜色丢失或可读性问题
- [x] 6.2 验证主题切换即时生效，刷新后保持
- [x] 6.3 验证 `prefers-color-scheme` 跟随模式在 macOS 切换系统外观后实时响应
- [x] 6.4 验证 FOUC：在慢网络下刷新页面，确认无白屏闪烁
