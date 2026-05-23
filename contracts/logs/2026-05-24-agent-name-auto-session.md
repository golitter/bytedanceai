# Session 新增 agent_name 字段 & CreateTask 自动创建 session

## 变更原因

前端创建对话时后端不会自动创建 session，导致刷新后对话消失。本次修改让 CreateTask 接受 `agents` 数组并自动创建 session，同时 Session 新增 `agent_name` 字段支持用户自定义名称。

## 变更文件

无 schema 文件变更（`agent_name` 不穿透到 agentend，仅影响前端展示和后端存储）。

## 对比结果

### Backend

- `model/session.go` Session 新增 `AgentName string` 字段 `gorm:"size:128" json:"agent_name"`
- `handler/task.go` 新增 `AgentConfig` 结构体（`type` + `name`），`CreateTaskReq` 新增 `Agents []AgentConfig` 字段
- `CreateTask` handler 循环 `req.Agents` 自动创建 sessions

### Frontend

- `lib/api.ts` `createTask` 新增 `agents` 参数；`createConversation` 通过 `createTask` + `fetchTask` 从后端获取 session，移除 `crypto.randomUUID()` 前端生成 ID
- `lib/api.ts` `Conversation` 接口新增 `agentName` 字段
- `hooks/use-conversations.ts` 恢复 `invalidateQueries`，移除 `setQueryData` 乐观更新
- `NewChatDialog` 支持为 agent 输入自定义名称
- `ConversationItem` / `ChatArea` / `ImPage` 显示 `agentName`（fallback `agentType`）

## 跨端影响

- **Frontend**: `Conversation` 接口新增 `agentName`，对话创建流程改为后端驱动
- **Backend**: Session 模型新增 `agent_name` 列（GORM AutoMigrate），CreateTask 新增 `agents` 参数
- **Agentend**: 无影响 — `agent_name` 不穿透，RunTask 仍只传 `agent_type`

## 契约变更

- 无 schema 变更（`contracts/schemas/*.yaml` 未修改）
- `agents` 为空时后端不创建 session，完全向后兼容
