package handler

import (
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/db"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct{}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{}
}

type PatchSessionReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *SessionHandler) PatchSession(c *gin.Context) {
	var req PatchSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "status is required")
		return
	}

	if req.Status != "inactive" {
		vo.BadRequest(c, "status must be \"inactive\"")
		return
	}

	sessionID := c.Param("sessionId")
	result := db.GetDB().Model(&model.Session{}).Where("session_id = ?", sessionID).Update("status", "inactive")
	if result.RowsAffected == 0 {
		vo.NotFound(c, "session not found")
		return
	}

	vo.OK(c, gin.H{"session_id": sessionID, "status": "inactive"})
}
