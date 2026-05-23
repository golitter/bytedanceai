from pydantic import Field

from src.generated.request import AgentRequest as _AgentRequest
from src.generated.request import AgentType


class AgentRequest(_AgentRequest):
    rules: list[str] = Field(default_factory=list)
    config: dict | None = None


__all__ = ["AgentType", "AgentRequest"]
