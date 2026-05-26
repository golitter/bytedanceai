# Handlers — HTTP 处理器

## 实现了什么

基于 Gin 框架实现了 8 组 HTTP 处理器，覆盖 Task CRUD、Session 管理、消息查询、Agent 类型枚举、头像上传、SSE 流式订阅、Diff 快照和工作区代理，构成完整的 RESTful API 层。

## 怎么实现的

### Task CRUD (`internal/handler/task.go`)

`TaskHandler` 通过闭包注入 `agentend_client.Client`，提供任务的创建、列表、详情、删除和运行。

构造函数与请求结构体：

```go
type TaskHandler struct {
	agentClient *agentend_client.Client
}

func NewTaskHandler(agentClient *agentend_client.Client) *TaskHandler {
	return &TaskHandler{agentClient: agentClient}
}

type CreateTaskReq struct {
	Title    string        `json:"title" binding:"required"`
	RepoPath string        `json:"repo_path"`
	Agents   []AgentConfig `json:"agents"`
}
```

**CreateTask** — 创建 Task 并可选地同时创建 Session 和 SessionAgent：

```go
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "title is required")
		return
	}
	t := model.Task{
		TaskID:   uuid.New().String(),
		Title:    req.Title,
		RepoPath: req.RepoPath,
		Status:   "active",
	}
	if err := db.GetDB().Create(&t).Error; err != nil {
		vo.InternalError(c, err.Error())
		return
	}
	for _, agent := range req.Agents {
		sid := uuid.New().String()
		s := model.Session{
			SessionID: sid,
			TaskID:    t.TaskID,
			AgentType: agent.Type,
			AgentName: agent.Name,
			Status:    "active",
		}
		sa := model.SessionAgent{
			SessionID: sid,
			AgentType: agent.Type,
			AgentName: agent.Name,
		}
		db.GetDB().Create(&s)
		db.GetDB().Create(&sa)
	}
	vo.Created(c, t)
}
```

**RunTask** — 核心 handler，启动 Agent 会话并返回流式 message_id：

```go
type RunTaskReq struct {
	Message   string `json:"message" binding:"required"`
	AgentType string `json:"agent_type"`
	SessionID string `json:"session_id" binding:"required"`
}
```

流程：验证 Task 存在 -> 保存用户消息 -> 查找或创建 Session -> 创建 agent Message（status: streaming）-> 启动后台 goroutine 调用 AgentEnd SSE -> 返回 `202 Accepted` + `message_id`。

后台 goroutine 核心逻辑：

```go
go func() {
	agentReq := &generated.AgentRequest{
		TaskId:    taskID,
		SessionId: req.SessionID,
		Message:   req.Message,
		AgentType: generated.AgentType(agentType),
		Stream:    true,
	}
	resp, err := h.agentClient.StreamAgent(agentReq)
	// ...
	sw := stream.NewStreamWriter(context.Background(), taskID, req.SessionID, messageID, agentType)
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	sw.Run(func(fn func(string)) {
		for scanner.Scan() {
			fn(scanner.Text())
		}
	})
}()
```

**ValidateRepoPath** — 转发仓库路径校验到 AgentEnd：

```go
func (h *TaskHandler) ValidateRepoPath(c *gin.Context) {
	var req ValidateRepoPathReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "repo_path is required")
		return
	}
	result, err := h.agentClient.ValidateRepoPath(req.RepoPath)
	// ...
	vo.OK(c, result)
}
```

### Session 管理 (`internal/handler/session.go`)

`SessionHandler` 提供 Session 状态的 PATCH 更新，当前仅支持停用（`inactive`）。

```go
type SessionHandler struct{}

type PatchSessionReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *SessionHandler) PatchSession(c *gin.Context) {
	// 仅允许 status = "inactive"
	sessionID := c.Param("sessionId")
	result := db.GetDB().Model(&model.Session{}).Where("session_id = ?", sessionID).Update("status", "inactive")
	// ...
}
```

### 消息查询 (`internal/handler/message.go`)

`MessageHandler` 按 Task 查询历史消息，支持 cursor 分页加载，按创建时间升序排列。

