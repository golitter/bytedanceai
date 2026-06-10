## 1. 数据模型修正

- [x] 1.1 重写 `backend/internal/model/task.go`：Task 变为顶层实体（task_id UUID, title, repo_path, status, created_at, updated_at）
- [x] 1.2 重写 `backend/internal/model/session.go`：Session 变为 Agent 会话（session_id 由调用方传入, task_id FK, agent_type, status, created_at, updated_at）

## 2. Handler 层重写

- [x] 2.1 删除 `backend/internal/handler/session.go`，创建 `backend/internal/handler/task.go`（CreateTask, ListTasks, GetTask, DeleteTask）
- [x] 2.2 重写 `backend/internal/handler/task.go` 的 RunTask：接受 session_id 参数，按需创建 Session 记录，传 session_id + task_id 给 AgentEnd
- [x] 2.3 重写 `backend/internal/handler/agent.go`：路径改为 `/api/agent-types`

## 3. 路由注册更新

- [x] 3.1 修改 `backend/cmd/server/main.go`：AutoMigrate 新模型，注册 `/api/tasks` 路由组，调整 `/api/agent-types` 路径

## 4. 验证

- [x] 4.1 `go build` 编译通过
- [x] 4.2 curl 测试 `POST /api/tasks` 创建 task
- [x] 4.3 curl 测试 `POST /api/tasks/:taskId/run` 传入 session_id，收到 SSE 流
- [x] 4.4 验证同一 task_id + session_id 二次 run 能 resume（AgentEnd 不报错）
