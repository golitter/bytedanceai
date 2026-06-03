import asyncio
import json
from pathlib import Path

from src.app.config import settings


class SessionMappingStore:
    def __init__(self, path: Path | None = None) -> None:
        # 默认路径来自 config.yaml 的 session.store_path
        self._path = path or Path(settings.session.store_path)
        self._mappings: dict[str, str] = {}
        self._load()

    def _load(self) -> None:
        if self._path.exists():
            try:
                self._mappings = json.loads(self._path.read_text())
            except (json.JSONDecodeError, OSError):
                self._mappings = {}
        else:
            self._mappings = {}

    def _save(self) -> None:
        self._path.parent.mkdir(parents=True, exist_ok=True)
        self._path.write_text(json.dumps(self._mappings, indent=2))

    @staticmethod
    def _key(session_id: str, task_id: str) -> str:
        return f"{session_id}::{task_id}"

    def get_cli_session_id(self, session_id: str, task_id: str = "") -> str | None:
        return self._mappings.get(self._key(session_id, task_id))

    async def set_cli_session_id(self, session_id: str, cli_session_id: str, task_id: str = "") -> None:
        self._mappings[self._key(session_id, task_id)] = cli_session_id
        await asyncio.to_thread(self._save)

    async def delete(self, session_id: str, task_id: str = "") -> None:
        self._mappings.pop(self._key(session_id, task_id), None)
        await asyncio.to_thread(self._save)
