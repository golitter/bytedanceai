## Why

当前聊天区域中，Agent 消息已使用 MarkdownRenderer 渲染富文本，但用户消息仍以纯文本展示，无法呈现 Markdown 格式（标题、列表、代码块、表格等）。同时输入栏是普通 textarea，用户粘贴 Markdown 后无法预览渲染效果，体验割裂。需要统一消息渲染能力并为输入栏增加双栏实时预览。

## What Changes

- 用户消息气泡内从纯文本改为使用 `MarkdownRenderer` 渲染，与 Agent 消息一致
- 输入栏上方新增 Markdown 模式切换按钮（默认关闭，保持现有行为）
- 开启 Markdown 模式后，输入区变为左右双栏布局：左栏 textarea 编辑 / 右栏 MarkdownRenderer 实时预览 / 最右侧发送按钮
- 双栏输入区自动增高（min 120px，max 60vh），内容超出时预览区按滚动比例同步跟随光标位置
- 用户气泡内 Markdown 样式需在 `bg-primary-soft` 背景上保持可读性

## Capabilities

### New Capabilities
- `markdown-input-preview`: 输入栏 Markdown 双栏实时预览模式 — 切换按钮、双栏布局、自动增高、光标同步滚动

### Modified Capabilities
- `chat-ui`: 用户消息渲染从纯文本改为 MarkdownRenderer
- `message-rendering`: 用户气泡内需适配 Markdown 样式（代码块、链接、引用块等在 `bg-primary-soft` 背景上的可读性）

## Impact

- **前端组件**：`MessageInput.tsx`（改为受控组件 + 双栏布局）、`MessageRenderer.tsx`（用户消息用 MarkdownRenderer）、`MessageBubble.tsx`（用户气泡 MD 样式微调）
- **依赖**：复用现有 `react-markdown`、`remark-gfm`、`shiki`，无新依赖
- **无后端/Agent 端变更**
