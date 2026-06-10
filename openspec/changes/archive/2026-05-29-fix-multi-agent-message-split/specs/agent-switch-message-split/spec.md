## ADDED Requirements

### Requirement: StreamWriter 检测 Agent 切换并创建新 Message
StreamWriter SHALL 解析 SSE TEXT 事件中的 `agent_type` 和 `agent` 字段。当 `agent_type` 与当前值不同时，SHALL 执行 Agent 切换流程：flush 当前 buffer 到当前 Message，将当前 Message status 设为 `completed`（除非是原始 Message），在同 session 下创建新 Message（status=`streaming`，agent_type/agent_name 为新值），更新内部 messageID。

#### Scenario: Agent 从 orchestrator 切换到 claude-code
- **WHEN** StreamWriter 收到 TEXT 事件，content 中 agent_type 从 `orchestrator` 变为 `claude-code`
- **THEN** flush 当前内容到当前 Message，创建新 Message（agent_type=`claude-code`），后续文本写入新 Message

#### Scenario: Agent 从 claude-code 切换回 orchestrator
- **WHEN** StreamWriter 收到 TEXT 事件，content 中 agent_type 从 `claude-code` 变为 `orchestrator`
- **THEN** finalize claude-code 的 Message（status=`completed`），创建新 Message（agent_type=`orchestrator`）

#### Scenario: 连续相同 agent_type 不触发切换
- **WHEN** 连续收到多个 TEXT 事件，agent_type 均为 `orchestrator`
- **THEN** 不触发切换，内容持续追加到当前 Message

#### Scenario: TEXT 事件无 agent_type 字段
- **WHEN** 收到 TEXT 事件但 content 中无 `agent_type` 字段
- **THEN** 不触发切换，按当前 agent_type 处理

#### Scenario: 切换时 buffer 为空
- **WHEN** Agent 切换触发时当前 buffer 无内容（streamingContent 为空）
- **THEN** 跳过 finalize 和创建新 Message，仅更新 agent_type/agent_name

### Requirement: 原始 Message 保持 streaming 直到整轮结束
StreamWriter SHALL 将 RunTask 创建的第一条 Message 标记为原始 Message。原始 Message 在 Agent 切换时 SHALL NOT 被 finalize。仅在 `finish()` 时（整轮 SSE 流结束），SHALL 将原始 Message 和最后一条子 Message 的 status 设为 `completed` 或 `failed`。

#### Scenario: 整轮正常结束
- **WHEN** SSE 流正常消费完毕（收到 done 事件）
- **THEN** 原始 Message status 设为 `completed`，最后一条子 Message status 设为 `completed`

#### Scenario: 整轮异常结束
- **WHEN** SSE 流消费过程中遇到错误
- **THEN** 原始 Message status 设为 `failed`，当前子 Message status 设为 `failed`

#### Scenario: ServeStream 通过原始 Message 判断活跃状态
- **WHEN** 前端连接到原始 messageID 的 SSE 端点
- **THEN** `IsActive(originalMessageID)` 在整轮运行期间返回 true，`ServeStream` 持续推送 Redis 事件

### Requirement: Redis stream key 在 Agent 切换期间保持不变
StreamWriter SHALL 在 Agent 切换创建新 Message 时，继续使用原始 Message 的 Redis stream key（`agent:{sessionID}:{originalMessageID}`）。新创建的子 Message SHALL NOT 创建独立的 Redis stream。

#### Scenario: 事件持续写入原始 stream
- **WHEN** Agent 切换后，新 Agent 的 TEXT 事件到来
- **THEN** 事件仍写入原始 messageID 的 Redis stream，前端从同一个 SSE 连接接收

### Requirement: ListMessages 支持按 session_id 过滤
`GET /api/tasks/:taskId/messages` SHALL 支持可选 `session_id` query param。传入时 SHALL 只返回匹配 session_id 的消息。不传入时行为不变（返回 task 下所有消息）。

#### Scenario: 按 session_id 过滤
- **WHEN** 前端请求 `GET /api/tasks/:taskId/messages?session_id=xxx&limit=20`
- **THEN** 只返回 session_id=xxx 的消息，分页逻辑不变

#### Scenario: 不传 session_id 向后兼容
- **WHEN** 前端请求 `GET /api/tasks/:taskId/messages?limit=20`（无 session_id）
- **THEN** 返回 task 下所有 session 的消息（当前行为不变）