```go
type MessageHandler struct{}

type ListMessagesResponse struct {
	Data    []model.Message `json:"data"`
	HasMore bool            `json:"has_more"`
}

func (h *MessageHandler) ListMessages(c *gin.Context) {
	taskID := c.Param("taskId")
	var task model.Task
	if err := db.GetDB().Where("task_id = ?", taskID).First(&task).Error; err != nil {
		vo.NotFound(c, "task not found")
		return
	}

	limitStr := c.Query("limit")
	beforeStr := c.Query("before")

	// No pagination params: return all messages, has_more=false
	if limitStr == "" && beforeStr == "" {
		var messages []model.Message
		db.GetDB().Where("task_id = ?", taskID).Order("created_at ASC").Find(&messages)
		vo.OK(c, ListMessagesResponse{Data: messages, HasMore: false})
		return
	}

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	query := db.GetDB().Where("task_id = ?", taskID)
	if beforeStr != "" {
		if beforeID, err := strconv.ParseUint(beforeStr, 10, 64); err == nil {
			query = query.Where("id < ?", beforeID)
		}
	}

	var messages []model.Message
	query.Order("created_at ASC").Limit(limit + 1).Find(&messages)

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	vo.OK(c, ListMessagesResponse{Data: messages, HasMore: hasMore})
}
```

- 无分页参数（`limit` / `before`）时返回全部消息，`has_more=false`
- `limit` 参数控制每页数量，默认 20；`before` 参数为 cursor（自增 ID），用于向前翻页
- 响应格式为 `{data: [...], has_more: bool}`

### Agent 类型枚举 (`internal/handler/agent.go`)

返回硬编码的四种 Agent 类型列表。

```go
type AgentHandler struct{}

var agentTypes = []string{"claude-code", "opencode", "orchestrator", "codex"}

func (h *AgentHandler) ListAgentTypes(c *gin.Context) {
	vo.OK(c, agentTypes)
}
```

### 头像上传 (`internal/handler/avatar.go`)

`AvatarHandler` 通过闭包注入七牛云 `Uploader`，处理 multipart 文件上传并更新 Session 元数据。

```go
type AvatarHandler struct {
	uploader *qiniu.Uploader
}

const maxAvatarSize = 2 << 20 // 2MB

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}
```

**UploadAvatar** — 文件校验 -> 上传七牛云 -> 返回 CDN URL：

```go
func (h *AvatarHandler) UploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("avatar")
	// 文件大小校验（2MB）、扩展名校验
	key := "avatars/" + uuid.New().String() + ext
	data, _ := io.ReadAll(file)
	avatarURL, err := h.uploader.UploadBytes(c.Request.Context(), key, data)
	vo.OK(c, gin.H{"avatar_url": avatarURL})
}
```

**UpdateSession** — 更新 SessionAgent 的 `agent_name` 和 `avatar_url`（Upsert 到 session_agents 表）：

```go
type UpdateSessionReq struct {
	AgentName string `json:"agent_name"`
	AvatarURL string `json:"avatar_url"`
}

func (h *AvatarHandler) UpdateSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	updates := map[string]interface{}{}
	if req.AgentName != "" { updates["agent_name"] = req.AgentName }
	if req.AvatarURL != "" { updates["avatar_url"] = req.AvatarURL }

	var sa model.SessionAgent
	result := db.GetDB().Where("session_id = ?", sessionID).First(&sa)
	if result.Error != nil {
		// No existing record — create one
		sa = model.SessionAgent{
			SessionID: sessionID,
			AgentName: req.AgentName,
			AvatarURL: req.AvatarURL,
		}
		db.GetDB().Create(&sa)
	} else {
		db.GetDB().Model(&sa).Updates(updates)
	}
}
```

### SSE 流式订阅 (`internal/handler/stream.go`)

`StreamHandler` 实现前端 SSE 订阅，支持 streaming / completed / failed 三种消息状态的分发。

```go
type StreamHandler struct{}

func (h *StreamHandler) ServeStream(c *gin.Context) {
	sessionID := c.Query("session_id")
	messageID := c.Query("message_id")
	// ...
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	switch msg.Status {
	case "streaming":
		h.serveStreaming(c, &msg)
	case "completed":
		h.serveCompleted(c, &msg)
	case "failed":
		h.serveFailed(c, &msg)
	}
}
```

