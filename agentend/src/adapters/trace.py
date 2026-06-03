"""LangSmith RunTree trace for CLI Adapter StreamEvent lifecycle.

Wraps an AsyncIterator[StreamEvent] with LangSmith RunTree tracing.
Only traces the StreamEvent lifecycle — does not capture the CLI's internal
LLM prompts (which are not exposed by CLI adapters).

Usage::

    from src.adapters.trace import trace_stream_events

    traced = trace_stream_events(
        adapter.stream_chat(session_id, message, **kwargs),
        run_name=f"claude-code session_id={session_id}",
        inputs={"message": message, "session_id": session_id},
    )
    async for event in traced:
        ...
"""

import asyncio
import logging
import os
from collections.abc import AsyncIterator

from langsmith.run_trees import RunTree

from src.schemas.events import EventType, StreamEvent

logger = logging.getLogger(__name__)


def _is_tracing_enabled() -> bool:
    """Check if LangSmith tracing is enabled.

    Mirrors LangChain SDK behavior: tracing is active when LANGSMITH_API_KEY
    is present, regardless of LANGSMITH_TRACING value.
    """
    return bool(os.getenv("LANGSMITH_API_KEY"))


async def trace_stream_events(
    events: AsyncIterator[StreamEvent],
    *,
    run_name: str,
    inputs: dict,
) -> AsyncIterator[StreamEvent]:
    """Wrap an AsyncIterator[StreamEvent] with LangSmith RunTree tracing.

    Creates a root RunTree, then maps StreamEvent lifecycle to child runs:

    - **INIT** → create an ``init`` child run (cli_session_id, agent_type).
    - **TEXT** events → aggregated into ``text_parts``; flushed as a single
      child run at TOOL_CALL / DONE / ERROR boundaries.
    - **TOOL_CALL** → flush text, then create a pending ``tool:{name}`` child run.
    - **TOOL_RESULT** → LIFO-match the pending tool run and ``end()`` it.
    - **DONE** → flush remaining text, create a ``done`` child run (usage).
    - **ERROR** → flush text + mark ``root.error``.

    The original StreamEvent objects are yielded unchanged — this wrapper
    is transparent to downstream consumers.

    Args:
        events: Raw event stream from ``adapter.stream_chat()``.
        run_name: Name for the root RunTree (e.g. ``"claude-code session_id=xxx"``).
        inputs: Input metadata for the root run (e.g. ``{"message": "...", "session_id": "..."}``).

    Yields:
        The original StreamEvent objects, unchanged.
    """
    if not _is_tracing_enabled():
        async for event in events:
            yield event
        return

    root = RunTree(
        name=run_name,
        run_type="chain",
        inputs=inputs,
    )
    # post() sends the root run to LangSmith server.
    # Without this, create_child() + patch() operate on a server-side orphan.
    await asyncio.to_thread(root.post)

    text_parts: list[str] = []
    # Track full output text across all TEXT events for root run outputs.
    output_parts: list[str] = []
    # Local variable — each call stack is independent, no concurrency issues.
    pending_tool_runs: list[tuple[str, RunTree]] = []

    async def _flush_text() -> None:
        """Flush accumulated text chunks into a single child run."""
        if not text_parts:
            return
        child = root.create_child(
            name="text",
            run_type="llm",
            inputs={},
            outputs={"text": "".join(text_parts)},
        )
        child.end()
        await asyncio.to_thread(child.post)
        text_parts.clear()

    try:
        async for event in events:
            etype = event.type
            content = event.content

            if etype == EventType.INIT.value:
                # Create a visible child run for adapter init.
                init_outputs: dict = {}
                cli_sid = content.get("cli_session_id", "")
                agent_type = content.get("agent_type", "")
                if cli_sid:
                    init_outputs["cli_session_id"] = cli_sid
                if agent_type:
                    init_outputs["agent_type"] = agent_type
                init_run = root.create_child(
                    name="init",
                    run_type="chain",
                    inputs={},
                    outputs=init_outputs,
                )
                init_run.end()
                await asyncio.to_thread(init_run.post)

            elif etype == EventType.TEXT.value:
                # Aggregate — do NOT create one child run per chunk.
                text = content.get("text", "")
                text_parts.append(text)
                output_parts.append(text)

            elif etype == EventType.TOOL_CALL.value:
                await _flush_text()
                tool_name = content.get("tool", "unknown")
                tool_run = root.create_child(
                    name=f"tool:{tool_name}",
                    run_type="tool",
                    inputs={"args": content.get("args", {})},
                )
                await asyncio.to_thread(tool_run.post)
                pending_tool_runs.append((tool_name, tool_run))

            elif etype == EventType.TOOL_RESULT.value:
                tool_name = content.get("tool", "unknown")
                # LIFO match: same tool may be called multiple times.
                matched_idx = None
                for i in range(len(pending_tool_runs) - 1, -1, -1):
                    if pending_tool_runs[i][0] == tool_name:
                        matched_idx = i
                        break
                if matched_idx is not None:
                    _, tool_run = pending_tool_runs.pop(matched_idx)
                    tool_run.end(outputs={"result": content.get("result", "")})
                    await asyncio.to_thread(tool_run.patch)
                else:
                    logger.warning("TOOL_RESULT for unmatched tool: %s", tool_name)

            elif etype == EventType.DONE.value:
                await _flush_text()
                usage = content.get("usage", {})
                if usage:
                    done_run = root.create_child(
                        name="done",
                        run_type="chain",
                        inputs={},
                        outputs={"usage": usage},
                    )
                    done_run.end()
                    await asyncio.to_thread(done_run.post)

            elif etype == EventType.ERROR.value:
                await _flush_text()
                root.error = content.get("error", "unknown error")

            yield event

    except Exception as e:
        await _flush_text()
        root.error = str(e)
        raise
    finally:
        # Clean up any remaining pending tool runs.
        for tool_name, tool_run in pending_tool_runs:
            tool_run.end(outputs={}, error=f"stream ended before result: {tool_name}")
            await asyncio.to_thread(tool_run.patch)

        root.end(
            outputs={
                "status": "completed" if not root.error else "error",
                "output": "".join(output_parts),
            }
        )
        await asyncio.to_thread(root.patch)
