# API Endpoints — HTTP 端点

## 实现了什么

FastAPI HTTP 端点，作为 Go Backend 调用 Runtime 的入口。

## 怎么实现的

### 依赖注入 (`src/api/dependencies.py`)

通过 FastAPI 的 `Request.app.state` 获取组件实例：

```python
def get_session_manager(request) -> SessionManager
def get_adapter_registry(request) -> AdapterRegistry
def get_rule_engine(request) -> RuleEngine
def get_session_store(request) -> SessionMappingStore
def get_workspace_manager(request) -> WorkspaceManager
```

### Health Check (`src/api/v1/health.py`)

```
GET /health → {"status": "ok", "version": "<config.yaml app.version>"}
```

### Session CRUD (`src/api/v1/session.py`)

| 端点 | 方法 | 说明 |
|------|------|------|
| `GET /v1/session` | GET | 列出所有会话 |
| `GET /v1/session/{id}` | GET | 获取会话详情（不存在返回 404） |
| `POST /v1/session/{id}/interrupt` | POST | 中断运行中的会话（非运行中返回提示） |
| `DELETE /v1/session/{id}` | DELETE | 销毁会话并清理进程（不存在返回 404） |

### Agent 执行 (`src/api/v1/agent.py`)

#### `POST /v1/agent/stream` — SSE 流式

请求体为 `AgentRequest`（stream=True），返回 SSE 流。

执行流程：
1. `_resolve_workspace()` — 有 workspace_path 直接使用，有 repo_path 自动创建 Git worktree
2. Rule Engine 评估请求 → 失败返回 HTTP 400
3. 从 AdapterRegistry 获取 Adapter 实例
4. `_resolve_session()` — 获取或创建 Session，查询 SessionMappingStore 获取 CLI session 映射
5. 启动 CLI 子进程，逐行解析 stdout → StreamEvent → SSE 事件推送
6. INIT 事件触发 CLI session_id 回写到 SessionMappingStore
7. 执行完成后 Session 状态更新为 COMPLETED

#### `POST /v1/agent/execute` — 同步

请求体为 `AgentRequest`（stream=False），阻塞等待完成后返回 `AgentResponse`。

执行流程：
1. `_resolve_workspace()` — 同 stream 路径
2. Rule Engine 评估请求 → 失败返回 HTTP 400
3. `_resolve_session()` — 同 stream 路径
4. 内联 `_collect()` 流式收集文本 + INIT 事件回写 mapping
5. `asyncio.wait_for` 设置超时（来自 `config.yaml` 的 `execution.timeout`）
6. 超时 → 中断进程，返回 HTTP 408
7. 成功 → 返回 AgentResponse

### Workspace 管理 (`src/api/v1/workspace.py`)

提供工作区 CRUD、文件操作、diff、commit、merge、preview 等端点：

| 端点 | 方法 | 说明 |
|------|------|------|
| `/v1/workspace/create` | POST | 创建 workspace（含 worktree + 技能分发） |
| `/v1/workspace/{id}/files/{path}` | GET | 读取文件（防路径穿越） |
| `/v1/workspace/{id}/files/{path}` | PUT | 写入文件（防路径穿越） |
| `/v1/workspace/{id}/diff` | GET | 获取 `git diff HEAD` |
| `/v1/workspace/{id}/commit` | POST | 提交变更（`git add -A && git commit`） |
| `/v1/workspace/{id}/revert` | POST | 撤销变更（`git checkout HEAD -- .`） |
| `/v1/workspace/{id}/merge` | POST | 合并分支 |
| `/v1/workspace/{id}/preview/start` | POST | 启动预览服务器（aiohttp 静态文件） |
| `/v1/workspace/{id}/preview/stop` | POST | 停止预览服务器 |
| `/v1/workspace/task/{task_id}/merge-to-main` | POST | 合并 task branch 到 main |
| `/v1/workspace/{id}` | DELETE | 清理 workspace（worktree + branch） |
| `/v1/workspace` | GET | 列出所有 workspace |
| `/v1/workspace/by-session/{session_id}` | GET | 按 session_id 查找活跃 workspace |

### Pin 管理 (`src/api/v1/pin.py`)

Pin Memory 上下文管理端点，允许用户将关键约束"钉住"供 Orchestrator 使用。

### Validate Repo Path (`src/api/v1/validate.py`)

校验仓库路径是否有效，供前端新建对话时使用。

### 完整请求生命周期

```
Go Backend 发送 POST /v1/agent/stream
  ↓
FastAPI 路由匹配
  ↓
依赖注入获取 SessionManager / AdapterRegistry / RuleEngine / SessionMappingStore / WorkspaceManager
  ↓
_resolve_workspace() → workspace_path（有 repo_path 则自动创建 worktree）
  ↓
RuleEngine.evaluate(request_context)
  ├─ 失败 → HTTP 400 {"error": "...", "rule": "..."}
  └─ 通过 ↓
AdapterRegistry.get(AgentType) → Adapter()
  ↓
_resolve_session() → (internal_session_id, cli_session_id, is_resume)
  ↓
SessionManager.update_state(RUNNING)
  ↓
Adapter.stream_chat(**stream_kwargs)
  → asyncio.create_subprocess_exec("claude", "-p", ..., "--output-format", "stream-json", "--verbose", "--include-partial-messages")
  → 逐行读取 stdout → _parse_stream_line() → StreamEvent
  → INIT 事件: session_store.set_cli_session_id() 回写 mapping
  → SSE: event: <type>\ndata: <json>\n\n
  ↓
SessionManager.update_state(COMPLETED)
```
