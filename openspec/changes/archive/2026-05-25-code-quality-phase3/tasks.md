## 1. 清理前端死代码

- [x] 1.1 从 package.json 移除 `radix-ui` 依赖，执行 `pnpm install`
- [x] 1.2 删除 `src/components/ui/card.tsx`、`button.tsx`、`input.tsx`
- [x] 1.3 检查 `@radix-ui/react-slot` 和 `class-variance-authority` 是否仍被引用，无引用则从 package.json 移除
- [x] 1.4 删除 `src/generated/response.ts` 和 `src/generated/session.ts`
- [x] 1.5 从 `src/lib/api.ts` 移除 `deleteTask()`、`patchSession()` 函数和 `StreamEvent` 重导出
- [x] 1.6 删除 `src/assets/hero.png`、`src/assets/react.svg`、`src/assets/vite.svg`
- [x] 1.7 运行 `pnpm build` 验证编译通过

## 2. 修复后端错误处理

- [x] 2.1 `stream.go`: 将 `gin.H{"msg": ...}` 错误响应替换为 `vo.BadRequest()`/`vo.NotFound()`
- [x] 2.2 `task.go` CreateTask: Session 创建失败时返回 `vo.InternalError()` 而非仅 Warn
- [x] 2.3 `task.go` RunTask: 用户消息保存失败时返回 `vo.InternalError()` 而非继续执行
- [x] 2.4 `writer.go`: Redis XADD 失败和 MySQL flush 失败日志从 `slog.Warn` 升级为 `slog.Error`
- [x] 2.5 `writer.go` updateStatus: 修改为返回 error，调用者处理返回值
- [x] 2.6 `redis.go` Init: 添加 `client.Ping(ctx)` 连接验证
- [x] 2.7 `redis.go`: 添加 `Close()` 函数
- [x] 2.8 `main.go`: 添加 `defer redis.Close()` 优雅关闭
- [x] 2.9 运行 `go build ./...` 验证编译通过

## 3. 修复 Agent 端异步阻塞

- [x] 3.1 `pin_memory.py`: `_generate_summary` 改为 async，使用 `await llm.ainvoke()`
- [x] 3.2 `pin_memory.py`: `pin()` 和 `pin_existing()` 中 `await _generate_summary()`
- [x] 3.3 `aggregator.py`: `aggregate` 改为 async，使用 `await llm.ainvoke()`
- [x] 3.4 检查并修复 `aggregate()` 所有上游调用方，确保使用 `await`
- [x] 3.5 `graph.py` `_extract_json`: 添加 try/except json.JSONDecodeError，失败返回 None
- [x] 3.6 `graph.py` `plan_node`: 添加 try/except 捕获 LLM 和 Pydantic 异常，返回 fallback `{"plan": None}`

## 4. 完善契约代码生成

- [x] 4.1 `generate_contracts.py`: 在 OUTPUT_MAP 添加 `message` 条目
- [x] 4.2 `generate_contracts.py`: 在 OUTPUT_MAP 添加 `validate-repo-path` 条目
- [x] 4.3 运行 `make generate` 生成三端类型，检查无命名冲突
- [x] 4.4 `session/models.py`: 在 `_VALID_TRANSITIONS` 添加 `SessionState.INACTIVE: set()`

## 5. 验证

- [x] 5.1 前端 `pnpm build` 通过
- [x] 5.2 后端 `go build ./...` 通过
- [x] 5.3 Agent 端 Python 导入检查通过（`python -c "import src"` 或类似）
- [x] 5.4 契约 `make generate` 通过，生成文件无冲突