**serveStreaming** 的两阶段分发：

```go
func (h *StreamHandler) serveStreaming(c *gin.Context, msg *model.Message) {
	// Phase 1: 从 MySQL 读取已刷写的历史文本
	if msg.Content != "" {
		chunks := splitContent(msg.Content, 500)
		for _, chunk := range chunks {
			fmt.Fprintf(c.Writer, "%s\n\n", stream.FormatSSE(chunk))
		}
	}
	// Phase 2: 从 Redis Stream 阻塞读取实时事件
	// 先非阻塞排空存量消息，再 5s 超时阻塞等待新消息
	// 当 StreamWriter goroutine 结束后发送 done/error 事件并退出
}
```

`splitContent` 将长文本按 500 字符分块，优先在换行符处断开：

```go
func splitContent(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}
	var chunks []string
	for len(text) > 0 {
		end := maxLen
		if end > len(text) { end = len(text) }
		if end < len(text) {
			if idx := lastIndexByte(text[:end], '\n'); idx > end/2 {
				end = idx + 1
			}
		}
		chunks = append(chunks, text[:end])
		text = text[end:]
	}
	return chunks
}
```

### Diff 快照 (`internal/handler/diff_snapshot.go`)

`DiffSnapshotHandler` 管理 Diff 快照的读取和 Upsert。

```go
type DiffSnapshotHandler struct{}
```

**GetDiffSnapshot** — 按 snapshot_id 查询快照，前端 DiffCard 加载时调用：

```go
func (h *DiffSnapshotHandler) GetDiffSnapshot(c *gin.Context) {
	snapshotID := c.Param("snapshotId")
	var snap model.DiffSnapshot
	if err := db.GetDB().Where("snapshot_id = ?", snapshotID).First(&snap).Error; err != nil {
		vo.NotFound(c, "snapshot not found")
		return
	}
	vo.OK(c, snap)
}
```

**SaveDiffSnapshot** — Upsert 保存快照。终态（committed/reverted/cancelled）的快照不可覆盖（返回 409）。同一 session 新建 pending 快照时自动取消同 session 的其他 pending 快照：

```go
if req.Status == "pending" {
	d.Model(&model.DiffSnapshot{}).
		Where("session_id = ? AND snapshot_id != ? AND status = ?", req.SessionID, snapshotID, "pending").
		Update("status", "cancelled")
}
```

### 工作区代理 (`internal/handler/workspace.go`)

`WorkspaceHandler` 通过闭包注入 `agentend_client.Client`，将前端请求代理到 AgentEnd 的工作区 API。提供两类路由：

**直接工作区路由**（`/api/workspace/:id/...`）— 直接使用 workspace ID：

```go
ws.GET("/:id/files/*filepath", workspaceHandler.ReadFile)
ws.PUT("/:id/files/*filepath", workspaceHandler.WriteFile)
ws.GET("/:id/diff", workspaceHandler.GetDiff)
ws.POST("/:id/commit", workspaceHandler.Commit)
ws.POST("/:id/revert", workspaceHandler.Revert)
ws.POST("/:id/preview/start", workspaceHandler.StartPreview)
ws.POST("/:id/preview/stop", workspaceHandler.StopPreview)
```

**Session 路由**（`/api/session/:sessionId/...`）— 先通过 `resolveWorkspaceID` 查询 AgentEnd 获取 workspace ID，再代理：

```go
func (h *WorkspaceHandler) resolveWorkspaceID(sessionID string) (string, error) {
	url := fmt.Sprintf("%s/v1/workspace/by-session/%s", h.agentClient.BaseURL(), sessionID)
	resp, err := http.DefaultClient.Get(url)
	// ...
	var ws struct { ID string `json:"id"` }
	json.NewDecoder(resp.Body).Decode(&ws)
	return ws.ID, nil
}
```

`proxy` 方法为通用代理函数，转发请求到 AgentEnd 并流式回传响应（通过 `io.Copy`）。
