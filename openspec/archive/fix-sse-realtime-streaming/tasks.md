## 1. 前端 SSE 客户端重写

- [x] 1.1 重写 `frontend/src/lib/sse.ts`：用 `EventSource` 替换 `fetch` + `ReadableStream`，保持 `connectSSE` 函数签名不变（url, params, onEvent, onError, reconnect, maxBackoff），返回 `AbortController`
- [x] 1.2 EventSource 实例通过 `onmessage` 解析 JSON 并调用 `onEvent`；通过 `onerror` 处理连接错误（reconnect=true 时不关闭，由 EventSource 自动重连）
- [x] 1.3 AbortController 的 `abort` 事件监听调用 `es.close()` 关闭连接

## 2. 前端乐观更新

- [x] 2.1 修改 `frontend/src/hooks/use-chat-stream.ts` 的 `sendMessage`：将 `store.sendMessage()` 移到 `await submitMessage()` 之前
- [x] 2.2 用 try/catch 包裹 `submitMessage` 调用，失败时调用 `store.streamError()` 回滚状态

## 3. 后端 serveStreaming 竞态修复

- [x] 3.1 修改 `backend/internal/handler/stream.go` 的 `serveStreaming` 阻塞循环：将 `if !stream.IsActive()` 内的 `else` 分支改为 `switch fresh.Status`，当 status 为 "streaming" 时不 return，fall through 到 XRead 阻塞等待
- [x] 3.2 确认 "completed" 和 "failed" 状态的行为不变（分别发送 done 和 error 事件后 return）

## 4. 验证

- [ ] 4.1 启动三端服务（`make run-frontend` / `make run-backend` / `make run-agentend`），发送消息验证 Agent 回复实时流式显示
- [ ] 4.2 验证断线重连：Agent 流式中断开网络，恢复后 SSE 自动重连继续接收
- [ ] 4.3 验证用户消息即时显示：发送后消息立刻出现在聊天界面
- [ ] 4.4 验证错误处理：关闭 agentend 服务后发送消息，前端显示错误提示
