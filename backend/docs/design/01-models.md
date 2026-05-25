# Models — 数据模型

## 实现了什么

使用 GORM 定义了三个核心数据模型（Task、Session、Message），构成 Task 1:N Session、Session 1:N Message 的层级关系，支撑多 Agent 会话管理。

## 怎么实现的

### Task — 顶层任务实体 (`internal/model/task.go`)

Task 是顶层实体，代表一个项目任务。`task_id` 为 UUID，供 AgentEnd 决定 git branch 和 worktree 隔离。

```go
type Task struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TaskID    string    `gorm:"uniqueIndex;size:36" json:"task_id"`
	Title     string    `gorm:"size:255" json:"title"`
	RepoPath  string    `gorm:"size:512" json:"repo_path"`
	Status    string    `gorm:"size:32;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

- `TaskID`：后端通过 `google/uuid` v4 生成，唯一索引
- `RepoPath`：仓库路径，运行时注入 AgentRequest
- `Status`：默认 `"active"`

### Session — Agent 会话 (`internal/model/session.go`)

Session 从属于 Task，代表一个 Agent 的会话。`session_id` 由调用方传入，与 `task_id` 组合映射到 AgentEnd 的 `cli_session_id`。

```go
type Session struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	SessionID string    `gorm:"uniqueIndex;size:128" json:"session_id"`
	TaskID    string    `gorm:"index;size:36" json:"task_id"`
	AgentType string    `gorm:"size:64" json:"agent_type"`
	AgentName string    `gorm:"size:128" json:"agent_name"`
	AvatarURL string    `gorm:"size:512" json:"avatar_url,omitempty"`
	Status    string    `gorm:"size:32;default:running" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

- `TaskID`：索引字段，关联 Task
- `AgentType`：Agent 类型（claude-code / opencode / orchestrator）
- `AgentName` / `AvatarURL`：Agent 的显示名称和头像，通过 `PUT /api/sessions/:sessionId` 更新
- `Status`：`active` -> `running` -> `completed` / `failed` / `inactive`

### Message — 消息记录 (`internal/model/message.go`)

Message 从属于 Session，记录用户和 Agent 的每条消息。流式场景下 `status` 为 `streaming`，结束后变为 `completed` 或 `failed`。

```go
type Message struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	MessageID string    `gorm:"uniqueIndex;size:36" json:"message_id"`
	TaskID    string    `gorm:"index;size:36" json:"task_id"`
	SessionID string    `gorm:"size:128" json:"session_id"`
	Role      string    `gorm:"size:16" json:"role"`
	Content   string    `gorm:"type:longtext" json:"content"`
	Status    string    `gorm:"size:16;default:completed" json:"status"`
	LastSeq   string    `gorm:"size:64;default:''" json:"last_seq"`
	AgentType string    `gorm:"size:64" json:"agent_type,omitempty"`
	AgentName string    `gorm:"size:128" json:"agent_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
```

- `MessageID`：UUID，唯一索引
- `Role`：`"user"` 或 `"agent"`
- `Content`：`longtext` 类型，Agent 消息由 StreamWriter 批量刷写
- `LastSeq`：Redis Stream 的最后消费位置，用于断线重连时从 MySQL 历史恢复后跳过已消费事件
- `Status`：`streaming`（流式中） / `completed` / `failed`

### 实体关系

```
Task 1:N Session 1:N Message
  │         │          │
  ├─ task_id ◄─────────┤ (FK)
  │    (uniqueIndex)   │
  │         │          │
  │    session_id ◄────┤ (FK)
  │    (uniqueIndex)   │
  │                    │
  └────────────────────┘  通过 task_id / session_id 字段关联，
                          未使用 GORM 外键约束（软关联）
```
