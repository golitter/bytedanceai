## 1. Backend StreamWriter Agent 切换检测

- [x] 1.1 StreamWriter 新增 `currentAgentType`、`currentAgentName`、`originalMessageID` 字段，`NewStreamWriter` 中初始化
- [x] 1.2 `Run()` 中解析 TEXT 事件的 `agent_type`/`agent` 字段，变化时调用 `switchAgent()`
- [x] 1.3 实现 `switchAgent()`：flush buffer → finalize 当前子 Message（status=completed）→ 创建新 Message（同 session，新 agent_type，status=streaming）→ 更新内部 messageID → 重置 buffer。buffer 为空时仅更新 agent 信息
- [x] 1.4 修改 `finish()`：将最后一条子 Message 和原始 Message 均设为 completed/failed，registry.Delete 使用 originalMessageID
- [x] 1.5 修改 `updateStatus()` → `updateMessageStatus(messageID, status)`，按指定 messageID 更新状态

## 2. Backend SSE 重放携带 agent 元数据

- [x] 2.1 新增 `FormatSSEWithMeta(text, agentType, agentName)` 函数，生成含 agent_type/agent 的 SSE 事件
- [x] 2.2 `serveCompleted` 使用消息自身 AgentType/AgentName 构造 SSE 事件（替换 `FormatSSE`）
- [x] 2.3 `serveStreaming` Phase 1 使用消息自身 AgentType/AgentName 构造 SSE 事件（替换 `FormatSSE`）

## 3. Backend ListMessages session_id 过滤

- [x] 3.1 `ListMessages` 新增可选 `session_id` query param，传入时添加 `WHERE session_id = ?` 条件

## 4. Frontend 传 session_id

- [x] 4.1 `getTaskMessages` 调用时传入当前 `sessionId` 参数
- [x] 4.2 `use-chat-stream.ts` 的 `useEffect` 历史加载中传入 sessionId

## 5. 验证

- [ ] 5.1 启动三端服务，创建 orchestrator + claude-code 的 task，发送消息触发多 Agent 流
- [ ] 5.2 实时流验证：确认前端正确显示 orchestrator → claude-code → orchestrator 三条独立消息
- [ ] 5.3 历史加载验证：刷新页面，确认消息按 agent_type 正确拆分显示
- [ ] 5.4 单 Agent 验证：纯 claude-code 会话（无 orchestrator）不受影响
