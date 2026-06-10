## Why

当前创建对话时后端不会自动创建 session，前端只能用乐观更新 workaround，导致刷新后对话消失。同时 session 只有 `agent_type` 字段（固定枚举），没有 `agent_name`（用户自定义名称），无法支持未来同类型多 agent 的群聊场景。

## What Changes

- 后端 `CreateTask` 接口新增 `agents` 参数（数组，每项包含 `type` 和 `name`），创建 task 时自动创建对应 sessions
- 后端 Session 模型新增 `agent_name` 字段
- 前端移除乐观更新 workaround，恢复正常 `invalidateQueries` 流程
- 前端 `NewChatDialog` 支持为 agent 输入自定义名称
- 前端对话列表和聊天区显示 `agent_name` 而非 `agent_type`
- `agent_name` 不穿透到 agentend 端，仅影响前端展示和后端存储

## Capabilities

### New Capabilities

- `auto-session`: 创建 task 时自动创建 session，前端无需 workaround
- `agent-naming`: Session 支持 agent_name 字段，前端展示自定义名称

### Modified Capabilities

## Impact

- `backend/internal/handler/task.go` — CreateTask 请求结构体和逻辑
- `backend/internal/model/session.go` — 新增 agent_name 字段
- `frontend/src/lib/api.ts` — createTask/createConversation 签名，fetchConversations 映射
- `frontend/src/hooks/use-conversations.ts` — 移除乐观更新，恢复 invalidateQueries
- `frontend/src/components/im/NewChatDialog.tsx` — 支持 agent name 输入
- `frontend/src/components/im/ConversationItem.tsx` — 显示 agent_name
- `frontend/src/components/chat/ChatArea.tsx` — header 显示 agent_name
- 不影响 agentend 端和 contracts/schemas/
