package mockdao

type SessionDao struct {
	DeactivateSessionFunc func(sessionID string) (bool, error)
}

func NewSessionDao() *SessionDao {
	return &SessionDao{}
}

func (dao *SessionDao) DeactivateSession(sessionID string) (bool, error) {
	if dao.DeactivateSessionFunc != nil {
		return dao.DeactivateSessionFunc(sessionID)
	}

	return false, nil
}
