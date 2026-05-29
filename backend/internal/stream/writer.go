package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"agenthub/backend/internal/generated"
	"agenthub/backend/internal/model"
	"agenthub/backend/pkg/db"
	pkgredis "agenthub/backend/pkg/redis"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	flushInterval    = 2 * time.Second
	flushThreshold   = 500
	maxStreamLen     = 10000
	streamExpireTTL  = 600 * time.Second
	goroutineTimeout = 30 * time.Minute
)

// Registry tracks active StreamWriter goroutines by messageID.
var registry sync.Map

// IsActive returns true if a goroutine is running for the given messageID.
func IsActive(messageID string) bool {
	_, ok := registry.Load(messageID)
	return ok
}

// StreamWriter consumes an agentend SSE stream, publishes events to Redis Stream,
// and batch-flushes text content to MySQL.
// When agent_type changes in SSE events, it finalizes the current Message and creates
// a new one under the same session, keeping the original Message in streaming status
// until the entire round finishes.
type StreamWriter struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	messageID string // current (latest) message ID — updated on agent switch
	sessionID string
	taskID    string
	streamKey string

	originalMessageID string // first message ID — never changes, used for registry and Redis stream
	currentAgentType  string // tracks the current agent type from SSE events
	currentAgentName  string // tracks the current agent name from SSE events

	buf        strings.Builder
	bufLen     int
	flushedLen int
	lastSeq    string
	lastFlush  time.Time
	mu         sync.Mutex
}

// NewStreamWriter creates a new StreamWriter and registers it.
func NewStreamWriter(ctx context.Context, taskID, sessionID, messageID, agentType string) *StreamWriter {
	childCtx, cancel := context.WithTimeout(ctx, goroutineTimeout)
	key := pkgredis.StreamKey(sessionID, messageID)
	sw := &StreamWriter{
		ctx:               childCtx,
		cancel:            cancel,
		messageID:         messageID,
		sessionID:         sessionID,
		taskID:            taskID,
		streamKey:         key,
		originalMessageID: messageID,
		currentAgentType:  agentType,
		currentAgentName:  "",
	}
	registry.Store(messageID, sw)
	return sw
}

