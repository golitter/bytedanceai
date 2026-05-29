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

## 经验教训

异步生成器（`async for ... yield`）中直接调用可能无限等待的外部进程时，必须加超时保护。Python 3.10 没有 `asyncio.timeout()` 上下文管理器，需要手动用 Queue + Task 实现。
