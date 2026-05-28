# Bug Fix: 子 Agent 消息持久化 — SSE 挂起 + 回复错乱

> 日期: 2026-05-29
> 关联 Change: `openspec/changes/sub-agent-message-persistence/`

## 现象

实现子 Agent 消息持久化（ExecutionEngine 改为回调后端 RunTask API）后，出现三个问题：

1. **SSE 永远「正在回复」** — 前端订阅的 orchestrator SSE 流永远不结束
2. **「管理者」消息显示内部规划** — orchestrator 回复内容是 LLM 的原始规划推理，而非汇总结果
3. **子 Agent 的 dispatch content 被当作用户消息显示** — "执行 pwd 命令..." 等内部指令出现在聊天中

---

## 根因分析

### Bug 1: SSE 永远挂起

**根因**: `BackendClient.stream_result()` 使用 `httpx.aiter_lines()` 消费后端 SSE 长连接。httpx 的 `aiter_lines()` 对 SSE 协议的 `\n\n` 事件分隔处理不可靠，导致 async generator 永远挂起。

**调用链**:

```
OrchestratorAdapter.stream_chat()
  → ExecutionEngine._execute_task()
    → BackendClient.stream_result()   ← 卡在这里
      → httpx.aiter_lines()           ← 永远不返回
```

`stream_result()` 挂起 → `_execute_task()` 无法 yield `RUNTIME_COMPLETED` → `stream_chat()` 卡在 Phase 3 → 前端永远收不到 DONE 事件。

**修复** (`agentend/src/clients/backend_client.py`):

- 改用独立 `httpx.AsyncClient`（长 read timeout）+ `aiter_text()` + 手动缓冲区拆行
- `_execute_task()` 中用 `asyncio.wait_for(_collect_stream(), timeout=timeout)` 加超时保底

```python
# Before (不可靠)
async with self._client.stream("GET", url, params=params) as resp:
    async for line in resp.aiter_lines():  # ← 可能永远挂起
        ...

# After (可靠)
sse_client = httpx.AsyncClient(timeout=httpx.Timeout(connect=10, read=600, write=10, pool=10))
async with sse_client.stream("GET", url, params=params) as resp:
    buf = ""
    async for text_chunk in resp.aiter_text():  # ← 逐块接收
        buf += text_chunk
        while "\n" in buf:
            line, buf = buf.split("\n", 1)
            # ... 解析 SSE data 行
```

### Bug 2: Orchestrator 回复显示内部规划

**根因**: `OrchestratorAdapter` Phase 4 的 TEXT 事件使用了 `overview`（LLM 的原始规划推理），而非 `aggregated`（汇总后的干净结果）。

```python
# Before
yield StreamEvent.create(EventType.TEXT, text=overview)
# overview 内容类似: "用户询问所有可用 agent，因此规划一个任务让该 agent 执行..."

# After
yield StreamEvent.create(EventType.TEXT, text=aggregated or overview)
# aggregated 是 Aggregator 合并子 Agent 结果后的干净摘要
```

**文件**: `agentend/src/adapters/orchestrator.py` Phase 4

### Bug 3: 子 Agent 的 dispatch content 被当作用户消息

**根因**: ExecutionEngine 调用 `POST /api/tasks/:taskId/run` 时，后端 RunTask 总是为子 Agent session 创建 user message。前端 `getTaskMessages` 返回所有 session 的所有消息，包括这条内部指令。

**修复**:

1. **后端** (`backend/internal/handler/task.go`): `RunTaskReq` 新增 `SkipUserMessage bool` 字段，当为 `true` 时跳过创建 user message

```go
type RunTaskReq struct {
    Message         string `json:"message" binding:"required"`
    AgentType       string `json:"agent_type"`
    SessionID       string `json:"session_id" binding:"required"`
    Cwd             string `json:"cwd"`
    SkipUserMessage bool   `json:"skip_user_message"`
}

// RunTask handler 中:
if !req.SkipUserMessage {
    // 创建 user message
}
```

2. **AgentEnd** (`agentend/src/clients/backend_client.py`): `run_task()` 请求体中传 `"skip_user_message": True`

---

## 修改文件清单

| 文件 | 变更 |
|------|------|
| `agentend/src/clients/backend_client.py` | SSE 改用 `aiter_text()` + 手动拆行；传 `skip_user_message` |
| `agentend/src/orchestrator/execution/engine.py` | SSE 消费加 `asyncio.wait_for` 超时保底 |
| `agentend/src/adapters/orchestrator.py` | TEXT 事件从 `overview` 改为 `aggregated or overview` |
| `backend/internal/handler/task.go` | RunTaskReq 新增 `SkipUserMessage` + `Cwd` 字段 |

---

## 教训

1. **httpx 不适合直接消费 SSE**: `aiter_lines()` 对长连接 SSE 的处理不可靠，应使用 `aiter_text()` 或 `aiter_bytes()` 手动解析
2. **SSE 消费必须加超时**: 任何 async generator 消费外部 SSE 的地方都应用 `asyncio.wait_for` 包裹，防止无限挂起
3. **内部 API 调用应与前端调用区分**: orchestrator 内部回调 RunTask 的语义不同于前端发起的 RunTask，需要通过参数（如 `skip_user_message`）区分行为
