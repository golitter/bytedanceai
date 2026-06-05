package impl

import (
	"errors"

	"agenthub/backend/internal/service"
	"agenthub/backend/internal/vo"

	"github.com/gin-gonic/gin"
)

func handleBizError(c *gin.Context, err error) {
	var bizErr *service.BizError
	if errors.As(err, &bizErr) {
		switch bizErr.Code {
		case 400:
			vo.BadRequest(c, bizErr.Message)
		case 401:
			vo.Unauthorized(c, bizErr.Message)
		case 403:
			vo.Forbidden(c, bizErr.Message)
		case 404:
			vo.NotFound(c, bizErr.Message)
		case 409:
			vo.Conflict(c, bizErr.Message)
		case 503:
			vo.ServiceUnavailable(c, bizErr.Message)
		default:
			vo.InternalError(c, bizErr.Message)
		}
		return
	}

	vo.InternalError(c, err.Error())
}
