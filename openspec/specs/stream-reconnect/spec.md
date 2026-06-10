## ADDED Requirements

### Requirement: GET /stream SSE 端点
系统 SHALL 提供 `GET /api/tasks/:taskId/stream?session_id=&message_id=` 端点，返回 SSE 流。该端点 SHALL 根据 Message 的 `status` 字段决定行为：

- `streaming`：先发 MySQL 已积累内容，再从 Redis Stream 续接实时事件
- `completed`：发全部 MySQL 内容 + Done event，结束
- `failed`：发已有内容 + Error event，结束

#### Scenario: 连接正在 streaming 的消息
- **WHEN** 前端请求 `GET /stream?message_id=xxx`，该 Message `status=streaming`
- **THEN** 后端从 MySQL 读取已存内容，拆成 text events 发送；然后 XREAD BLOCK 从 Redis Stream `last_seq` 之后读取实时事件并转发

#### Scenario: 连接已完成的消息
- **WHEN** 前端请求 `GET /stream?message_id=xxx`，该 Message `status=completed`
- **THEN** 后端发送全部内容作为 text events，然后发送 Done event，关闭 SSE 连接

#### Scenario: 连接已失败的消息
- **WHEN** 前端请求 `GET /stream?message_id=xxx`，该 Message `status=failed`
- **THEN** 后端发送已有内容作为 text events，然后发送 Error event，关闭 SSE 连接

#### Scenario: message_id 不存在
- **WHEN** 前端请求不存在的 `message_id`
- **THEN** 返回 HTTP 404

### Requirement: 历史内容伪装为 text events
系统 SHALL 将 MySQL 中已积累的 agent 输出内容拆分成多个 SSE text event 发送给前端，格式与实时 text event 完全一致。前端 SHALL 无需区分历史回放和实时流。

#### Scenario: 历史内容拆分发送
- **WHEN** GET /stream 连接到 `status=streaming` 的 Message，MySQL 中已有 2000 字符内容
- **THEN** 后端将 2000 字符拆成若干 text events（每个约 500 字符），逐个发送，格式为 `data: {"type":"text","content":{"text":"..."}}`

#### Scenario: 无历史内容直接续接
- **WHEN** GET /stream 连接到 `status=streaming` 的 Message，MySQL 中尚无内容
- **THEN** 后端跳过历史发送阶段，直接 XREAD BLOCK 等待 Redis Stream 实时事件

### Requirement: Redis Stream 续接
系统 SHALL 在发送完 MySQL 历史内容后，从 `last_seq` 之后的位置开始 `XREAD BLOCK` 读取 Redis Stream 事件并转发给前端，实现零间隙衔接。

#### Scenario: MySQL + Redis 无缝衔接
- **WHEN** MySQL 历史内容发送完毕，Redis Stream `last_seq` 之后还有未消费的事件
- **THEN** 后端先通过 `XREAD` 读取积压事件（非阻塞），逐个转发；积压消费完后切换到 `XREAD BLOCK` 等待实时事件

#### Scenario: Agent 流在重连时刚好结束
- **WHEN** 前端重连，MySQL 内容已全部发送，但 goroutine 刚好设置了 `status=completed`
- **THEN** 后端发送 Done event，关闭 SSE 连接

### Requirement: 前端断线自动重连
前端 SHALL 在 SSE 连接断开时自动重连同一 `message_id`。重连后从 MySQL + Redis 恢复完整内容，用户无感知。

#### Scenario: 网络断开后自动重连
- **WHEN** SSE 连接因网络问题断开
- **THEN** 前端在短暂延迟后（指数退避，初始 1s，最大 10s）使用同一 `message_id` 重新连接 GET /stream

#### Scenario: 页面刷新后自动重连
- **WHEN** 用户刷新页面，store 中有 `status=streaming` 的活跃消息
- **THEN** 页面加载后从 MySQL 读取历史消息，发现 streaming 消息，自动连接 GET /stream 恢复
