# 契约变更：新增 runtime + coordination 事件类型

## 变更原因

Phase 5 多 Agent 协作运行时需要新事件类型来支持：Agent 执行状态追踪（runtime_executing/runtime_completed）和协调通道消息（coordination_start/message/done）。

## 变更文件

- `contracts/schemas/event-types.yaml`

## 对比结果

EventType 枚举新增 5 个值：

| 新增值 | 说明 |
|--------|------|
| `runtime_executing` | Agent 开始执行任务，content 含 task_id、agent、title、status |
| `runtime_completed` | Agent 完成任务，content 含 task_id、agent、success、duration、status |
| `coordination_start` | 协调通道开启，content 含 round、agents |
| `coordination_message` | 协调消息，content 含 from、to、text、round |
| `coordination_done` | 协调通道关闭，content 含 rounds、decisions |

## 跨端影响

- **AgentEnd**: EventType 枚举扩展，OrchestratorAdapter 的 ExecutionEngine 和 CoordinationChannel 使用新类型
- **Backend**: EventType 常量扩展，SSE 透传新事件类型（无需特殊处理）
- **Frontend**: EventType 枚举扩展，use-chat-stream 和 block-reducer 新增对应处理逻辑
