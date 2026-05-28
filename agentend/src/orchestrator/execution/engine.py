from __future__ import annotations

import asyncio
import logging
import time
from collections.abc import AsyncIterator

from src.clients.backend_client import BackendClient
from src.generated.events import EventType
from src.orchestrator.models import DispatchResult, TaskResult
from src.schemas.events import StreamEvent
from src.schemas.request import AgentType
from src.workspace.manager import WorkspaceManager

logger = logging.getLogger(__name__)


class ExecutionEngine:
    def __init__(
        self,
        backend_client: BackendClient,
        workspace_mgr: WorkspaceManager | None = None,
        repo_path: str = "",
        task_id: str = "",
        shared_dir: str = "",
        cwd: str = "",
    ) -> None:
        self._backend_client = backend_client
        self._workspace_mgr = workspace_mgr
        self._repo_path = repo_path
        self._task_id = task_id
        self._shared_dir = shared_dir
        self._cwd = cwd

    async def execute(
        self,
        dispatches: list[DispatchResult],
        timeout_per_task: float = 300.0,
    ) -> AsyncIterator[tuple[StreamEvent, TaskResult | None]]:
        if len(dispatches) <= 1:
            for dispatch in dispatches:
                async for item in self._execute_task(dispatch, timeout_per_task):
                    yield item
            return

        queue: asyncio.Queue[tuple[StreamEvent, TaskResult | None]] = asyncio.Queue()

        async def _run(dispatch: DispatchResult) -> None:
            async for item in self._execute_task(dispatch, timeout_per_task):
                await queue.put(item)

        tasks = [asyncio.create_task(_run(d)) for d in dispatches]

        async def _drain() -> None:
            await asyncio.gather(*tasks)
            await queue.put(None)  # sentinel

        drain_task = asyncio.create_task(_drain())

        while True:
            item = await queue.get()
            if item is None:
                break
            yield item

        await drain_task

    async def _ensure_worktree(self, dispatch: DispatchResult) -> str:
        """为子 agent 创建独立 worktree，返回 worktree 路径作为 cwd。"""
        if not self._workspace_mgr or not self._repo_path:
            return self._cwd or dispatch.workspace_path

        real_session_id = dispatch.real_session_id
        if not real_session_id:
            logger.warning("ExecutionEngine: no real_session_id for task=%s, fallback to shared cwd", dispatch.task_id)
            return self._cwd or dispatch.workspace_path

        try:
            agent_type_str = dispatch.agent_type or dispatch.agent
            try:
                agent_type = AgentType(agent_type_str)
            except ValueError:
                agent_type = AgentType.CLAUDE_CODE

            ws = await self._workspace_mgr.create(
                repo_path=self._repo_path,
                task_id=self._task_id,
                agent_name=dispatch.agent,
                session_id=real_session_id,
                agent_type=agent_type,
            )
            logger.info(
                "ExecutionEngine: created worktree for agent=%s session=%s path=%s",
                dispatch.agent,
                real_session_id,
                ws.worktree_path,
            )
            return ws.worktree_path
        except Exception:
            logger.exception(
                "ExecutionEngine: failed to create worktree for agent=%s session=%s",
                dispatch.agent,
                real_session_id,
            )
            return self._cwd or dispatch.workspace_path

    async def _execute_task(
        self,
        dispatch: DispatchResult,
        timeout: float,
    ) -> AsyncIterator[tuple[StreamEvent, TaskResult | None]]:
        task_id = dispatch.task_id
        agent_name = dispatch.agent
        agent_type = dispatch.agent_type or agent_name
        start = time.monotonic()

        yield (
            StreamEvent.create(
                EventType.RUNTIME_EXECUTING,
                task_id=task_id,
                agent=agent_name,
                title=dispatch.content[:80],
                status="running",
            ),
            None,
        )

        session_id = dispatch.real_session_id or f"orch-{task_id}"
        success = False
        collected: list[str] = []

        try:
            agent_cwd = await self._ensure_worktree(dispatch)

            logger.info(
                "ExecutionEngine: dispatching agent=%s type=%s task=%s session=%s cwd=%s via backend RunTask",
                agent_name,
                agent_type,
                task_id,
                session_id,
                agent_cwd,
            )

            message_id = await asyncio.wait_for(
                self._backend_client.run_task(
                    task_id=self._task_id,
                    session_id=session_id,
                    message=dispatch.content,
                    agent_type=agent_type,
                    cwd=agent_cwd,
                ),
                timeout=30.0,
            )

            # Collect SSE events with an overall timeout as safety net
            async def _collect_stream() -> None:
                nonlocal success
                async for event in self._backend_client.stream_result(
                    task_id=self._task_id,
                    message_id=message_id,
                    session_id=session_id,
                ):
                    event_type = event.get("type", "")
                    if event_type == "text":
                        content = event.get("content", {})
                        text = content.get("text", "") if isinstance(content, dict) else str(content)
                        if text:
                            collected.append(text)
                    elif event_type == "done":
                        break
                    elif event_type == "error":
                        content = event.get("content", {})
                        msg = content.get("message", "unknown error") if isinstance(content, dict) else str(content)
                        collected.append(f"[Error] {msg}")
                        break

            await asyncio.wait_for(_collect_stream(), timeout=timeout)

            success = True
            logger.info(
                "ExecutionEngine: completed agent=%s task=%s collected=%d chars",
                agent_name,
                task_id,
                len("".join(collected)),
            )
        except asyncio.TimeoutError:
            msg = f"[Timeout] Task {task_id} exceeded {timeout}s"
            logger.warning("ExecutionEngine: %s", msg)
            collected.append(msg)
            yield (StreamEvent.create(EventType.ERROR, error=msg), None)
        except Exception as exc:
            msg = f"[Error] Task {task_id} agent={agent_name} failed: {exc}"
            logger.error("ExecutionEngine: %s", msg, exc_info=True)
            collected.append(msg)
            yield (StreamEvent.create(EventType.ERROR, error=msg), None)

        duration = time.monotonic() - start
        result = TaskResult(
            task_id=task_id,
            agent=agent_name,
            success=success,
            content="".join(collected),
            duration=round(duration, 2),
        )

        yield (
            StreamEvent.create(
                EventType.RUNTIME_COMPLETED,
                task_id=task_id,
                agent=agent_name,
                success=success,
                duration=result.duration,
                status="completed" if success else "failed",
            ),
            result,
        )
