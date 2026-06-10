## ADDED Requirements

### Requirement: StreamAgent 调用 AgentEnd SSE 流
系统 SHALL 提供 `StreamAgent(req AgentRequest)` 方法，向 AgentEnd 发起 `POST /v1/agent/stream` 请求，返回原始 `*http.Response` 供上层逐行读取 SSE 数据。

#### Scenario: 成功建立 SSE 连接
- **WHEN** 调用 `StreamAgent` 并传入合法的 AgentRequest
- **THEN** 返回 HTTP 200 的 response body，Content-Type 为 `text/event-stream`

#### Scenario: AgentEnd 不可用
- **WHEN** AgentEnd 服务未启动或网络不通
- **THEN** 返回 error，调用方可根据 error 判断连接失败

### Requirement: HealthCheck 检测 AgentEnd 状态
系统 SHALL 提供 `HealthCheck()` 方法，向 AgentEnd 发起 `GET /health` 请求。

#### Scenario: AgentEnd 健康
- **WHEN** 调用 `HealthCheck()` 且 AgentEnd 正常运行
- **THEN** 返回 nil（无错误）

#### Scenario: AgentEnd 不可达
- **WHEN** 调用 `HealthCheck()` 且 AgentEnd 不可达
- **THEN** 返回非 nil error

### Requirement: Client 从配置初始化
系统 SHALL 支持通过 `New(host string, port int)` 构造函数创建 Client 实例，base URL 为 `http://{host}:{port}`。

#### Scenario: 构造 Client 实例
- **WHEN** 调用 `New("http://localhost", 8001)`
- **THEN** 创建 Client 实例，后续请求的 base URL 为 `http://localhost:8001`
