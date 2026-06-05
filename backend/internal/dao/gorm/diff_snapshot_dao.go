package gormdao

import (
	"errors"

	"agenthub/backend/internal/model"
	"agenthub/backend/pkg/db"

	"gorm.io/gorm"
)

type DiffSnapshotDao struct{}

func NewDiffSnapshotDao() *DiffSnapshotDao {
	return &DiffSnapshotDao{}
}

func (dao *DiffSnapshotDao) GetBySnapshotID(snapshotID string) (*model.DiffSnapshot, error) {
	var snapshot model.DiffSnapshot
	if err := db.GetDB().Where("snapshot_id = ?", snapshotID).First(&snapshot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &snapshot, nil
}

func (dao *DiffSnapshotDao) CancelPendingBySession(sessionID, excludedSnapshotID string) error {
	return db.GetDB().
		Model(&model.DiffSnapshot{}).
		Where("session_id = ? AND snapshot_id != ? AND status = ?", sessionID, excludedSnapshotID, "pending").
		Update("status", "cancelled").
		Error
}

func (dao *DiffSnapshotDao) Upsert(snapshot model.DiffSnapshot) (*model.DiffSnapshot, error) {
	result := db.GetDB().
		Where("snapshot_id = ?", snapshot.SnapshotID).
		Assign(model.DiffSnapshot{
			SessionID:   snapshot.SessionID,
			DiffContent: snapshot.DiffContent,
			Status:      snapshot.Status,
		}).
		FirstOrCreate(&snapshot)
	if result.Error != nil {
		return nil, result.Error
	}

	return &snapshot, nil
}

func (dao *DiffSnapshotDao) UpsertPending(snapshotID, sessionID, diff string) error {
	snapshot := model.DiffSnapshot{
		SnapshotID:  snapshotID,
		SessionID:   sessionID,
		DiffContent: diff,
		Status:      "pending",
	}
	_, err := dao.Upsert(snapshot)
	return err
}