// Run consumes the agentend response body (SSE lines), publishes to Redis, and flushes to MySQL.
// This should be called in a goroutine.
func (sw *StreamWriter) Run(scanFunc func(func(line string))) {
	defer sw.finish()

	sw.wg.Add(1)
	go sw.flushLoop()

	sawError := false

	scanFunc(func(line string) {
		if sw.ctx.Err() != nil {
			return
		}
		// Publish every SSE line to Redis Stream
		sw.publishToRedis(line)

		// Detect error/done events to set final status correctly
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var event generated.StreamEvent
			if err := json.Unmarshal([]byte(data), &event); err == nil {
				switch event.Type {
				case generated.EventTypeText:
					if text, ok := event.Content["text"].(string); ok {
						// Check for agent switch
						if newAgentType, ok := event.Content["agent_type"].(string); ok && newAgentType != "" && newAgentType != sw.currentAgentType {
							newAgentName, _ := event.Content["agent"].(string)
							sw.switchAgent(newAgentType, newAgentName)
						}
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

	// Final flush
	sw.doFlush()

	status := "completed"
	if sawError {
		status = "failed"
	}
	// Finalize the current (last) sub-message
	sw.updateMessageStatus(sw.messageID, status)
	// Finalize the original message (may be the same if no agent switch happened)
	if sw.messageID != sw.originalMessageID {
		sw.updateMessageStatus(sw.originalMessageID, status)
	}
}

// switchAgent handles agent type transition: flushes buffer, finalizes current Message,
// and creates a new Message under the same session with the new agent_type.
func (sw *StreamWriter) switchAgent(newAgentType, newAgentName string) {
	sw.mu.Lock()
	hasContent := sw.bufLen > 0
	sw.mu.Unlock()

	if hasContent {
		// Flush current buffer to the current Message
		sw.doFlush()

		// Finalize current Message (not the original — it stays streaming)
		sw.updateMessageStatus(sw.messageID, "completed")

		// Create new Message in MySQL
		newMsgID := uuid.New().String()
		newMsg := model.Message{
			MessageID: newMsgID,
			TaskID:    sw.taskID,
			SessionID: sw.sessionID,
			Role:      "agent",
			Content:   "",
			Status:    "streaming",
			AgentType: newAgentType,
			AgentName: newAgentName,
		}
		if err := db.GetDB().Create(&newMsg).Error; err != nil {
			slog.Error("create sub-message failed", "error", err)
			return
		}

		// Switch internal state to new Message
		sw.mu.Lock()
		sw.messageID = newMsgID
		sw.buf.Reset()
		sw.bufLen = 0
		sw.flushedLen = 0
		sw.mu.Unlock()
	}

	// Always update agent tracking
	sw.mu.Lock()
	sw.currentAgentType = newAgentType
	sw.currentAgentName = newAgentName
	sw.mu.Unlock()
}

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

func (sw *StreamWriter) appendText(text string) {
	sw.mu.Lock()
	sw.buf.WriteString(text)
	sw.bufLen += len(text)
	shouldFlush := sw.bufLen-sw.flushedLen >= flushThreshold
	sw.mu.Unlock()

	if shouldFlush {
		sw.doFlush()
	}
}

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

func (sw *StreamWriter) doFlush() {
	sw.mu.Lock()
	content := sw.buf.String()
	seq := sw.lastSeq
	if content == "" {
		sw.mu.Unlock()
		return
	}
	sw.flushedLen = sw.bufLen
	sw.lastFlush = time.Now()
	sw.mu.Unlock()

	err := db.GetDB().Model(&model.Message{}).
		Where("message_id = ?", sw.messageID).
		Updates(map[string]interface{}{
			"content":  content,
			"last_seq": seq,
		}).Error
	if err != nil {
		slog.Error("flush to MySQL failed", "message_id", sw.messageID, "error", err)
	}
}

func (sw *StreamWriter) updateMessageStatus(messageID, status string) error {
	err := db.GetDB().Model(&model.Message{}).
		Where("message_id = ?", messageID).
		Updates(map[string]interface{}{
			"status": status,
		}).Error
	if err != nil {
		slog.Error("update message status failed", "message_id", messageID, "error", err)
	}
	return err
}

func (sw *StreamWriter) finish() {
	sw.cancel()
	sw.wg.Wait()

	// Set Redis EXPIRE on the stream
	rdb := pkgredis.GetClient()
	if rdb != nil {
		if err := rdb.Expire(context.Background(), sw.streamKey, streamExpireTTL).Err(); err != nil {
			slog.Warn("redis EXPIRE failed", "key", sw.streamKey, "error", err)
		}
	}

	registry.Delete(sw.originalMessageID)
}

// Fail marks a StreamWriter's message as failed (e.g., on context cancellation).
func (sw *StreamWriter) Fail() {
	sw.doFlush()
	sw.updateMessageStatus(sw.messageID, "failed")
	if sw.messageID != sw.originalMessageID {
		sw.updateMessageStatus(sw.originalMessageID, "failed")
	}
}

// PublishErrorAndFail writes an error event to Redis Stream, then marks the message as failed.
// Used when agentend fails before/during streaming so the frontend can see the error.
func PublishErrorAndFail(messageID, sessionID, errMsg string) {
	key := pkgredis.StreamKey(sessionID, messageID)
	rdb := pkgredis.GetClient()
	if rdb != nil {
		event := map[string]interface{}{
			"type": "error",
			"content": map[string]string{
				"message": errMsg,
			},
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

// CleanupStaleMessages marks all streaming messages as failed (called at startup).
func CleanupStaleMessages() {
	result := db.GetDB().Model(&model.Message{}).
		Where("status = ?", "streaming").
		Update("status", "failed")
	if result.RowsAffected > 0 {
		slog.Info("cleaned up stale streaming messages", "count", result.RowsAffected)
	}
}

// FormatSSE formats a text chunk as an SSE data line matching the StreamEvent contract.
func FormatSSE(text string) string {
	event := map[string]interface{}{
		"type": "text",
		"content": map[string]string{
			"text": text,
		},
	}
	data, _ := json.Marshal(event)
	return fmt.Sprintf("data: %s", string(data))
}

// FormatSSEWithMeta formats a text chunk as an SSE data line with agent metadata.
func FormatSSEWithMeta(text, agentType, agentName string) string {
	content := map[string]string{
		"text": text,
	}
	if agentType != "" {
		content["agent_type"] = agentType
	}
	if agentName != "" {
		content["agent"] = agentName
	}
	event := map[string]interface{}{
		"type":    "text",
		"content": content,
	}
	data, _ := json.Marshal(event)
	return fmt.Sprintf("data: %s", string(data))
}
