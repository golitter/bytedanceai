## 1. 配置与依赖

- [x] 1.1 修改 `backend/internal/conf/conf.go` 添加 `AgentEndConfig` 结构体（Host + Port），加入 `Config.AgentEnd` 字段
- [x] 1.2 修改 `backend/configs/config.yaml` 添加 agentend 配置项（host: http://localhost, port: 8000）
- [x] 1.3 添加 `github.com/google/uuid` 依赖到 go.mod

## 2. GORM 数据模型

- [x] 2.1 创建 `backend/internal/model/session.go`（Session struct：ID, SessionID, Title, Status, CreatedAt, UpdatedAt）
- [x] 2.2 创建 `backend/internal/model/task.go`（Task struct：ID, TaskID, SessionID, AgentType, Status, Message, Result, CreatedAt, UpdatedAt）

## 3. AgentEnd HTTP Client

- [x] 3.1 创建 `backend/pkg/agentend_client/client.go`（New 构造函数 + StreamAgent 方法 + HealthCheck 方法）
- [x] 3.2 StreamAgent 向 AgentEnd POST /v1/agent/stream，返回 *http.Response 供逐行读取
- [x] 3.3 HealthCheck 向 AgentEnd GET /v1/health，返回 error

## 4. Handler 层

- [x] 4.1 创建 `backend/internal/handler/session.go`（CreateSession, ListSessions, GetSession, DeleteSession）
- [x] 4.2 创建 `backend/internal/handler/task.go`（RunTask：解析请求 → 写入 DB → 调 StreamAgent → SSE 透传 → 更新状态）
- [x] 4.3 创建 `backend/internal/handler/agent.go`（ListAgents：硬编码返回 agent 类型列表）

## 5. 路由注册与启动

- [x] 5.1 修改 `backend/cmd/server/main.go`：注册 API 路由组、AutoMigrate、初始化 AgentEnd Client 并注入 handler
- [x] 5.2 验证 `make run-backend` 启动成功，`/ping` 正常

## 6. 端到端验证

- [x] 6.1 启动 AgentEnd，用 curl 测试 `POST /api/sessions` 创建 session
- [x] 6.2 用 curl 测试 `GET /api/sessions/{sid}/tasks/run` 收到 SSE 流
- [x] 6.3 验证 SSE 流结束后 task.status = completed，AgentEnd 断开时 task.status = failed
