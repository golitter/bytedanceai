## MODIFIED Requirements

### Requirement: Chat 状态管理迁移到 Zustand 全局 store
系统 SHALL 使用 Zustand 全局 store 管理 chat 状态，替代 `useReducer`。Store SHALL 按 `sessionId` 组织状态，跨组件共享，切 Tab 不丢失。每个 session 的状态包含 messages、streamingContent、streamingAgentType、status、activeStream（messageId + 重连信息）。

#### Scenario: 发送消息更新 store
- **WHEN** 用户发送消息
- **THEN** 通过 store action 添加 user message，创建 activeStream，连接 GET /stream SSE

#### Scenario: SSE event 更新 store
- **WHEN** 收到 text event
- **THEN** 追加内容到对应 session 的 streamingContent

#### Scenario: 流结束更新 store
- **WHEN** 收到 done event
- **THEN** 将 streamingContent 转为完整 agent message 添加到 messages，清除 activeStream

#### Scenario: 切 Tab 保持状态
- **WHEN** 用户从 Session A 切换到 Session B，再切回 Session A
- **THEN** Session A 的 messages 和 streamingContent 在 store 中完整保留，UI 秒恢复

### Requirement: SSE 连接改为 GET 请求
系统 SHALL 将 SSE 连接从 `POST /api/tasks/:taskId/run` 改为 `GET /api/tasks/:taskId/stream?message_id=`。消息提交和流订阅分离为两步：先 POST 获取 message_id，再 GET 订阅 SSE。

#### Scenario: 提交消息并订阅流
- **WHEN** 用户发送消息
- **THEN** 先调用 `POST /run` 获取 `{ message_id }`，再连接 `GET /stream?message_id=xxx` 接收 SSE events

#### Scenario: POST /run 返回 202
- **WHEN** 前端发送 `POST /api/tasks/:taskId/run`
- **THEN** 后端返回 HTTP 202，body 为 `{ message_id: "uuid", status: "streaming" }`

### Requirement: 断线自动重连
系统 SHALL 在 SSE 连接断开时自动使用同一 `message_id` 重连 `GET /stream`。重连 SHALL 使用指数退避策略（初始 1 秒，最大 10 秒，抖动 ±0.5 秒）。

#### Scenario: 网络断开重连
- **WHEN** SSE 连接因网络错误断开
- **THEN** 前端等待退避时间后重连 GET /stream，恢复历史内容 + 实时流

#### Scenario: 页面刷新重连
- **WHEN** 用户刷新页面
- **THEN** store 从 MySQL 加载历史消息；发现 status=streaming 的 agent message，自动连接 GET /stream 恢复

### Requirement: StreamEvent parsing
系统 SHALL 解析每个 SSE `data:` 行为符合 StreamEvent 契约的 JSON 对象（`type`, `content`, `timestamp`）。支持的 event 类型：`init`, `text`, `tool_call`, `tool_result`, `artifact`, `planning`, `done`, `error`。

#### Scenario: Text event received
- **WHEN** 收到 `{ "type": "text", "content": { "text": "..." } }` event
- **THEN** 追加文本内容到当前 streaming agent message

#### Scenario: Done event received
- **WHEN** 收到 `{ "type": "done" }` event
- **THEN** 完成 streaming agent message，状态转为 `done`

#### Scenario: Error event received
- **WHEN** 收到 `{ "type": "error", "content": { "message": "..." } }` event
- **THEN** 显示错误信息，状态转为 `error`

#### Scenario: Malformed SSE data
- **WHEN** SSE `data:` 行包含无效 JSON
- **THEN** 记录解析错误并跳过该 event，不中断流

### Requirement: Chat state machine
系统 SHALL 使用 Zustand store 管理聊天生命周期，状态包括：`idle`, `loading`, `streaming`, `tool_running`, `done`, `error`。

#### Scenario: State transitions on happy path
- **WHEN** 用户发送消息
- **THEN** 状态转换: `idle` → `loading` → `streaming`（收到第一个 text event）→ `done`（收到 done event）

#### Scenario: State transition on error
- **WHEN** 收到 error event 或网络失败
- **THEN** 状态转为 `{ status: 'error', error: Error }`，显示错误 UI

#### Scenario: State transition on tool_call
- **WHEN** streaming 过程中收到 `tool_call` event
- **THEN** 状态转为 `{ status: 'tool_running', toolName: string }`，保持 abort controller

### Requirement: Message accumulation during streaming
系统 SHALL 将流式文本累积到单个 agent message 对象中，不按 chunk 创建新消息。最终消息只在 streaming 完成时添加到消息列表。每次 chunk 只更新最后一条消息的 DOM，不刷新整个列表。

#### Scenario: Progressive text accumulation
- **WHEN** 连续收到多个 `text` event
- **THEN** 内容拼接为单个 streaming message

#### Scenario: Streaming message finalization
- **WHEN** 收到 `done` event
- **THEN** 累积的 streaming message 被终结并追加到消息历史作为完整的 agent message
