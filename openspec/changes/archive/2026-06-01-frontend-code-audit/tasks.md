## 1. 常量收敛（Magic Strings）

- [x] 1.1 在 `lib/constants.ts` 中补充 `AGENT_TYPES` 枚举（`CLAUDE_CODE`, `OPENCODE`, `ORCHESTRATOR`, `CODEX`）、`MESSAGE_ROLES`（`USER`, `AGENT`, `SYSTEM`）、`CHAT_STATUSES`（`IDLE`, `LOADING`, `STREAMING`, `TOOL_RUNNING`, `DONE`, `ERROR`, `RETRYING`）、`CURRENT_USER_NAME` 常量
- [x] 1.2 替换 `components/chat/ChatArea.tsx` 中的 agent type / status / role 硬编码字符串为常量引用
- [x] 1.3 替换 `components/chat/MessageRenderer.tsx` 中的 agent type / role 硬编码字符串
- [x] 1.4 替换 `components/chat/AskAgentCard.tsx` 中的 agent type 硬编码字符串
- [x] 1.5 替换 `components/im/NewChatDialog.tsx` 和 `AgentSelectList.tsx` 中的 agent type 硬编码字符串
- [x] 1.6 替换 `lib/api.ts` 中的 `'orchestrator'` 硬编码字符串
- [x] 1.7 替换 `components/chat/MembersSection.tsx` 中的 status 比较和硬编码用户名
- [x] 1.8 替换 `components/layout/IconSidebar.tsx` 和 `AdminMenu.tsx` 中的硬编码用户名
- [x] 1.9 替换 `pages/admin/UserManagementPage.tsx` 中的硬编码用户名
- [x] 1.10 替换 `components/chat/MessageList.tsx`、`RightSidebar.tsx`、`HistorySearch.tsx`、`AnnouncementsSection.tsx` 中的 role 硬编码字符串
- [x] 1.11 替换 `hooks/use-chat-stream.ts` 中的 role 硬编码字符串

## 2. 色值迁移 — 核心组件

- [x] 2.1 修复 `components/chat/RightSidebar.tsx`：将 `color: '#fff'` → `text-foreground`，`background: 'rgba(0,0,0,0.75)'` → `bg-black/75` Tailwind 类
- [x] 2.2 修复 `components/cards/FinalSummaryCard.tsx`：将 `text-emerald-300` / `bg-emerald-500/10` / `text-red-300` / `bg-red-500/10` → 语义类
- [x] 2.3 修复 `components/cards/TaskFailureCard.tsx`：将 `border-red-500/25` / `bg-red-500/[0.08]` / `text-red-*` → `border-destructive` / `bg-destructive/*` / `text-destructive-*`
- [x] 2.4 修复 `components/cards/PlanCard.tsx`：将 `text-red-500` → `text-destructive`
- [x] 2.5 修复 `components/cards/RuntimeStatus.tsx`：将 `bg-red-500/10` / `text-red-500` → `bg-destructive/10` / `text-destructive`
- [x] 2.6 修复 `components/diff/DiffFileTabs.tsx`：将 `bg-green-500/15` / `bg-red-500/15` / `bg-blue-500/15` / `bg-purple-500/15` → 语义色类
- [x] 2.7 修复 `components/diff/DiffFileInfo.tsx`：将 `text-green-500` / `text-red-500` → 语义类
- [x] 2.8 修复 `components/diff/DiffHeader.tsx`：将 `bg-green-500/10` / `text-green-600` / `text-red-500` → 语义类

## 3. 色值迁移 — Admin 页面

- [x] 3.1 修复 `pages/admin/AgentOverviewPage.tsx`：将 `background: 'rgba(0,0,0,0.5)'` → `bg-black/50`
- [x] 3.2 修复 `pages/admin/WorkspacePage.tsx`：将 `rgba(34, 197, 94, 0.1)` → `bg-success/10`，统一其余硬编码色值
- [x] 3.3 修复 `pages/admin/SessionCleanupPage.tsx`：将 `rgba(34, 197, 94, 0.1)` → `bg-success/10`
- [x] 3.4 修复 `pages/admin/UserManagementPage.tsx`：将 `bg-black/40` → `bg-black/40`（确认是否已有 Tailwind 等价类）
- [x] 3.5 修复 `pages/AgentProfilePage.tsx`：将 `bg-black/40` → 确认合规

## 4. 内联样式迁移 — Admin 页面

- [x] 4.1 迁移 `pages/admin/StatisticsPage.tsx`：将所有 `style={{ color: 'var(--text-primary)', border: '...' }}` → Tailwind 类
- [x] 4.2 迁移 `pages/admin/WorkspacePage.tsx`：将所有 `style={{ border, color }}` → Tailwind 类
- [x] 4.3 迁移 `pages/admin/SessionCleanupPage.tsx`：将所有 `style={{ border, ... }}` → Tailwind 类
- [x] 4.4 迁移 `pages/admin/AgentOverviewPage.tsx`：将所有 `style={{ border, ... }}` → Tailwind 类
- [x] 4.5 迁移 `pages/admin/ServiceHealthPage.tsx`：将所有 `style={{ border, ... }}` → Tailwind 类
- [x] 4.6 迁移 `pages/admin/UserManagementPage.tsx`：将所有 `style={{ border, ... }}` → Tailwind 类
- [x] 4.7 迁移 `pages/admin/DashboardPage.tsx`：将所有 `style={{ background: 'var(--border)' }}` → Tailwind 类

