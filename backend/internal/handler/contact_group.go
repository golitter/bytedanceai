package handler

import (
	"errors"

	"agenthub/backend/internal/model"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactGroupHandler struct{}

func NewContactGroupHandler() *ContactGroupHandler {
	return &ContactGroupHandler{}
}

// ── Request types ──

type CreateGroupReq struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGroupReq struct {
	Name string `json:"name" binding:"required"`
}

type AddItemReq struct {
	TaskID string `json:"task_id" binding:"required"`
}

// ── Response types ──

type GroupItem struct {
	TaskID    string `json:"task_id"`
	SortOrder int    `json:"sort_order"`
}

type GroupWithItems struct {
	GroupID   string      `json:"group_id"`
	Name      string      `json:"name"`
	SortOrder int         `json:"sort_order"`
	Items     []GroupItem `json:"items"`
}

type ListGroupsResponse struct {
	Groups           []GroupWithItems `json:"groups"`
	UngroupedTaskIDs []string         `json:"ungrouped_task_ids"`
}

// ── Handlers ──

func (h *ContactGroupHandler) ListGroups(c *gin.Context) {
	var groups []model.ContactGroup
	db.GetDB().Order("sort_order ASC, created_at ASC").Find(&groups)

	// Collect all grouped task IDs
	groupedSet := make(map[string]bool)
	result := make([]GroupWithItems, 0, len(groups))

	for _, g := range groups {
		var items []model.ContactGroupItem
		db.GetDB().Where("group_id = ?", g.GroupID).Order("sort_order ASC").Find(&items)

		groupItems := make([]GroupItem, 0, len(items))
		for _, item := range items {
			groupItems = append(groupItems, GroupItem{TaskID: item.TaskID, SortOrder: item.SortOrder})
			groupedSet[item.TaskID] = true
		}

		result = append(result, GroupWithItems{
			GroupID:   g.GroupID,
			Name:      g.Name,
			SortOrder: g.SortOrder,
			Items:     groupItems,
		})
	}

	// Find ungrouped tasks
	var allTaskIDs []string
	db.GetDB().Model(&model.Task{}).Where("status = ?", "active").Pluck("task_id", &allTaskIDs)
	ungrouped := make([]string, 0)
	for _, tid := range allTaskIDs {
		if !groupedSet[tid] {
			ungrouped = append(ungrouped, tid)
		}
	}

	vo.OK(c, ListGroupsResponse{Groups: result, UngroupedTaskIDs: ungrouped})
}

func (h *ContactGroupHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "name is required")
		return
	}

	g := model.ContactGroup{
		GroupID: uuid.New().String()[:8],
		Name:    req.Name,
	}
	if err := db.GetDB().Create(&g).Error; err != nil {
		vo.InternalError(c, "failed to create group")
		return
	}
	vo.Created(c, g)
}

func (h *ContactGroupHandler) UpdateGroup(c *gin.Context) {
	groupID := c.Param("groupId")

	var req UpdateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "name is required")
		return
	}

	result := db.GetDB().Model(&model.ContactGroup{}).Where("group_id = ?", groupID).Update("name", req.Name)
	if result.RowsAffected == 0 {
		vo.NotFound(c, "group not found")
		return
	}
	vo.OK(c, nil)
}

func (h *ContactGroupHandler) DeleteGroup(c *gin.Context) {
	groupID := c.Param("groupId")

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Delete items first
		tx.Where("group_id = ?", groupID).Delete(&model.ContactGroupItem{})
		// Delete group
		result := tx.Where("group_id = ?", groupID).Delete(&model.ContactGroup{})
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			vo.NotFound(c, "group not found")
			return
		}
		vo.InternalError(c, "failed to delete group")
		return
	}
	vo.OK(c, nil)
}

func (h *ContactGroupHandler) AddItem(c *gin.Context) {
	groupID := c.Param("groupId")

	var req AddItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "task_id is required")
		return
	}

	item := model.ContactGroupItem{
		GroupID: groupID,
		TaskID:  req.TaskID,
	}
	if err := db.GetDB().Create(&item).Error; err != nil {
		vo.InternalError(c, "failed to add item to group")
		return
	}
	vo.Created(c, item)
}

func (h *ContactGroupHandler) RemoveItem(c *gin.Context) {
	groupID := c.Param("groupId")
	taskID := c.Param("taskID")

	result := db.GetDB().Where("group_id = ? AND task_id = ?", groupID, taskID).Delete(&model.ContactGroupItem{})
	if result.RowsAffected == 0 {
		vo.NotFound(c, "item not found")
		return
	}
	vo.OK(c, nil)
}
