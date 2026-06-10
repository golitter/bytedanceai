## Context

三端 monorepo 代码审计产出 97 个发现，其中代码质量类问题 46 个 Medium + 24 个 Low。本次变更聚焦第三阶段代码质量修复，涉及前端（React/TS）、后端（Go/Gin）、Agent 端（Python/FastAPI）三端。当前状态：前端存在未使用依赖和死代码，后端错误处理不一致（部分绕过 vo.Response），Agent 端同步 LLM 调用阻塞事件循环，契约生成器缺少 2 个 schema。

## Goals / Non-Goals

**Goals:**
- 清理前端所有已确认的死代码和未使用依赖
- 统一后端错误响应格式，修复被吞没的错误
- 消除 Agent 端事件循环阻塞，补全 LLM 调用错误处理
- 补全契约生成器缺失的 schema，修复状态机转换
- 所有改动不引入破坏性 API 变更

**Non-Goals:**
- 不做安全加固（Phase 1）和可靠性修复（Phase 2）
- 不添加新的测试覆盖（后续 Phase 4）
- 不修改 API 接口定义或路由结构
- 不处理性能优化

## Decisions

### 前端清理策略：直接删除 vs 保留
- **决定**: 直接删除所有确认未使用的文件和依赖
- **理由**: `radix-ui`、`card.tsx`/`button.tsx`/`input.tsx`、`response.ts`/`session.ts` 经 grep 确认无任何导入。删除后 `pnpm build` 即可验证无破坏。
- **替代方案**: 保留并标记 `@deprecated` — 增加维护噪音，不如直接删除。

### 后端错误处理：vo.Response 统一
- **决定**: 所有 HTTP 错误响应统一使用 `vo.BadRequest()`/`vo.NotFound()`/`vo.InternalError()` 等辅助函数
- **理由**: 已有 `internal/vo/response.go` 提供标准格式 `{code, data, msg}`，只需修复绕过点（stream.go）。
- **注意**: SSE 流式响应中的错误仍需保持 SSE 格式（`event: error`），不强制 JSON。

### Agent 端异步化：ainvoke vs to_thread
- **决定**: 使用 `await llm.ainvoke()` 替代 `llm.invoke()`
- **理由**: LangChain 的 `ChatOpenAI` 原生支持 `ainvoke()`，比 `asyncio.to_thread()` 更高效（无需线程切换）。调用方 `pin()` 和 `aggregate()` 已在 async 上下文中，改为 `async` 顺理成章。
- **替代方案**: `asyncio.to_thread(llm.invoke, ...)` — 可行但不如原生 async 高效。

### 契约生成器补全
- **决定**: 在 `OUTPUT_MAP` 中添加 `message` 和 `validate-repo-path` 条目
- **理由**: `contracts/schemas/` 中已有 YAML 定义但未被处理。添加后运行 `make generate` 即可生成三端类型。如果生成类型与现有手写类型冲突，优先使用生成类型。

## Risks / Trade-offs

- **删除文件可能遗漏隐式引用** → 删除后立即运行 `pnpm build` 验证编译通过
- **Agent 端 async 改动可能影响调用链** → `pin()` 和 `aggregate()` 的调用方需同步改为 `await`，需要检查完整调用链
- **契约生成可能产生与手写代码冲突的类型** → 先检查生成输出是否与现有代码命名冲突
- **后端错误处理改动可能改变 API 响应格式** → stream.go 的 SSE 错误保持 SSE 格式，仅 JSON 错误点改用 vo.Response
