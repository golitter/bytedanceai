## Why

全端代码审计发现 97 个问题，其中代码质量类（Medium/Low）占 70 个。第一阶段（安全加固）和第二阶段（可靠性修复）暂缓，第三阶段聚焦代码质量：死代码清理、错误处理模式修复、异步阻塞消除、契约层补全、错误响应格式统一。在功能继续迭代前清理技术债，避免后续维护成本指数增长。

## What Changes

- 移除前端未使用的依赖（`radix-ui`）、未使用的 shadcn/ui 组件（card/button/input）、未使用的生成类型文件、死代码导出和 Vite 脚手架资源
- 修复后端 `stream.go` 错误响应绕过 `vo.Response` 的问题，统一使用标准响应格式
- 修复后端 `task.go` 中被吞没的错误（Session 创建失败、用户消息保存失败），改为正确返回错误
- 升级后端 `writer.go` 中 Redis/MySQL 失败日志级别从 Warn 到 Error，`updateStatus` 返回错误
- 后端 `redis.go` 添加 Ping 连接验证和 Close 优雅关闭
- 将 Agent 端 `pin_memory.py` 和 `aggregator.py` 中的同步 LLM 调用改为异步（`ainvoke`）
- 修复 Agent 端 `graph.py` 中 `_extract_json` 和 `plan_node` 缺失的错误处理
- 将 `message.yaml` 和 `validate-repo-path.yaml` 加入契约代码生成器 `OUTPUT_MAP` 并重新生成类型
- 修复 Agent 端 `session/models.py` 中 `_VALID_TRANSITIONS` 缺少 `inactive` 状态
- 统一三端错误响应格式：后端全部使用 `vo.Response`，Agent 端创建共享错误模型

## Capabilities

### New Capabilities
- `cleanup-frontend-dead-code`: 清理前端死代码 — 移除未使用依赖、组件、生成文件、导出函数和脚手架资源
- `fix-backend-error-handling`: 修复后端错误处理 — 统一响应格式、修复吞没错误、升级日志级别、Redis 连接验证
- `fix-agentend-async-blocking`: 修复 Agent 端异步阻塞 — 同步 LLM 调用改异步、JSON 解析错误处理
- `complete-contract-generator`: 完善契约代码生成 — 补全缺失 schema、修复状态机转换

### Modified Capabilities
- `api-gateway`: 后端错误响应格式统一 — stream.go 错误响应改用 vo.Response 辅助函数
- `contract-workflow`: 契约生成器补全 — OUTPUT_MAP 新增 message 和 validate-repo-path schema
- `stream-protocol`: Agent 端流式调用异步化 — pin_memory 和 aggregator 同步改异步

## Impact

- **前端**: `package.json`（移除依赖）、`src/components/ui/`（删除文件）、`src/generated/`（删除文件）、`src/lib/api.ts`（移除导出）、`src/assets/`（删除文件）
- **后端**: `internal/handler/stream.go`、`internal/handler/task.go`、`internal/stream/writer.go`、`pkg/redis/redis.go`、`cmd/server/main.go`
- **Agent 端**: `src/orchestrator/pin_memory.py`、`src/orchestrator/aggregator.py`、`src/orchestrator/graph.py`、`src/session/models.py`
- **契约层**: `scripts/generate_contracts.py`（修改 OUTPUT_MAP）、三端 `generated/` 目录（新增文件）
- **无破坏性 API 变更** — 所有改动为内部质量提升，不影响外部接口
