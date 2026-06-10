## 1. Backend — 基础设施 & 数据模型

- [x] 1.1 在 `configs/config.yaml` 添加 Redis 连接配置（host、port、password、db）
- [x] 1.2 添加 `go-redis` 依赖到 `go.mod`，创建 Redis 客户端单例 `pkg/redis/redis.go`
- [x] 1.3 `model/message.go` 新增字段：`message_id`（VARCHAR(36) UNIQUE）、`status`（VARCHAR(16) DEFAULT "completed"）、`last_seq`（VARCHAR(64) DEFAULT ""）
- [x] 1.4 验证 DB migration 正确执行，旧数据兼容（status 默认 completed，message_id 默认空）

## 2. Backend — Stream 包（核心持久化逻辑）

- [x] 2.1 创建 `stream/writer.go`：StreamWriter 结构体，封装 goroutine 生命周期管理（context、WaitGroup、registry）
- [x] 2.2 实现 BatchWriter：buffer 累积 + 500 字符/2 秒混合 flush 策略 + final flush
- [x] 2.3 实现从 agentend 流读取 → `XADD` 写入 Redis Stream（`agent:{session_id}:{message_id}`，MAXLEN ~10000）
- [x] 2.4 BatchWriter flush 时同步 UPDATE MySQL（content + last_seq）
- [x] 2.5 流结束时：final flush + UPDATE status 为 completed/failed + Redis EXPIRE 600s
- [x] 2.6 实现 goroutine registry（`sync.Map` 按 message_id 追踪活跃 goroutine），支持超时取消（30 分钟）

## 3. Backend — API 端点改造

- [x] 3.1 改造 `POST /run`：创建 user Message → 创建 agent Message（status=streaming, message_id=UUID）→ 启动 StreamWriter goroutine → 返回 202 `{ message_id, status }`
- [x] 3.2 新建 `handler/stream.go`：`GET /stream?session_id=&message_id=` SSE 端点
- [x] 3.3 GET /stream 实现：根据 Message.status 分支处理（streaming → MySQL 历史 + Redis 续接，completed → 全量 + Done，failed → 已有 + Error）
- [x] 3.4 实现 MySQL 历史内容拆分为 text events（~500 字符/chunk），格式与实时事件一致
- [x] 3.5 实现 Redis Stream 续接：先 XREAD 消费积压 → 再 XREAD BLOCK 等待实时事件
- [x] 3.6 Backend 启动时扫描 `status=streaming` 的 Message，标记为 `failed`

## 4. Frontend — Zustand Store & API 层

- [x] 4.1 新建 `stores/chat.ts`：Zustand store，按 sessionId 组织状态（messages、streamingContent、status、activeStream）
- [x] 4.2 `lib/api.ts` 新增 `submitMessage(taskId, body)` 函数，POST /run 返回 `{ message_id }`
- [x] 4.3 `lib/sse.ts` 改造：支持 GET 请求方式连接 SSE（`connectSSE` 接受 URL + params）
- [x] 4.4 `lib/sse.ts` 新增断线重连逻辑：指数退避（1s → 2s → 4s → ... → 10s max），同一 message_id 重连 GET /stream

## 5. Frontend — Hook 改造 & 组件适配

- [x] 5.1 改造 `hooks/use-chat-stream.ts`：从 Zustand store 读写状态，替代 useReducer
- [x] 5.2 `sendMessage` 流程改为：submitMessage() → 获取 message_id → connectSSE() → 通过 store 更新状态
- [x] 5.3 mount 时加载历史消息，检测 streaming 状态的 agent message，自动重连 GET /stream
- [x] 5.4 `ChatArea.tsx` 和 `MessageList.tsx` 适配新 store 接口（如果数据获取方式有变化）
- [x] 5.5 `MessageInput.tsx` 适配新的 `onSend` 签名（如 submitMessage 变为 async）

## 6. 集成验证

- [ ] 6.1 端到端测试：发送消息 → Agent 流式输出 → 前端实时显示 → 流结束 → Message 完整持久化
- [ ] 6.2 切 Tab 测试：发送消息 → 切到其他 session → 切回 → streaming 状态秒恢复
- [ ] 6.3 刷新测试：发送消息 → 刷新页面 → 历史消息恢复 → streaming 消息自动重连
- [ ] 6.4 Agent 长时间运行测试：Agent 跑 5+ 分钟 → 中途刷新 → 内容完整恢复
- [ ] 6.5 异常测试：Agent 失败 → 前端显示 Error；Backend 重启 → streaming 消息标记 failed
