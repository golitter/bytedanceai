## 1. 后端 Session 模型

- [x] 1.1 `backend/internal/model/session.go` 新增 `AgentName` 字段 `gorm:"size:128" json:"agent_name"`

## 2. 后端 CreateTask 处理

- [x] 2.1 `backend/internal/handler/task.go` 新增 `AgentConfig` 结构体和 `CreateTaskReq.Agents` 字段
- [x] 2.2 `CreateTask` handler 中循环 `req.Agents` 创建 sessions（每个 agent 一个 session，`agent_name` 可选）
- [x] 2.3 验证 `agents` 为空时行为不变（向后兼容）

## 3. 前端 API 层

- [x] 3.1 `frontend/src/lib/api.ts` — `createTask` 新增 `agents` 参数
- [x] 3.2 `frontend/src/lib/api.ts` — `createConversation` 传 `agents` 数组，移除 `crypto.randomUUID()` 逻辑，恢复 `fetchTask` 取 session
- [x] 3.3 `frontend/src/lib/api.ts` — `Conversation` 接口新增 `agentName` 字段，`fetchConversations` 映射 `agent_name`
- [x] 3.4 `frontend/src/lib/api.ts` — `fetchAgentTypes` 保持对 `string[]` 的兼容转换

## 4. 前端 Hooks

- [x] 4.1 `frontend/src/hooks/use-conversations.ts` — `useCreateConversation` 恢复 `invalidateQueries`，移除 `setQueryData` 乐观更新

## 5. 前端组件

- [x] 5.1 `frontend/src/components/im/NewChatDialog.tsx` — 每个 agent 按钮旁增加名称输入框（可选），传递 `name` 到 `createConversation`
- [x] 5.2 `frontend/src/components/im/ConversationItem.tsx` — 显示 `agentName`（fallback `agentType`）
- [x] 5.3 `frontend/src/components/chat/ChatArea.tsx` — header 显示 `agentName`
- [x] 5.4 `frontend/src/pages/ImPage.tsx` — 传递 `agentName` 到 ChatArea

## 6. 验证

- [ ] 6.1 重启后端，验证 `POST /api/tasks` 带 `agents` 时自动创建 sessions
- [ ] 6.2 重启前端，验证创建对话 → 刷新页面 → 对话仍在
- [ ] 6.3 验证 `agent_name` 为空时 fallback 到 `agent_type` 显示
