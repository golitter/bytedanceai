package impl

import (
	gormdao "agenthub/backend/internal/dao/gorm"
	"agenthub/backend/internal/service"
	svcimpl "agenthub/backend/internal/service/impl"
	"agenthub/backend/internal/vo"

	"github.com/gin-gonic/gin"
)

type SessionController struct {
	service service.SessionService
}

func NewSessionController() *SessionController {
	sessionDao := gormdao.NewSessionDao()
	sessionService := svcimpl.NewSessionService(sessionDao)

	return &SessionController{service: sessionService}
}

type PatchSessionReq struct {
	Status string `json:"status" binding:"required"`
}

func (ctrl *SessionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.PATCH("/sessions/:sessionId", ctrl.PatchSession)
}

func (ctrl *SessionController) PatchSession(c *gin.Context) {
	var req PatchSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.BadRequest(c, "status is required")
		return
	}

	result, err := ctrl.service.PatchSessionStatus(c.Param("sessionId"), req.Status)
	if err != nil {
		handleBizError(c, err)
		return
	}

	vo.OK(c, result)
}
