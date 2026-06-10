## Context

Go Backend (`agenthub/backend`) 目前有 Gin 路由骨架、GORM MySQL 连接、JWT 配置、统一响应包装器 (`vo.Response`)。模块名 `agenthub/backend`。无业务 handler、无数据模型。

AgentEnd (Python FastAPI) 已就绪，暴露：
- `POST /v1/agent/stream` — SSE 流式响应
- `GET /v1/health` — 健康检查

Contracts 已生成 Go 类型到 `backend/internal/generated/`（`AgentRequest`, `StreamEvent`, `AgentType`, `EventType`, `SessionState` 等）。

## Goals / Non-Goals

**Goals:**
- Go 层作为薄壳代理，透传 AgentEnd SSE 流
- Session / Task 的基础 CRUD 存入 MySQL
- `curl` 能完整走通：创建 session → 运行 task → 收到 SSE 事件流
- 为 Phase 2 前端提供可对接的真实 API

**Non-Goals:**
- 不做 JWT 认证（Phase 3 再加 auth 中间件）
- 不做 EventLog 持久化（直接透传，不解析 SSE 内容）
- 不做断线重连、重放、事件过滤
- 不改 AgentEnd 代码
- 不碰前端

## Decisions

### D1: SSE 纯透传（不解析 AgentEnd 返回内容）
Go 用 `bufio.Scanner` 逐行读 AgentEnd response body，逐行写回 `gin.ResponseWriter` + `Flush()`。不解析 JSON、不校验事件类型。

**替代方案**: 解析为 `StreamEvent` 再序列化转发——灵活但增加延迟和复杂度，Phase 1 不需要。

### D2: AgentEnd Client 在 main.go 注入
`main.go` 中 `agentend_client.New(host, port)` 创建 client 实例，通过闭包传入 handler。不引入 DI 框架。

**理由**: Phase 1 只有 1 个 client，闭包注入够用，简单直接。

### D3: UUID 生成 SessionID / TaskID
使用 `github.com/google/uuid` 生成 v4 UUID 作为业务 ID。数据库自增 `ID` 字段保留为主键。

### D4: GORM AutoMigrate
在 `main.go` 启动时执行 `db.AutoMigrate(&model.Session{}, &model.Task{})`，开发阶段不需要手动 SQL migration。

### D5: 类型复用 generated 包
Handler 和 Client 直接引用 `agenthub/backend/internal/generated` 中的类型，不重复定义。

### D6: 流断开处理
SSE 读取过程中如果 AgentEnd 连接断开或出错，将 `task.status` 更新为 `"failed"`，前端通过连接关闭感知错误。

## Risks / Trade-offs

- **[SSE 透传不可观测]** → Phase 1 可接受，Phase 4 可加 EventLog 持久化
- **[无认证]** → Phase 1 不暴露公网，Phase 3 加 JWT 中间件
- **[AutoMigrate 不可逆]** → 开发阶段可接受，生产前需改用版本化 migration
- **[AgentEnd 不可用时无降级]** → client.HealthCheck() 可做启动检查，但 Phase 1 不做重试
