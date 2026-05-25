# 2026-05-26: Diff 快照持久化 + SessionAgent 拆表

## 变更类型: ADD

## 影响范围
- Backend: 新增 `diff_snapshots`、`session_agents` 表及对应 handler
- Frontend: DiffCard 改为 snapshotId 驱动
- Agent: card_diff.go 输出 snapshotId

## 说明
- 新增 `diff_snapshots` 表存储每个 diff block 的快照状态（pending/committed/reverted/cancelled）
- 新增 `session_agents` 表从 Session 拆分 agent 信息（1:1 关系），后端 JOIN 保持 API 兼容
- Agent render skill 输出 diff block 时携带 snapshotId（UUID v4）
- 前端 DiffCard 首次渲染创建 pending snapshot，用户操作后更新状态
- 现有契约 schema 无需修改（API 响应结构不变）
