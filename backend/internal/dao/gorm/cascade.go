package gormdao

import (
	"agenthub/backend/internal/model"

	"gorm.io/gorm"
)

func cascadeDeleteBySessionIDs(tx *gorm.DB, sessionIDs []string) {
	if len(sessionIDs) == 0 {
		return
	}
	tx.Where("session_id IN ?", sessionIDs).Delete(&model.Message{})
	tx.Where("session_id IN ?", sessionIDs).Delete(&model.SessionAgent{})
	tx.Where("session_id IN ?", sessionIDs).Delete(&model.DiffSnapshot{})
}

func cascadeDeleteByTaskID(tx *gorm.DB, taskID string) {
	var sessionIDs []string
	tx.Model(&model.Session{}).Where("task_id = ?", taskID).Pluck("session_id", &sessionIDs)
	cascadeDeleteBySessionIDs(tx, sessionIDs)
	tx.Where("task_id = ?", taskID).Delete(&model.Session{})
	tx.Where("task_id = ?", taskID).Delete(&model.Announcement{})
}
