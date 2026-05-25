# Session 级工作区代理路由

## 变更原因

前端卡片组件（ImageCard、AttachmentCard、DiffCard）通过 workspaceId 访问文件，但前端只有 sessionId。之前硬编码 `_placeholder` 导致 404。需要新增 session → workspace 解析链路，使前端通过 sessionId 即可访问工作区文件和 diff。

## 变更文件

无 schema 文件修改。变更仅涉及内部路由：
- `agentend/src/api/v1/workspace.py` — 新增 `GET /v1/workspace/by-session/{session_id}`
- `backend/internal/handler/workspace.go` — 新增 `resolveWorkspaceID` + Session 级代理 handler
- `backend/cmd/server/main.go` — 注册 `/api/session/:sessionId/files/*` 路由组
- `frontend/src/components/cards/` — 三卡片改用 `/api/session/{sessionId}/files/*` URL
- `frontend/src/components/chat/` — sessionId 从 ChatArea 透传至 BlockRenderer

## 对比结果

无契约变更。所有 SSE 事件、请求/响应格式保持不变。

## 跨端影响

- **AgentEnd**：新增 `by-session` 查询端点，返回 `{ workspace_id: string }`，不修改现有 schema
- **Backend**：新增 `/api/session/{sessionId}/files/*` 和 `/api/session/{sessionId}/diff` 代理路由，内部解析 sessionId → workspaceId 后转发至 AgentEnd，不修改现有 handler
- **Frontend**：卡片组件 URL 从 `/api/workspace/{workspaceId}/files/*` 改为 `/api/session/{sessionId}/files/*`，sessionId 通过组件 props 透传，无需额外 API 调用
