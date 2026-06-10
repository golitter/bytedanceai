## MODIFIED Requirements

### Requirement: GET /stream SSE 端点
系统 SHALL 提供 `GET /api/tasks/:taskId/stream?session_id=&message_id=` 端点，返回 SSE 流。该端点 SHALL 根据 Message 的 `status` 字段决定行为：

- `streaming`：先发 MySQL 已积累内容，replay Redis gap，再从 RuntimeHub 订阅实时事件
- `completed`：发全部 MySQL 内容 + Done event，结束
- `failed`：发已有内容 + Error event，结束

#### Scenario: 连接正在 streaming 的消息（三阶段握手）
- **WHEN** 前端请求 `GET /stream?message_id=xxx`，该 Message `status=streaming`
- **THEN** 后端执行三阶段：(1) 从 MySQL 读取已存内容，拆成 text events 发送；(2) 调用 `hub.Subscribe()` 获取 currentSeq，replay Redis Stream 中 `(lastSeq, currentSeq]` 范围的 gap 事件；(3) 从 hub channel 消费实时事件直到收到 Done 或连接断开

#### Scenario: 连接已完成的消息
- **WHEN** 前端请求 `GET /stream?message_id=xxx`，该 Message `status=completed`
- **THEN** 后端发送全部内容作为 text events，然后发送 Done event，关闭 SSE 连接

#### Scenario: 连接已失败的消息
- **WHEN** 前端请求 `GET /stream?message_id=xxx`，该 Message `status=failed`
- **THEN** 后端发送已有内容作为 text events，然后发送 Error event，关闭 SSE 连接

#### Scenario: message_id 不存在
- **WHEN** 前端请求不存在的 `message_id`
- **THEN** 返回 HTTP 404

### Requirement: Redis Stream gap replay
系统 SHALL 在 subscribe RuntimeHub 后，使用 hub 返回的 `currentSeq` 作为切割点，从 Redis Stream replay `(lastSeq, currentSeq]` 范围内的事件，确保 MySQL 历史和实时流之间无 gap。

#### Scenario: 有 gap 需要 replay
- **WHEN** MySQL 历史发送完毕（lastSeq=5），hub.Subscribe() 返回 currentSeq=8
- **THEN** 后端从 Redis XREAD 读取 seq 5 到 seq 8 之间的事件，逐个发送给前端

#### Scenario: 无 gap 直接进入实时消费
- **WHEN** MySQL 历史发送完毕（lastSeq=10），hub.Subscribe() 返回 currentSeq=10
- **THEN** 后端跳过 Redis replay，直接从 hub channel 消费实时事件

#### Scenario: hub 尚未创建（StreamWriter 还在连接 agentend）
- **WHEN** SSE handler 调用 hub.Subscribe()，但该 stream 尚未创建
- **THEN** hub.Subscribe() 返回空 channel 和 currentSeq=0，SSE handler 等待 channel 事件或定期检查 MySQL status

### Requirement: 前端断线自动重连
前端 SHALL 在 SSE 连接断开时自动重连同一 `message_id`。重连后从 MySQL + Redis replay + hub 恢复完整内容，用户无感知。

#### Scenario: 重连时 gap replay 补齐
- **WHEN** 前端断线后重连，MySQL 中已有到 lastSeq=20 的内容，hub 当前在 seq=25
- **THEN** 后端发送 MySQL 历史，replay Redis gap (20,25]，然后从 hub 消费 seq>25 的实时事件
