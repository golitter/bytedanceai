# 契约变更记录 — RuntimeHub 低延迟实时流

## 变更原因

SSE token 流式传输延迟 ~250-700ms，根因是 token 事件经两层 batching（StreamWriter 500ms/2KB + StreamHandler Redis XRead 200ms block）才到达前端。引入内存 RuntimeHub 将实时传输与持久化解耦，延迟降至 ~10ms。

## 变更文件

无 schema 文件变更。

## 对比结果

不涉及 schema 对比。`event-types.yaml` 中 EventType 枚举和 StreamEvent 结构体均未修改，SSE 事件格式与前端完全兼容。

## 跨端影响

| 端 | 影响 |
|---|---|
| **Frontend** | 无影响。SSE 事件格式不变，前端代码零改动。 |
| **Backend** | 核心变更。新增 `internal/stream/hub.go`（RuntimeHub 内存 pub/sub），修改 `writer.go`（双写 hub + Redis），修改 `handler/stream.go`（三阶段握手替代 Redis 轮询）。 |
| **AgentEnd** | 无影响。Agent 产出的事件格式不变。 |

## 契约变更

无。本次变更为后端内部传输架构优化，不涉及跨端协议。

### 架构变更说明

- **HOT PATH**（新增）：StreamWriter → `RuntimeHub.Publish()` → Go channel → StreamHandler → SSE flush（~1-10ms）
- **COLD PATH**（不变）：StreamWriter → Redis XADD（异步）→ MySQL flush（500ms/2KB 批处理）
- **Replay**（不变）：断线重连仍从 MySQL 历史 + Redis Stream gap replay 恢复
- Redis Stream 保留用于 replay/reconnect/fanout，退出实时关键路径
