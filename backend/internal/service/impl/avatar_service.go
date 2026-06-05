package impl

import (
	"context"
	"path/filepath"
	"strings"

	"agenthub/backend/internal/dao"
	"agenthub/backend/internal/service"
	"agenthub/backend/pkg/qiniu"

	"github.com/google/uuid"
)

type AvatarService struct {
	sessionDao dao.SessionDao
	uploader   *qiniu.Uploader
}

func NewAvatarService(sessionDao dao.SessionDao, uploader *qiniu.Uploader) *AvatarService {
	return &AvatarService{sessionDao: sessionDao, uploader: uploader}
}

func (svc *AvatarService) UploadAvatar(filename string, size int64, data []byte) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	key := "avatars/" + uuid.New().String() + ext
	avatarURL, err := svc.uploader.UploadBytes(context.Background(), key, data)
	if err != nil {
		return "", service.ErrInternal("failed to upload file")
	}
	return avatarURL, nil
}

func (svc *AvatarService) UpdateSession(sessionID, agentName, avatarURL string) error {
	updates := map[string]interface{}{}
	if agentName != "" {
		updates["agent_name"] = agentName
	}
	if avatarURL != "" {
		updates["avatar_url"] = avatarURL
	}
	if len(updates) == 0 {
		return service.ErrBadRequest("at least one field (agent_name or avatar_url) is required")
	}

	updated, err := svc.sessionDao.UpdateFields(sessionID, updates)
	if err != nil {
		return err
	}
	if !updated {
		return service.ErrNotFound("session not found")
	}
	return nil
}
