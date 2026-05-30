package handler

import (
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/db"

	"github.com/gin-gonic/gin"
)

const adminAvatarKey = "admin_avatar_url"

type AvatarRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *AdminHandler) GetAvatar(c *gin.Context) {
	var setting model.AdminSetting
	if err := db.GetDB().Where("`key` = ?", adminAvatarKey).First(&setting).Error; err != nil {
		vo.OK(c, gin.H{"url": "https://api.dicebear.com/9.x/notionists/svg?seed=tln&backgroundColor=c0aede"})
		return
	}
	if setting.Value == "" {
		vo.OK(c, gin.H{"url": "https://api.dicebear.com/9.x/notionists/svg?seed=tln&backgroundColor=c0aede"})
		return
	}
	vo.OK(c, gin.H{"url": setting.Value})
}

func (h *AdminHandler) UpdateAvatar(c *gin.Context) {
	var req AvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "url is required")
		return
	}

	db.GetDB().Where("`key` = ?", adminAvatarKey).Delete(&model.AdminSetting{})
	db.GetDB().Create(&model.AdminSetting{Key: adminAvatarKey, Value: req.URL})

	vo.OK(c, gin.H{"success": true})
}
