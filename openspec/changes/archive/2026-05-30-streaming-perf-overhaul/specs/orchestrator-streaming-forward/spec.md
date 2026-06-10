## ADDED Requirements

### Requirement: ExecutionEngine 流式透传子 Agent token
ExecutionEngine._execute_task SHALL 在收到子 Agent 的每个 TEXT 事件时，立即 yield 一个 `RUNTIME_TEXT` StreamEvent，不再缓冲全部文本后才发出 RUNTIME_COMPLETED。同时 SHALL 维护 collected 列表用于最终 TaskResult.content。

#### Scenario: 子 Agent 输出 token 实时透传
- **WHEN** 子 Agent adapter yield 一个 TEXT 事件，内容为 "你好"
- **THEN** ExecutionEngine SHALL 立即 yield `StreamEvent(type=RUNTIME_TEXT, content={"task_id": "task-001", "agent": "claude-code", "text": "你好"})`

#### Scenario: 多个 token 连续透传
- **WHEN** 子 Agent 连续输出三个 TEXT 事件："你"、"好"、"世界"
- **THEN** ExecutionEngine SHALL 依次 yield 三个独立的 RUNTIME_TEXT 事件

#### Scenario: 子 Agent 执行完成后仍发出 RUNTIME_COMPLETED
- **WHEN** 子 Agent 执行完成（done event 或 timeout）
- **THEN** ExecutionEngine SHALL yield `StreamEvent(type=RUNTIME_COMPLETED, ...)` 且 TaskResult.content 包含完整累积文本

#### Scenario: 子 Agent 出错时发出 ERROR 事件
- **WHEN** 子 Agent 执行遇到异常或超时
- **THEN** ExecutionEngine SHALL yield `StreamEvent(type=ERROR, ...)` 并跳过 RUNTIME_COMPLETED

### Requirement: Orchestrator 进程内短路调用 adapter
Orchestrator 调度子 Agent 时 SHALL 直接在 agentend 进程内调用对应 adapter 的 stream_chat() 方法，而非通过 BackendClient HTTP 回环。短路路径 SHALL 仍通过 BackendClient.run_task() 创建 agent message 记录。

#### Scenario: 短路调用 adapter 流式输出
- **WHEN** Orchestrator 调度一个 claude-code 类型的子 Agent
- **THEN** ExecutionEngine SHALL 直接调用 ClaudeAdapter.stream_chat() 获取 AsyncIterator[StreamEvent]，逐事件 yield RUNTIME_TEXT

#### Scenario: 短路调用仍创建 message 记录
- **WHEN** 短路路径执行子 Agent
- **THEN** SHALL 先调用 BackendClient.run_task() 创建 agent message（获取 message_id），然后直接调用 adapter.stream_chat()

#### Scenario: 不支持短路的 agent_type 回退到 HTTP 路径
- **WHEN** 子 Agent 的 agent_type 没有可用的 adapter 实例
- **THEN** SHALL 回退到 BackendClient.run_task() + stream_result() 的原有 HTTP 路径，并仍然流式透传 RUNTIME_TEXT

### Requirement: RUNTIME_TEXT 事件类型定义
系统 SHALL 支持 `runtime_text` 事件类型，content 格式为 `{"task_id": string, "agent": string, "text": string}`。该事件类型 SHALL 注册到 contracts/schemas/ 的 StreamEvent 契约中。

#### Scenario: RUNTIME_TEXT 事件序列化
- **WHEN** ExecutionEngine yield 一个 RUNTIME_TEXT 事件
- **THEN** 序列化为 `{"type": "runtime_text", "content": {"task_id": "task-001", "agent": "claude-code", "text": "增量文本"}, "timestamp": "..."}`

#### Scenario: RUNTIME_TEXT 事件经过 Backend StreamWriter 正常持久化
- **WHEN** RUNTIME_TEXT 事件到达 Backend StreamWriter
- **THEN** SHALL 被正常 XADD 到 Redis Stream，前端可通过 SSE 连接消费
