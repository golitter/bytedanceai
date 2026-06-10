## ADDED Requirements

### Requirement: SSE Handler 批量读取 Redis Stream
Backend SSE Handler 的实时阶段 SHALL 使用 `XREAD Count:100 Block:200ms` 替代当前的 `Count:1 Block:5s`，一次读取多条消息后批量写入 HTTP 响应。

#### Scenario: Redis 中有多条待读消息
- **WHEN** Redis Stream 中有 50 条未读消息
- **THEN** XREAD 一次返回全部 50 条，SSE Handler 遍历写出 50 个 data 行，最后一次 Flush

#### Scenario: Redis 中无新消息
- **WHEN** XREAD 200ms 内没有新消息到达
- **THEN** 返回空结果，SSE Handler 继续下一次 XREAD 循环

#### Scenario: 流结束后仍能检测
- **WHEN** StreamWriter goroutine 结束（IsActive 返回 false）
- **THEN** SSE SHALL 按现有逻辑检查 message status 并发送 done/error 事件

### Requirement: StreamWriter 事件合并写入
StreamWriter.Run 中的 scanFunc 回调 SHALL 将高频小 TEXT 事件合并：累积不超过 500ms 或 2KB 文本后执行一次 XADD，而非逐条 XADD。

#### Scenario: 高频小 token 合并
- **WHEN** 100ms 内连续收到 10 个 TEXT 事件，每个 20 字节
- **THEN** 合并为一次 XADD，包含 200 字节文本

#### Scenario: 低频事件立即写入
- **WHEN** 收到一个 TEXT 事件后 500ms 内无新事件
- **THEN** 立即执行 XADD，不等累积

#### Scenario: 大文本立即写入
- **WHEN** 单个 TEXT 事件文本超过 2KB
- **THEN** 立即执行 XADD，不等累积

#### Scenario: 非文本事件不合并
- **WHEN** 收到 done、error、tool_call 等非 TEXT 事件
- **THEN** 先 flush 已累积的文本，再立即 XADD 该事件
