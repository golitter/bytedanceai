## Context

当前架构中，Agent 流式输出的持久化依赖三层临时存储：

1. **Frontend**: `useReducer` 组件内存（切 Tab 即丢）
2. **Backend**: `contentBuilder` 函数局部变量 + 流结束后写入 MySQL
3. **Agentend**: `Session.history` 内存（重启即丢）

Backend 的 `RunTask` handler 是同步 SSE 代理——在一个 HTTP 请求内同时读 agentend 流并写给前端。前端断连后 goroutine 虽然继续运行（`fmt.Fprintf` 错误未检查），但内容只在流结束后才 `db.Create`，中途无法查询。

现有 Message 表只有自增 `ID`，无对外标识、无状态字段、无流式追踪。

## Goals / Non-Goals

**Goals:**
- Agent 输出从第一条 event 起就持久化，前端随时可恢复
- 前端切 Tab 不丢失流式状态
- 前端刷新/断网后重连，能看到已输出的全部内容 + 续接实时流
- Backend 与前端连接解耦——前端断开不影响 Agent 执行和持久化

**Non-Goals:**
- Agent 端改动（agentend 不变）
- Agent 执行队列化/调度（仍是即时启动）
- 历史消息的全文搜索
- WebSocket 替代 SSE

## Decisions

### D1: Redis Stream 作为实时事件总线

**选择**: Redis Stream (`XADD` / `XREAD BLOCK`)
**替代方案**:
- A) goroutine + Go channel：进程内通信，backend 重启丢失，无法跨实例
- B) Redis Pub/Sub：无持久化，晚连接的订阅者收不到历史消息
- C) Kafka：过重，单实例部署不需要

**理由**: Redis Stream 兼具持久化和实时性。Stream 内消息在 TTL 内可回放（`XREAD` 指定起始 ID），支持前端断线重连后无缝续接。部署成本低（单 Redis 实例），且为未来多 backend 实例水平扩展留空间。

### D2: Redis Stream key 格式

**选择**: `agent:{session_id}:{message_id}`（每条回复独立 Stream）
**替代方案**:
- `agent:{session_id}`（Session 级复用）：需要 message_index 区分轮次，清理复杂

**理由**: 每条回复独立生命周期——创建、TTL、清理互不干扰。Stream 内无需额外字段标识轮次。前端重连时只需一个 message_id 即可定位。

### D3: 增量持久化策略——混合批处理

**选择**: 每 500 字符 **或** 每 2 秒 flush 一次 MySQL，先到先触发。流结束时强制 final flush。
**替代方案**:
- A) 每个 Text event 都 UPDATE：DB 压力大（Agent 可达 100+ event/s）
- B) 仅流结束写入：前端断连后无法查询部分内容（当前行为）

**理由**: 混合策略在持久化完整性和 DB 压力间取得平衡。极端情况（backend 崩溃）最多丢失 500 字符或 2 秒内容，Redis Stream 中有全量事件可回补。

### D4: POST /run 异步化

**选择**: `POST /run` 立即返回 `202 Accepted` + `message_id`，后台 goroutine 独立消费 agentend 流
**替代方案**:
- A) 保持同步 SSE 代理 + 另加 reconnect 端点：双路径逻辑复杂
- B) WebSocket：双向通信，但需要改动 agentend

**理由**: 统一所有流式输出通过 `GET /stream`，无论首次连接还是重连走同一路径。前端逻辑简化——一个 `submitMessage()` + 一个 `connectSSE()`。

### D5: 前端 Zustand 全局 store

**选择**: Zustand store 管理所有 session 的 chat 状态，`useChatStream` 从 store 读写
**替代方案**:
- A) 保持 `useReducer` + Context：需要 Provider 包裹，性能差（全树 re-render）
- B) localStorage 持久化：序列化延迟，不适合高频 streaming 更新

**理由**: Zustand 是项目已有的状态管理方案。全局 store 天然跨组件共享，切 Tab 时组件卸载但 store 保持。刷新时从 MySQL 重新加载 + 自动重连 SSE。

### D6: 重连时历史内容伪装为 text events

**选择**: GET /stream 先将 MySQL 已存内容拆成多个 text event 发给前端，格式与实时事件完全一致
**替代方案**:
- A) 单独 `history` event type：前端需要区分两种模式，逻辑复杂
- B) 一个大 text event 包含全部历史：JSON 载荷可能很大

**理由**: 前端无需区分"历史回放"和"实时流"，复用同一套 event handler。拆成合理大小的 chunk（按 SSE event 边界）避免单个过大的 JSON 载荷。

## Risks / Trade-offs

**[Risk] Redis 成为单点故障** → Backend 启动时检查 Redis 连接，Redis 不可用时降级为同步模式（与当前行为一致）。后续可加 Redis Sentinel/Cluster。

**[Risk] goroutine 泄漏** → 使用 `context.WithTimeout` 设置上限（如 30 分钟）。Backend 启动时扫描 `status=streaming` 的 Message，标记为 `failed`。使用 `sync.WaitGroup` 或 registry 追踪活跃 goroutine。

**[Risk] 并发 UPDATE MySQL 冲突** → BatchWriter 是单 goroutine 串行写入，无并发冲突。最终 `status` 更新在 goroutine 退出时执行。

**[Risk] 前端 POST /run 和 GET /stream 之间有时间窗口** → goroutine 可能在 GET /stream 连接前已产出部分事件。但 Redis Stream 会保留，GET /stream 从 MySQL + Redis 衔接，无丢失。

**[Trade-off] Redis Stream MAXLEN 截断** → 设 MAXLEN ~10000 防止超长 Agent 回复撑爆内存。MySQL 作为 ground truth 保留完整内容，Redis 只是实时中转。

## Migration Plan

1. **部署 Redis** 实例，Backend 连接配置加入 `configs/config.yaml`
2. **DB Migration**: ALTER messages 表新增 `message_id`, `status`, `last_seq` 字段（无破坏性，默认值兼容旧行为）
3. **Backend**: 新增 `stream/` 包和 `GET /stream` 端点，`POST /run` 改为异步。两阶段发布：先上线 GET /stream（向后兼容），再切换 POST /run
4. **Frontend**: 新增 Zustand store，`useChatStream` 切换到新 API。`connectSSE` 改用 GET 请求
5. **回滚**: 如果新流程出问题，Frontend 可回退到直接连 POST /run（同步模式），Backend 保留旧逻辑作为 fallback

## Open Questions

- Redis Stream 的 TTL 策略：流结束后保留多久？当前建议 10 分钟，是否需要更长？
- BatchWriter 的 flush 间隔是否需要可配置（写入 config.yaml）？
