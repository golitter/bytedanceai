package handler

import (
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

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
	sessionID := c.Query("session_id")
	mode := c.Query("mode")
	primarySessionID := c.Query("primary_session_id")
	if primarySessionID == "" {
		primarySessionID = sessionID
	}
	if mode == "group" && primarySessionID == "" {
		primarySessionID = findPrimaryGroupSessionID(taskID)
	}

	// No pagination params: return all messages, has_more=false
	if limitStr == "" && beforeStr == "" {
		query := db.GetDB().Where("task_id = ?", taskID)
		if mode == "group" {
			query = applyGroupMessageVisibility(query, primarySessionID)
		} else if sessionID != "" {
			query = query.Where("session_id = ?", sessionID)
		}
		var messages []model.Message
		query.Order("created_at ASC").Order("id ASC").Find(&messages)
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
	if mode == "group" {
		query = applyGroupMessageVisibility(query, primarySessionID)
	} else if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	if beforeStr != "" {
		if beforeID, err := strconv.ParseUint(beforeStr, 10, 64); err == nil {
			query = query.Where("id < ?", beforeID)
		}
	}

	var messages []model.Message
	query.Order("id DESC").Limit(limit + 1).Find(&messages)

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}
	reverseMessages(messages)

	vo.OK(c, ListMessagesResponse{Data: messages, HasMore: hasMore})
}

func findPrimaryGroupSessionID(taskID string) string {
	var session model.Session
	if err := db.GetDB().
		Where("task_id = ? AND agent_type = ?", taskID, "orchestrator").
		First(&session).Error; err == nil {
		return session.SessionID
	}
	return ""
}

func applyGroupMessageVisibility(query *gorm.DB, primarySessionID string) *gorm.DB {
	if primarySessionID == "" {
		return query.Where("role = ? OR role = ?", "user", "agent")
	}
	return query.Where(
		"role = ? OR (role = ? AND (session_id <> ? OR (session_id = ? AND (agent_type = ? OR agent_type = '' OR agent_type IS NULL))))",
		"user",
		"agent",
		primarySessionID,
		primarySessionID,
		"orchestrator",
	)
}

func reverseMessages(messages []model.Message) {
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
}

// WindowMessages returns the group chat window messages for a given session.
func (h *MessageHandler) WindowMessages(c *gin.Context) {
	taskID := c.Param("taskId")
	sessionID := c.Query("session_id")

	if sessionID == "" {
		vo.BadRequest(c, "session_id is required")
		return
	}

	result := fetchGroupChatWindow(taskID, sessionID)
	vo.OK(c, result)
}