## 5. 内联样式迁移 — 核心组件

- [x] 5.1 迁移 `components/layout/AdminMenu.tsx`：将 `style={{ borderRight, background }}` → Tailwind 类
- [x] 5.2 迁移 `components/layout/IconSidebar.tsx`：将 `style={{ fontSize: 11 }}` → `text-[11px]`，`style={{ background }}` → Tailwind 类
- [x] 5.3 迁移 `components/chat/MembersSection.tsx`：将 `style={{ background: 'transparent', boxShadow }}` → Tailwind 类
- [x] 5.4 迁移 `components/im/ConversationItem.tsx`：将 active indicator 的 inline style → Tailwind 类

## 6. 过渡合规

- [x] 6.1 替换 `components/chat/AnnouncementsSection.tsx` 中 4 处 `transition-colors` → `transition-[transform,opacity]`
- [x] 6.2 替换 `components/chat/MembersSection.tsx` 中 2 处 `transition-colors`
- [x] 6.3 替换 `components/chat/MessageBubble.tsx` 中 1 处 `transition-colors`
- [x] 6.4 替换 `components/chat/HistorySearch.tsx` 中 2 处 `transition-colors`
- [x] 6.5 替换 `components/chat/RightSidebar.tsx` 中 7 处 `transition-colors`
- [x] 6.6 替换 `components/chat/MessageList.tsx` 中 1 处 `transition-colors`
- [x] 6.7 替换 `components/layout/IconSidebar.tsx`、`AdminMenu.tsx`、`AdminPasswordDialog.tsx` 中 4 处 `transition-colors`
- [x] 6.8 替换 `components/cards/PreviewCard.tsx`、`AttachmentCard.tsx` 中 2 处 `transition-colors`
- [x] 6.9 替换 `components/diff/DiffFileTabs.tsx`、`DiffHeader.tsx`、`DiffFileEditorInner.tsx` 中 8 处 `transition-colors`
- [x] 6.10 替换 `components/im/NewChatDialog.tsx`、`AgentSelectList.tsx` 中 2 处 `transition-colors`
- [x] 6.11 替换 `pages/AgentProfilePage.tsx` 中 1 处 `transition-colors`
- [x] 6.12 替换 admin 页面中 `transition-all`（`StatisticsPage` 2 处、`DashboardPage` 1 处）及约 10 处 `transition-colors`

## 7. 字号修正

- [x] 7.1 替换 `components/chat/AnnouncementsSection.tsx` 中 2 处 `text-[10px]` → `text-[11px]`
- [x] 7.2 替换 `components/chat/HistorySearch.tsx` 中 1 处 `text-[10px]` → `text-[11px]`
- [x] 7.3 替换 `components/diff/DiffFileTabs.tsx` 中 1 处 `text-[10px]` → `text-[11px]`
- [x] 7.4 替换 `pages/admin/StatisticsPage.tsx` 中 2 处 `text-[10px]` → `text-[11px]`

## 8. 阴影修正

- [x] 8.1 移除 `components/chat/RightSidebar.tsx` toggle 按钮上的 `shadow-sm`，改用背景色差表达层级

## 9. Admin 页面数据获取迁移（TanStack Query）

- [x] 9.1 在 `lib/api.ts` 中补充 admin API 函数（`fetchStatistics`、`fetchWorkspaces`、`fetchServiceHealth`、`fetchSessions`、`fetchAgents`、`fetchUsers`），每个函数返回类型标注
- [x] 9.2 迁移 `pages/admin/StatisticsPage.tsx`：`useState+useEffect` → `useQuery`
- [x] 9.3 迁移 `pages/admin/DashboardPage.tsx`：`useState+useEffect` → `useQuery`
- [x] 9.4 迁移 `pages/admin/WorkspacePage.tsx`：`useState+useEffect` → `useQuery`
- [x] 9.5 迁移 `pages/admin/ServiceHealthPage.tsx`：`useState+useEffect` → `useQuery`
- [x] 9.6 迁移 `pages/admin/UserManagementPage.tsx`：`useState+useEffect` → `useQuery`
- [x] 9.7 迁移 `pages/admin/SessionCleanupPage.tsx`：`useState+useEffect` → `useQuery`
- [x] 9.8 迁移 `pages/admin/AgentOverviewPage.tsx`：`useState+useEffect` → `useQuery`

## 10. Chat Store 拆分

- [x] 10.1 从 `stores/chat.ts` 提取导航状态到 `stores/navigation-store.ts`（activeTab, sidebarVisible 等）
- [x] 10.2 从 `stores/chat.ts` 提取会话状态到 `stores/session-store.ts`（sessions, activeSession, session CRUD）
- [x] 10.3 从 `stores/chat.ts` 提取消息状态到 `stores/message-store.ts`（messages, streaming, runtimeBlocks, blockReducer）
- [x] 10.4 更新所有 import 引用：将 `stores/chat` 拆分为对应的新 store 路径
- [x] 10.5 验证拆分后所有组件功能正常（无状态丢失、无循环依赖）
