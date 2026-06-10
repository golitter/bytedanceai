## 1. CSS 变量体系统一

- [x] 1.1 在 `index.css` 的 `@theme inline` 块中添加自定义扩展 token 的 Tailwind 桥接：`--color-tertiary` → `var(--text-tertiary)`、`--color-code-bg` → `var(--code-bg)`、`--color-danger-bg` → `var(--color-danger-bg)`、`--color-agent-claude` / `--color-agent-opencode` / `--color-agent-orchestrator`
- [x] 1.2 从 `.dark` 块中移除与 shadcn token 重复的 8 个自定义 app token（`--bg-canvas`、`--bg-sidebar`、`--bg-hover`、`--text-primary`、`--text-secondary`、`--color-brand`、`--color-error`、`--divider`）
- [x] 1.3 移除未使用的 CSS 变量：`--bg-active`、`--color-success`、`--color-warning`
- [x] 1.4 为需要基于品牌色透明变体的场景新增 CSS 变量（如 `--primary-soft: rgba(99,102,241,0.08)` 和 `--primary-border: rgba(99,102,241,0.15)`）并加入 `@theme inline`

## 2. 组件样式迁移 — chat 模块

- [x] 2.1 `ChatArea.tsx`：将 6 处内联样式替换为 Tailwind 类（`var(--text-primary)` → `text-foreground`，`var(--bg-canvas)` → `bg-background`，`var(--color-brand)` → `text-primary`，`var(--text-secondary)` → `text-muted-foreground`，`var(--text-tertiary)` → `text-tertiary`，硬编码边框 → `border-border`）
- [x] 2.2 `AgentAvatar.tsx`：将 4 处硬编码状态色替换为 CSS 变量（`#22C55E` → `var(--color-success)`，`#F59E0B` → `var(--color-warning)`，`#5A6070` → `var(--text-tertiary)`，`#EF4444` → `var(--destructive)`），Agent 标识色保持使用 `--agent-*` 变量
- [x] 2.3 `AgentEditDialog.tsx`：将 9 处内联样式替换为 Tailwind 类，4 处硬编码 `rgba(255,255,255,0.06)` 边框替换为 `border-border`
- [x] 2.4 `MessageBubble.tsx`：将 4 处内联样式替换为 Tailwind 类，`rgba(99,102,241,0.08/0.15)` 替换为 `var(--primary-soft)` / `var(--primary-border)`，修复 `text-brand` 为 `text-primary`
- [x] 2.5 `MessageInput.tsx`：将 4 处内联样式替换为 Tailwind 类，`placeholder:text-[var(--text-tertiary)]` 替换为 `placeholder:text-tertiary`
- [x] 2.6 `MessageList.tsx`：将 4 处内联样式替换为 Tailwind 类

## 3. 组件样式迁移 — im 模块

- [x] 3.1 `ConversationList.tsx`：将 9 处内联样式替换为 Tailwind 类，1 处硬编码 `rgba(255,255,255,0.06)` 替换为 `border-border`
- [x] 3.2 `ConversationItem.tsx`：将 4 处内联样式替换为 Tailwind 类
- [x] 3.3 `NewChatDialog.tsx`：将 14 处内联样式替换为 Tailwind 类，5 处硬编码 `rgba(255,255,255,0.06)` 替换为 `border-border`，2 处硬编码状态色 `#EF4444` / `#22C55E` 替换为 CSS 变量

## 4. 组件样式迁移 — markdown 模块

- [x] 4.1 `MarkdownRenderer.tsx`：将 4 处内联样式替换为 Tailwind 类，3 处硬编码 `rgba(255,255,255,0.06)` 替换为 `border-border`
- [x] 4.2 `CodeBlock.tsx`：将 2 处内联样式替换为 Tailwind 类，1 处硬编码 `#5A6070` 替换为 `var(--text-tertiary)`

## 5. 组件样式迁移 — pages 与 hooks

- [x] 5.1 `ImPage.tsx`：将 3 处内联样式替换为 Tailwind 类
- [x] 5.2 `use-hover-style.ts`：更新 `--bg-hover` 引用为 `--accent`（或对应的 Tailwind 类）

## 6. 死代码与冗余清理

- [x] 6.1 删除 `frontend/src/App.tsx`（返回 null，从未被渲染）
- [x] 6.2 删除 `frontend/src/hooks/use-sessions.ts`（从未被引用）
- [x] 6.3 检查所有组件文件的未使用 import 并清理

## 7. 文档更新

- [x] 7.1 更新 `frontend/AGENTS.md`：反映当前目录结构（移除 `stores/app.ts` 引用，补充 chat/im/markdown 组件、hooks、stores 描述）

## 8. 验证

- [x] 8.1 运行 `npm run build` 确认无编译错误
- [x] 8.2 启动前端并逐页面验证视觉表现与迁移前一致（无颜色偏差、无布局错乱）
- [x] 8.3 确认所有硬编码颜色值已消除（grep 验证无残留 hex/rgba）
