## ADDED Requirements

### Requirement: 后台 goroutine 消费 agentend 流并持久化
系统 SHALL 在收到 `POST /run` 请求后启动后台 goroutine，独立于前端 HTTP 连接消费 agentend SSE 流。goroutine SHALL 将每个 SSE event 同时：(1) 通过 `hub.Publish()` 立即推送到 RuntimeHub（实时通道），(2) 发布到 Redis Stream `agent:{session_id}:{message_id}`（持久通道）。Redis 发布 SHALL 在 hub 发布之后执行。

#### Scenario: goroutine 独立于前端连接运行
- **WHEN** 前端发送 `POST /run`，后端返回 `202 Accepted`
- **THEN** 后台 goroutine 开始消费 agentend 流，不受前端连接状态影响

#### Scenario: 前端断连后 goroutine 继续
- **WHEN** 前端在 Agent 流式输出中断开连接
- **THEN** goroutine 继续读取 agentend 流，继续 hub.Publish + Redis XADD，继续 flush 到 MySQL

#### Scenario: 每个 SSE event 先推 hub 再写 Redis
- **WHEN** goroutine 从 agentend 收到一个 SSE event
- **THEN** 先调用 `hub.Publish(key, event)` 立即推送到所有订阅者，再通过 `XADD` 发布到 Redis Stream

#### Scenario: goroutine 超时保护
- **WHEN** goroutine 运行超过 30 分钟
- **THEN** context 取消，goroutine 退出，调用 `hub.Close(key)`，Message status 设为 `failed`，Redis Stream 设 EXPIRE

### Requirement: BatchWriter 混合批处理策略
系统 SHALL 使用 BatchWriter 在 goroutine 内批量将文本内容 flush 到 MySQL Message 表。触发条件为以下三者先到先触发：①累积文本 >= 2048 字符 ②距上次 flush >= 500ms ③流结束（Done/Error event）。此批处理策略仅影响 MySQL 持久化路径，不影响 hub 实时推送路径。

#### Scenario: 文本累积触发 MySQL flush
- **WHEN** goroutine 收到 Text event，累积内容达到 2048 字符
- **THEN** 执行 `UPDATE messages SET content=?, last_seq=? WHERE message_id=?`，清空 buffer

#### Scenario: 定时触发 MySQL flush
- **WHEN** 距上次 flush 已过 500ms，且有新内容
- **THEN** 执行 UPDATE，清空 buffer，重置 timer

#### Scenario: 流结束强制 flush
- **WHEN** goroutine 收到 Done 或 Error event
- **THEN** 强制 flush 全部剩余内容到 MySQL，调用 `hub.Close(key)`，设置 `status` 为 `completed` 或 `failed`

#### Scenario: TEXT 事件即时推送到 hub 不受批处理影响
- **WHEN** StreamWriter 收到一个 TEXT 事件，内容为 "你好"
- **THEN** 立即调用 `hub.Publish()` 推送（不等待批处理），同时异步进入 Redis/MySQL 批处理路径

### Requirement: Redis Stream 生命周期管理
系统 SHALL 在流开始时创建 Redis Stream，流结束后设置 EXPIRE。Stream key 格式为 `agent:{session_id}:{message_id}`。SHALL 设置 MAXLEN ~10000 防止内存溢出。

#### Scenario: 流结束后设置 TTL
- **WHEN** goroutine 完成流消费（Done 或 Error）
- **THEN** 对 Redis Stream key 执行 `EXPIRE` 设置 600 秒 TTL

#### Scenario: Stream 长度上限
- **WHEN** Stream 中的 event 数量超过 10000
- **THEN** 通过 `XADD MAXLEN ~10000` 自动截断最早的 event

### Requirement: Backend 启动时清理残留 streaming 状态
系统 SHALL 在 Backend 启动时扫描 `status=streaming` 的 Message 记录，将其标记为 `failed`。对应的 goroutine 已丢失，无法恢复。

#### Scenario: 启动时发现残留 streaming 消息
- **WHEN** Backend 启动，MySQL 中存在 `status=streaming` 的 Message
- **THEN** 将这些 Message 的 `status` 更新为 `failed`
