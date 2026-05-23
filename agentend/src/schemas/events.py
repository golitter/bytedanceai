import time

from pydantic import Field

from src.generated.events import EventType
from src.generated.events import StreamEvent as _StreamEvent


class StreamEvent(_StreamEvent):
    type: str
    content: dict = Field(default_factory=dict)
    timestamp: float = Field(default_factory=time.time)

    @staticmethod
    def create(event_type: EventType, agent_type: str | None = None, **kwargs) -> "StreamEvent":
        if agent_type:
            kwargs["agent_type"] = agent_type
        return StreamEvent(type=event_type.value, content=kwargs)


__all__ = ["EventType", "StreamEvent"]
