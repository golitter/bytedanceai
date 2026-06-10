## 1. 项目初始化

- [x] 1.1 创建目录结构 `cmd/server`, `internal/{conf,controller/impl,dao/gorm,dao/mock,model,service/impl,middleware,vo}`, `pkg/db`, `configs`, `docs/api`
- [x] 1.2 执行 `go mod init agenthub/backend` 并安装核心依赖
- [x] 1.3 创建 `configs/config.yaml`（mysql + jwt 配置）

## 2. 基础设施层

- [x] 2.1 实现 `internal/conf/conf.go` — YAML 配置加载（MySQLConfig + JWTConfig → Config struct）
- [x] 2.2 实现 `pkg/db/mysql.go` — MySQL 单例连接（sync.Once + DSN 构建）
- [x] 2.3 实现 `internal/vo/response.go` — 统一响应（OK / Created / BadRequest / NotFound / InternalError）

## 3. 中间件

- [x] 3.1 实现 `internal/middleware/cors.go` — CORS 中间件（允许 localhost:5173）
- [x] 3.2 实现 `internal/middleware/logger.go` — 请求日志（slog，按状态码分级）
- [x] 3.3 实现 `internal/middleware/auth.go` — JWT Auth 中间件（GenerateToken + 校验）

## 4. 入口与接口

- [x] 4.1 实现 `cmd/server/main.go` — 加载配置 → 连 DB → 注册中间件 → 注册路由 → 启动
- [x] 4.2 实现 GET /ping 接口，返回 `{"code":0,"data":{"message":"pong"}}`

## 5. 构建工具

- [x] 5.1 创建 `Makefile`（build / run / fmt / tidy）
- [x] 5.2 验证：`make run` 启动后 `curl localhost:8080/ping` 返回正确响应
