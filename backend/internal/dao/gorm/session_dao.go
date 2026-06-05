package gormdao

import (
	"errors"

	"agenthub/backend/internal/model"
	"agenthub/backend/pkg/db"

	"gorm.io/gorm"
)

type SessionDao struct{}

func NewSessionDao() *SessionDao {
	return &SessionDao{}
}

func (dao *SessionDao) DeactivateSession(sessionID string) (bool, error) {
	result := db.GetDB().
		Model(&model.Session{}).
		Where("session_id = ?", sessionID).
		Update("status", "inactive")
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (dao *SessionDao) GetBySessionID(sessionID string) (*model.Session, error) {
	var session model.Session
	if err := db.GetDB().Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (dao *SessionDao) GetByTaskAndSessionID(taskID, sessionID string) (*model.Session, error) {
	var session model.Session
	if err := db.GetDB().Where("task_id = ? AND session_id = ?", taskID, sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (dao *SessionDao) ListByTaskID(taskID string) ([]model.Session, error) {
	var sessions []model.Session
	if err := db.GetDB().Where("task_id = ?", taskID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (dao *SessionDao) ListAll() ([]model.Session, error) {
	var sessions []model.Session
	if err := db.GetDB().Order("created_at DESC").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (dao *SessionDao) FindPrimaryGroupSessionID(taskID string) (string, error) {
	var session model.Session
	if err := db.GetDB().
		Where("task_id = ? AND agent_type = ?", taskID, "orchestrator").
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return session.SessionID, nil
}

func (dao *SessionDao) UpdateFields(sessionID string, updates map[string]interface{}) (bool, error) {
	result := db.GetDB().
		Model(&model.Session{}).
		Where("session_id = ?", sessionID).
		Updates(updates)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (dao *SessionDao) UpdateSoul(sessionID, soulMD string) (bool, error) {
	result := db.GetDB().
		Model(&model.Session{}).
		Where("session_id = ?", sessionID).
		Update("soul_md", soulMD)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (dao *SessionDao) UpdateStatusByTask(sessionID, taskID, status string) error {
	return db.GetDB().
		Model(&model.Session{}).
		Where("session_id = ? AND task_id = ?", sessionID, taskID).
		Update("status", status).
		Error
}
