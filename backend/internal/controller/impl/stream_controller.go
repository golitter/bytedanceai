package impl

import (
	"net/http"

	gormdao "agenthub/backend/internal/dao/gorm"
	"agenthub/backend/internal/service"
	svcimpl "agenthub/backend/internal/service/impl"

	"github.com/gin-gonic/gin"
)

type StreamController struct {
	service service.StreamService
}

func NewStreamController() *StreamController {
	messageDao := gormdao.NewMessageDao()
	streamService := svcimpl.NewStreamService(messageDao)
	return &StreamController{service: streamService}
}

func (ctrl *StreamController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tasks/:taskId/stream", ctrl.ServeStream)
}

func (ctrl *StreamController) ServeStream(c *gin.Context) {
	if c.Query("message_id") == "" || c.Query("session_id") == "" {
		handleBizError(c, service.ErrBadRequest("session_id and message_id are required"))
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Flush()

	if err := ctrl.service.ServeStream(c.Request.Context(), c.Query("session_id"), c.Query("message_id"), c.Writer, c.Writer); err != nil {
		handleBizError(c, err)
	}
}
