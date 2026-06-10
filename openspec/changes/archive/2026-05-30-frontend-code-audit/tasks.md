## 1. Border-radius 修正（visual-radius-fix）

- [x] 1.1 修正 MessageInput 发送按钮：`rounded-lg` → `rounded-[6px]`
- [x] 1.2 修正 MessageInput textarea：`rounded-lg` → `rounded-[8px]`
- [x] 1.3 修正 ConversationList 搜索输入框：`rounded-lg` → `rounded-[8px]`
- [x] 1.4 修正 ConversationList 新建按钮：`rounded-md` → `rounded-[6px]`
- [x] 1.5 修正 IconSidebar NavItem：`rounded-[10px]` → `rounded-[6px]`
- [x] 1.6 修正 IconSidebar hover card：`rounded-xl` → `rounded-[12px]`
- [x] 1.7 修正 IconSidebar 用户头像：`rounded-[10px]` → `rounded-full`
- [x] 1.8 修正 MessageBubble 用户头像：`rounded-lg` → `rounded-full`
- [x] 1.9 修正 AskAgentCard 状态 Badge：`rounded-md` → `rounded-full` + `text-[10px]` → `text-[11px]`
- [x] 1.10 全量扫描 TSX 文件中 `rounded-md`/`rounded-lg`/`rounded-xl` 用法，分类为"保留"（Agent 头像 8px、卡片 10px 合规）和"修正"

## 2. 色彩 Token 合规（color-token-enforcement）

- [x] 2.1 修正 AgentMeta 描述文字：`text-foreground/50` → `text-secondary`，`tracking-wider` → `tracking-wide`
- [x] 2.2 修正 SkillCard 描述文字：`text-foreground/75` → `text-secondary`
- [x] 2.3 修正 SkillCard 元信息：`text-foreground/55` → `text-tertiary`
- [x] 2.4 修正 MessageBubble streaming 光标：`text-primary`（品牌色）→ `text-foreground`
- [x] 2.5 修正 ChatArea streaming 状态文字：`text-primary`（品牌色）→ `text-tertiary`
- [x] 2.6 修正 AskAgentCard 边框/背景：`border-primary/20` → `border`，`bg-primary/[0.03]` → `bg-card`，`bg-primary/[0.06]` → `bg-hover`
- [x] 2.7 修正 AskAgentCard hover 态：`hover:border-primary/40` → `hover:bg-hover`，`border-primary/10` → `border`
- [x] 2.8 修正 AskAgentCard 状态 Badge 颜色：`text-emerald-400` → `text-success`，`bg-emerald-500/10` → `bg-success/10`
- [x] 2.9 修正 AskAgentCard 状态 Badge 字号：`text-[10px]` → `text-[11px]`（已在 1.9 中完成）
- [x] 2.10 修正 IconSidebar NavItem label 字号：`fontSize: 9` → `fontSize: 11`
- [x] 2.11 全量扫描 TSX 文件中 opacity 调色用法（核心聊天组件已修正，AgentProfilePage 残留待后续处理）

## 3. 过渡动画合规（transition-compliance）

- [x] 3.1 修正 MessageInput 发送按钮：移除 `transition-colors`
- [x] 3.2 修正 ConversationList 搜索/按钮：移除 `transition-colors duration-120`
- [x] 3.3 修正 ConversationItem：移除 `transition-colors duration-120`
- [x] 3.4 修正 IconSidebar hover card：`transition-all` → `transition-[transform,opacity]`
- [x] 3.5 修正 IconSidebar hover card 阴影：`0 8px 32px` → `0 4px 24px`
- [x] 3.6 修正 IconSidebar 在线状态指示器边框：`border-2` → `border`
- [x] 3.7 全量扫描 `transition-colors`/`transition-all`（核心组件已修正，admin/diff 残留待后续处理）

## 4. Error Boundary 补充（error-boundary）

- [x] 4.1 创建 `components/ui/error-boundary.tsx`：class component，支持 `fallback` prop，默认降级 UI 显示错误描述 + 重试按钮
- [x] 4.2 在 ImPage 中为 ChatArea 包裹 ErrorBoundary
- [x] 4.3 在 ImPage 中为 AdminContent 包裹 ErrorBoundary
- [x] 4.4 在 ImPage 中为 ConversationList 包裹 ErrorBoundary
- [x] 4.5 AdminContent 已在模块顶层（无需提取）

## 5. API 健壮性（api-resilience）

- [x] 5.1 在 `lib/api.ts` 中创建 `ApiError` 类和 `handleResponse<T>(res)` 辅助函数
- [x] 5.2 将 `fetchTasks` 改用 `handleResponse`
- [x] 5.3 将 `fetchTask` 改用 `handleResponse`
- [x] 5.4 将 `createTask` 改用 `handleResponse`
- [x] 5.5 将 `fetchAgentTypes` 改用 `handleResponse`（自定义逻辑保留，补充 res.ok 检查）
- [x] 5.6 将 `getTaskMessages` 改用 `handleResponse`
- [x] 5.7 其余 API 函数（submitMessage、uploadAvatar、updateSession、validateRepoPath、fetchAgentProfile、fetchAgentDetail、admin 系列）已有 res.ok 检查
- [x] 5.8 修正 ChatArea `loadMoreMessages` catch 块：添加 loadError 状态 + 错误 banner

## 6. 验证与回归检查

- [x] 6.1 运行 ESLint 检查 + TypeScript 类型检查 + Vite 构建，全部通过
- [x] 6.2 ESLint 零错误，TypeScript 零错误，Vite build 成功
- [x] 6.3 Error Boundary 组件已创建，ImPage 中 ChatArea/AdminContent/ConversationList 已包裹
