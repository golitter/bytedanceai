from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    CLAUDE_CLI_PATH: str = "claude"
    DEFAULT_MAX_TURNS: int = 20
    EXECUTION_TIMEOUT: int = 300
    HOST: str = "0.0.0.0"
    PORT: int = 8001

    model_config = {"env_prefix": "", "case_sensitive": True}


settings = Settings()
