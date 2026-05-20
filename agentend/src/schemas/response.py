from pydantic import BaseModel, Field


class AgentResponse(BaseModel):
    session_id: str
    content: str
    artifacts: list[dict] = Field(default_factory=list)
    usage: dict = Field(default_factory=dict)
