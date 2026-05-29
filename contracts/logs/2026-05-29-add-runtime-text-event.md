# 2026-05-29-add-runtime-text-event

## 变更原因
Orchestrator 多 Agent 模式下，子 Agent 的 token 经过四层网络回环后被全量缓冲，用户在子 Agent 执行期间看不到任何输出。新增 `runtime_text` 事件类型用于 Orchestrator 实时透传子 Agent 的增量文本。

## 变更文件
- `contracts/schemas/event-types.yaml`: EventType 枚举新增 `runtime_text`

## 对比结果
```diff
 enum:
   - init
   - text
   - tool_call
   - tool_result
   - artifact
   - planning
   - done
   - error
   - runtime_executing
+  - runtime_text          # Orchestrator 透传子 Agent 增量文本
   - runtime_completed
   - coordination_start
   - coordination_message
   - coordination_done
```

## 跨端影响
- **AgentEnd**: ExecutionEngine 需 yield `RUNTIME_TEXT` 事件，content 包含 task_id/agent/text
- **Backend**: StreamWriter 需透传 `runtime_text` 事件到 Redis Stream
- **Frontend**: Chat store 需识别 `runtime_text` 事件并归属到子 Agent 消息块

## 契约变更
- EventType 枚举新增值: `runtime_text`
- StreamEvent.content 在 `runtime_text` 类型下结构为 `{ task_id: string, agent: string, text: string }`
