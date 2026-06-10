## Why

前端代码经过多轮迭代后，累积了 92 处视觉规范违规（硬编码颜色、圆角超限、阴影滥用、图标 strokeWidth 不统一、className 拼接未用 cn() 等），导致实际 UI 与 `visual-style-guide.md` 定义的设计系统产生偏差。同时 chat store 中 messages/streaming 等 Server State 数据直接存在 Zustand 中，违反 `development-strategy.md` 的状态管理策略。需要一次性审计并修正所有违规，确保代码实现与设计规范完全对齐。

## What Changes

- **修复硬编码颜色**：3 处 `#fff` inline style → `var(--primary-foreground)`；2 处 `text-white` → `text-primary-foreground`
- **统一图标 strokeWidth**：将 AskAgentCard / FinalSummaryCard / TaskFailureCard / MessageBubble 中的 1.5/1.6/1.7 统一为 `strokeWidth={1.25}`
- **修正圆角违规**：5 处 `rounded-xl` → 按视觉规范降级（头像 8px `rounded-lg`，面板 12px `rounded-[12px]`）
- **清理阴影滥用**：2 处 `shadow-lg` → 仅弹出菜单保留，卡片/tooltip 改用背景色阶表达深度
- **统一 className 拼接**：18 处模板字面量 `className={\`...\${cond}\`} → `cn()` 调用
- **CSS 变量补齐**：确保 `--color-agent-codex` 在 `.dark` 中已定义（当前已存在）
- **标记 Zustand 中的 Server State**：chat store 的 messages/streaming 属于 Server State，当前因 SSE streaming 架构原因暂存 Zustand，记录为已知技术债务

## Capabilities

### New Capabilities
- `style-violation-fixes`: 一次性修复全部 30 处视觉规范违规（颜色、圆角、阴影、strokeWidth、className）
- `classname-cleanup`: 统一 className 拼接方式，18 处模板字面量迁移到 cn()

### Modified Capabilities
<!-- 无已有 spec 的行为变更 -->

## Impact

- **受影响文件**（30+ 组件文件）：
  - `components/chat/` — AskAgentCard, MessageBubble, AgentAvatar, AgentHoverCard, AgentMeta
  - `components/cards/` — FinalSummaryCard, RuntimeStatus, DiffCard, PlanCard, TaskFailureCard
  - `components/diff/` — DiffFileTabs, DiffHeader
  - `components/layout/` — AdminPasswordDialog, IconSidebar
  - `components/ui/` — dialog.tsx（shadow 审查）
  - `pages/` — AgentProfilePage, ImPage
  - `pages/admin/` — UserManagementPage, SessionCleanupPage, AgentOverviewPage, ServiceHealthPage, StatisticsPage, DashboardPage, WorkspacePage
- **无 API 变更**：纯前端样式和代码风格修改
- **无依赖变更**：所有修复使用现有工具链（Tailwind + cn() + Lucide）
- **风险等级**：低 — 全部为视觉效果和代码规范修改，不涉及业务逻辑
