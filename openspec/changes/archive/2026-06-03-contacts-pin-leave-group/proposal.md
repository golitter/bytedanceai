## Why

当前 IM 系统的会话列表平铺展示所有会话，缺乏分组管理和排序控制。RightSidebar 中"置顶会话"和"退出群聊"按钮均为 TODO 状态。IconSidebar 中"通讯录"按钮已存在但 disabled。需要实现这三个功能来补全 IM 体验，使项目的会话管理能力达到可演示水平。

## What Changes

- **置顶会话**：复用已有 `Task.pinned_at` 字段和 `updateTaskPin` API，前端接入置顶/取消置顶交互，ConversationList 按置顶优先排序，ConversationItem 显示置顶图标
- **退出群聊（彻底删除）**：三端联动清理——AgentEnd 暴露 task 级清理 API（worktree + Git 分支），Backend 新增 `DELETE /tasks/:taskId/leave` 协调 AgentEnd 清理 + 数据库级联删除，前端二次确认对话框
- **通讯录页面**：新增 `contact_group` + `contact_group_item` 数据库表，Backend 提供分组 CRUD API，前端新建 ContactsPage 组件（分组管理 + 未分组 + 新建分组），启用 IconSidebar 通讯录导航

## Capabilities

### New Capabilities
- `pin-conversation`: 会话置顶/取消置顶，列表排序调整，置顶图标展示
- `leave-group`: 退出群聊（彻底删除），三端联动清理 worktree/分支/数据库记录
- `contacts-page`: 通讯录页面，会话分组管理（创建/删除/重命名分组，分配会话到分组）

### Modified Capabilities
<!-- 无需修改现有 spec -->

## Impact

- **Backend (Go)**: 新增 `contact_group.go` model + handler，修改 `task.go`（ListTasks 排序 + LeaveTask），修改 `main.go`（路由注册 + AutoMigrate），新增 `agentend_client` 方法（DestroySession + CleanupByTask）
- **AgentEnd (Python)**: 新增 `DELETE /v1/workspace/task/{task_id}` 端点
- **Frontend (React)**: 新增 ContactsPage + use-contact-groups hook，修改 ConversationItem（置顶图标）、RightSidebar（置顶/退群按钮接入）、ImPage（路由 + Props）、IconSidebar（启用通讯录），api.ts 新增 leaveTask + contact group API
- **数据库**: 新增 `contact_groups` + `contact_group_items` 两张表
- **契约层**: 无需修改 contracts/schemas/*.yaml（这些功能不涉及三端共享协议变更），仅需在 contracts/logs/ 添加变更记录
