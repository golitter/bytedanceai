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
```

### Health Check (`src/api/v1/health.py`)

```
GET /health → {"status": "ok", "version": "0.1.0"}
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
1. Rule Engine 评估请求 → 失败返回 HTTP 400
2. 从 AdapterRegistry 获取 Adapter 实例
3. 获取或创建 Session
4. 启动 CLI 子进程，逐行解析 stdout → StreamEvent → SSE 事件推送
5. 执行完成后 Session 状态更新为 COMPLETED

#### `POST /v1/agent/execute` — 同步

请求体为 `AgentRequest`（stream=False），阻塞等待完成后返回 `AgentResponse`。

执行流程：
1. Rule Engine 评估请求 → 失败返回 HTTP 400
2. 启动 Adapter 执行，`asyncio.wait_for` 设置超时
3. 超时 → 中断进程，返回 HTTP 408
4. 成功 → 返回 AgentResponse

### 完整请求生命周期

```
Go Backend 发送 POST /v1/agent/stream
  ↓
FastAPI 路由匹配
  ↓
依赖注入获取 SessionManager / AdapterRegistry / RuleEngine
  ↓
RuleEngine.evaluate(request_context)
  ├─ 失败 → HTTP 400 {"error": "...", "rule": "..."}
  └─ 通过 ↓
AdapterRegistry.get("claude-code") → ClaudeCodeAdapter()
  ↓
SessionManager.get_or_create() → session_id
  ↓
SessionManager.update_state(RUNNING)
  ↓
ClaudeCodeAdapter.stream_chat()
  → asyncio.create_subprocess_exec("claude", "-p", ..., "--output-format", "stream-json")
  → 逐行读取 stdout → _parse_stream_line() → StreamEvent
  → SSE: event: <type>\ndata: <json>\n\n
  ↓
SessionManager.update_state(COMPLETED)
```
