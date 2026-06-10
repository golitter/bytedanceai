## ADDED Requirements

### Requirement: Go 模块初始化与依赖安装
系统 SHALL 在 backend/ 目录下初始化 Go 模块 `agenthub/backend` 并安装核心依赖（gin, gorm, mysql driver, yaml, jwt, cors）。

#### Scenario: 模块可用
- **WHEN** 执行 `cd backend && go mod tidy`
- **THEN** `go.mod` 和 `go.sum` 生成且无错误

### Requirement: 分层目录结构
系统 SHALL 按照 controller/dao/service/model/middleware/vo/conf 分层组织代码，每层接口与实现分离。

#### Scenario: 目录存在
- **WHEN** 查看 backend/ 目录
- **THEN** 存在 `internal/{conf,controller/impl,dao/gorm,dao/mock,model,service/impl,middleware,vo}` 和 `pkg/db` 目录

### Requirement: YAML 配置加载
系统 SHALL 从 `configs/config.yaml` 读取 MySQL 和 JWT 配置，解析为 Go struct。

#### Scenario: 配置加载成功
- **WHEN** 调用 `conf.Load("configs/config.yaml")`
- **THEN** 返回包含 MySQL 连接信息和 JWT 配置的 Config 对象

#### Scenario: 配置文件不存在
- **WHEN** 配置文件路径不存在
- **THEN** 返回错误

### Requirement: MySQL 单例连接
系统 SHALL 使用 sync.Once 确保 MySQL 连接只初始化一次，通过 `db.Init(cfg)` 创建，`db.GetDB()` 全局获取。

#### Scenario: 数据库连接成功
- **WHEN** 调用 `db.Init(&cfg.MySQL)` 且 MySQL 服务可用
- **THEN** 返回可用的 `*gorm.DB` 实例

#### Scenario: 数据库不可用
- **WHEN** MySQL 服务未启动
- **THEN** 返回连接错误

### Requirement: 统一响应格式
系统 SHALL 所有 API 响应使用 `{code, data, msg}` 格式，通过 vo 包的辅助函数返回。

#### Scenario: 成功响应
- **WHEN** 调用 `vo.OK(ctx, data)`
- **THEN** 返回 HTTP 200，body 为 `{"code":0,"data":...}`

#### Scenario: 错误响应
- **WHEN** 调用 `vo.BadRequest(ctx, "invalid")`
- **THEN** 返回 HTTP 400，body 为 `{"code":400,"msg":"invalid"}`

### Requirement: CORS 中间件
系统 SHALL 提供 CORS 中间件，允许 `http://localhost:5173` 跨域访问。

#### Scenario: 预检请求
- **WHEN** 发送 OPTIONS 请求
- **THEN** 返回 204 并设置正确的 CORS 头

### Requirement: Logger 中间件
系统 SHALL 记录每个请求的 method、path、status、latency、client_ip，按状态码分级（>=500 ERROR, >=400 WARN, 其余 INFO）。

#### Scenario: 请求日志
- **WHEN** 收到 GET /ping 请求
- **THEN** 输出 INFO 级别日志包含 status=200

### Requirement: JWT Auth 中间件
系统 SHALL 提供 JWT Bearer Token 校验中间件，校验失败返回 401。

#### Scenario: 无 Token
- **WHEN** 请求不含 Authorization 头
- **THEN** 返回 HTTP 401

#### Scenario: 有效 Token
- **WHEN** 请求携带有效 Bearer Token
- **THEN** 将 userID 和 userName 写入 gin.Context 并放行

### Requirement: Ping 接口
系统 SHALL 提供 GET /ping 接口，返回 `{"code":0,"data":{"message":"pong"}}`。

#### Scenario: Ping 成功
- **WHEN** 发送 GET /ping
- **THEN** 返回 HTTP 200，body 为 `{"code":0,"data":{"message":"pong"}}`

### Requirement: Makefile
系统 SHALL 提供 Makefile 支持 build、run、fmt、tidy 命令。

#### Scenario: 构建并运行
- **WHEN** 执行 `make run`
- **THEN** 编译并启动服务在 :8080
