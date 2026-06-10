## ADDED Requirements

### Requirement: CreateTask SHALL auto-create sessions
当 `POST /api/tasks` 请求包含 `agents` 数组时，后端 MUST 为每个 agent 配置创建一条 Session 记录。

#### Scenario: 单聊创建 session
- **WHEN** 前端发送 `POST /api/tasks { title: "Chat", agents: [{ type: "claude-code" }] }`
- **THEN** 后端创建 1 个 task 和 1 个 session，session 的 `agent_type` 为 `"claude-code"`

#### Scenario: 多 agent 创建多个 session
- **WHEN** 前端发送 `POST /api/tasks { title: "Group", agents: [{ type: "claude-code" }, { type: "opencode" }] }`
- **THEN** 后端创建 1 个 task 和 2 个 session

#### Scenario: agents 为空时向后兼容
- **WHEN** 前端发送 `POST /api/tasks { title: "Task" }` 不含 agents 字段
- **THEN** 后端只创建 task，不创建 session（向后兼容）

### Requirement: createConversation SHALL 从后端获取 session
前端 `createConversation` MUST 通过 `createTask` + `fetchTask` 获取后端创建的 session，不使用前端生成 ID。

#### Scenario: 创建对话后立即可用
- **WHEN** 用户在 NewChatDialog 选择一个 agent 创建对话
- **THEN** 对话立即出现在列表中，ChatArea 正确显示

#### Scenario: 刷新页面后对话仍在
- **WHEN** 用户创建对话后刷新页面
- **THEN** 对话仍然出现在列表中（因为后端已有 session 记录）
