## Context

当前聊天区为三栏布局：`IconSidebar(56px)` → `ConversationList(280px)` → `ChatArea(flex-1)`。群聊会话已有基础能力（`isGroupChat` + `groupAgentTypes`/`groupAgentNames`/`groupSessions`），但缺少群信息展示、公告、成员列表和消息搜索。

前端使用 React + Zustand 状态管理 + Tailwind CSS，后端使用 Go/Gin + GORM + MySQL。数据库迁移通过 `main.go` 中的 `AutoMigrate` 自动处理。

视觉参考：[chat-enhanced-demo.html](../../docs/common/dev-plan/phase5-2-chat-enhanced/chat-enhanced-demo.html)，CSS 变量与 `frontend/src/index.css` dark 主题完全对齐。

## Goals / Non-Goals

**Goals:**
- 在群聊会话右侧新增 280px 常驻侧栏，包含搜索、公告、成员、操作四个区块
- 公告系统完整 CRUD：后端 API + 前端组件，支持置顶
- 历史消息前端搜索：过滤已加载消息，高亮匹配，点击跳转
- 成员列表从现有 props 提取，显示角色 + 在线状态
- 所有可折叠区块状态持久化到 localStorage

**Non-Goals:**
- 后端搜索 API（消息量 <1000 时前端过滤足够，后续再考虑）
- 公告编辑/更新（只做创建和删除）
- 成员管理（踢人/加人）—— 仅展示
- 单聊右侧栏——仅群聊显示
- AgentEnd 改动

## Decisions

### D1: RightSidebar 挂载位置

**决策**: 在 `ImPage.tsx` 中挂载 RightSidebar，而非 ChatArea 内部。

**理由**: RightSidebar 需要访问 `taskId`、`groupSessions` 等顶层 props，且需要影响 ChatArea 的宽度。放在 ImPage 层级可以避免 props 层层传递。ChatArea 保持 `flex-1`，RightSidebar 使用 `flex-shrink-0 w-[280px]`。

**替代方案**: 放在 ChatArea 内部——会增加 ChatArea 复杂度，且 ChatArea 已承担消息渲染职责。

### D2: 公告数据模型

**决策**: 使用 `task_announcements` 表，关联 `task_id`，字段包括 `sender_id`、`sender_name`、`content`、`pinned`、`created_at`。

**理由**: 公告是 task（群聊会话）级别的资源，一个群聊可以有多条公告。使用独立表而非 JSON 字段，便于排序、分页和单独删除。

**替代方案**: 存为 session 的 system message——会与普通消息混淆，不利于独立管理。

### D3: 历史搜索实现方式

**决策**: 前端过滤 `chatStore` 中已加载的消息，匹配 `message.content` 和 `block.content`。

**理由**: 当前消息量通常 <1000 条，前端过滤足够快，无需额外后端 API。避免增加后端复杂度。

**替代方案**: 后端搜索 API——过早优化，后续消息量 >1000 时再引入。

### D4: 成员在线状态判断

**决策**: 通过 SSE 连接状态推断——有活跃 stream 的成员显示 online，否则 offline。

**理由**: 无需额外心跳或 WebSocket 连接，复用现有 SSE 基础设施。

### D5: 数据库迁移策略

**决策**: 使用 GORM AutoMigrate，新增 `task_announcements` 表和 `tasks.pinned_at` 字段。

**理由**: 项目已有 AutoMigrate 模式，保持一致。`pinned_at` 使用 `DATETIME DEFAULT NULL`，NULL 表示未置顶。

## Risks / Trade-offs

- **[性能]** 消息量 >1000 时前端搜索变慢 → 使用 debounce 300ms + 限制结果 50 条，后续迁移到后端搜索
- **[布局]** 280px 侧栏在小屏幕上挤压 ChatArea → 设置最小屏幕宽度阈值，或添加侧栏折叠按钮
- **[SSE 状态]** SSE 连接状态不完全等于成员在线 → 接受近似判断，后续可引入更精确的心跳机制
- **[数据一致性]** AutoMigrate 在生产环境可能不够安全 → 当前为开发阶段，可接受；后续引入版本化迁移
