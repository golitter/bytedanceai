package impl

import (
	"agenthub/backend/internal/dao"
	"agenthub/backend/internal/service"
)

type SessionService struct {
	dao dao.SessionDao
}

func NewSessionService(sessionDao dao.SessionDao) *SessionService {
	return &SessionService{dao: sessionDao}
}

func (svc *SessionService) PatchSessionStatus(sessionID, status string) (*service.SessionStatus, error) {
	if status != "inactive" {
		return nil, service.ErrBadRequest("status must be \"inactive\"")
	}

	found, err := svc.dao.DeactivateSession(sessionID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, service.ErrNotFound("session not found")
	}

	return &service.SessionStatus{
		SessionID: sessionID,
		Status:    "inactive",
	}, nil
}
