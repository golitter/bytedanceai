# 契约变更：TEXT 事件新增 agent 元数据字段

## 变更原因

多 Agent 编排场景下，前端需根据 TEXT 事件的来源 agent 进行消息分离显示（不同 agent 头像、名称、独立消息气泡）。原 TEXT 事件 content 仅含 `text` 字段，无法区分消息来源。

## 变更文件

本次变更 **未修改** `contracts/schemas/*.yaml`。`StreamEvent.content` 已为 `type: object` + `additionalProperties: true`，新增字段在现有 schema 下合法。

## 新增 TEXT 事件 content 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `agent` | `string` | 来源 agent 名称（如 `"Orchestrator"`、子 agent 名称） |
| `agent_type` | `string` | 来源 agent 类型（如 `"orchestrator"`、`"claude-code"`） |

前端 `streamAgentUpdate` 检测 `agent` + `agent_type` 变化时，自动将当前 streaming content 终结为独立消息，切换到新 agent 继续接收。

## Orchestrator 事件流时序

```
TEXT(text=overview,          agent="Orchestrator",  agent_type="orchestrator")  → 消息 1
TEXT(text=sub-agent-1-结果,   agent="Alice",         agent_type="claude-code")   → 消息 2
TEXT(text=sub-agent-2-结果,   agent="Bob",           agent_type="claude-code")   → 消息 3
TEXT(text=aggregated-总结,    agent="Orchestrator",  agent_type="orchestrator")  → 消息 4
```

## 跨端影响

- **AgentEnd**: `OrchestratorAdapter.stream_chat()` Phase 1/3/4 yield TEXT 事件时携带 `agent` + `agent_type`
- **Backend**: 透传 SSE 事件，无特殊处理
- **Frontend**: `use-chat-stream.ts` 的 TEXT 事件处理已有 `streamAgentUpdate` 逻辑，当 `event.content.agent` 和 `event.content.agent_type` 均存在时触发 agent 切换
- **Contracts**: 无 schema 变更
