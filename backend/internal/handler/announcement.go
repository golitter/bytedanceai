package handler

import (
	"context"
	"log/slog"
	"path/filepath"

	"agenthub/backend/internal/model"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/agentend_client"
	"agenthub/backend/pkg/db"

	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	agentClient *agentend_client.Client
}

func NewAnnouncementHandler(agentClient *agentend_client.Client) *AnnouncementHandler {
	return &AnnouncementHandler{agentClient: agentClient}
}

type CreateAnnouncementReq struct {
	SenderID   string `json:"sender_id" binding:"required"`
	SenderName string `json:"sender_name" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Pinned     bool   `json:"pinned"`
}

// ListAnnouncements returns all announcements for a task, pinned first then by time descending.
// Supports optional ?pinned=true query parameter to filter only pinned announcements.
func (h *AnnouncementHandler) ListAnnouncements(c *gin.Context) {
	taskID := c.Param("taskId")
	pinnedOnly := c.Query("pinned") == "true"

	query := db.GetDB().Where("task_id = ?", taskID)
	if pinnedOnly {
		query = query.Where("pinned = ?", true)
	}

	var announcements []model.Announcement
	if err := query.
		Order("pinned DESC, created_at DESC").
		Find(&announcements).Error; err != nil {
		vo.InternalError(c, "failed to fetch announcements")
		return
	}

	vo.OK(c, announcements)
}

// CreateAnnouncement creates a new announcement for a task.
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	taskID := c.Param("taskId")

	var req CreateAnnouncementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "content is required")
		return
	}

	announcement := model.Announcement{
		TaskID:     taskID,
		SenderID:   req.SenderID,
		SenderName: req.SenderName,
		Content:    req.Content,
		Pinned:     req.Pinned,
	}

	if err := db.GetDB().Create(&announcement).Error; err != nil {
		vo.InternalError(c, "failed to create announcement")
		return
	}

	vo.Created(c, announcement)
}

// DeleteAnnouncement deletes an announcement by ID.
// If the announcement was pinned, notifies agentend to write an unpin event.
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")

	// Fetch before deleting so we can check pinned state
	var announcement model.Announcement
	if err := db.GetDB().Where("id = ?", id).First(&announcement).Error; err != nil {
		vo.NotFound(c, "announcement not found")
		return
	}

	if err := db.GetDB().Delete(&announcement).Error; err != nil {
		vo.InternalError(c, "failed to delete announcement")
		return
	}

	if announcement.Pinned {
		h.notifyUnpin(c.Request.Context(), announcement)
	}

	vo.OK(c, nil)
}

// notifyUnpin looks up the task's repo_path, derives shared_dir, and fires an async
// notification to agentend so the Orchestrator knows this constraint was cancelled.
func (h *AnnouncementHandler) notifyUnpin(_ context.Context, announcement model.Announcement) {
	go func() {
		var task model.Task
		if err := db.GetDB().Where("task_id = ?", announcement.TaskID).First(&task).Error; err != nil {
			slog.Warn("failed to find task for announcement unpin notification",
				"task_id", announcement.TaskID, "error", err)
			return
		}
		if task.RepoPath == "" {
			slog.Warn("task has no repo_path, skipping unpin notification",
				"task_id", announcement.TaskID)
			return
		}

		absRepoPath, err := filepath.Abs(task.RepoPath)
		if err != nil {
			slog.Warn("failed to resolve repo_path", "repo_path", task.RepoPath, "error", err)
			return
		}
		sharedDir := filepath.Join(filepath.Dir(absRepoPath), "worktrees", announcement.TaskID, "shared", ".agent")

		err = h.agentClient.NotifyAnnouncementUnpin(agentend_client.AnnouncementUnpinRequest{
			SharedDir:  sharedDir,
			Content:    announcement.Content,
			SenderName: announcement.SenderName,
		})
		if err != nil {
			slog.Warn("failed to notify agentend of announcement unpin",
				"task_id", announcement.TaskID, "announcement_id", announcement.ID, "error", err)
		}
	}()
}
