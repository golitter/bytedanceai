from pathlib import Path

from pydantic import BaseModel
from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic_settings.sources import YamlConfigSettingsSource

# 基于 __file__ 定位 config.yaml，不依赖 CWD，任何目录启动都能找到
_CONFIG_PATH = Path(__file__).resolve().parent.parent.parent / "config.yaml"


class CorsConfig(BaseModel):
    origins: list[str]
    credentials: bool
    methods: list[str]
    headers: list[str]


class ServerConfig(BaseModel):
    host: str
    port: int
    reload: bool
    cors: CorsConfig


class AppConfig(BaseModel):
    title: str
    version: str


class CliConfig(BaseModel):
    claude_path: str
    opencode_path: str


class WorkspaceConfig(BaseModel):
    base_dir: str
    ttl_seconds: int
    ttl_check_interval: int
    store_path: str
    git_default_branch: str


class SessionConfig(BaseModel):
    store_path: str


class ExecutionConfig(BaseModel):
    max_turns: int
    timeout: int
    process_terminate_timeout: float


class SkillsConfig(BaseModel):
    builtin_dir: str
    manifest: dict


class Settings(BaseSettings):
    server: ServerConfig
    app: AppConfig
    cli: CliConfig
    workspace: WorkspaceConfig
    session: SessionConfig
    execution: ExecutionConfig
    skills: SkillsConfig

    model_config = SettingsConfigDict(yaml_file=str(_CONFIG_PATH))

    @classmethod
    def settings_customise_sources(
        cls,
        settings_cls,
        init_settings,
        env_settings,
        dotenv_settings,
        file_secret_settings,
    ):
        # 纯 YAML 驱动，禁用环境变量 / dotenv / file secret 等默认源
        return (init_settings, YamlConfigSettingsSource(settings_cls))


# 模块加载时立即实例化，config.yaml 缺失或字段校验失败会直接崩溃退出（fail-fast）
settings = Settings()
