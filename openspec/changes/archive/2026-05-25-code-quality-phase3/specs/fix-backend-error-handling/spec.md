## ADDED Requirements

### Requirement: stream.go 错误响应使用 vo.Response 辅助函数
`internal/handler/stream.go` 中的 JSON 错误响应 SHALL 使用 `vo.BadRequest()`、`vo.NotFound()` 等辅助函数，替换直接构造 `gin.H{"msg": "..."}` 的方式。

#### Scenario: 无效 session ID 返回标准格式
- **WHEN** 客户端请求 SSE 流但 sessionId 无效
- **THEN** 响应体格式为 `{"code": 400, "msg": "...", "data": null}`，使用 `vo.BadRequest()`

#### Scenario: 未找到 session 返回标准格式
- **WHEN** 客户端请求 SSE 流但 session 不存在
- **THEN** 响应体格式为 `{"code": 404, "msg": "...", "data": null}`，使用 `vo.NotFound()`

### Requirement: CreateTask 中 Session 创建失败返回错误
`internal/handler/task.go` 的 `CreateTask` 函数 SHALL 在 Session 创建失败时返回 HTTP 500 错误给客户端，而非仅记录 Warn 日志后继续。

#### Scenario: Session 创建失败时返回 500
- **WHEN** 数据库 Session 插入失败
- **THEN** 返回 `vo.InternalError(c, "failed to create session")`，HTTP 状态码 500

### Requirement: 用户消息保存失败返回错误
`internal/handler/task.go` 的 `RunTask` 函数 SHALL 在用户消息保存到数据库失败时返回错误给客户端。

#### Scenario: 消息保存失败时返回 500
- **WHEN** 用户消息 DB 插入失败
- **THEN** 返回 `vo.InternalError(c, "failed to save user message")`，不继续执行 Agent 调用

### Requirement: 基础设施故障日志升级为 Error 级别
`internal/stream/writer.go` 中 Redis XADD 失败和 MySQL flush 失败 SHALL 使用 `slog.Error` 而非 `slog.Warn`。

#### Scenario: Redis XADD 失败记录 Error
- **WHEN** Redis XADD 操作返回错误
- **THEN** 使用 `slog.Error("redis xadd failed", ...)` 记录日志

#### Scenario: MySQL flush 失败记录 Error
- **WHEN** MySQL 批量写入失败
- **THEN** 使用 `slog.Error("mysql flush failed", ...)` 记录日志

### Requirement: updateStatus 返回错误
`internal/stream/writer.go` 的 `updateStatus` 函数 SHALL 返回 error，允许调用者感知状态更新失败。

#### Scenario: 状态更新失败传播给调用者
- **WHEN** 数据库状态更新失败
- **THEN** `updateStatus` 返回非 nil error，调用者可据此决定后续行为

### Requirement: Redis Init 添加连接验证
`pkg/redis/redis.go` 的 `Init` 函数 SHALL 在创建客户端后调用 `Ping()` 验证连接可用性。

#### Scenario: Redis 不可用时启动报错
- **WHEN** Redis 服务不可用
- **THEN** `Init()` 返回明确的错误信息，服务启动失败

### Requirement: Redis Close 在服务关闭时调用
`cmd/server/main.go` SHALL 在服务关闭时调用 `redis.Close()` 释放连接。

#### Scenario: 优雅关闭释放 Redis 连接
- **WHEN** 服务收到关闭信号
- **THEN** `defer redis.Close()` 被执行，Redis 连接被正确释放
