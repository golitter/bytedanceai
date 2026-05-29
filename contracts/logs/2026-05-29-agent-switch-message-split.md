# 契约变更：Backend 多 Agent 消息拆分 + ListMessages session_id 过滤

## 变更原因

Orchestrator 一次 RunTask 中多个 Agent 的输出合并为一条 MySQL Message，`agent_type` 固定为初始值。页面刷新后所有内容混在一个气泡下。同时 ExecutionEngine 为子 Agent 单独调用 RunTask 产生跨 session 重复消息。

## 变更文件

本次变更 **未修改** `contracts/schemas/*.yaml`。`Message` schema 已有 `session_id`、`agent_type`、`agent_name` 字段，`StreamEvent.content` 已为 `additionalProperties: true`，新增行为在现有 schema 下合法。

## 变更内容

### 1. StreamWriter Agent 切换检测

`backend/internal/stream/writer.go` 的 `StreamWriter` 新增 `switchAgent()` 方法：
- 解析 SSE TEXT 事件中的 `agent_type`/`agent` 字段
- `agent_type` 变化时：flush buffer → finalize 当前 Message（status=completed）→ 创建新 Message（同 session，新 agent_type，status=streaming）→ 更新内部 messageID → 重置 buffer
- 原始 Message 保持 streaming 直到整轮结束（`finish()` 时标记 completed/failed）
- buffer 为空时仅更新 agent 信息，不创建空 Message
- Redis stream key 不变（所有事件继续写入原始 Message 的 stream）

### 2. SSE 重放携带 agent 元数据

`backend/internal/handler/stream.go` 的 `serveCompleted`、`serveStreaming` Phase 1、`serveFailed` 改用 `FormatSSEWithMeta(chunk, msg.AgentType, msg.AgentName)`，SSE 事件格式：

```json
{"type":"text","content":{"text":"...","agent_type":"claude-code","agent":"代码助手"}}
```

`agent_type` 为空时不包含 agent 字段，向后兼容。

### 3. ListMessages session_id 过滤

`backend/internal/handler/message.go` 的 `ListMessages` 新增可选 `session_id` query param：
- 传入时：`WHERE session_id = ?`，仅返回当前 session 的消息
- 不传入时：行为不变（返回 task 下所有 session 消息）

### 4. Frontend 传 session_id

`frontend/src/lib/api.ts` 的 `getTaskMessages` 新增 `sessionId` 参数。
`use-chat-stream.ts` 和 `ChatArea.tsx` 的历史加载和 loadMore 均传入 `sessionId`。

## 跨端影响

- **AgentEnd**: 无改动（SSE 事件格式已在 `2026-05-29-text-event-agent-metadata.md` 中确认）
- **Backend**: StreamWriter 自动拆分 Message；SSE 重放带 agent 元数据；ListMessages 支持 session_id 过滤
- **Frontend**: API 调用传入 sessionId，消除重复消息；SSE 重放事件中的 agent 元数据使前端正确归属消息
- **Contracts**: 无 schema 变更
