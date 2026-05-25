# Phase 1: Go 胶水层实现

## 概述

Go Backend 作为薄壳代理，透传 AgentEnd SSE 流，提供 Task/Session CRUD API。
curl 端到端验证通过：创建 task → 运行 agent session → 收到 SSE 事件流。

## 架构

```
curl / 前端
    │
    ▼
Go Backend (Gin :8080)
    │
    ├── /api/tasks              → Task CRUD (MySQL)
    ├── /api/tasks/:taskId/run  → 按需创建 Session → SSE 透传
    ├── /api/agent-types        → 硬编码 Agent 类型列表
    │
    └── AgentEnd Client ──► AgentEnd (FastAPI :8001)
                              POST /v1/agent/stream (SSE)
                              GET  /health
```

## 文件清单

| 文件 | 说明 |
|------|------|
| `internal/conf/conf.go` | 新增 `AgentEndConfig` 结构体 |
| `internal/model/task.go` | Task GORM 模型（群组任务，顶层实体） |
| `internal/model/session.go` | Session GORM 模型（Agent 会话，从属 Task） |
| `internal/handler/task.go` | Task CRUD + Session 运行 + SSE 透传 handler |
| `internal/handler/agent.go` | Agent 类型列表 handler |
| `pkg/agentend_client/client.go` | AgentEnd HTTP 客户端 |
| `cmd/server/main.go` | 路由注册、AutoMigrate、DI |
| `configs/config.yaml` | agentend 配置项 |

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/tasks` | 创建 Task（需 `title`，可选 `repo_path`，后端生成 task_id） |
| GET | `/api/tasks` | 列出所有 Task |
| GET | `/api/tasks/:taskId` | 获取 Task 详情（含关联 Session 列表） |
| DELETE | `/api/tasks/:taskId` | 删除 Task |
| POST | `/api/tasks/:taskId/run` | 运行 Agent Session（SSE 流式返回，后端生成 session_id） |
| GET | `/api/tasks/:taskId/stream` | SSE 流式读取（Redis Stream → MySQL 批量刷写） |
| GET | `/api/tasks/:taskId/messages` | 获取 Task 的历史消息 |
| GET | `/api/agent-types` | 获取 Agent 类型列表 |
| PATCH | `/api/sessions/:sessionId` | 更新 Session 状态（如停用为 inactive） |
| PUT | `/api/sessions/:sessionId` | 更新 Session 元数据（agent_name、avatar_url） |
| POST | `/api/agents/avatar` | 上传 Agent 头像（七牛云存储） |
| POST | `/api/validate-repo-path` | 验证仓库路径合法性 |

## 数据模型

> **Phase 1 为单聊模式**：Task : Session : Agent = 1 : 1 : 1。
> 一个 Task 只绑定一个 Agent、创建一个 Session，不支持多 Agent 协作。

### Task（任务 — 顶层实体）
- `id` (uint, PK)
- `task_id` (string, UUID, 唯一索引) — AgentEnd 用此决定 git branch 和 worktree
- `title` (string)
- `repo_path` (string, 仓库路径，运行时注入 AgentRequest)
- `status` (string, default: "active")
- `created_at`, `updated_at`

### Session（Agent 会话 — 从属 Task，1:1）
- `id` (uint, PK)
- `session_id` (string, 后端生成 UUID, 唯一索引) — AgentEnd 用 `session_id::task_id` 映射 cli_session_id
- `task_id` (string, 索引, FK → Task)
- `agent_type` (string)
- `agent_name` (string, Agent 显示名称)
- `avatar_url` (string, Agent 头像 URL)
- `status` (string: active → running/completed/failed/inactive)
- `created_at`, `updated_at`

### 单聊流程

```
1. POST /api/tasks          → 创建 Task（传入 title、repo_path，后端生成 task_id）
2. POST /api/tasks/:taskId/run → 运行 Session（传入 message、agent_type）
   - 后端自动生成 session_id
   - 首次 run → 自动创建 Session 记录
   - 后续 run → 复用同一 Session（相同 session_id 即 resume）
3. SSE 流通过 Redis Stream 中转，消息批量刷写到 MySQL
```

## SSE 流式机制

RunTask handler：
1. 验证 Task 存在
2. 生成 session_id，创建 Session 记录
3. 调用 `agentend_client.StreamAgent()` 获取 AgentEnd 的 SSE 响应
4. 将 SSE 事件写入 Redis Stream（`stream:task:{taskId}`）
5. 后台 goroutine 从 Redis Stream 消费，批量刷写到 MySQL Message 表
6. 前端通过 `GET /api/tasks/:taskId/stream` 订阅 Redis Stream 接收实时事件
7. 流结束/中断时更新 `session.status`

- 流正常结束 → `session.status = "completed"`
- 流中断/出错 → `session.status = "failed"`

## 头像上传

Agent 头像通过七牛云存储：
1. `POST /api/agents/avatar` 接收 multipart 文件上传
2. 上传至七牛云，返回 CDN URL
3. 通过 `PUT /api/sessions/:sessionId` 更新 Session 的 `avatar_url` 和 `agent_name`

## 设计决策

- **Task 为顶层，Session 从属**: 与 AgentEnd 的 `task_id`（workspace 隔离）和 `session_id`（CLI session 映射）语义对齐
- **单聊模式（1:1）**: Phase 1 中一个 Task 只有一个 Session、一个 Agent，不支持多 Agent 协作
- **session_id 由后端生成**: 后端使用 `google/uuid` 生成 session_id，简化调用方逻辑
- **Session 按需创建**: 调用 `/api/tasks/:taskId/run` 时自动创建，无需单独端点
- **Redis Stream 中转**: SSE 事件先入 Redis Stream，再批量刷写 MySQL，支持断线重连和历史回放
- **七牛云头像**: 头像文件上传至七牛云 CDN，Session 记录保存 URL
- **闭包注入**: AgentEnd Client 在 main.go 创建，通过构造函数传入 handler
- **UUID 业务 ID**: `task_id` 和 `session_id` 均由后端使用 `google/uuid` v4 生成
- **AutoMigrate**: 开发阶段自动建表，生产前需改为版本化 migration
- **复用 generated 包**: 直接引用 `internal/generated` 中的类型

## 配置

```yaml
# backend/configs/config.yaml
agentend:
  host: http://localhost
  port: 8001
```

## 端到端验证示例

```bash
# 创建 task（群组任务，绑定仓库路径）
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"my-project", "repo_path": "/path/to/repo"}'

# 运行 agent session（SSE 流，后端自动生成 session_id）
curl -N -X POST http://localhost:8080/api/tasks/{taskId}/run \
  -H "Content-Type: application/json" \
  -d '{"message":"hello","agent_type":"claude-code"}'

# 获取 agent 类型列表
curl http://localhost:8080/api/agent-types
```

## Non-Goals (Phase 1 不做)

- 断线重连、重放、事件过滤
- 修改 AgentEnd 代码
- 前端改动
