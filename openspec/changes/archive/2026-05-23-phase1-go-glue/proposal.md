## Why

Go Backend 目前只有骨架（Gin + GORM + JWT），无业务代码。前端需要一个真实的 API 来对接开发。Phase 1 的目标是让 curl 能走通 Go → AgentEnd SSE 流，为 Phase 2 前端聊天界面提供真实后端。

## What Changes

- 新增 AgentEnd HTTP Client，封装对 AgentEnd `/v1/agent/stream` 的 SSE 调用
- 新增 Session / Task GORM 数据模型，支持 MySQL 自动迁移
- 新增 RESTful API：Session CRUD + Task 运行（SSE 流式响应）+ Agent 类型列表
- 修改 `main.go` 注册路由、注入依赖、执行 DB migration
- 修改 `config.yaml` 增加 AgentEnd 地址配置

## Capabilities

### New Capabilities
- `agentend-client`: 封装对 AgentEnd 的 HTTP 调用（SSE 流 + 健康检查）
- `session-api`: Session 的 CRUD REST API
- `task-run-api`: Task 创建与 SSE 流式运行 API
- `agent-list-api`: 返回可用 Agent 类型列表

### Modified Capabilities
（无已有 spec 需要修改）

## Impact

- `backend/` 目录新增 ~355 行代码，修改 ~50 行
- API 端点：`/api/sessions`、`/api/sessions/:sid/tasks/run`、`/api/agents`
- 依赖：无新外部依赖（仅使用标准库 `net/http` + `bufio`）
- 数据库：新增 `sessions` 和 `tasks` 两张表
- 不涉及前端和 AgentEnd 的改动
