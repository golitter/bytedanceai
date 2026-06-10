## 1. 硬编码颜色修复

- [x] 1.1 修复 AdminPasswordDialog.tsx — inline style `color: '#fff'` → `color: 'var(--primary-foreground)'`
- [x] 1.2 修复 SessionCleanupPage.tsx — inline style `color: '#fff'` → `color: 'var(--primary-foreground)'`
- [x] 1.3 修复 AgentOverviewPage.tsx — inline style `color: '#fff'` → `color: 'var(--primary-foreground)'`
- [x] 1.4 修复 AgentProfilePage.tsx — Camera 图标 `text-white` → `text-primary-foreground`
- [x] 1.5 修复 UserManagementPage.tsx — Camera 图标 `text-white` → `text-primary-foreground`

## 2. Lucide strokeWidth 统一

- [x] 2.1 修复 AskAgentCard.tsx — 5 处 strokeWidth 从 1.5/1.6 → 1.25
- [x] 2.2 修复 MessageBubble.tsx — 1 处 strokeWidth 从 1.5 → 1.25
- [x] 2.3 修复 FinalSummaryCard.tsx — 3 处 strokeWidth 从 1.7 → 1.25
- [x] 2.4 修复 TaskFailureCard.tsx — 1 处 strokeWidth 从 1.7 → 1.25

## 3. 圆角规范修正

- [x] 3.1 修复 AgentProfilePage.tsx — 3 处 `rounded-xl` → `rounded-lg`（头像容器、头像图片、hover overlay）
- [x] 3.2 修复 UserManagementPage.tsx — 2 处 `rounded-xl` → `rounded-lg`（头像容器、hover overlay）

## 4. className 迁移 cn()

- [x] 4.1 修复 AgentAvatar.tsx — 模板字面量 → `cn()`
- [x] 4.2 修复 AgentHoverCard.tsx — 模板字面量 → `cn()`
- [x] 4.3 修复 AgentMeta.tsx — 模板字面量 → `cn()`
- [x] 4.4 修复 FinalSummaryCard.tsx — 模板字面量 → `cn()`
- [x] 4.5 修复 RuntimeStatus.tsx — 2 处模板字面量 → `cn()`
- [x] 4.6 修复 DiffFileTabs.tsx — 模板字面量 → `cn()`
- [x] 4.7 修复 DiffCard.tsx — 模板字面量 → `cn()`
- [x] 4.8 修复 PlanCard.tsx — 模板字面量 → `cn()`
- [x] 4.9 修复 DiffHeader.tsx — 模板字面量 → `cn()`
- [x] 4.10 修复 AgentProfilePage.tsx — badge 模板字面量 → `cn()`
- [x] 4.11 修复 WorkspacePage.tsx — loading 图标模板字面量 → `cn()`
- [x] 4.12 修复 StatisticsPage.tsx — loading 图标模板字面量 → `cn()`
- [x] 4.13 修复 SessionCleanupPage.tsx — loading 图标模板字面量 → `cn()`
- [x] 4.14 修复 AgentOverviewPage.tsx — loading 图标模板字面量 → `cn()`
- [x] 4.15 修复 ServiceHealthPage.tsx — 2 处模板字面量 → `cn()`（loading 图标 + 状态指示器）
- [x] 4.16 修复 DashboardPage.tsx — loading 图标模板字面量 → `cn()`

## 5. 技术债务标记

- [x] 5.1 在 chat.ts store 文件头部添加 Server State 技术债务注释，说明 messages/streaming 数据暂存 Zustand 的原因和后续迁移方向

## 6. 验证

- [x] 6.1 全局搜索确认无残留硬编码 `#fff` / `text-white` 违规
- [x] 6.2 全局搜索确认无残留 `strokeWidth` 非 1.25 值（排除 shadcn/ui 内部）
- [x] 6.3 全局搜索确认无残留模板字面量 className 拼接
- [x] 6.4 启动前端确认无编译错误和视觉回归
