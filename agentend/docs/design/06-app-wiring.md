# App Wiring — 应用组装与启动

## 实现了什么

FastAPI 应用入口，负责组件初始化、路由注册、CORS 配置和生命周期管理。

## 怎么实现的

### 配置管理 (`src/app/config.py`)

使用 `pydantic-settings` 管理配置，支持环境变量覆盖：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `CLAUDE_CLI_PATH` | `"claude"` | Claude CLI 可执行文件路径 |
| `DEFAULT_MAX_TURNS` | `20` | 默认最大执行轮数 |
| `EXECUTION_TIMEOUT` | `300` | 同步执行超时（秒） |
| `HOST` | `"0.0.0.0"` | 监听地址 |
| `PORT` | `8001` | 监听端口 |

### DI 容器 (`src/app/dependencies.py`)

集中创建各组件实例：

```python
def create_adapter_registry() -> AdapterRegistry:
    registry = AdapterRegistry()
    registry.register("claude-code", ClaudeCodeAdapter)
    return registry

def create_session_manager() -> SessionManager:
    return SessionManager()

def create_rule_engine() -> RuleEngine:
    rules = [SafetyRule(), ScopeRule()]
    return RuleEngine(rules)
```

### 应用入口 (`src/app/main.py`)

#### Lifespan 生命周期

```python
@asynccontextmanager
async def lifespan(app: FastAPI):
    # 启动时初始化
    app.state.adapter_registry = create_adapter_registry()
    app.state.session_manager = create_session_manager()
    app.state.rule_engine = create_rule_engine()

    yield

    # 关闭时清理所有活跃 Session
    for session in mgr.list():
        await mgr.destroy(session.id)
```

#### 路由注册

```python
app.include_router(health_router)     # GET /health
app.include_router(session_router)    # /v1/session/*
app.include_router(agent_router)      # /v1/agent/*
```

#### CORS

允许所有来源（内部服务间调用，不做认证）。

#### 启动

```bash
uv run python -m src.app.main
# 等同于
uvicorn src.app.main:app --host 0.0.0.0 --port 8001 --reload
```
