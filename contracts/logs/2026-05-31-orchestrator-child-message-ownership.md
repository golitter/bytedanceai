# Orchestrator 子 Agent 消息所有权与群聊展示顺序修复

## 变更原因

Orchestrator 群聊中，子 Agent 回复会先在子 Agent session 持久化，随后又被 Orchestrator SSE 透传并写入 Orchestrator session，导致刷新后出现重复消息、ask-agent 卡片挂到错误 Agent 气泡、以及群聊窗口上下文重复注入。

本次修复明确消息所有权：

- 子 Agent 正文由子 Agent session 持久化。
- Orchestrator session 只持久化 Orchestrator 自己的规划、总结与 ask-agent 卡片。
- Orchestrator 透传子 Agent 正文时复用子 Agent 原始 `message_id`，Backend 识别为外部 session 已持久化消息后只做 SSE/Redis 转发。

## 变更文件

本次变更 **未修改** `contracts/schemas/*.yaml`。

相关契约仍使用现有开放字段：

- `contracts/schemas/event-types.yaml` 中 `StreamEvent.content.additionalProperties: true`
- `text.content.message_id` 已在既有文本事件元数据约定中使用
- `ask_card_start` / `ask_card_done` 已在既有事件枚举中定义

## 对比结果

### 变更前

- Orchestrator 透传子 Agent `text` 时，Backend `StreamWriter` 可能在 Orchestrator session 中创建重复 Message。
- 前端 live store 在 subagent streaming 气泡期间收到新的 `ask_card_start` 时，会把 ask 卡片追加到当前 subagent 气泡。
- 刷新历史时，旧 Orchestrator session 内的重复子 Agent 转发行会被当成真实 aa/god 气泡展示。
- `fetchGroupChatWindow` 会把 user / Orchestrator / 旧重复转发行都注入给下一个 Agent。

### 变更后

- `TaskResult` 记录子 Agent 原始 `message_id`，Orchestrator 转发正文时携带该 ID。
- Backend 对属于其他 session 的 `text.content.message_id` 只转发，不重复写 MySQL。
- 子 Agent 回复后，下一段 Orchestrator 正文或新的 ask-card 会强制拆成新 Message，支持多轮 `Orchestrator → SubAgent → Orchestrator → SubAgent` 顺序。
- ask-card 按 `question_id` 回写完成态，避免 `ask_card_done` 写错消息。
- 群聊历史过滤 Orchestrator session 中旧的非 Orchestrator agent 转发行。
- 群聊窗口保留 `role` / `agent_type` 全量上下文，只按 task、session、状态和时间窗裁剪，并去重旧内容。

## 跨端影响

- **AgentEnd**: `ExecutionEngine` 将 Backend 返回的子 Agent `message_id` 写入 `TaskResult`；Orchestrator 执行阶段输出 ask-card 并携带原始 `message_id` 转发子 Agent 正文；ask_agent 不再伪造子 Agent session 的 user 消息。
- **Backend**: `StreamWriter` 根据 `text.content.message_id` 判定外部已持久化消息，只转发不落库；ask-card 按 `question_id` 维护消息归属；群聊窗口保留 role / agent_type 全量上下文并去重旧重复内容。
- **Frontend**: 群聊历史读取 task 可见消息并过滤旧重复转发行；live store 在新 ask-card 到达时按 speaker 切换结算上一条 subagent 消息。
- **Contracts**: 无 schema 变更；本次仅记录既有开放 `StreamEvent.content` 字段的使用语义和跨端消息所有权约定。

## 契约变更

无 `contracts/schemas/*.yaml` 变更。

现有约定说明：

| 事件/字段 | 类型 | 说明 |
|-----------|------|------|
| `text.content.message_id` | `string` | 可指向原始子 Agent Message；当该 Message 属于其他 session 时，Backend 只转发不重复落库 |
| `ask_card_start.content.question_id` | `string` | ask-card 的稳定归属 ID |
| `ask_card_done.content.question_id` | `string` | 用于回写同一张 ask-card 的完成态 |

兼容性：

- 旧客户端忽略未知 `content` 字段时仍可显示文本。
- 新 Backend 遇到旧历史重复行时通过群聊历史/窗口过滤降低影响。
