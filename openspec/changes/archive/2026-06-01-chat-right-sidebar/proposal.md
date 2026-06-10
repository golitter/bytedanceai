## Why

群聊会话当前缺少信息展示区域——用户无法查看群公告、成员列表，也无法搜索历史消息。聊天区为三栏布局（IconSidebar → ConversationList → ChatArea），缺少群聊场景下的信息侧栏，导致群聊体验不完整。

## What Changes

- 新增第四栏 `RightSidebar(280px)`，仅群聊会话时显示，包含：历史消息搜索、群公告、群成员列表、更多操作
- 后端新增 `task_announcements` 表及 CRUD API（GET/POST/DELETE `/api/tasks/{taskId}/announcements`）
- 后端 `tasks` 表新增 `pinned_at` 字段，支持会话置顶
- 前端新增 4 个组件：`RightSidebar`、`HistorySearch`、`AnnouncementsSection`、`MembersSection`
- 前端 `ImPage.tsx` 从三栏变四栏布局
- 前端 `ChatArea.tsx` 集成右侧栏

## Capabilities

### New Capabilities

- `right-sidebar`: 聊天右侧栏容器组件，仅群聊时显示，包含搜索、公告、成员、操作四个区块
- `history-search`: 前端历史消息搜索功能，过滤已加载消息，高亮匹配关键词，点击跳转定位
- `announcements`: 群公告系统，后端 CRUD API + 前端展示组件，支持置顶、折叠、发布
- `members-section`: 群成员列表展示，从现有 props 提取数据，显示角色与在线状态，支持点击跳转

### Modified Capabilities

- `chat-ui`: 新增第四栏右侧栏，ImPage 布局从三栏扩展为四栏

## Impact

- **Frontend**: `ImPage.tsx` 布局变更、`ChatArea.tsx` 集成右侧栏、新增 4 个组件、`chat.ts` store 新增公告状态、`api.ts` 新增公告 API
- **Backend**: 新增 `announcement` model/handler、`task.go` 新增 pinned_at 更新、`main.go` 注册新路由、数据库新增 `task_announcements` 表和 `tasks.pinned_at` 字段
- **Contracts**: 无改动——公告走 REST API，不走 SSE 事件
- **AgentEnd**: 无改动
