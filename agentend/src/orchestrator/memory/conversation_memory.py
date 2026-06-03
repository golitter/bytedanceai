"""Persist Orchestrator's memory_messages across conversation turns.

Follows the same file-based pattern as EvolutionStore and PinMemory.
Stores serialized LangChain messages in JSON format within shared_dir.
"""

from __future__ import annotations

import json
import logging
from pathlib import Path

from langchain_core.messages import messages_from_dict, messages_to_dict

logger = logging.getLogger(__name__)

_MAX_TURNS = 10


class ConversationMemoryStore:
    """Persist Orchestrator's memory_messages (including dynamic context messages).

    Storage location: ``{shared_dir}/memory/conversation_memory.json``

    Uses ``langchain_core.messages.messages_to_dict`` /
    ``messages_from_dict`` for lossless serialization of HumanMessage,
    AIMessage (with tool_calls), ToolMessage, and SystemMessage.
    """

    def __init__(self, shared_dir: str | Path) -> None:
        self.memory_dir = Path(shared_dir) / "memory"
        self.memory_dir.mkdir(parents=True, exist_ok=True)

    @property
    def memory_path(self) -> Path:
        return self.memory_dir / "conversation_memory.json"

    # ------------------------------------------------------------------
    # Public API
    # ------------------------------------------------------------------

    def save_messages(self, messages: list) -> None:
        """Serialize *messages*, append to file, trim to retention limit."""
        existing = self._load_raw()
        new_entries = messages_to_dict(messages)
        combined = existing + new_entries
        trimmed = self._trim_to_turns(combined, _MAX_TURNS)
        self._write(trimmed)

    def replace_messages(self, messages: list) -> None:
        """Replace the store with exactly *messages* (after trimming).

        Unlike :meth:`save_messages` which reads existing entries and appends,
        this method writes *messages* directly — no duplication.
        """
        entries = messages_to_dict(messages)
        trimmed = self._trim_to_turns(entries, _MAX_TURNS)
        self._write(trimmed)

    def load_messages(self) -> list:
        """Deserialize stored messages back into LangChain message objects."""
        raw = self._load_raw()
        if not raw:
            return []
        try:
            return messages_from_dict(raw)
        except Exception:
            logger.warning("Failed to deserialize conversation memory; starting fresh", exc_info=True)
            return []

    # ------------------------------------------------------------------
    # Internal helpers
    # ------------------------------------------------------------------

    def _load_raw(self) -> list[dict]:
        if not self.memory_path.exists():
            return []
        try:
            data = self.memory_path.read_text(encoding="utf-8")
            return json.loads(data)
        except (json.JSONDecodeError, OSError):
            logger.warning("Failed to read conversation memory file; treating as empty", exc_info=True)
            return []

    def _write(self, entries: list[dict]) -> None:
        self.memory_path.write_text(
            json.dumps(entries, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

    @staticmethod
    def _trim_to_turns(entries: list[dict], max_turns: int) -> list[dict]:
        """Keep the last *max_turns* complete turns.

        A turn starts wherever ``type == "human"`` appears.  We never
        split a turn — we keep every message from the chosen start index
        to the end.
        """
        human_indices = [i for i, e in enumerate(entries) if e.get("type") == "human"]
        if len(human_indices) <= max_turns:
            return entries
        start = human_indices[-max_turns]
        return entries[start:]
