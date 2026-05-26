# Stream — SSE 流式中转

## 实现了什么

`StreamWriter` 消费 AgentEnd 的 SSE 响应流，将每个事件实时写入 Redis Stream，同时将文本内容定时批量刷写到 MySQL Message 表。实现了 Redis Stream key 管理、全局 goroutine 注册表、启动时残留消息清理等机制。

## 怎么实现的

### StreamWriter 结构体 (`internal/stream/writer.go`)

```go
type StreamWriter struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	messageID string
	sessionID string
	taskID    string
	agentType string
	streamKey string

	buf       strings.Builder
	bufLen    int
	lastSeq   string
	lastFlush time.Time
	mu        sync.Mutex
}
```

常量配置：

```go
const (
	flushInterval    = 2 * time.Second   // 定时刷写间隔
	flushThreshold   = 500               // 缓冲区字节数阈值
	maxStreamLen     = 10000             // Redis Stream 最大长度
	streamExpireTTL  = 600 * time.Second // 流结束后 Redis key TTL
	goroutineTimeout = 30 * time.Minute  // 单次 goroutine 超时
)
```

### 创建与注册

通过 `NewStreamWriter` 创建实例并注册到全局 `sync.Map` 注册表，设置 30 分钟超时 context：

```go
var registry sync.Map

func IsActive(messageID string) bool {
	_, ok := registry.Load(messageID)
	return ok
}

func NewStreamWriter(ctx context.Context, taskID, sessionID, messageID, agentType string) *StreamWriter {
	childCtx, cancel := context.WithTimeout(ctx, goroutineTimeout)
	key := pkgredis.StreamKey(sessionID, messageID)
	sw := &StreamWriter{
		ctx: childCtx, cancel: cancel,
		messageID: messageID, sessionID: sessionID,
		taskID: taskID, agentType: agentType, streamKey: key,
	}
	registry.Store(messageID, sw)
	return sw
}
```

Redis Stream key 格式为 `agent:{sessionId}:{messageId}`，由 `pkg/redis.StreamKey()` 生成：

```go
func StreamKey(sessionID, messageID string) string {
	return fmt.Sprintf("agent:%s:%s", sessionID, messageID)
}
```

### Run — 主消费循环

`Run` 接收一个扫描函数，将每行 SSE 数据同时发送到 Redis 和本地缓冲区：

```go
func (sw *StreamWriter) Run(scanFunc func(func(line string))) {
	defer sw.finish()
	sw.wg.Add(1)
	go sw.flushLoop()

	sawError := false
	scanFunc(func(line string) {
		if sw.ctx.Err() != nil { return }
		sw.publishToRedis(line)

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var event generated.StreamEvent
			if err := json.Unmarshal([]byte(data), &event); err == nil {
				switch event.Type {
				case generated.EventTypeText:
					if text, ok := event.Content["text"].(string); ok {
						sw.appendText(text)
					}
				case generated.EventTypeDone:
					// normal end
				case generated.EventTypeError:
					sawError = true
					if errMsg, ok := event.Content["error"].(string); ok && errMsg != "" {
						sw.appendText("[Error] " + errMsg)
					}
				}
			}
		}
	})

	sw.doFlush()
	if sawError {
		sw.updateStatus("failed")
	} else {
		sw.updateStatus("completed")
	}
}
```

### Redis Stream 写入

每行 SSE 数据通过 `XADD` 写入 Redis Stream，自动截断到 `maxStreamLen` 条：

```go
func (sw *StreamWriter) publishToRedis(line string) {
	rdb := pkgredis.GetClient()
	if rdb == nil {
		return
	}
	seq, err := rdb.XAdd(sw.ctx, &redis.XAddArgs{
		Stream: sw.streamKey,
		MaxLen: maxStreamLen,
		Approx: true,
		Values: map[string]interface{}{
			"data": line,
		},
	}).Result()
	if err != nil {
		slog.Error("redis XADD failed", "key", sw.streamKey, "error", err)
		return
	}
	sw.mu.Lock()
	sw.lastSeq = seq
	sw.mu.Unlock()
}
```

