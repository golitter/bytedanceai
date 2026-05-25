# 架构概览

## 设计思路

将原来的 Task 管理视角改为 IM 聊天体验。左侧为对话列表（搜索 + 新建），右侧为聊天区（Agent 头像 + 消息列表 + 输入框）。整体采用双栏布局，暗色主题。

## 目录结构

```
src/
├── main.tsx                          # 应用入口：StrictMode + QueryClient + BrowserRouter
├── index.css                         # 全局样式：Tailwind + CSS 变量暗色主题
│
├── pages/
│   └── ImPage.tsx                    # 主页面：ConversationList + ChatArea 双栏布局
│
├── components/
│   ├── im/                           # 对话列表侧栏
│   │   ├── ConversationList.tsx      # 侧栏容器：Header + 搜索 + 列表 + 新建弹窗
│   │   ├── ConversationItem.tsx      # 单条对话：头像 + 名称 + 时间 + 最后消息预览
│   │   └── NewChatDialog.tsx         # 新建对话弹窗：选择 Agent 类型
│   │
│   ├── chat/                         # 聊天区
│   │   ├── ChatArea.tsx              # 聊天容器：Header + 消息列表 + 输入框
│   │   ├── MessageList.tsx           # 消息列表（支持虚拟滚动）
│   │   ├── MessageBubble.tsx         # 消息气泡（user / agent / system 三种变体）
│   │   ├── MessageInput.tsx          # 输入框（自动高度 + Enter 发送 + Shift+Enter 换行）
│   │   ├── AgentAvatar.tsx           # Agent 头像（颜色 + 状态指示灯）
│   │   └── AgentEditDialog.tsx       # Agent 编辑弹窗（修改名称 + 上传头像）
│   │
│   ├── markdown/                     # Markdown 渲染
│   │   ├── MarkdownRenderer.tsx      # react-markdown + remark-gfm + 自定义组件
│   │   └── CodeBlock.tsx             # 代码块（Shiki 高亮 + 行号）
│   │
│   └── ui/                           # shadcn/ui 基础组件
│       ├── button.tsx
│       ├── card.tsx
│       ├── dialog.tsx
│       └── input.tsx
│
├── hooks/
│   ├── use-chat-stream.ts            # 聊天流：useReducer 状态机 + SSE 连接
│   ├── use-conversations.ts          # 对话列表查询 + 新建 mutation
│   └── use-hover-style.ts            # 悬停样式工具 hook
│
├── lib/
│   ├── api.ts                        # REST API 封装（Task / Session / Avatar 上传）
│   ├── sse.ts                        # SSE 客户端（fetch ReadableStream 解析）
│   ├── constants.ts                  # 常量定义（AGENT_NAMES / AGENT_DESCRIPTIONS）
│   └── utils.ts                      # cn() 工具函数
│
├── stores/
│   └── chat.ts                       # 聊天导航状态（currentSessionId）
│
├── generated/                        # 契约生成的类型文件（勿手改）
│   ├── events.ts                     # StreamEvent 类型
│   ├── request.ts                    # AgentType 等请求类型
│   ├── response.ts                   # AgentResponse 类型
│   └── session.ts                    # SessionState 类型
│
└── assets/                           # 静态资源
```

## 页面路由

| 路径 | 组件 | 说明 |
|------|------|------|
| `/*` | `ImPage` | 当前只有一个页面，所有路径都渲染 IM 聊天界面 |

路由在 `main.tsx` 中通过 `BrowserRouter` + `Routes` 配置，后续可根据需要添加 `/settings`、`/history` 等路由。

## 构建配置

**Vite** (`vite.config.ts`)：
- 插件：`@vitejs/plugin-react` + `@tailwindcss/vite`
- 路径别名：`@` → `./src`
- 开发代理：`/api` → `http://localhost:8080`（后端 Go 服务）

## 依赖总览

| 类别 | 库 | 用途 |
|------|------|------|
| 框架 | React 19 + React Router 7 | UI 渲染 + 路由 |
| 样式 | Tailwind CSS 4 + shadcn/ui | 原子化 CSS + 组件库 |
| 状态 | Zustand | 全局轻量状态 |
| 数据 | TanStack React Query | 服务端状态管理 |
| Markdown | react-markdown + remark-gfm | Markdown 渲染 |
| 代码高亮 | Shiki | VS Code 级别语法高亮 |
| 虚拟滚动 | @tanstack/react-virtual | 大量消息时的性能优化 |
| 图标 | Lucide React | 统一图标系统 |
| 字体 | @fontsource-variable/geist | Geist Variable 字体 |
