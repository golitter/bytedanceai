# 2026-06-02 — Orchestrator Merge Review Diff

## 变更原因

Orchestrator 需要在决定将 `task/{task_id}` 合入 `main` 前先展示 CodeDiff，并等待用户确认后再执行实际 git merge。此前 `plan_review` 事件只能表达普通规划审查，缺少合并审查类型、源/目标分支和 diff snapshot 关联字段，前端无法复用现有 DiffCard 渲染 task-to-main diff。

## 变更文件

- `agentend/src/adapters/orchestrator.py`
- `agentend/src/orchestrator/models.py`
- `agentend/src/orchestrator/planning/graph.py`
- `agentend/src/orchestrator/planning/prompts.py`
- `agentend/src/orchestrator/planning/tools.py`
- `agentend/src/workspace/git_ops.py`
- `agentend/src/workspace/manager.py`
- `backend/internal/stream/writer.go`
- `frontend/src/components/cards/PlanReviewCard.tsx`
- `frontend/src/components/chat/MessageBubble.tsx`
- `frontend/src/hooks/use-chat-stream.ts`
- `frontend/src/lib/block-reducer.ts`
- `frontend/src/lib/block-types.ts`

## 对比结果

- `PlanOutput` 新增 `merge_to_main` 决策字段，由 Orchestrator 规划是否请求合入 `main`。
- `plan_and_dispatch` 工具新增 `merge_to_main` 参数，允许“只请求合并 main”的用户输入生成零任务合并规划。
- AgentEnd 新增 task-to-main diff 获取能力，合并前发出 `plan_review` 审批事件。
- Backend 在收到带 `diff_snapshot_id` 和 `diff` 的 `plan_review` 事件时写入 `diff_snapshots`，消息块只持久化 snapshot id 等元数据。
- Frontend 的 `PlanReviewCard` 对 `review_type=merge_to_main` 使用现有 `DiffCard` 按 `diff_snapshot_id` 渲染 CodeDiff。

## 跨端影响

- **AgentEnd**: Orchestrator 只在 `merge_to_main=true` 时请求合并审查；用户通过后才调用 workspace merge。
- **Backend**: `plan_review` 持久化支持合并审查元数据，并复用 diff snapshot 存储。
- **Frontend**: 历史消息和实时 SSE 都能恢复 `review_type/source_branch/target_branch/diff_snapshot_id`，合并审批卡展示 task-to-main CodeDiff。
- **Contracts**: 现有 `contracts/schemas/event-types.yaml` 的 `StreamEvent.content.additionalProperties=true` 可承载本次 payload 扩展，无需修改 YAML schema 或运行 `make generate`。

## 契约变更

无 `contracts/schemas/*.yaml` 变更。

实际扩展的 `plan_review.content` 字段如下：

```json
{
  "review_key": "external_review:{session_id}",
  "review_type": "merge_to_main",
  "source_branch": "task/{task_id}",
  "target_branch": "main",
  "diff_snapshot_id": "{uuid_v4}",
  "diff": "{git diff main...task/{task_id}}"
}
```

约束：

- `review_type` 为空或缺省时表示普通规划审查。
- `review_type=merge_to_main` 时，`source_branch`、`target_branch`、`diff_snapshot_id` 必须存在。
- `diff_snapshot_id` 使用 UUID v4 字符串，不添加业务前缀，避免超过 `diff_snapshots.snapshot_id` 字段长度。
- `diff` 仅用于 Backend 创建 diff snapshot，不写入持久化消息块。
