## 1. Project Setup (uv)

- [x] 1.1 Run `cd agentend && uv init` 初始化项目，生成 pyproject.toml 和 .venv
- [x] 1.2 Run `uv add fastapi uvicorn pydantic pydantic-settings sse-starlette httpx` 添加运行时依赖
- [x] 1.3 Run `uv add --dev pytest pytest-asyncio` 添加开发依赖
- [x] 1.4 Create `agentend/src/__init__.py` and top-level package structure
- [x] 1.5 Create `agentend/src/app/config.py` — pydantic-settings 配置类，支持环境变量覆盖（CLAUDE_CLI_PATH, DEFAULT_MAX_TURNS, EXECUTION_TIMEOUT, HOST, PORT）
- [x] 1.6 Create `agentend/src/app/main.py` — FastAPI 应用入口，lifespan 中初始化 AdapterRegistry / SessionManager / RuleEngine
- [x] 1.7 Create `agentend/.gitignore` 忽略 .venv/, __pycache__/, *.pyc, uv.lock 等

## 2. Schemas (Data Models)

- [x] 2.1 Create `agentend/src/schemas/__init__.py` with public exports
- [x] 2.2 Create `agentend/src/schemas/request.py` — `AgentRequest` Pydantic model (task_id, conversation_id, session_id?, message, agent_type, stream, system_prompt?, rules, workspace_path?, config?)
- [x] 2.3 Create `agentend/src/schemas/response.py` — `AgentResponse` Pydantic model (session_id, content, artifacts, usage)
- [x] 2.4 Create `agentend/src/schemas/events.py` — `StreamEvent` Pydantic model (type, content, timestamp) + EventType enum (text, tool_call, tool_result, artifact, done, error)

## 3. Adapter Layer

- [x] 3.1 Create `agentend/src/adapters/__init__.py` with public exports (`BaseAgentAdapter`, `AdapterRegistry`, `ClaudeCodeAdapter`)
- [x] 3.2 Create `agentend/src/adapters/base.py` — `BaseAgentAdapter` ABC with abstract methods: create_session, chat, stream_chat, interrupt, destroy_session
- [x] 3.3 Create `agentend/src/adapters/registry.py` — `AdapterRegistry` with register/get/list methods, startup auto-registration
- [x] 3.4 Create `agentend/src/adapters/claude.py` — `ClaudeCodeAdapter` implementing `_build_command` (组装 CLI 参数), `_parse_stream_line` (解析 stream-json), `stream_chat` (AsyncIterator[StreamEvent]), `chat` (同步), `interrupt` (SIGTERM→SIGKILL), `destroy_session`

## 4. Session Manager

- [x] 4.1 Create `agentend/src/session/__init__.py` with public exports
- [x] 4.2 Create `agentend/src/session/models.py` — `Session` dataclass + `SessionState` enum (IDLE, RUNNING, COMPLETED, INTERRUPTED, ERROR)
- [x] 4.3 Create `agentend/src/session/manager.py` — `SessionManager` with create/get/list/update_state/destroy methods, state transition validation, history tracking

## 5. Rule Engine

- [x] 5.1 Create `agentend/src/rules/__init__.py` with public exports
- [x] 5.2 Create `agentend/src/rules/base.py` — `BaseRule` ABC with name, description, phase, priority, check(), enforce()
- [x] 5.3 Create `agentend/src/rules/registry.py` — `RuleRegistry` with register/list methods
- [x] 5.4 Create `agentend/src/rules/engine.py` — `RuleEngine` with evaluate() method (按 priority 降序执行 check→enforce，合并约束结果)
- [x] 5.5 Create `agentend/src/rules/builtin.py` — 内置 `SafetyRule` (priority=10, 注入安全提示 + 限制危险工具) 和 `ScopeRule` (priority=5, 验证 workspace_path)

## 6. API Endpoints

- [x] 6.1 Create `agentend/src/api/__init__.py` with router exports
- [x] 6.2 Create `agentend/src/api/v1/health.py` — `GET /health` 返回 {"status": "ok", "version": "0.1.0"}
- [x] 6.3 Create `agentend/src/api/v1/session.py` — `GET /v1/session`, `GET /v1/session/{id}`, `POST /v1/session/{id}/interrupt`, `DELETE /v1/session/{id}`
- [x] 6.4 Create `agentend/src/api/v1/agent.py` — `POST /v1/agent/stream` (SSE 流式) 和 `POST /v1/agent/execute` (同步)，集成 Rule Engine 评估 + Adapter 执行
- [x] 6.5 Create `agentend/src/api/dependencies.py` — FastAPI 依赖注入函数，从 app state 获取 SessionManager / AdapterRegistry / RuleEngine

## 7. App Wiring & Entry Point

- [x] 7.1 Wire all routers into FastAPI app in `main.py`，配置 CORS、lifespan
- [x] 7.2 Create `agentend/src/app/dependencies.py` — DI 容器，统一创建和提供各组件实例
- [x] 7.3 Add `if __name__ == "__main__"` entry point in main.py 启动 uvicorn

## 8. Tests

- [x] 8.1 Create `agentend/tests/test_schemas.py` — 测试 AgentRequest/AgentResponse/StreamEvent 序列化和校验
- [x] 8.2 Create `agentend/tests/test_session.py` — 测试 SessionManager CRUD、状态转移、history
- [x] 8.3 Create `agentend/tests/test_rules.py` — 测试 RuleEngine evaluate、priority 排序、check 失败
- [x] 8.4 Create `agentend/tests/test_adapter.py` — 测试 ClaudeCodeAdapter._build_command 参数组装
- [x] 8.5 Create `agentend/tests/test_api.py` — 测试 health endpoint、session CRUD endpoints、agent execute endpoint (mock adapter)
