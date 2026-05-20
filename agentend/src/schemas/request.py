from pydantic import BaseModel, Field


class AgentRequest(BaseModel):
    task_id: str
    conversation_id: str
    session_id: str | None = None
    message: str
    agent_type: str = "claude-code"
    stream: bool = True
    system_prompt: str | None = None
    rules: list[str] = Field(default_factory=list)
    workspace_path: str | None = None
    config: dict | None = None