### MySQL 批量刷写

双触发机制：缓冲区满 500 字节或 2 秒定时器到期。

定时刷写 goroutine：

```go
func (sw *StreamWriter) flushLoop() {
	defer sw.wg.Done()
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			sw.mu.Lock()
			hasContent := sw.bufLen > 0
			sw.mu.Unlock()
			if hasContent {
				sw.doFlush()
			}
		case <-sw.ctx.Done():
			return
		}
	}
}
```

刷写逻辑 — 替换式写入（每次刷写写入完整内容，非追加）：

```go
func (sw *StreamWriter) doFlush() {
	sw.mu.Lock()
	content := sw.buf.String()
	seq := sw.lastSeq
	if content == "" {
		sw.mu.Unlock()
		return
	}
	sw.buf.Reset()
	sw.bufLen = 0
	sw.mu.Unlock()

	db.GetDB().Model(&model.Message{}).
		Where("message_id = ?", sw.messageID).
		Updates(map[string]interface{}{
			"content":  content,
			"last_seq": seq,
		})
}
```

### 结束与清理

流结束后设置 Redis key 600 秒 TTL，从注册表移除：

```go
func (sw *StreamWriter) finish() {
	sw.cancel()
	sw.wg.Wait()

	rdb := pkgredis.GetClient()
	if rdb != nil {
		if err := rdb.Expire(context.Background(), sw.streamKey, streamExpireTTL).Err(); err != nil {
			slog.Warn("redis EXPIRE failed", "key", sw.streamKey, "error", err)
		}
	}
	registry.Delete(sw.messageID)
}
```

### 辅助函数

**PublishErrorAndFail** — 当 AgentEnd 在流式之前或中途失败时，向 Redis Stream 写入 error 事件并标记 Message 为 failed：

```go
func PublishErrorAndFail(messageID, sessionID, errMsg string) {
	key := pkgredis.StreamKey(sessionID, messageID)
	rdb := pkgredis.GetClient()
	if rdb != nil {
		event := map[string]interface{}{
			"type": "error",
			"content": map[string]string{"message": errMsg},
		}
		data, _ := json.Marshal(event)
		rdb.XAdd(context.Background(), &redis.XAddArgs{
			Stream: key,
			MaxLen: maxStreamLen,
			Approx: true,
			Values: map[string]interface{}{
				"data": fmt.Sprintf("data: %s", string(data)),
			},
		}).Result()
		rdb.Expire(context.Background(), key, streamExpireTTL)
	}
	db.GetDB().Model(&model.Message{}).Where("message_id = ?", messageID).Update("status", "failed")
}
```

**CleanupStaleMessages** — 服务启动时将所有残留的 `streaming` 状态消息标记为 `failed`：

```go
func CleanupStaleMessages() {
	result := db.GetDB().Model(&model.Message{}).
		Where("status = ?", "streaming").
		Update("status", "failed")
	if result.RowsAffected > 0 {
		slog.Info("cleaned up stale streaming messages", "count", result.RowsAffected)
	}
}
```

**FormatSSE** — 将文本块格式化为 SSE data 行，供 StreamHandler 回放历史内容：

```go
func FormatSSE(text string) string {
	event := map[string]interface{}{
		"type": "text",
		"content": map[string]string{"text": text},
	}
	data, _ := json.Marshal(event)
	return fmt.Sprintf("data: %s", string(data))
}
```

### 事件流全景

```
AgentEnd (FastAPI :8001)
    │ SSE response body
    ▼
RunTask handler goroutine
    │ bufio.Scanner 逐行读取
    ├──────────────────────┐
    ▼                      ▼
publishToRedis()       appendText()
    │ XADD                │ 缓冲区 >= 500B?
    ▼                     ▼
Redis Stream          doFlush()
agent:{sid}:{mid}         │ UPDATE content, last_seq
    │                      ▼
    ▼                  MySQL messages 表
StreamHandler.serveStreaming()
    │ XREAD (阻塞 5s)
    ▼
前端 EventSource
```
