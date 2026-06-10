## 1. Backend — 数据模型与迁移

- [x] 1.1 创建 `backend/internal/model/announcement.go`，定义 `Announcement` 结构体（id, task_id, sender_id, sender_name, content, pinned, created_at）
- [x] 1.2 修改 `backend/cmd/server/main.go`，将 `&model.Announcement{}` 加入 AutoMigrate 列表
- [x] 1.3 修改 `backend/internal/model/task.go`，在 Task 结构体新增 `PinnedAt *time.Time` 字段

## 2. Backend — 公告 CRUD Handler

- [x] 2.1 创建 `backend/internal/handler/announcement.go`，实现 `AnnouncementHandler` 结构体
- [x] 2.2 实现 `ListAnnouncements` 方法 — `GET /api/tasks/:taskId/announcements`，按置顶优先+时间倒序返回
- [x] 2.3 实现 `CreateAnnouncement` 方法 — `POST /api/tasks/:taskId/announcements`，校验必填字段
- [x] 2.4 实现 `DeleteAnnouncement` 方法 — `DELETE /api/tasks/:taskId/announcements/:id`

## 3. Backend — 路由注册与 Task 更新

- [x] 3.1 修改 `backend/cmd/server/main.go`，注册公告路由（GET/POST/DELETE）
- [x] 3.2 修改 `backend/internal/handler/task.go`，在 UpdateTask 中支持 `pinned_at` 字段更新

## 4. Frontend — API 层

- [x] 4.1 在 `frontend/src/lib/api.ts` 新增 `fetchAnnouncements(taskId)` 函数
- [x] 4.2 新增 `createAnnouncement(taskId, data)` 函数
- [x] 4.3 新增 `deleteAnnouncement(taskId, announcementId)` 函数
- [x] 4.4 新增 `updateTaskPin(taskId, pinnedAt)` 函数

## 5. Frontend — Store 层

- [x] 5.1 在 `frontend/src/stores/chat.ts` 新增公告相关状态（announcements map、loading 状态）
- [x] 5.2 新增 `loadAnnouncements(taskId)` action
- [x] 5.3 新增 `addAnnouncement(taskId, data)` / `removeAnnouncement(taskId, id)` actions

## 6. Frontend — RightSidebar 容器组件

- [x] 6.1 创建 `frontend/src/components/chat/RightSidebar.tsx`，实现 280px 侧栏容器，接收 taskId、isGroupChat、groupAgentTypes、groupAgentNames、groupSessions props
- [x] 6.2 实现可折叠区块通用逻辑（标题栏点击折叠/展开，localStorage 持久化）

## 7. Frontend — 历史消息搜索组件

- [x] 7.1 创建 `frontend/src/components/chat/HistorySearch.tsx`，实现搜索输入框 + 300ms debounce
- [x] 7.2 实现前端消息过滤逻辑（匹配 message.content 和 block.content）
- [x] 7.3 实现搜索结果下拉列表（Agent 头像 + 名称 + 时间 + 摘要，`<mark>` 高亮关键词）
- [x] 7.4 实现点击结果跳转定位（滚动到消息 + 800ms 高亮闪烁动画）
- [x] 7.5 实现下拉关闭逻辑（点击外部、清空输入）
- [x] 7.6 实现结果上限 50 条 + 超出提示

## 8. Frontend — 群公告组件

- [x] 8.1 创建 `frontend/src/components/chat/AnnouncementsSection.tsx`，展示公告列表
- [x] 8.2 实现置顶 badge + 按置顶优先/时间倒序排序
- [x] 8.3 实现「+ 发布新公告」按钮（Owner/Admin 可见，弹出输入框 → 调用创建 API）
- [x] 8.4 实现公告删除（长按或右键菜单 → 调用删除 API）

## 9. Frontend — 群成员组件

- [x] 9.1 创建 `frontend/src/components/chat/MembersSection.tsx`，从 props 提取成员数据
- [x] 9.2 实现成员行（头像 + 名称 + 角色标签 + 在线状态圆点）
- [x] 9.3 实现角色 badge 样式（Owner: primary, Admin: warning）
- [x] 9.4 实现在线状态判断（基于 SSE stream 活跃状态）
- [x] 9.5 实现点击成员跳转（切换 activeSessionId 到对应 session）

## 10. Frontend — 布局集成

- [x] 10.1 修改 `frontend/src/pages/ImPage.tsx`，在群聊时渲染 RightSidebar，实现四栏布局
- [x] 10.2 实现更多操作按钮（导出聊天记录为 Markdown、置顶会话、退出群聊）

## 11. 验证

- [ ] 11.1 验证单聊不显示右侧栏、群聊显示右侧栏
- [ ] 11.2 验证公告 CRUD：创建、删除、置顶排序
- [ ] 11.3 验证历史搜索：关键词匹配、高亮、点击跳转
- [ ] 11.4 验证成员列表：角色显示、在线状态、点击跳转
- [ ] 11.5 验证折叠状态持久化（刷新页面后恢复）
