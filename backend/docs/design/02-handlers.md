# Handlers — HTTP 处理器

## 实现了什么

基于 Gin 框架实现了 6 组 HTTP 处理器，覆盖 Task CRUD、Session 管理、消息查询、Agent 类型枚举、头像上传和 SSE 流式订阅，构成完整的 RESTful API 层。

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

**CreateTask** — 创建 Task 并可选地同时创建 Session：

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
		s := model.Session{
			SessionID: uuid.New().String(),
			TaskID:    t.TaskID,
			AgentType: agent.Type,
			AgentName: agent.Name,
			Status:    "active",
		}
		db.GetDB().Create(&s)
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

`MessageHandler` 按 Task 查询所有历史消息，按创建时间升序排列。

```go
type MessageHandler struct{}

func (h *MessageHandler) ListMessages(c *gin.Context) {
	taskID := c.Param("taskId")
	var task model.Task
	if err := db.GetDB().Where("task_id = ?", taskID).First(&task).Error; err != nil {
		vo.NotFound(c, "task not found")
		return
	}
	var messages []model.Message
	db.GetDB().Where("task_id = ?", taskID).Order("created_at ASC").Find(&messages)
	vo.OK(c, messages)
}
```

### Agent 类型枚举 (`internal/handler/agent.go`)

返回硬编码的三种 Agent 类型列表。

```go
type AgentHandler struct{}

var agentTypes = []string{"claude-code", "opencode", "orchestrator"}

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

**UpdateSession** — 更新 Session 的 `agent_name` 和 `avatar_url`：

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
	db.GetDB().Model(&model.Session{}).Where("session_id = ?", sessionID).Updates(updates)
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
