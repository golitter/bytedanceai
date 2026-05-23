import logging

import aiomysql

logger = logging.getLogger(__name__)


class DBReader:
    """Read-only DB accessor for inactive cleanup queries."""

    def __init__(self, host: str, port: int, user: str, password: str, db: str) -> None:
        self._host = host
        self._port = port
        self._user = user
        self._password = password
        self._db = db
        self._pool: aiomysql.Pool | None = None

    async def connect(self) -> None:
        self._pool = await aiomysql.create_pool(
            host=self._host,
            port=self._port,
            user=self._user,
            password=self._password,
            db=self._db,
            autocommit=True,
            minsize=1,
            maxsize=2,
        )

    async def close(self) -> None:
        if self._pool:
            self._pool.close()
            await self._pool.wait_closed()

    async def query_inactive_sessions(self) -> list[tuple[str, str]]:
        """Return (session_id, task_id) pairs where status == 'inactive'."""
        assert self._pool is not None
        async with self._pool.acquire() as conn:
            async with conn.cursor() as cur:
                await cur.execute("SELECT session_id, task_id FROM sessions WHERE status = 'inactive'")
                rows = await cur.fetchall()
                return [(r[0], r[1]) for r in rows]

    async def query_task_session_statuses(self) -> dict[str, set[str]]:
        """Return {task_id: {status1, status2, ...}} for all sessions."""
        assert self._pool is not None
        async with self._pool.acquire() as conn:
            async with conn.cursor() as cur:
                await cur.execute("SELECT task_id, status FROM sessions")
                rows = await cur.fetchall()
                result: dict[str, set[str]] = {}
                for task_id, status in rows:
                    result.setdefault(task_id, set()).add(status)
                return result
