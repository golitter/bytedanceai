package impl

import (
	"agenthub/backend/internal/dao"
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/service"
)

type DiffSnapshotService struct {
	dao dao.DiffSnapshotDao
}

func NewDiffSnapshotService(diffSnapshotDao dao.DiffSnapshotDao) *DiffSnapshotService {
	return &DiffSnapshotService{dao: diffSnapshotDao}
}

func (svc *DiffSnapshotService) GetDiffSnapshot(snapshotID string) (*model.DiffSnapshot, error) {
	snapshot, err := svc.dao.GetBySnapshotID(snapshotID)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, service.ErrNotFound("snapshot not found")
	}

	return snapshot, nil
}

func (svc *DiffSnapshotService) SaveDiffSnapshot(snapshotID string, input service.SaveDiffSnapshotInput) (*model.DiffSnapshot, error) {
	existing, err := svc.dao.GetBySnapshotID(snapshotID)
	if err != nil {
		return nil, err
	}
	if existing != nil && isTerminalDiffSnapshotStatus(existing.Status) {
		return nil, service.ErrConflict("snapshot is in terminal state")
	}

	if input.Status == "pending" {
		if err := svc.dao.CancelPendingBySession(input.SessionID, snapshotID); err != nil {
			return nil, err
		}
	}

	return svc.dao.Upsert(model.DiffSnapshot{
		SnapshotID:  snapshotID,
		SessionID:   input.SessionID,
		DiffContent: input.Diff,
		Status:      input.Status,
	})
}

func isTerminalDiffSnapshotStatus(status string) bool {
	switch status {
	case "committed", "reverted", "cancelled":
		return true
	default:
		return false
	}
}
