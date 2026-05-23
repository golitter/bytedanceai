from pydantic import Field

from src.generated.response import AgentResponse as _AgentResponse


class AgentResponse(_AgentResponse):
    artifacts: list[dict] = Field(default_factory=list)
    usage: dict = Field(default_factory=dict)


__all__ = ["AgentResponse"]
