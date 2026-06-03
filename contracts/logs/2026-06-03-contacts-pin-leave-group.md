# 通讯录分组 + 置顶会话 + 退群 — 契约影响评估

> 日期: 2026-06-03
> 结论: **无需修改 schemas/**，仅记录跨端影响评估

## 变更原因

实现三个 IM 功能补全用户体验：
1. **通讯录分组** — 自定义分组管理会话
2. **置顶会话** — 置顶/取消置顶，排序优先
3. **退群（Leave Task）** — 三端联动彻底删除会话及关联资源

## 变更文件

无。`contracts/schemas/` 未修改。

## 跨端影响评估

### 1. 通讯录分组

| 端 | 变更 | 影响 |
|----|------|------|
| Backend | 新增 `contact_group` + `contact_group_item` 模型，6 个 REST 端点 | 仅 Backend ↔ Frontend |
| Frontend | 新增 `ContactsPage` 组件、`use-contact-groups` hooks、API 函数 | 纯 UI 层 |
| AgentEnd | 无变更 | — |

**不涉及 schemas**：通讯录分组是 Backend 独有数据模型，不经 SSE/Agent 请求协议。

### 2. 置顶会话

| 端 | 变更 | 影响 |
|----|------|------|
| Backend | `ListTasks` 排序改为 `pinned_at IS NULL, pinned_at DESC` | 复用已有 `Task.pinned_at` 字段 |
| Frontend | `Conversation.pinnedAt` 映射 + 排序 + Pin 图标 | 纯前端展示 |
| AgentEnd | 无变更 | — |

**不涉及 schemas**：`pinned_at` 字段已存在于 Task 模型，`updateTaskPin` API 已有。

### 3. 退群（Leave Task）

| 端 | 变更 | 影响 |
|----|------|------|
| Backend | 新增 `DELETE /tasks/:taskId/leave`，调用 AgentEnd 清理 API | 编排层 |
| Backend | `agentend_client` 新增 `DestroySession` / `CleanupByTask` / `CleanupTaskBranches` | HTTP 客户端 |
| AgentEnd | 新增 `DELETE /v1/workspace/task/{task_id}` | 封装已有 `cleanup_by_task()` |
| AgentEnd | 新增 `POST /v1/workspace/task/{task_id}/cleanup-branches` | 强制清理分支 |
| Frontend | 新增 `leaveTask()` API + RightSidebar 按钮接入 | 纯 UI 层 |

**不涉及 schemas**：新增的 AgentEnd 端点是 REST 操作型 API（清理 worktree/分支），不是 SSE 事件类型、Agent 请求/响应协议或会话状态变更。这些端点不引入新的枚举值或跨端共享类型。

## 设计文档

- [docs/design/06-contacts-pin-leave-group.md](../../docs/design/06-contacts-pin-leave-group.md)
