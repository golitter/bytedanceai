from __future__ import annotations

import asyncio
import logging
import time
from collections.abc import AsyncIterator

from src.adapters.registry import AdapterRegistry
from src.generated.events import EventType
from src.orchestrator.models import DispatchResult, TaskResult
from src.schemas.events import StreamEvent
from src.schemas.request import AgentType
from src.workspace.manager import WorkspaceManager

logger = logging.getLogger(__name__)


class ExecutionEngine:
    def __init__(
        self,
        registry: AdapterRegistry,
        workspace_mgr: WorkspaceManager | None = None,
        repo_path: str = "",
        task_id: str = "",
        shared_dir: str = "",
        cwd: str = "",
    ) -> None:
        self._registry = registry
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
        for dispatch in dispatches:
            async for item in self._execute_task(dispatch, timeout_per_task):
                yield item

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

        adapter = None
        session_id = dispatch.real_session_id or f"orch-{task_id}"
        success = False
        collected: list[str] = []

        try:
            agent_cwd = await self._ensure_worktree(dispatch)

            adapter_cls = self._registry.get(agent_type)
            adapter = adapter_cls()
            await adapter.create_session(session_id)

            stream_kwargs: dict = {
                "task_id": task_id,
                "shared_dir": self._shared_dir,
            }
            if agent_cwd:
                stream_kwargs["cwd"] = agent_cwd

            logger.info(
                "ExecutionEngine: starting agent=%s type=%s task=%s session=%s cwd=%s",
                agent_name,
                agent_type,
                task_id,
                session_id,
                agent_cwd,
            )

            async for event in adapter.stream_chat(session_id, dispatch.content, **stream_kwargs):
                # Inject agent identity into sub-agent text events
                if event.type in (EventType.TEXT.value, "text"):
                    event.content["agent"] = agent_name
                    event.content["agent_type"] = agent_type
                yield (event, None)
                if event.type in (EventType.TEXT.value, "text"):
                    text = event.content.get("text", "")
                    if text:
                        collected.append(text)

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
        finally:
            if adapter:
                try:
                    await adapter.destroy_session(session_id)
                except Exception:
                    pass

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
