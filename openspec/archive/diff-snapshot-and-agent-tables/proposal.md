## Why

当前 DiffCard 在用户接受/拒绝变更后调用 `refresh()` 获取 workspace diff，但 git 操作后 diff 为空，导致原始 diff 内容丢失、卡片变空白。此外，agent 信息（agent_type、agent_name、avatar_url）直接耦合在 Session 表中，违反关注点分离。需要一个独立的 diff 快照机制来持久化每个 diff 块的状态，以及将 agent 信息拆到独立表中。

## What Changes

- **新建 `diff_snapshots` 表**：以 agent 生成的 `snapshotId`（UUID）为主键，独立存储每个 diff 块的 unified diff 全文和状态（pending → committed/reverted/cancelled）
- **同一 session 同时只允许一个 pending diff**：新 diff 出现时，后端自动将同 session 下其他 pending 快照置为 cancelled
- **Agent 端 `card_diff.go` 输出 snapshotId**：diff block 从 `{ type: diff }` 变为 `{ type: diff, snapshotId: uuid }`
- **新建 `session_agents` 表**：从 Session 提取 agent_type、agent_name、avatar_url 为独立 1:1 表
- **Session 表瘦身**：移除 agent_type、agent_name、avatar_url 以及之前临时加的 settled_diff、diff_status 字段
- **DiffCard 用 snapshotId 驱动**：首次渲染时创建 pending snapshot 并展示 workspace diff；用户操作后更新为 committed/reverted；页面刷新后从 snapshot 恢复
- **后端 JOIN session_agents**：Session API 响应结构不变，前端对 agent 信息拆表无感知

## Capabilities

### New Capabilities
- `diff-snapshot`: diff 快照的创建、查询、状态流转（pending/committed/reverted/cancelled）及同 session 自动取消逻辑
- `session-agent-table`: agent 信息独立存储表，从 Session 拆分，后端 JOIN 保持 API 兼容

### Modified Capabilities
- `message-rendering`: diff block 类型从 `{ type: 'diff' }` 扩展为 `{ type: 'diff'; snapshotId: string }`
- `workspace-management`: agent render skill 输出 diff block 时携带 snapshotId

## Impact

- **Backend（Go）**: 新增 2 个 model（DiffSnapshot、SessionAgent），新增 diff-snapshot handler，修改 session/task/avatar handler 使用 JOIN，main.go 注册新路由
- **Frontend（React/TS）**: block-types.ts、block-reducer.ts 扩展 diff block，DiffCard 重写为 snapshotId 驱动，MessageBubble 传递 snapshotId
- **Agent（Go render skill）**: card_diff.go 生成 UUID 并输出到 diff block
- **Database**: AutoMigrate 新增两张表，Session 表删除 5 个字段（需确认数据迁移策略）
