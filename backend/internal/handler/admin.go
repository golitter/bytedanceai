package handler

import (
	"agenthub/backend/internal/conf"
	"agenthub/backend/internal/middleware"
	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/agentend_client"
	"agenthub/backend/pkg/qiniu"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	cfg         *conf.Config
	uploader    *qiniu.Uploader
	agentClient *agentend_client.Client
}

func NewAdminHandler(cfg *conf.Config, uploader *qiniu.Uploader, agentClient *agentend_client.Client) *AdminHandler {
	return &AdminHandler{cfg: cfg, uploader: uploader, agentClient: agentClient}
}

func (h *AdminHandler) RegisterRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		authLimiter := middleware.NewIPRateLimiter(5, time.Minute)
		admin.POST("/auth", authLimiter.Middleware(), h.Auth)
		admin.GET("/health", h.HealthCheck)
		admin.GET("/avatar", h.GetAvatar)

		protected := admin.Group("")
		protected.Use(middleware.AdminAuth(h.cfg.JWT.Secret))
		{
			protected.GET("/resources", h.GetResources)
			protected.DELETE("/sessions", h.DeleteSessions)
			protected.GET("/workspaces", h.GetWorkspaces)
			protected.DELETE("/workspaces/:id", h.DeleteWorkspace)
			protected.GET("/agents", h.GetAgents)
			protected.GET("/services", h.GetServices)
			protected.GET("/statistics", h.GetStatistics)
			protected.PUT("/avatar", h.UpdateAvatar)
		}
	}
}

type AuthRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) Auth(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "password is required")
		return
	}

	if !middleware.VerifyAdminPassword(req.Password, h.cfg.Admin.Password) {
		vo.Unauthorized(c, "密码错误")
		return
	}

	token, err := middleware.GenerateAdminToken(h.cfg.JWT.Secret)
	if err != nil {
		vo.InternalError(c, "failed to generate token")
		return
	}

	vo.OK(c, gin.H{
		"token":      token,
		"expires_in": 3600,
	})
}

func (h *AdminHandler) HealthCheck(c *gin.Context) {
	vo.OK(c, gin.H{"status": "ok"})
}
