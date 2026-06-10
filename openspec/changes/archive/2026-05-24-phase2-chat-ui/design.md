## Context

AgentHub 前端基于 React 19 + Vite 8 + TypeScript 6 + Tailwind CSS 4 + shadcn/ui + Zustand 5 + TanStack React Query 5 + React Router 7。当前代码全是脚手架验证，可直接重写。

后端 API 已完整（Go Gin）：
- `POST /api/tasks` / `GET /api/tasks` / `GET /api/tasks/:taskId` / `DELETE /api/tasks/:taskId`
- `POST /api/tasks/:taskId/run` — SSE 流式聊天（text/event-stream）
- `GET /api/agent-types` / `PATCH /api/sessions/:sessionId`

数据模型：Task (1) → Session (N)，每次发消息走 `POST /tasks/:taskId/run` + session_id。

视觉规范：Dark Utilitarian — 5 级灰阶背景、单一品牌色 #6366F1、3 个 Agent 标识色、禁止渐变/毛玻璃/纯白文字。

开发实践约束（development-strategy.md）：
- 三层组件模型：Page → Smart → Dumb
- 三类状态归属：Server State (TanStack Query)、Global Client State (Zustand)、Local State (useState/useReducer)
- 不把 Server State 塞进 Zustand
- 抽象是负债，第三次重复再抽象
- streaming 更新只追加最后一条消息的 DOM，不重建整个列表

## Goals / Non-Goals

**Goals:**
- 完整 Chat 界面：发送消息 + SSE 流式回复 + 会话管理
- Markdown 渲染 + 代码语法高亮（shiki Tokyo Night）
- 虚拟列表（>50 条消息时启用）
- 视觉风格严格遵循 visual-style-guide.md
- 三栏布局（侧边栏 + 聊天区 + 预留详情面板）

**Non-Goals:**
- 工具调用卡片（tool_call/tool_result）富渲染 — Phase 3
- Artifact 文件卡片 — Phase 3
- Planning 思考过程展示 — Phase 3
- 多 Agent 协作时间线 — Phase 3
- 消息编辑/重新生成/搜索 — Phase 3
- 离线支持 / Event Replay — Phase 3
- 亮色主题切换 UI — 保留 CSS 变量基础，不提供切换按钮

## Decisions

### D1: CSS 变量映射 — 改值不改名

保留 shadcn 变量名（`--background`、`--card` 等），把 `.dark` 块中的值替换为 visual-style-guide 色值。同时新增 visual-style-guide 特有变量（`--bg-canvas`、`--color-brand`、Agent 标识色等）供自定义组件使用。

映射：
```
--background  → #0A0B0E (canvas)    --card  → #1A1D24 (card)
--popover     → #22262F (hover)     --border → rgba(255,255,255,0.06)
--foreground  → #E8EBF0 (primary)   --muted-foreground → #8B91A0 (secondary)
--primary     → #6366F1 (brand)     --ring  → #6366F1
```

`:root` 保留合理的亮色值作为占位。`<html>` 硬编码 `.dark` class。

**替代方案**：创建全新变量体系 → 拒绝，会导致所有 shadcn 组件失效。

### D2: SSE 解析 — fetch + ReadableStream

原生 `fetch` + `ReadableStream` + `TextDecoder`。因为后端 `POST /api/tasks/:taskId/run` 需要 POST body，`EventSource` 只支持 GET。

关键点：
- 用 `\n\n` 分割 SSE 事件，维护 buffer 处理跨 chunk 边界
- 解析 `data: {...}` 行，JSON.parse 得到 StreamEvent
- 返回 AbortController 供取消

### D3: Chat 状态管理 — useReducer + Zustand 分层

遵循 development-strategy.md 的三类状态归属原则：
- **Chat 状态机**（streaming 状态、消息累加）→ `useReducer`，封装在 `useChatStream` hook 中（Smart 层）
- **全局导航状态**（currentSessionId、currentTaskId）→ Zustand
- **消息不放进 Zustand**，因为每个 token 更新会触发所有订阅者 re-render

