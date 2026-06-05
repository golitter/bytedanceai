package impl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"agenthub/backend/pkg/agentend_client"

	"github.com/gin-gonic/gin"
)

type WorkspaceController struct {
	agentClient *agentend_client.Client
	httpClient  *http.Client
}

func NewWorkspaceController(agentClient *agentend_client.Client) *WorkspaceController {
	return &WorkspaceController{
		agentClient: agentClient,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (ctrl *WorkspaceController) RegisterRoutes(rg *gin.RouterGroup) {
	ws := rg.Group("/workspace")
	{
		ws.GET("/task/:taskId/git-info", ctrl.TaskGitInfo)
		ws.POST("/task/:taskId/merge-to-main", ctrl.MergeTaskToMain)
		ws.GET("/:id/files/*filepath", ctrl.ReadFile)
		ws.PUT("/:id/files/*filepath", ctrl.WriteFile)
		ws.GET("/:id/diff", ctrl.GetDiff)
		ws.POST("/:id/commit", ctrl.Commit)
		ws.POST("/:id/revert", ctrl.Revert)
		ws.POST("/:id/preview/start", ctrl.StartPreview)
		ws.POST("/:id/preview/stop", ctrl.StopPreview)
	}

	ss := rg.Group("/session")
	{
		ss.GET("/:sessionId/files/*filepath", ctrl.SessionFileRead)
		ss.PUT("/:sessionId/files/*filepath", ctrl.SessionFileWrite)
		ss.GET("/:sessionId/diff", ctrl.SessionGetDiff)
		ss.POST("/:sessionId/commit", ctrl.SessionCommit)
		ss.POST("/:sessionId/revert", ctrl.SessionRevert)
	}
}

func sanitizePath(p string) (string, bool) {
	cleaned := path.Clean(p)
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", false
	}
	for _, seg := range strings.Split(cleaned, "/") {
		if seg == ".." {
			return "", false
		}
	}
	return cleaned, true
}

func (ctrl *WorkspaceController) resolveWorkspaceID(sessionID string) (string, error) {
	url := fmt.Sprintf("%s/v1/workspace/by-session/%s", ctrl.agentClient.BaseURL(), sessionID)
	resp, err := ctrl.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("agentend unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("workspace not found for session %s", sessionID)
	}

	var ws struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil {
		return "", fmt.Errorf("invalid workspace response")
	}
	return ws.ID, nil
}

func (ctrl *WorkspaceController) withResolvedWorkspace(c *gin.Context, fn func(wsID string)) {
	wsID, err := ctrl.resolveWorkspaceID(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fn(wsID)
}

func (ctrl *WorkspaceController) SessionFileRead(c *gin.Context) {
	filePath, ok := sanitizePath(c.Param("filepath"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file path"})
		return
	}
	ctrl.withResolvedWorkspace(c, func(wsID string) {
		ctrl.proxy(c, "GET", fmt.Sprintf("/v1/workspace/%s/files%s", wsID, filePath), nil)
	})
}

func (ctrl *WorkspaceController) SessionFileWrite(c *gin.Context) {
	filePath, ok := sanitizePath(c.Param("filepath"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file path"})
		return
	}
	ctrl.withResolvedWorkspace(c, func(wsID string) {
		ctrl.proxy(c, "PUT", fmt.Sprintf("/v1/workspace/%s/files%s", wsID, filePath), c.Request.Body)
	})
}

func (ctrl *WorkspaceController) SessionGetDiff(c *gin.Context) {
	ctrl.withResolvedWorkspace(c, func(wsID string) {
		ctrl.proxy(c, "GET", fmt.Sprintf("/v1/workspace/%s/diff", wsID), nil)
	})
}

func (ctrl *WorkspaceController) SessionCommit(c *gin.Context) {
	ctrl.withResolvedWorkspace(c, func(wsID string) {
		ctrl.proxy(c, "POST", fmt.Sprintf("/v1/workspace/%s/commit", wsID), c.Request.Body)
	})
}

func (ctrl *WorkspaceController) SessionRevert(c *gin.Context) {
	ctrl.withResolvedWorkspace(c, func(wsID string) {
		ctrl.proxy(c, "POST", fmt.Sprintf("/v1/workspace/%s/revert", wsID), nil)
	})
}

func (ctrl *WorkspaceController) ReadFile(c *gin.Context) {
	filePath, ok := sanitizePath(c.Param("filepath"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file path"})
		return
	}
	ctrl.proxy(c, "GET", fmt.Sprintf("/v1/workspace/%s/files%s", c.Param("id"), filePath), nil)
}

func (ctrl *WorkspaceController) WriteFile(c *gin.Context) {
	filePath, ok := sanitizePath(c.Param("filepath"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file path"})
		return
	}
	ctrl.proxy(c, "PUT", fmt.Sprintf("/v1/workspace/%s/files%s", c.Param("id"), filePath), c.Request.Body)
}

func (ctrl *WorkspaceController) GetDiff(c *gin.Context) {
	ctrl.proxy(c, "GET", fmt.Sprintf("/v1/workspace/%s/diff", c.Param("id")), nil)
}

func (ctrl *WorkspaceController) Commit(c *gin.Context) {
	ctrl.proxy(c, "POST", fmt.Sprintf("/v1/workspace/%s/commit", c.Param("id")), c.Request.Body)
}

func (ctrl *WorkspaceController) Revert(c *gin.Context) {
	ctrl.proxy(c, "POST", fmt.Sprintf("/v1/workspace/%s/revert", c.Param("id")), nil)
}

func (ctrl *WorkspaceController) TaskGitInfo(c *gin.Context) {
	ctrl.proxy(c, "GET", fmt.Sprintf("/v1/workspace/task/%s/git-info", c.Param("taskId")), nil)
}

func (ctrl *WorkspaceController) MergeTaskToMain(c *gin.Context) {
	ctrl.proxy(c, "POST", fmt.Sprintf("/v1/workspace/task/%s/merge-to-main", c.Param("taskId")), c.Request.Body)
}

func (ctrl *WorkspaceController) StartPreview(c *gin.Context) {
	ctrl.proxy(c, "POST", fmt.Sprintf("/v1/workspace/%s/preview/start", c.Param("id")), nil)
}

func (ctrl *WorkspaceController) StopPreview(c *gin.Context) {
	ctrl.proxy(c, "POST", fmt.Sprintf("/v1/workspace/%s/preview/stop", c.Param("id")), nil)
}

func (ctrl *WorkspaceController) proxy(c *gin.Context, method, path string, body io.Reader) {
	url := ctrl.agentClient.BaseURL() + path

	req, err := http.NewRequestWithContext(c.Request.Context(), method, url, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if body != nil && c.ContentType() != "" {
		req.Header.Set("Content-Type", c.ContentType())
	}

	resp, err := ctrl.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "agentend unavailable"})
		return
	}
	defer resp.Body.Close()

	for k, vs := range resp.Header {
		for _, v := range vs {
			c.Writer.Header().Add(k, v)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(c.Writer, resp.Body)
}
