## 1. 色彩类名修正（color-usage-normalization）

- [x] 1.1 修正 `MarkdownRenderer.tsx` 中 `<em>` 标签的 `text-secondary` → `text-text-secondary`（约 2 处）
- [x] 1.2 修正 `MarkdownRenderer.tsx` 中 `<th>` 表头的 `text-secondary` → `text-text-secondary`（约 1 处）
- [x] 1.3 修正 `ToolCard.tsx` 中 `bg-secondary` → `bg-bg-card`（约 1 处）
- [x] 1.4 修正 `GitGraphPanel.tsx` 中非当前分支的 `bg-secondary` → `bg-bg-hover`（约 1 处）
- [x] 1.5 修正 `index.css` 中 `--prose-heading-h1: #FFFFFF` → `--prose-heading-h1: #F0F2F7`
- [x] 1.6 全局 grep 验证：确认 `text-secondary`（非 `text-text-secondary`）不再出现在 .tsx 文件中

## 2. UI 字符串集中（ui-strings-centralization）

- [x] 2.1 创建 `src/lib/ui-text.ts`，定义 UI_ACTIONS / UI_STATUS / UI_MESSAGES / UI_LABELS / UI_PLACEHOLDERS 五组常量（使用 `as const`）
- [x] 2.2 替换 `RightSidebar.tsx` 中散落的中文 UI 字符串 → 引用 ui-text 常量（约 10 处：展开侧栏、Agent 信息、路径信息、仓库路径、任务路径、导出、置顶、退出群聊、删除会话、复制成功）
- [x] 2.3 替换 `MessageBubble.tsx` 中的中文 UI 字符串 → 引用 ui-text 常量（约 3 处：放大、点击放大查看、消息详情）
- [x] 2.4 替换 `GitGraphPanel.tsx` 中的中英文字符串 → 引用 ui-text 常量（约 4 处：Git Graph、commits、branches、git ref 不存在）
- [x] 2.5 替换 `use-chat-stream.ts` 中的错误提示字符串 → 引用 ui-text 常量（约 1 处：发送失败）
- [x] 2.6 替换其余组件中的散落字符串 → 引用 ui-text 常量（ImPage、ChatArea、MessageInput、ConversationList、NewChatDialog、ContactsPage、AgentProfilePage、SkillsHubPage、AdminPages 等）
- [x] 2.7 全局 grep 验证：`grep -rn '[\x{4e00}-\x{9fff}]' frontend/src/components/ frontend/src/hooks/` 确认无遗漏（注释和 console.log 除外）

## 3. 组件拆分（component-split-refactor）— MessageBubble

- [x] 3.1 从 `MessageBubble.tsx` 提取 `BlockRenderer` 组件到 `components/chat/BlockRenderer.tsx`（约 110 行）
- [x] 3.2 从 `MessageBubble.tsx` 提取 `AgentMessageContent` 组件到 `components/chat/AgentMessageContent.tsx`（约 125 行）
- [x] 3.3 简化 `MessageBubble.tsx` 主文件，仅保留编排逻辑（目标 ≤160 行）
- [x] 3.4 验证 MessageBubble 拆分后渲染输出与拆分前完全一致

## 4. 组件拆分 — RightSidebar

- [x] 4.1 从 `RightSidebar.tsx` 提取 `useCollapsible` Hook 到 `components/chat/useCollapsible.ts`（约 22 行）
- [x] 4.2 从 `RightSidebar.tsx` 提取 `AgentInfoSection` 到 `components/chat/AgentInfoSection.tsx`（约 57 行）
- [x] 4.3 从 `RightSidebar.tsx` 提取 `SidebarPathSection` 到 `components/chat/SidebarPathSection.tsx`（约 58 行）
- [x] 4.4 从 `RightSidebar.tsx` 提取 `SidebarActions` 到 `components/chat/SidebarActions.tsx`（约 56 行）
- [x] 4.5 简化 `RightSidebar.tsx` 主文件，仅保留编排逻辑（目标 ≤230 行）
- [x] 4.6 验证 RightSidebar 拆分后交互行为与拆分前完全一致

## 5. 组件拆分 — GitGraphPanel

- [x] 5.1 从 `GitGraphPanel.tsx` 提取 `GraphRenderer` 到 `components/chat/GraphRenderer.tsx`（约 90 行）
- [x] 5.2 从 `GitGraphPanel.tsx` 提取 `GraphBranchLabels` 到 `components/chat/GraphBranchLabels.tsx`（约 42 行）
- [x] 5.3 从 `GitGraphPanel.tsx` 提取 `GraphTooltip` 到 `components/chat/GraphTooltip.tsx`（约 20 行）
- [x] 5.4 简化 `GitGraphPanel.tsx` 主文件，仅保留编排逻辑（目标 ≤170 行）
- [x] 5.5 验证 GitGraphPanel 拆分后渲染和交互与拆分前完全一致

## 6. 最终验证

- [x] 6.1 运行 `pnpm build` 确认零编译错误
- [x] 6.2 运行 `pnpm lint` 确认无新增 lint 警告
- [ ] 6.3 启动前端服务，暗色模式下逐页面验证文字可见性和交互无回归
