## Context

当前系统的三端数据流：

- 前端创建对话 → 后端 `POST /api/tasks` 只创建 task，不创建 session
- session 只在 `POST /api/tasks/:id/run` 被调用时按需创建
- 前端用乐观更新（`setQueryData`）临时展示新对话，刷新后消失
- Session 模型只有 `agent_type`（固定枚举），没有用户自定义名称

## Goals / Non-Goals

**Goals:**
- 创建 task 时自动创建 session，消除前端乐观更新 workaround
- Session 支持 `agent_name` 字段（用户自定义，仅前端展示）
- 单聊（1 agent）和群聊（多 agent）使用同一个 `agents` 数组参数
- 后端到 agentend 的通信不变，只传 `agent_type`

**Non-Goals:**
- 不改 agentend 端代码
- 不改 contracts/schemas/ 契约文件
- 不实现群聊 UI（本次只做单聊 + 数据模型预留）

## Decisions

### D1: CreateTask 接收 `agents` 数组

```go
type AgentConfig struct {
    Type string `json:"type" binding:"required"`  // agent_type: claude-code / opencode / orchestrator
    Name string `json:"name"`                      // 用户自定义名称，不传则用 type 作为默认值
}

type CreateTaskReq struct {
    Title    string       `json:"title" binding:"required"`
    RepoPath string       `json:"repo_path"`
    Agents   []AgentConfig `json:"agents"`           // 空则不创建 session（向后兼容）
}
```

单聊传 `[{ type: "claude-code", name: "代码审查助手" }]`，群聊传多个。后端循环创建 sessions。

### D2: Session 模型新增 `agent_name`

```go
type Session struct {
    // ... 现有字段 ...
    AgentName string `gorm:"size:128" json:"agent_name"`  // 新增
}
```

GORM AutoMigrate 会自动加列。`agent_name` 为空时前端 fallback 到 `agent_type`。

### D3: 前端 createConversation 恢复正常流程

```
createConversation(agentType, agentName)
  → createTask(title, [{ type: agentType, name: agentName }])
  → fetchTask(taskId)
  → sessions 不再为空 → 取第一个 session
  → 返回 Conversation
```

移除 `crypto.randomUUID()` 和 `setQueryData` 乐观更新，恢复 `invalidateQueries`。

### D4: agent_name 不穿透到 agentend

后端 `RunTask` 转发请求给 agentend 时，只传 `agent_type`。`agent_name` 纯前端展示用。

## Risks / Trade-offs

- [数据库 migration] → GORM AutoMigrate 自动加列，无需手动 migration，安全
- [向后兼容] → `agents` 数组为空时后端不创建 session，不影响已有 API 调用
- [agent_name 可选] → 不传时 fallback 到 agent_type，前端无破坏性变更
