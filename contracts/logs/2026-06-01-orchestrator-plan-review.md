# Orchestrator Plan Review

- **变更原因**：Orchestrator 需要在规划完成后、任务执行前暂停，等待用户批准或反馈修改意见。
- **变更文件**：`contracts/schemas/session-state.yaml`、`contracts/schemas/event-types.yaml`。
- **对比结果**：新增 `awaiting_review` 会话状态；新增 `plan_review` SSE 事件类型。
- **跨端影响**：AgentEnd 可发出规划审查事件并等待审查结果；Backend 可标记审查等待状态并转发审查请求；Frontend 可渲染规划审查卡片。
- **契约变更**：`SessionState.awaiting_review`；`EventType.plan_review`。
