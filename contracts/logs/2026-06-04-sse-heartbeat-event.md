# 契约变更：SSE heartbeat 事件入契约

## 变更原因

Backend SSE stream handler 会发送 `heartbeat` 保活事件。此前该事件没有出现在 `contracts/schemas/event-types.yaml` 中，Frontend 只能通过 default 分支静默忽略，三端契约不一致。

本次将 `heartbeat` 加入 EventType 枚举，并重新生成 Python / TypeScript / Go 类型。保留 data event 形式而不是改成 SSE comment，是因为 Frontend 的 stale watchdog 依赖 `onmessage` 更新活跃时间。

## 变更文件

- `contracts/schemas/event-types.yaml`
- `agentend/src/generated/events.py`
- `frontend/src/generated/events.ts`
- `backend/internal/generated/events.go`

## 跨端影响

- **Backend**：`heartbeat` 成为合法 SSE 事件类型。
- **Frontend**：可显式处理 `EventTypeValues.Heartbeat`，保持连接活跃且不触发业务渲染。
- **AgentEnd**：生成类型同步，无行为变更。
