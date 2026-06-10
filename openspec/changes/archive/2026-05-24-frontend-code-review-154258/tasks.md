## 1. 文档矛盾修复

- [x] 1.1 更新 `visual-style-guide.md` 图标部分：将推荐从 Phosphor Icons 改为 Lucide React，说明通过 `strokeWidth={1.25}` 模拟细线视觉
- [x] 1.2 在 `lib/api.ts` 中与 `generated/` 重叠的手写类型添加 TODO 注释

## 2. CSS 变量系统补齐

- [x] 2.1 在 `index.css` 的 `.dark` 中补充 `--code-bg: #0D0F14` 变量
- [x] 2.2 检查并补充其他缺失的设计令牌（如 `--color-danger-bg` 等语义色背景）

## 3. CSS 变量迁移 — 核心组件

- [x] 3.1 SessionList.tsx：将所有 `bg-gray-*`、`text-green-*` 等 Tailwind 原始色阶替换为 CSS 变量
- [x] 3.2 AgentAvatar.tsx：将硬编码的 Agent 标识色 hex 值替换为 `var(--agent-xxx)` CSS 变量
- [x] 3.3 AgentAvatar.tsx：将状态点尺寸从 8px 修正为 4px
- [x] 3.4 AgentAvatar.tsx：删除内联 `AGENT_LABELS`，改从 `lib/constants.ts` 导入 `AGENT_NAMES`
- [x] 3.5 ChatArea.tsx：将错误横幅背景色从浅色 `#FEF2F2` 替换为暗色语义色 `rgba(239, 68, 68, 0.1)`

## 4. 纯白文字清除

- [x] 4.1 AgentEditDialog.tsx：将 `text-white` 替换为 `text-[var(--text-primary)]`
- [x] 4.2 NewChatDialog.tsx：将 `text-white` 替换为 `text-[var(--text-primary)]`
- [x] 4.3 MessageInput.tsx：将 `text-white` 替换为 `text-[var(--text-primary)]`
- [x] 4.4 button.tsx（shadcn）：确认默认变体的文字色为 `var(--text-primary)` 而非 `#FFFFFF`

## 5. 边框和圆角一致性

- [x] 5.1 全局替换 `rgba(255, 255, 255, 0.1)` 为 `var(--border)` 或 `rgba(255, 255, 255, 0.06)`
- [x] 5.2 滚动到底部按钮：将 `rounded-full` 替换为 `rounded-xl`（12px）
- [x] 5.3 CodeBlock/MarkdownRenderer：将硬编码 `#0D0F14` 替换为 `var(--code-bg)`

## 6. 魔法字符串修复

- [x] 6.1 在 `lib/constants.ts` 中定义 SSE 事件类型枚举（或从 `generated/` 导出）
- [x] 6.2 `use-chat-stream.ts`：将所有字符串字面量事件类型替换为常量引用

## 7. Bug 修复

- [x] 7.1 MessageList.tsx：非虚拟化路径的 `MessageRenderer` 补充 `avatarUrl` 和 `agentName` props

## 8. 死代码清理

- [x] 8.1 删除 `stores/app.ts`（空 store，确认无导入）
- [x] 8.2 删除 `SessionList.tsx`（确认无导入）
- [x] 8.3 清理 `use-sessions.ts` 中未使用的 mutation
- [x] 8.4 清理 `stores/chat.ts` 中未使用的 `currentTaskId`

## 9. 验证

- [x] 9.1 运行 `pnpm build` 确认无编译错误
- [x] 9.2 启动 dev server 视觉验证：检查 SessionList、AgentAvatar、ChatArea、对话框等组件显示正常
- [x] 9.3 确认无 `text-white`、`bg-gray-*`（非变量）残留（全局 grep 验证）
