## ADDED Requirements

### Requirement: FastAPI application with lifespan
系统 SHALL 创建 FastAPI 应用实例，在 lifespan 中初始化 `AdapterRegistry`、`SessionManager`、`RuleEngine`，并通过依赖注入提供给端点使用。

#### Scenario: App startup initializes components
- **WHEN** FastAPI 应用启动
- **THEN** SHALL 初始化 `AdapterRegistry`（注册 ClaudeCodeAdapter）、`SessionManager`、`RuleEngine`（加载内置 Rule），并存储在 app state 中

#### Scenario: App shutdown cleans up
- **WHEN** FastAPI 应用关闭
- **THEN** SHALL 遍历所有活跃 Session，调用 `interrupt` 和 `destroy` 清理进程资源

### Requirement: POST /v1/agent/stream endpoint
系统 SHALL 提供 `POST /v1/agent/stream` 端点，接收 `AgentRequest`，通过 SSE 流式返回 `StreamEvent` 序列。MUST 先经过 Rule Engine 评估，通过后启动 Adapter 执行。

#### Scenario: Successful stream execution
- **WHEN** 发送有效的 `AgentRequest`（stream=True）
- **THEN** SHALL 返回 `Content-Type: text/event-stream` 响应，逐个推送 StreamEvent SSE 事件，最终推送 `type="done"` 事件

#### Scenario: Rule evaluation blocks execution
- **WHEN** Rule Engine 评估失败
- **THEN** SHALL 返回 HTTP 400 错误，body 包含失败 Rule 名称和错误原因

#### Scenario: Invalid request body
- **WHEN** 发送缺少必填字段的请求
- **THEN** SHALL 返回 HTTP 422 错误，body 包含 Pydantic 验证错误详情

### Requirement: POST /v1/agent/execute endpoint
系统 SHALL 提供 `POST /v1/agent/execute` 端点，接收 `AgentRequest`（stream=False），同步等待 Adapter 执行完成，返回 `AgentResponse`。

#### Scenario: Successful sync execution
- **WHEN** 发送有效的 `AgentRequest`（stream=False）
- **THEN** SHALL 等待 CLI 进程完成，返回包含完整结果的 `AgentResponse` JSON

#### Scenario: Execution timeout
- **WHEN** 同步执行超过配置的超时时间（默认 300 秒）
- **THEN** SHALL 终止进程，返回 HTTP 408 错误，包含 "execution timeout" 信息

### Requirement: GET /v1/session endpoint
系统 SHALL 提供 `GET /v1/session` 端点，返回所有活跃 Session 的列表。

#### Scenario: List sessions
- **WHEN** 调用 `GET /v1/session`
- **THEN** SHALL 返回 JSON 数组，每个元素包含 session_id、agent_type、state、created_at、last_active

### Requirement: GET /v1/session/{id} endpoint
系统 SHALL 提供 `GET /v1/session/{id}` 端点，返回指定 Session 的详细信息。

#### Scenario: Get existing session
- **WHEN** 调用 `GET /v1/session/exist-id`
- **THEN** SHALL 返回该 Session 的完整信息 JSON

#### Scenario: Get non-existent session
- **WHEN** 调用 `GET /v1/session/non-existent`
- **THEN** SHALL 返回 HTTP 404 错误

### Requirement: POST /v1/session/{id}/interrupt endpoint
系统 SHALL 提供 `POST /v1/session/{id}/interrupt` 端点，中断指定 Session 的正在执行的 Agent 进程。

#### Scenario: Interrupt running session
- **WHEN** Session 状态为 RUNNING，调用 `POST /v1/session/{id}/interrupt`
- **THEN** SHALL 发送中断信号，Session 状态变为 INTERRUPTED，返回 HTTP 200

#### Scenario: Interrupt idle session
- **WHEN** Session 状态为 IDLE 或 COMPLETED
- **THEN** SHALL 返回 HTTP 200，body 包含 `"message": "session not running"`

### Requirement: DELETE /v1/session/{id} endpoint
系统 SHALL 提供 `DELETE /v1/session/{id}` 端点，销毁指定 Session 并清理进程资源。

#### Scenario: Delete existing session
- **WHEN** Session 存在，调用 `DELETE /v1/session/{id}`
- **THEN** SHALL 终止进程（如有）、移除 Session，返回 HTTP 200

#### Scenario: Delete non-existent session
- **WHEN** Session 不存在
- **THEN** SHALL 返回 HTTP 404 错误

### Requirement: GET /health endpoint
系统 SHALL 提供 `GET /health` 端点，返回服务健康状态。

#### Scenario: Health check
- **WHEN** 调用 `GET /health`
- **THEN** SHALL 返回 `{"status": "ok", "version": "0.1.0"}`

### Requirement: Configuration via environment variables
系统 SHALL 使用 `pydantic-settings` 管理配置，支持通过环境变量覆盖。MUST 支持的配置项：`CLAUDE_CLI_PATH`（CLI 路径）、`DEFAULT_MAX_TURNS`（默认最大轮数）、`EXECUTION_TIMEOUT`（执行超时秒数）、`HOST`（监听地址）、`PORT`（监听端口）。

#### Scenario: Load config from environment
- **WHEN** 设置环境变量 `CLAUDE_CLI_PATH=/usr/local/bin/claude`
- **THEN** 系统 SHALL 使用该路径作为 Claude Code CLI 的执行路径

#### Scenario: Use default config values
- **WHEN** 未设置任何环境变量
- **THEN** `CLAUDE_CLI_PATH` 默认为 `"claude"`，`DEFAULT_MAX_TURNS` 默认为 `20`，`EXECUTION_TIMEOUT` 默认为 `300`，`HOST` 默认为 `"0.0.0.0"`，`PORT` 默认为 `8001`
