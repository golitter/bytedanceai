## Context

项目是 Monorepo 三端架构（Frontend React + Backend Go + AgentEnd Python），支持单聊和群聊。Backend 当前是"自有逻辑 + 透传代理 + 协调编排"三种模式混合：
- Task/Session/Message CRUD 由 Backend 自有 MySQL 管理
- Workspace 文件操作透传到 AgentEnd
- CreateTask/RunTask/Stream 是双端协调

当前会话列表平铺展示，缺乏分组和排序控制。已有 `Task.pinned_at` 字段和 `updateTaskPin` API，但前端未接入。退出群聊按钮存在但为 TODO。通讯录导航已定义但 disabled。

## Goals / Non-Goals

**Goals:**
- 置顶会话：最小改动，纯复用已有 API，前端接入排序 + 图标
- 退出群聊：三端联动彻底删除，确保 worktree/分支/数据库不留残留
- 通讯录：会话分组管理，启用通讯录导航

**Non-Goals:**
- 不做拖拽排序（通讯录分组内排序使用 sortOrder 字段，UI 用下拉菜单分配）
- 不做多用户权限（单用户系统，无需 user_id）
- 不修改 contracts/schemas/*.yaml 共享契约（这些功能不涉及三端共享协议）
- 不做会话归档/恢复（退出即彻底删除）

## Decisions

### D1: 置顶排序 — 服务端排序 vs 客户端排序

**选择：服务端排序** — 修改 `ListTasks` 的 `ORDER BY pinned_at DESC NULLS LAST, created_at DESC`。

理由：前端 `fetchConversations` 基于 `fetchTasks` 构建，排序在数据源头完成最简单，避免每个客户端重复实现排序逻辑。同时兼容未来其他端（如移动端）直接消费排序结果。

### D2: 退出群聊 — 先清 AgentEnd 还是先清数据库

**选择：先清 AgentEnd，再清数据库** — 但 AgentEnd 调用失败不阻断数据库删除。

理由：AgentEnd 的 worktree/分支是运行时资源，尽早释放更好。但即使 AgentEnd 调用超时或失败，后台 `_inactive_cleanup_loop` 兜底清理，不应阻断用户操作。数据库删除是最终态，必须保证成功。

### D3: 通讯录分组 — 多对多关联 vs Task 字段

**选择：多对多关联表** — 新增 `contact_groups` + `contact_group_items` 两张表。

理由：一个会话可以只属于一个分组（当前需求），但关联表模式更灵活（未来支持一会对多组、排序、批量移动），且不需要修改现有 Task 模型。单用户系统无需 user_id，未出现在关联表中的 task 隐式属于"未分组"。

### D4: 通讯录前端 — 独立 Tab 页 vs 聊天列表内嵌

**选择：独立 Tab 页** — IconSidebar 已预留 contacts 导航，新建 ContactsPage 组件。

理由：通讯录有独立的交互模式（分组管理、分配操作），与聊天列表的"快速切换"定位不同。独立页面更清晰，也匹配已有导航架构（NavTab 已定义 'contacts'）。

### D5: AgentEnd 清理 API — 新增端点 vs 复用已有

**选择：新增 `DELETE /v1/workspace/task/{task_id}`** — 封装已有 `WorkspaceManager.cleanup_by_task()`。

理由：`cleanup_by_task` 已存在但仅被后台 `_inactive_cleanup_loop` 调用，无 HTTP 端点暴露。新增一个 DELETE 端点是最小改动，复用全部已有清理逻辑。

## Risks / Trade-offs

- **[退出群聊数据不可恢复]** → 确认对话框明确列出将被删除的内容，二次确认防误操作
- **[AgentEnd 调用失败导致 worktree 残留]** → best-effort 策略，不阻断数据库删除；后台 `_inactive_cleanup_loop` 兜底清理不活跃 session 的 workspace
- **[通讯录分组表数据量]** → 单用户系统，分组数和会话数均有限，无需担心性能
- **[MySQL NULLS LAST 语法]** → MySQL 8.0+ 支持 `NULLS LAST`，需确认生产环境 MySQL 版本。若为 MySQL 5.7，需改用 `ORDER BY pinned_at IS NULL, pinned_at DESC, created_at DESC`
