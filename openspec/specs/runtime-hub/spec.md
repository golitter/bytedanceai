## ADDED Requirements

### Requirement: RuntimeHub manages in-memory pub/sub streams
系统 SHALL 提供 `RuntimeHub` 单例，管理多个 `RuntimeStream` 实例。每个 RuntimeStream 由 stream key（格式 `sessionID:messageID`）标识，拥有独立的单调递增序列号和订阅者集合。

#### Scenario: 创建新 stream 并发布事件
- **WHEN** StreamWriter 对不存在的 key 调用 `hub.Publish(key, event)`
- **THEN** RuntimeHub 自动创建 RuntimeStream，分配 seq=1，将事件发送给所有订阅者

#### Scenario: 发布事件到已有 stream
- **WHEN** StreamWriter 对已有 key 调用 `hub.Publish(key, event)`，当前 seq=5
- **THEN** seq 递增为 6，事件携带 seq=6 发送给所有订阅者

#### Scenario: 无订阅者时发布
- **WHEN** 调用 `hub.Publish(key, event)` 但该 stream 无订阅者
- **THEN** Publish 为 no-op，不报错，seq 仍递增

### Requirement: RuntimeStream 订阅返回 channel 和 currentSeq
调用 `hub.Subscribe(key)` SHALL 返回一个 `<-chan HubEvent` 和当前 `currentSeq`（uint64）。订阅者通过 channel 接收实时事件。

#### Scenario: 订阅活跃的 stream
- **WHEN** SSE handler 调用 `hub.Subscribe("sess1:msg1")`，该 stream 已发布到 seq=10
- **THEN** 返回 channel 和 currentSeq=10，后续 Publish 的事件从 seq=11 开始通过 channel 接收

#### Scenario: 订阅不存在的 stream
- **WHEN** SSE handler 调用 `hub.Subscribe("sess1:msg999")` 但该 stream 尚未创建
- **THEN** 返回 channel 和 currentSeq=0，当 stream 被创建后事件开始流入 channel

### Requirement: HubEvent 携带序列号和事件数据
`HubEvent` 结构体 SHALL 包含 `Seq uint64`（单调递增序列号）和 `Data string`（SSE data line 内容，如 `data: {"type":"text","content":{"text":"..."}}`）。

#### Scenario: 事件序列号单调递增
- **WHEN** hub 连续发布 3 个事件
- **THEN** 订阅者收到 HubEvent{Seq:1}, HubEvent{Seq:2}, HubEvent{Seq:3}

#### Scenario: 事件数据与 SSE 格式一致
- **WHEN** StreamWriter 发布一个 TEXT 事件
- **THEN** HubEvent.Data 为 `data: {"type":"text","content":{"text":"token text","agent_type":"claude"}}`

### Requirement: RuntimeStream 关闭通知订阅者
调用 `hub.Close(key)` SHALL 向所有订阅者的 channel 发送一个 `HubEvent{Done: true}` 事件，然后关闭 channel 并清理订阅者。

#### Scenario: StreamWriter 完成时关闭 stream
- **WHEN** StreamWriter 的 `finish()` 调用 `hub.Close("sess1:msg1")`
- **THEN** 所有订阅者收到 HubEvent{Done:true}，channel 被关闭，订阅者从 map 中移除

#### Scenario: 无订阅者时关闭
- **WHEN** 调用 `hub.Close(key)` 但无订阅者
- **THEN** stream 从 RuntimeHub 中移除，不报错

### Requirement: Subscriber channel 有界缓冲
每个订阅者的 channel SHALL 有界缓冲区（256 个事件）。当缓冲区满时，旧事件 SHALL 被丢弃（非阻塞发送）。

#### Scenario: 缓冲区满时丢弃旧事件
- **WHEN** hub 发布第 257 个事件，某订阅者尚未消费前 256 个
- **THEN** 最旧的事件被丢弃，订阅者收到最新事件（无阻塞）

#### Scenario: 慢订阅者不断线
- **WHEN** 订阅者消费速度慢于生产速度
- **THEN** 订阅者不断线，但可能丢失中间事件（可从 Redis replay 恢复）

### Requirement: RuntimeHub 线程安全
RuntimeHub 的所有方法（Publish、Subscribe、Close）SHALL 线程安全，支持多 goroutine 并发调用。

#### Scenario: StreamWriter 和 SSE handler 并发访问
- **WHEN** StreamWriter goroutine 调用 Publish，同时 SSE handler goroutine 调用 Subscribe
- **THEN** 无 race condition，无 panic，操作正确完成
