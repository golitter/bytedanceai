package agentend_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"agenthub/backend/internal/generated"
)

type Client struct {
	baseURL      string
	httpClient   *http.Client
	streamClient *http.Client
}

func New(host string, port int) *Client {
	if !strings.Contains(host, "://") {
		host = "http://" + host
	}
	return &Client{
		baseURL:      fmt.Sprintf("%s:%d", host, port),
		httpClient:   &http.Client{Timeout: 60 * time.Second},
		streamClient: &http.Client{},
	}
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) StreamAgent(req *generated.AgentRequest) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	httpReq, err := http.NewRequest("POST", c.baseURL+"/v1/agent/stream", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	return c.streamClient.Do(httpReq)
}

type ReviewRequest struct {
	SessionID string `json:"session_id"`
	Action    string `json:"action"`
	Content   string `json:"content,omitempty"`
}

func (c *Client) ReviewAgent(req ReviewRequest) (map[string]interface{}, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal review request: %w", err)
	}
	httpReq, err := http.NewRequest("POST", c.baseURL+"/v1/agent/review", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create review request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("submit review: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if len(respBody) > 0 {
			return nil, fmt.Errorf("agent review failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
		}
		return nil, fmt.Errorf("agent review failed: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if len(respBody) == 0 {
		return map[string]interface{}{"status": "ok"}, nil
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode review response: %w", err)
	}
	return result, nil
}

type ValidateRepoPathResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

func (c *Client) ValidateRepoPath(repoPath string) (*ValidateRepoPathResult, error) {
	body, err := json.Marshal(map[string]string{"repo_path": repoPath})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	httpReq, err := http.NewRequest("POST", c.baseURL+"/v1/validate-repo-path", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("validate repo path: %w", err)
	}
	defer resp.Body.Close()

	var result ValidateRepoPathResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("health check: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetResources() (*http.Response, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/v1/resources")
	if err != nil {
		return nil, fmt.Errorf("get resources: %w", err)
	}
	return resp, nil
}

type AnnouncementUnpinRequest struct {
	SharedDir  string `json:"shared_dir"`
	Content    string `json:"content"`
	SenderName string `json:"sender_name"`
}

func (c *Client) NotifyAnnouncementUnpin(req AnnouncementUnpinRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal unpin request: %w", err)
	}
	httpReq, err := http.NewRequest("POST", c.baseURL+"/v1/pin/announcement-unpin", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create unpin request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("notify announcement unpin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("announcement unpin failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
