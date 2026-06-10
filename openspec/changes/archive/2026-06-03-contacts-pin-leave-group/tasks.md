## 1. Phase 1 — 置顶会话（纯前端 + 微量后端）

- [x] 1.1 Backend `task.go` ListTasks 排序改为 `ORDER BY pinned_at DESC NULLS LAST, created_at DESC`
- [x] 1.2 Frontend `api.ts` Task 接口添加 `pinned_at`，Conversation 接口添加 `pinnedAt`
- [x] 1.3 Frontend `api.ts` fetchConversations 映射 `pinnedAt` 字段，调整排序逻辑（置顶优先 + 时间降序）
- [x] 1.4 Frontend `ConversationItem.tsx` 添加 Pin 图标（`pinnedAt` 存在时显示旋转 45° 的 Lucide Pin）
- [x] 1.5 Frontend `RightSidebar.tsx` 置顶按钮 onClick 接入 `updateTaskPin`，动态切换"置顶会话"/"取消置顶"文案和高亮状态，Props 新增 `pinnedAt`
- [x] 1.6 Frontend `ImPage.tsx` 将当前会话的 `pinnedAt` 传递给 RightSidebar

## 2. Phase 2 — 退出群聊（三端联动）

- [x] 2.1 AgentEnd `workspace.py` 新增 `DELETE /v1/workspace/task/{task_id}` 端点，封装 `cleanup_by_task()`
- [x] 2.2 Backend `agentend_client/client.go` 新增 `DestroySession(sessionID)` 和 `CleanupByTask(taskID)` 方法
- [x] 2.3 Backend `task.go` 新增 `LeaveTask` handler（`DELETE /api/tasks/:taskId/leave`）：查 sessionID → 调 AgentEnd 清理 → 数据库事务级联删除
- [x] 2.4 Backend `main.go` 注册 `DELETE /tasks/:taskId/leave` 路由
- [x] 2.5 Frontend `api.ts` 新增 `leaveTask(taskId)` 函数
- [x] 2.6 Frontend `RightSidebar.tsx` 退出群聊按钮：二次确认对话框 → 调用 `leaveTask` → `clearNavigation` → `invalidateQueries`
- [x] 2.7 验证：创建群聊 → 发消息 → 退出群聊 → 确认数据库无残留、AgentEnd 无残留 worktree/分支

## 3. Phase 3 — 通讯录页面

- [x] 3.1 Backend `model/contact_group.go` 新增 ContactGroup 和 ContactGroupItem 模型
- [x] 3.2 Backend `main.go` AutoMigrate 追加 ContactGroup 和 ContactGroupItem
- [x] 3.3 Backend `handler/contact_group.go` 实现 CRUD handler（ListGroups、CreateGroup、UpdateGroup、DeleteGroup、AddItem、RemoveItem）
- [x] 3.4 Backend `main.go` 注册 `/api/contact-groups` 路由组（6 个端点）
- [x] 3.5 Frontend `api.ts` 新增 ContactGroup 接口和 6 个 API 函数
- [x] 3.6 Frontend `hooks/use-contact-groups.ts` 新增 useContactGroups query + 4 个 mutation hooks
- [x] 3.7 Frontend `components/im/ContactsPage.tsx` 新建通讯录页面组件（搜索栏 + 置顶分组 + 自定义分组 + 未分组 + 新建分组）
- [x] 3.8 Frontend `IconSidebar.tsx` 移除 contacts NavItem 的 `disabled` 属性
- [x] 3.9 Frontend `ImPage.tsx` 将 contacts 的 PlaceholderPage 替换为 `<ContactsPage />`
- [x] 3.10 验证：创建分组 → 分配会话 → 删除分组 → 从通讯录点击会话卡片跳转聊天
