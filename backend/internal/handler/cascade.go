package handler

import (
	"agenthub/backend/internal/model"

	"gorm.io/gorm"
)

// cascadeDeleteBySessionIDs deletes all child records (Message, SessionAgent, DiffSnapshot)
// for the given session IDs within the provided transaction.
func cascadeDeleteBySessionIDs(tx *gorm.DB, sessionIDs []string) {
	if len(sessionIDs) == 0 {
		return
	}
	tx.Where("session_id IN ?", sessionIDs).Delete(&model.Message{})
	tx.Where("session_id IN ?", sessionIDs).Delete(&model.SessionAgent{})
	tx.Where("session_id IN ?", sessionIDs).Delete(&model.DiffSnapshot{})
}

// cascadeDeleteByTaskID deletes all task-level and session-level child records
// (Message, SessionAgent, DiffSnapshot, Session, Announcement) within the provided transaction.
func cascadeDeleteByTaskID(tx *gorm.DB, taskID string) {
	var sessionIDs []string
	tx.Model(&model.Session{}).Where("task_id = ?", taskID).Pluck("session_id", &sessionIDs)
	cascadeDeleteBySessionIDs(tx, sessionIDs)
	tx.Where("task_id = ?", taskID).Delete(&model.Session{})
	tx.Where("task_id = ?", taskID).Delete(&model.Announcement{})
}