状态机（discriminated union）：
```
idle → loading → streaming → done
                  ↓              ↑
            tool_running        error
```

每个状态携带的字段由 TypeScript 保证存在性，不需要到处判空。

### D4: 组件架构 — 三层模型

遵循 development-strategy.md 的三层组件模型：

| 层级 | 组件 | 职责 |
|------|------|------|
| Page | ChatPage | 路由参数、三栏布局编排、数据组合（useQuery） |
| Smart | ChatArea | streaming 状态管理（useChatStream）、消息区+输入区编排 |
| Smart | ChatSidebar | 会话列表管理（TanStack Query CRUD） |
| Dumb | MessageBubble | 纯 props 渲染，user/agent/system 三态 |
| Dumb | AgentAvatar | 纯 props 渲染，标识色 + 状态灯 |
| Dumb | MessageInput | 纯 props + 本地 textarea 状态 |
| Dumb | MarkdownRenderer | 纯 props 渲染 Markdown + 代码高亮 |

MessageList 作为独立 Smart 组件隔离，因为它处理虚拟列表 + streaming diff + 滚动，这些高频逻辑不能污染整个页面。

### D5: Markdown + 代码高亮 — react-markdown + shiki

`react-markdown` + `remark-gfm` 渲染 Markdown，`shiki` 做代码语法高亮。

选择 shiki 而非 highlight.js：
- visual-style-guide 指定 Tokyo Night 主题，shiki 内置
- shiki 高亮质量 = VSCode 级别
- 通过 react-markdown 的 components prop 覆盖 `<pre>` 渲染

性能：shiki 动态 import 懒加载，先显示纯文本再替换为高亮结果。

### D6: 虚拟列表 — @tanstack/react-virtual

体积小（~3KB），headless 设计。支持动态高度项（消息长度不一）。与 streaming 更新配合：消息内容变化时 invalidate 测量缓存。

触发阈值：>50 条启用虚拟化，少于 50 条直接渲染。

### D7: 路由 — /tasks/:taskId 作为 Chat 入口

与后端 Task → Session 数据模型对齐。`/tasks` 列表页选择/创建 task，`/tasks/:taskId` 进入聊天。

### D8: 图标库 — 保留 Lucide React

shadcn/ui 默认搭配，已安装。`strokeWidth` 可调细线条弥补 visual-style-guide 的顾虑。

### D9: 数据请求策略

遵循 development-strategy.md：
- Task/Session 列表 → TanStack Query（缓存、自动刷新、窗口聚焦重验证）
- SSE streaming → useReducer（不用 TanStack Query 管理 streaming）
- Session CRUD → TanStack Query mutations（invalidateQueries 刷新列表）
- API 调用集中在 `lib/api.ts`

## Risks / Trade-offs

**[SSE 跨 chunk 边界]** → 用 buffer 累积未完成的 data 行，每次 append 新 chunk 后按 `\n\n` 分割。

**[shiki 首次加载慢]** → 动态 import + 先显示纯文本再替换。不阻塞渲染。

**[虚拟列表动态高度 + streaming]** → @tanstack/react-virtual 的 `measureElement` 在消息内容变化时自动重新测量。需在 streaming 更新后触 invalidate。

**[Streaming 只更新最后一条消息]** → useReducer 的 dispatch(text) 只修改当前 agent 消息的 content 字段，不重建 messages 数组。虚拟列表只重测最后一个元素。

**[Dark mode 固定]** → 先硬编码 `.dark`。亮色 CSS 变量值已定义在 `:root`，Phase 3 加切换按钮即可。

**[文件结构约定]** → 遵循 development-strategy.md：不在 shadcn/ui 上建 BaseXXX 层；抽象是负债，第三次重复再抽；纯计算函数放 `lib/`，不放 hooks。
