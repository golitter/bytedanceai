package gormdao

import (
	"errors"
	"time"

	"agenthub/backend/internal/model"
	"agenthub/backend/pkg/db"

	"gorm.io/gorm"
)

type TaskDao struct{}

func NewTaskDao() *TaskDao {
	return &TaskDao{}
}

func (dao *TaskDao) GetByTaskID(taskID string) (*model.Task, error) {
	var task model.Task
	if err := db.GetDB().Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (dao *TaskDao) FindRepoPathByTaskID(taskID string) (string, error) {
	var task model.Task
	if err := db.GetDB().Select("repo_path").Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return task.RepoPath, nil
}

func (dao *TaskDao) CreateTaskWithSessions(task *model.Task, sessions []model.Session, sessionAgents []model.SessionAgent) error {
	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(task).Error; err != nil {
			return err
		}
		for _, session := range sessions {
			if err := tx.Create(&session).Error; err != nil {
				return err
			}
		}
		for _, sessionAgent := range sessionAgents {
			if err := tx.Create(&sessionAgent).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *TaskDao) ListTasks() ([]model.Task, error) {
	var tasks []model.Task
	if err := db.GetDB().Order("pinned_at IS NULL, pinned_at DESC, created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (dao *TaskDao) ListSessionAgentsBySessionIDs(sessionIDs []string) ([]model.SessionAgent, error) {
	var agents []model.SessionAgent
	if len(sessionIDs) == 0 {
		return agents, nil
	}
	if err := db.GetDB().Where("session_id IN ?", sessionIDs).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func (dao *TaskDao) DeleteTaskCascade(taskID string) (bool, error) {
	found := true
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("task_id = ?", taskID).Delete(&model.Task{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			found = false
			return nil
		}
		cascadeDeleteByTaskID(tx, taskID)
		return nil
	})
	if err != nil {
		return false, err
	}
	return found, nil
}

func (dao *TaskDao) GetTaskAndSessionIDs(taskID string) (*model.Task, []string, error) {
	task, err := dao.GetByTaskID(taskID)
	if err != nil || task == nil {
		return task, nil, err
	}

	var sessionIDs []string
	if err := db.GetDB().Model(&model.Session{}).Where("task_id = ?", taskID).Pluck("session_id", &sessionIDs).Error; err != nil {
		return nil, nil, err
	}
	return task, sessionIDs, nil
}

func (dao *TaskDao) PatchTask(taskID string, updates map[string]interface{}) (bool, error) {
	result := db.GetDB().Model(&model.Task{}).Where("task_id = ?", taskID).Updates(updates)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (dao *TaskDao) EnsureSession(sessionID, taskID, agentType string) (*model.Session, bool, error) {
	var session model.Session
	result := db.GetDB().Where(model.Session{
		SessionID: sessionID,
		TaskID:    taskID,
	}).Attrs(model.Session{
		AgentType: agentType,
		Status:    "running",
	}).FirstOrCreate(&session)
	if result.Error != nil {
		return nil, false, result.Error
	}
	created := result.RowsAffected > 0 && session.CreatedAt.After(time.Time{})
	if !created {
		if err := db.GetDB().Model(&session).Update("status", "running").Error; err != nil {
			return nil, false, err
		}
	}
	return &session, created, nil
}

func (dao *TaskDao) CreateSessionAgent(agent model.SessionAgent) error {
	return db.GetDB().Create(&agent).Error
}
