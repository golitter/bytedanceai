# 2026-06-03 SSE 流式性能优化 — 消除三处同步阻塞

**类型**: agentend（无跨端契约变更）

## 变更原因

SSE 流式链路中存在三处同步阻塞点，导致首 token 延迟和流式卡顿：
1. `get_pinned_announcements` 串行 HTTP 阻塞（50-500ms）
2. `set_cli_session_id._save()` 同步磁盘 I/O 阻塞事件循环（10-50ms）
3. LangSmith trace `post()`/`patch()` 同步网络调用引入随机延迟尖刺

## 变更文件

| 文件 | 变更 |
|------|------|
| `agentend/src/session/store.py` | `set_cli_session_id` / `delete` 改为 `async def`，`_save()` 用 `asyncio.to_thread` 包裹 |
| `agentend/src/api/v1/agent.py` | `get_pinned_announcements` 改为 `create_task` 并行获取；两处 `set_cli_session_id` 加 `await` |
| `agentend/src/adapters/trace.py` | 所有 `post()`/`patch()` 调用改为 `await asyncio.to_thread(...)` |

## 对比结果

- API 接口签名不变（`/stream`、`/execute` 的请求/响应格式完全一致）
- SSE 事件类型、事件格式不变
- LangSmith tracing 行为不变（开启时正常上报，关闭时零开销）
- `get_pinned_announcements` 失败时保留原有优雅降级（返回 `[]`）

## 跨端影响

- **Frontend**: 无影响
- **Backend**: 无影响
- **contracts/schemas/**: 无变更
