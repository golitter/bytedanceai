package handler

import (
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeleteSessionsRequest struct {
	SessionIDs []string `json:"session_ids" binding:"required"`
}

func (h *AdminHandler) DeleteSessions(c *gin.Context) {
	var req DeleteSessionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "session_ids is required")
		return
	}

	deleted := 0
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, sid := range req.SessionIDs {
			// Verify session exists
			var count int64
			tx.Model(&model.Session{}).Where("session_id = ?", sid).Count(&count)
			if count == 0 {
				continue
			}
			tx.Where("session_id = ?", sid).Delete(&model.Message{})
			tx.Where("session_id = ?", sid).Delete(&model.SessionAgent{})
			tx.Where("session_id = ?", sid).Delete(&model.DiffSnapshot{})
			tx.Where("session_id = ?", sid).Delete(&model.Session{})
			deleted++
		}
		return nil
	})
	if err != nil {
		vo.InternalError(c, "failed to delete sessions")
		return
	}

	vo.OK(c, gin.H{"deleted": deleted})
}
