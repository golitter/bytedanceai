## Why

前端目前只有脚手架验证代码（counter demo、简单 Task/Session 列表），没有实际的 Chat 功能。Phase 2 需要实现完整的聊天界面：用户发送消息，通过 SSE 接收 Agent 流式回复，并具备 Markdown 渲染、代码语法高亮、虚拟列表能力。这是 AgentHub 作为 AI Runtime Workspace 的核心交互入口。

## What Changes

- **移除全部脚手架代码**：App.tsx counter demo、TaskList/TaskDetail/SessionList 验证页面
- **构建完整 Chat UI**：三栏布局（侧边栏 + 聊天区 + 预留详情面板），消息发送与 SSE 流式回复
- **新增 SSE 流式连接层**：fetch + ReadableStream 对接 `POST /api/tasks/:taskId/run`，解析 8 种 StreamEvent
- **新增 Chat 状态机**：discriminated union + useReducer，管理 idle → streaming → done/error 流转
- **消息渲染体系**：Markdown 渲染（react-markdown + remark-gfm）、代码语法高亮（shiki Tokyo Night）、虚拟列表（@tanstack/react-virtual）
- **适配 Visual Style Guide**：shadcn CSS 变量映射到 Dark Utilitarian 色值体系，暗色优先
- **路由整合**：`/tasks/:taskId` 作为 Chat 入口，与后端 Task → Session 数据模型对齐

## Capabilities

### New Capabilities

- `chat-streaming`: SSE 连接、StreamEvent 解析、Chat 状态机（idle/loading/streaming/tool_running/done/error）、消息累加逻辑
- `chat-ui`: 三栏布局、ChatSidebar（会话列表 + 新建）、ChatArea（消息区 + 输入区）、MessageBubble（user/agent/system 三态）、AgentAvatar（标识色 + 状态灯）、MessageInput（自适应高度 + 快捷键）
- `message-rendering`: Markdown 渲染、GFM 支持、代码块语法高亮（shiki Tokyo Night）、虚拟列表滚动、自动滚动行为
- `theme-adaptation`: CSS 变量映射 visual-style-guide、暗色优先双主题基础、字体/间距/动效约束

### Modified Capabilities

（无已有 spec 需要修改）

## Impact

- **重写文件**：App.tsx、main.tsx、index.css、lib/api.ts、stores/app.ts、pages/TaskList.tsx、pages/TaskDetail.tsx、components/SessionList.tsx
- **新增依赖**：react-markdown、remark-gfm、shiki、@tanstack/react-virtual
- **保留依赖**：shadcn/ui、Tailwind CSS 4、Zustand 5、TanStack React Query 5、React Router 7、Lucide React
- **后端 API 依赖**：POST /api/tasks、GET /api/tasks、GET /api/tasks/:taskId、POST /api/tasks/:taskId/run、GET /api/agent-types、PATCH /api/sessions/:sessionId
- **开发实践约束**：遵循 frontend/docs/development-strategy.md 的全部指导原则（三层组件模型、状态管理策略、性能策略、Hook 规范等）
