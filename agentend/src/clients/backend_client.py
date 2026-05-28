from __future__ import annotations

import json
import logging
from collections.abc import AsyncIterator

import httpx

logger = logging.getLogger(__name__)

_SSE_READ_TIMEOUT = 600.0


class BackendClient:
    """HTTP client for calling backend API (RunTask, SSE stream)."""

    def __init__(self, base_url: str, timeout: float = 30.0) -> None:
        self._base_url = base_url.rstrip("/")
        # run_task 用短超时; stream_result 用独立的 stream client
        self._client = httpx.AsyncClient(timeout=timeout)

    async def close(self) -> None:
        await self._client.aclose()

    async def run_task(
        self,
        task_id: str,
        session_id: str,
        message: str,
        agent_type: str,
        cwd: str = "",
    ) -> str:
        """POST /api/tasks/:taskId/run → returns message_id."""
        resp = await self._client.post(
            f"{self._base_url}/api/tasks/{task_id}/run",
            json={
                "message": message,
                "session_id": session_id,
                "agent_type": agent_type,
                "cwd": cwd,
                "skip_user_message": True,
            },
        )
        resp.raise_for_status()
        body = resp.json()
        data = body.get("data", body)
        message_id = data.get("message_id")
        if not message_id:
            raise ValueError(f"RunTask response missing message_id: {body}")
        logger.info(
            "BackendClient.run_task: task=%s session=%s agent=%s → message_id=%s",
            task_id,
            session_id,
            agent_type,
            message_id,
        )
        return message_id

    async def stream_result(
        self,
        task_id: str,
        message_id: str,
        session_id: str,
    ) -> AsyncIterator[dict]:
        """Subscribe to SSE GET /api/tasks/:taskId/stream and yield parsed events.

        Uses a dedicated httpx client with long read timeout to handle
        the long-lived SSE connection, and aiter_text() with manual line
        splitting for robust SSE parsing.
        """
        url = f"{self._base_url}/api/tasks/{task_id}/stream"
        params = {"message_id": message_id, "session_id": session_id}
        logger.info("BackendClient.stream_result: connecting %s params=%s", url, params)

        # Dedicated client with long read timeout for SSE
        sse_client = httpx.AsyncClient(
            timeout=httpx.Timeout(connect=10.0, read=_SSE_READ_TIMEOUT, write=10.0, pool=10.0),
        )
        try:
            async with sse_client.stream("GET", url, params=params) as resp:
                logger.info("BackendClient.stream_result: status=%d", resp.status_code)
                resp.raise_for_status()

                buf = ""
                async for text_chunk in resp.aiter_text():
                    buf += text_chunk
                    while "\n" in buf:
                        line, buf = buf.split("\n", 1)
                        line = line.strip()
                        if not line or not line.startswith("data: "):
                            continue
                        payload = line[6:]  # strip "data: " prefix
                        try:
                            event = json.loads(payload)
                        except json.JSONDecodeError:
                            logger.warning("BackendClient.stream_result: JSON parse error: %s", payload[:120])
                            continue
                        logger.debug(
                            "BackendClient.stream_result: event type=%s",
                            event.get("type"),
                        )
                        yield event
        finally:
            await sse_client.aclose()
