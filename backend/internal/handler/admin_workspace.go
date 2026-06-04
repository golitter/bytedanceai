package handler

import (
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/agentend_client"
	"agenthub/backend/pkg/db"

	"github.com/gin-gonic/gin"
)

type WorkspaceItem struct {
	ID     string  `json:"id"`
	Task   string  `json:"task"`
	Agent  string  `json:"agent"`
	Branch string  `json:"branch"`
	DiskMB float64 `json:"disk_mb"`
	Status string  `json:"status"`
}

func (h *AdminHandler) GetWorkspaces(c *gin.Context) {
	database := db.GetDB()

	// 从 MySQL 获取所有 session 作为数据底座
	var sessions []model.Session
	database.Order("created_at DESC").Find(&sessions)

	// 从 agentend 获取真实工作区数据，用于补充 branch / 真实 status
	wsMap := make(map[string]*agentend_client.WorkspaceInfo) // key: session_id
	if workspaces, err := h.agentClient.ListWorkspaces(); err == nil {
		for i := range workspaces {
			wsMap[workspaces[i].SessionID] = &workspaces[i]
		}
	}

	var items []WorkspaceItem
	var active, cleaned int
	for _, s := range sessions {
		wi := wsMap[s.SessionID]
		status := s.Status // 默认用 session 状态

		var branch string
		if wi != nil {
			branch = wi.BranchName
			status = wi.Status // 有真实工作区数据则用工作区状态
		}

		items = append(items, WorkspaceItem{
			ID:     s.SessionID,
			Task:   s.TaskID,
			Agent:  s.AgentName,
			Branch: branch,
			DiskMB: 0,
			Status: status,
		})

		switch status {
		case "active":
			active++
		case "cleaned":
			cleaned++
		}
	}

	if items == nil {
		items = []WorkspaceItem{}
	}

	vo.OK(c, gin.H{
		"workspaces": items,
		"total":      len(items),
		"active":     active,
		"cleaned":    cleaned,
		"totalDisk":  0,
	})
}

func (h *AdminHandler) DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")

	// 1. 通知 agentend 清理真实工作区（git worktree + branch）
	if err := h.agentClient.CleanupWorkspace(id); err != nil {
		// agentend 可能不可用或 workspace 已清理，继续更新数据库
	}

	// 2. 更新 MySQL session 状态
	database := db.GetDB()
	database.Model(&model.Session{}).Where("session_id = ?", id).Update("status", "cleaned")

	vo.OK(c, gin.H{"success": true})
}
