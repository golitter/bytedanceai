## Context

当前聊天系统前端中，Agent 消息通过 `MarkdownRenderer`（基于 react-markdown + remark-gfm + shiki）渲染为富文本，但用户消息在 `MessageRenderer.tsx:99` 以 `{msg.content}` 纯文本渲染。输入栏 `MessageInput.tsx` 是非受控 textarea，直接操作 `textareaRef.current.value`，无法驱动实时预览。

项目已有依赖：`react-markdown@^10.1.0`、`remark-gfm@^4.0.1`、`shiki@^4.1.0`，无需新增依赖。

设计参考 Demo HTML：`frontend/docs/payloads/markdown-chat-preview-demo.html`

## Goals / Non-Goals

**Goals:**
- 用户消息与 Agent 消息统一使用 MarkdownRenderer 渲染
- 输入栏新增 Markdown 模式切换按钮，开启后进入双栏实时预览
- 双栏输入区自动增高，超出最大高度时预览区跟随光标同步滚动
- 保持现有功能（@提及、Enter 发送、IME 兼容）不受影响

**Non-Goals:**
- 不改变后端消息存储格式（仍是纯文本 Markdown source）
- 不做消息编辑/重新渲染
- 不做 Markdown 工具栏（加粗、斜体等按钮）

## Decisions

### 1. MessageInput 改为受控组件

**决策**：引入 `useState` 管理 textarea 值，作为双栏预览的数据源。

**理由**：当前 `textareaRef.current.value` 直接赋值无法触发 React 重渲染，无法驱动预览区更新。受控组件是 React 标准模式，也是唯一能同步 `MarkdownRenderer` 的方式。

**影响**：`handleSend`、`insertMention` 从 state 读取值；所有 `textareaRef.current.value` 赋值改为 `setValue()`。

### 2. Markdown 模式切换按钮位于输入栏上方工具栏

**决策**：在 `border-t` 与 textarea 之间新增一个工具栏行，放置 Markdown 切换按钮。

**理由**：按钮不占用输入区宽度，视觉层级清晰。默认关闭时工具栏仅一行按钮，不干扰正常使用。

### 3. 双栏布局：flex 平分 + 固定发送按钮

**决策**：双栏模式下布局为 `[textarea flex:1] [divider 1px] [preview flex:1] [send-btn fixed]`，三栏在同一 flex 容器内。

**理由**：与 Demo HTML 验证过的布局一致，预览区与编辑区等宽，发送按钮在最右侧不随内容变化。

### 4. 光标同步滚动采用比例映射

**决策**：`previewScrollTop = (textareaScrollTop / textareaMaxScroll) * previewMaxScroll`，绑定 textarea 的 `scroll` / `click` / `keyup` 事件。

**理由**：Markdown source → rendered HTML 不是严格的行对行映射（代码块、列表等会膨胀），比例映射足够实用，实现简单。

**替代方案**：逐行 AST 映射 — 精确但复杂度高，收益有限，不采用。

### 5. 用户气泡 Markdown 样式通过 Tailwind 子选择器覆盖

**决策**：在 `MessageBubble.tsx` 的用户气泡 div 上添加 `[&_a]:`、`[&_pre]:` 等 CSS 覆盖，调整链接/代码块在 `bg-primary-soft` 上的颜色。

**理由**：不修改 MarkdownRenderer 通用组件，避免影响 Agent 消息渲染。子选择器局部作用域安全。

## Risks / Trade-offs

- **[受控组件性能]** textarea 每次输入触发 setState → 重渲染 → MarkdownRenderer 重解析 → 可能卡顿 → **缓解**：预览区使用 `useMemo` + `debounce(150ms)` 缓冲 Markdown 解析，textarea 本身不受影响
- **[用户气泡 MD 样式冲突]** prose 默认样式可能在 `bg-primary-soft` 上可读性差 → **缓解**：通过子选择器覆盖关键元素颜色，已在 Demo 中验证
- **[双栏模式下的 @提及]** 弹出列表的定位需要适配双栏布局 → **缓解**：mention popup 已是 absolute 定位，不受双栏影响
