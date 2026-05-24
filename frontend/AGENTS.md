# AGENTS.md — frontend

基于 React 19 + Vite + TypeScript 的前端项目，使用 Tailwind CSS + shadcn/ui 组件库，Zustand 状态管理，TanStack React Query 数据请求。包管理 pnpm。

## 目录结构

```
src/
├── main.tsx                # 应用入口（StrictMode + QueryClient + BrowserRouter）
├── index.css               # 全局样式（Tailwind + shadcn 语义 token + 自定义扩展 token）
├── assets/                 # 静态资源
├── components/
│   ├── chat/               # 聊天模块（ChatArea, MessageBubble, MessageInput, MessageList, AgentAvatar, AgentEditDialog）
│   ├── im/                 # 会话列表模块（ConversationList, ConversationItem, NewChatDialog）
│   ├── markdown/           # Markdown 渲染（MarkdownRenderer, CodeBlock）
│   └── ui/                 # shadcn/ui 基础组件（button, card, input, dialog）
├── pages/
│   └── ImPage.tsx          # 主页面（会话列表 + 聊天区布局）
├── hooks/
│   ├── use-chat-stream.ts  # 聊天 SSE 流式数据 hook
│   ├── use-conversations.ts # 会话列表 CRUD hook
│   └── use-hover-style.ts  # 通用 hover 样式 hook
├── stores/
│   └── chat.ts             # Zustand store（聊天导航 + 消息状态）
├── lib/
│   ├── api.ts              # API 请求封装
│   ├── constants.ts        # 常量（Agent 名称、描述）
│   ├── sse.ts              # SSE 连接工具
│   └── utils.ts            # cn() 工具函数
└── generated/              # 契约生成的 TypeScript 类型（events, request, response, session）
```

## 常用命令

> 通过根目录 Makefile 统一管理，需在项目根目录执行。

```bash
make run-frontend          # 启动（热重载）
make stop-frontend         # 停止
make restart-frontend      # 重启
make status                # 查看状态
```

如需手动启动：`cd frontend && pnpm dev`

- Makefile 完整说明：[docs/common/makefile-guide.md](../docs/common/makefile-guide.md)

## 详细文档

- 技术栈详情：[docs/tech-stack.md](docs/tech-stack.md)
