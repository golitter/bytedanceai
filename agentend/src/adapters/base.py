from abc import ABC, abstractmethod
from collections.abc import AsyncIterator

from src.schemas.events import StreamEvent
from src.schemas.response import AgentResponse


class BaseAgentAdapter(ABC):
    @abstractmethod
    async def create_session(self, session_id: str) -> None: ...

    @abstractmethod
    async def chat(self, session_id: str, message: str, **kwargs) -> AgentResponse: ...

    @abstractmethod
    async def stream_chat(self, session_id: str, message: str, **kwargs) -> AsyncIterator[StreamEvent]: ...

    @abstractmethod
    async def interrupt(self, session_id: str) -> bool: ...

    @abstractmethod
    async def destroy_session(self, session_id: str) -> None: ...
