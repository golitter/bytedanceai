package impl

import (
	"log/slog"
	"path/filepath"

	"agenthub/backend/internal/dao"
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/service"
	"agenthub/backend/pkg/agentend_client"
)

type AnnouncementService struct {
	announcementDao dao.AnnouncementDao
	taskDao         dao.TaskDao
	agentClient     *agentend_client.Client
}

func NewAnnouncementService(announcementDao dao.AnnouncementDao, taskDao dao.TaskDao, agentClient *agentend_client.Client) *AnnouncementService {
	return &AnnouncementService{
		announcementDao: announcementDao,
		taskDao:         taskDao,
		agentClient:     agentClient,
	}
}

func (svc *AnnouncementService) ListAnnouncements(taskID string, pinnedOnly bool) ([]model.Announcement, error) {
	return svc.announcementDao.ListByTaskID(taskID, pinnedOnly)
}

func (svc *AnnouncementService) CreateAnnouncement(taskID string, input service.CreateAnnouncementInput) (*model.Announcement, error) {
	return svc.announcementDao.CreateAnnouncement(model.Announcement{
		TaskID:     taskID,
		SenderID:   input.SenderID,
		SenderName: input.SenderName,
		Content:    input.Content,
		Pinned:     input.Pinned,
	})
}

func (svc *AnnouncementService) DeleteAnnouncement(id string) error {
	announcement, err := svc.announcementDao.DeleteAnnouncement(id)
	if err != nil {
		return err
	}
	if announcement == nil {
		return service.ErrNotFound("announcement not found")
	}
	if announcement.Pinned {
		svc.notifyUnpin(*announcement)
	}
	return nil
}

func (svc *AnnouncementService) notifyUnpin(announcement model.Announcement) {
	go func() {
		repoPath, err := svc.taskDao.FindRepoPathByTaskID(announcement.TaskID)
		if err != nil {
			slog.Warn("failed to find task for announcement unpin notification", "task_id", announcement.TaskID, "error", err)
			return
		}
		if repoPath == "" {
			slog.Warn("task has no repo_path, skipping unpin notification", "task_id", announcement.TaskID)
			return
		}

		absRepoPath, err := filepath.Abs(repoPath)
		if err != nil {
			slog.Warn("failed to resolve repo_path", "repo_path", repoPath, "error", err)
			return
		}
		sharedDir := filepath.Join(filepath.Dir(absRepoPath), "worktrees", announcement.TaskID, "shared", ".agent")

		if err := svc.agentClient.NotifyAnnouncementUnpin(agentend_client.AnnouncementUnpinRequest{
			SharedDir:  sharedDir,
			Content:    announcement.Content,
			SenderName: announcement.SenderName,
		}); err != nil {
			slog.Warn("failed to notify agentend of announcement unpin", "task_id", announcement.TaskID, "announcement_id", announcement.ID, "error", err)
		}
	}()
}
