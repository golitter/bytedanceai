# Architecture — 架构与目录结构

## 实现了什么

基于 React 19 + Vite 8 的单页应用，采用 IM 聊天体验的双栏布局。左侧为对话列表侧栏，右侧为聊天区。使用 Zustand 管理聊天导航与会话状态，TanStack React Query 管理服务端数据缓存。

## 怎么实现的

### 应用入口 (`src/main.tsx`)

顶层挂载 `StrictMode` + `QueryClientProvider` + `BrowserRouter`，所有路由指向 `ImPage`：

```tsx
const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/*" element={<ImPage />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
)
```

### 主页面 (`src/pages/ImPage.tsx`)

双栏布局编排：`ConversationList` 固定宽度侧栏 + `ChatArea` 弹性填充主区域。通过 `useConversations()` 获取对话列表，通过 `useChatNav()` 读取当前选中会话：

```tsx
export function ImPage() {
  const { data: conversations } = useConversations()
  const { currentSessionId } = useChatNav()

  const active = conversations?.find((c) => c.sessionId === currentSessionId)

  return (
    <div className="flex h-screen bg-background">
      <ConversationList />
      <div className="flex-1">
        {active ? (
          <ChatArea
            taskId={active.taskId}
            sessionId={active.sessionId}
            agentType={active.agentType}
            agentName={active.agentName || undefined}
            avatarUrl={active.avatarUrl}
            repoPath={active.repoPath}
          />
        ) : (
          <div className="flex h-full flex-col items-center justify-center gap-3">
            <MessageSquare className="h-10 w-10 text-tertiary" strokeWidth={1} />
            <p className="text-sm text-tertiary">选择一个对话开始聊天</p>
          </div>
        )}
      </div>
    </div>
  )
}
```

### 构建配置 (`vite.config.ts`)

Vite 插件链为 `@vitejs/plugin-react` + `@tailwindcss/vite`，路径别名 `@` 指向 `./src`，开发代理将 `/api` 转发到后端 Go 服务：

```ts
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

### 目录结构

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
│   │   ├── ConversationItem.tsx      # 单条对话：头像 + 名称 + 时间
│   │   └── NewChatDialog.tsx         # 新建对话弹窗：仓库路径校验 + 选择 Agent 类型
│   │
│   ├── chat/                         # 聊天区
│   │   ├── ChatArea.tsx              # 聊天容器：Header + 消息列表 + 输入框
│   │   ├── MessageList.tsx           # 消息列表（支持虚拟滚动）
│   │   ├── MessageBubble.tsx         # 消息气泡（user / agent / system 三种变体）
│   │   ├── MessageInput.tsx          # 输入框（自动高度 + Enter 发送）
│   │   ├── AgentAvatar.tsx           # Agent 头像（颜色 + 状态指示灯）
│   │   └── AgentEditDialog.tsx       # Agent 编辑弹窗（修改名称 + 上传头像）
│   │
│   ├── markdown/                     # Markdown 渲染
│   │   ├── MarkdownRenderer.tsx      # react-markdown + remark-gfm + 自定义组件
│   │   └── CodeBlock.tsx             # 代码块（Shiki 高亮 + 行号）
│   │
│   └── ui/                           # shadcn/ui 基础组件
│       └── dialog.tsx
│
├── hooks/
│   ├── use-chat-stream.ts            # 聊天流：SSE 连接 + store actions 驱动状态
│   ├── use-conversations.ts          # 对话列表查询 + 新建 mutation
│   └── use-hover-style.ts            # 悬停样式工具 hook
│
├── lib/
│   ├── api.ts                        # REST API 封装
│   ├── sse.ts                        # SSE 客户端（EventSource 封装）
│   ├── constants.ts                  # 常量定义（AGENT_NAMES / AGENT_DESCRIPTIONS）
│   └── utils.ts                      # cn() 工具函数
│
├── stores/
│   └── chat.ts                       # Zustand Store：聊天导航 + 各会话独立状态
│
└── generated/                        # 契约生成的类型文件（勿手改）
    ├── events.ts                     # StreamEvent 类型
    ├── message.ts                    # 消息类型
    ├── request.ts                    # AgentType 等请求类型
    ├── response.ts                   # AgentResponse 类型
    ├── session.ts                    # SessionState 类型
    └── validate-repo-path.ts         # 仓库路径校验请求类型
```

### 依赖总览

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
