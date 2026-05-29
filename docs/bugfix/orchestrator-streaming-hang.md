# Bug Fix: Orchestrator 流式透传后一直 streaming 不结束

> 日期: 2026-05-29
> 关联 Change: `openspec/changes/streaming-perf-overhaul/`

## 现象

`streaming-perf-overhaul` 变更后，Orchestrator 多 Agent 模式下发起对话，orchestrator 一直停留在 `streaming` 状态，前端始终收不到 `done` 事件。

---

## 根因分析

### Bug: 短路路径无超时保护，CLI 子进程挂起导致整个流永不结束

**变更引入的短路路径**：`ExecutionEngine._execute_task()` 在检测到 `adapter_registry` 中有对应 agent_type 的 adapter 时，直接在进程内调用 `adapter.stream_chat()`，跳过 BackendClient 的 HTTP 回环。

**问题**：`adapter.stream_chat()` 内部通过 `asyncio.create_subprocess_exec()` 启动 CLI 子进程，然后用 `async for line in process.stdout` 逐行读取。整个调用 **没有超时保护**：

```python
# 修复前 — 无超时，CLI 挂起则整个 orchestrator 卡死
async for event in adapter.stream_chat(
    session_id, dispatch.content, cwd=agent_cwd,
):
    ...
```

如果 CLI 子进程因以下原因挂起：
- `cwd` 路径不存在（worktree 创建失败后的 fallback 路径为空）
- CLI 可执行文件路径错误
- CLI 本身卡在等待用户输入
- 子进程 stdout 管道阻塞

`_execute_task()` 永远不会返回 → `engine.execute()` 永远不会结束 → `orchestrator.stream_chat()` 永远到不了最后的 `yield DONE` 事件 → Backend StreamWriter 的 goroutine 永远在等 SSE 响应 → 前端永远看到 `streaming` 状态。

**次要问题**：`success` 标记在 ERROR break 后仍被设为 `True`（`success = True` 写在 for 循环之后，无论是正常结束还是 break 都会执行）。

---

## 修复方案

### 1. 新增 `_iter_adapter_with_timeout()` 方法

用 `asyncio.Queue` + `asyncio.create_task()` 包裹 adapter 迭代，实现总超时保护：

```python
async def _iter_adapter_with_timeout(
    self, adapter, session_id, message, cwd, timeout
) -> AsyncIterator[StreamEvent]:
    queue: asyncio.Queue[StreamEvent | None] = asyncio.Queue()

    async def _consume() -> None:
        try:
            async for event in adapter.stream_chat(session_id, message, cwd=cwd):
                await queue.put(event)
        except Exception as exc:
            await queue.put(StreamEvent.create(EventType.ERROR, error=str(exc)))
        await queue.put(None)  # sentinel

    task = asyncio.create_task(_consume())
    deadline = time.monotonic() + timeout

    try:
        while True:
            remaining = deadline - time.monotonic()
            if remaining <= 0:
                raise asyncio.TimeoutError()
            event = await asyncio.wait_for(queue.get(), timeout=max(remaining, 0.1))
            if event is None:
                break
            yield event
    finally:
        if not task.done():
            task.cancel()
            try:
                await task
            except asyncio.CancelledError:
                pass
```

### 2. 修正 `success` 标记

- `DONE` 事件 → `success = True`
- `ERROR` break → `success` 保持 `False`（不再在循环后无条件设 True）
- 循环正常结束（CLI 退出但没发 DONE/ERROR）→ `for...else` 中设 `success = True`

---

## 影响范围

- `agentend/src/orchestrator/execution/engine.py` — 新增 `_iter_adapter_with_timeout()`，修改 `_execute_task()` 的短路路径分支

---

## 第二轮修复：Orchestrator 重复输出 + Sub-Agent 执行不可见

> 日期: 2026-05-29（续）

### 现象

超时修复后，多 Agent 流式对话仍存在两个问题：

1. **重复输出**：Orchestrator 的规划概述（overview）被输出两次
2. **Sub-Agent 执行不可见**：实时流中只看到 Orchestrator 的规划和总结，看不到 Sub-Agent 的执行内容；刷新页面后才正常

### 根因分析

#### Bug A: `flushTextBuffer()` 丢失 agent 元数据

`StreamWriter.flushTextBuffer()` 合并多个 TEXT 事件为一条 SSE 后发布到 Redis。原实现使用 `FormatSSE()` 只保留 `{"type":"text","content":{"text":"..."}}` —— **丢失了 `agent` 和 `agent_type` 字段**。

前端 `use-chat-stream.ts` 只在 TEXT 事件同时携带 `agent` 和 `agent_type` 时才调用 `streamAgentUpdate()` 触发 Agent 切换。元数据丢失 → 前端永远检测不到 Agent 切换 → 所有内容（规划 + 执行 + 总结）被合并到一条 streamingContent 中，显示为单一的 Orchestrator 消息。

刷新后正常是因为 `switchAgent()` 在 MySQL 中为每个 Agent 创建了独立消息记录。

> 此 Bug 已通过将 `flushTextBuffer()` 改用 `FormatSSEWithMeta()` 修复。

#### Bug B: `currentAgentName` 未随 TEXT 事件更新

TEXT 事件处理只检查 `agent_type` 变化来触发 `switchAgent()`。当第一个 Orchestrator TEXT 到达时，`agent_type` 与初始值相同（"orchestrator"），不触发 `switchAgent()`，导致 `currentAgentName` 始终为 `""`。即使 `FormatSSEWithMeta` 包含了 `agent_type`，`agent` 字段为空，前端仍不会调用 `streamAgentUpdate()`。

#### Bug C: 无 Sub-Agent 时 overview 被输出两次

`OrchestratorAdapter.stream_chat()` 中：
- Phase 1: `yield TEXT(text=overview)` —— 输出规划概述
- Phase 4: `yield TEXT(text=aggregated or overview)` —— 聚合报告

当 `Aggregator.aggregate([], overview)` 返回空字符串（无 Sub-Agent 结果），`aggregated or overview` 回退到 `overview`，**同一文本被 yield 两次**。

### 修复方案

#### 1. `currentAgentName` 同步更新（writer.go）

在 TEXT 事件处理中，即使 `agent_type` 未变，也更新 `currentAgentName`：

```go
// agent_type 未变但携带 agent 名称 — 同步更新
if newName, ok := event.Content["agent"].(string); ok && newName != "" {
    sw.mu.Lock()
    sw.currentAgentName = newName
    sw.mu.Unlock()
}
```

#### 2. 跳过空聚合的 Phase 4 TEXT（orchestrator.py）

当 `aggregated` 为空（无 Sub-Agent 结果）时，不 yield Phase 4 TEXT 事件：

```python
aggregated = await aggregator.aggregate(task_results, overview)
if aggregated:
    yield StreamEvent.create(EventType.TEXT, text=aggregated, ...)
```

### 影响范围（第二轮）

- `backend/internal/stream/writer.go` — TEXT 事件处理中增加 `currentAgentName` 同步更新
- `agentend/src/adapters/orchestrator.py` — Phase 4 聚合结果为空时跳过 TEXT 事件

### 经验教训

1. SSE 事件合并（batch/merge）时必须保留完整的业务元数据，不能只保留 payload
2. Agent 追踪（type + name）必须在每个 TEXT 事件上保持同步，不能只在 switchAgent 时更新

## 经验教训

异步生成器（`async for ... yield`）中直接调用可能无限等待的外部进程时，必须加超时保护。Python 3.10 没有 `asyncio.timeout()` 上下文管理器，需要手动用 Queue + Task 实现。
