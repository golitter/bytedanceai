## Context

`agentend` 的配置参数分布在多个文件中：`app/config.py`（pydantic-settings，环境变量驱动）、`workspace/store.py` 和 `session/store.py`（硬编码路径常量）、`app/main.py`（CORS、app 元信息硬编码）、`adapters/claude.py` / `adapters/opencode.py`（超时常量）、`skills/builtin/manifest.yaml`（独立 YAML 文件）。运维修改参数需要改动多处源码，缺乏统一入口。

## Goals / Non-Goals

**Goals:**
- 将所有可配置参数统一收敛到 `agentend/config.yaml`
- 保留 `pydantic-settings` 做类型校验，数据源从环境变量切换为 YAML
- 将 `skills/builtin/manifest.yaml` 合并进 `config.yaml`
- 启动时配置文件缺失或字段校验失败，立即报错退出

**Non-Goals:**
- 不改变 git branch prefix、`_DANGEROUS_TOOLS`、`SAFETY_SYSTEM_PROMPT` 等代码逻辑常量
- 不支持环境变量覆盖（纯 YAML 驱动）
- 不支持热加载 / 运行时配置变更
- 不改变 `AGENT_CONFIG_DIRS` 的定义位置

## Decisions

### 1. 使用 pydantic-settings 的 `yaml_file` 加载 YAML

**选择**: `pydantic-settings` 原生支持 `yaml_file` 配置源，在 `model_config` 中指定 YAML 文件路径即可。

**替代方案**: 手动 `yaml.safe_load()` + `BaseModel`。放弃，因为 pydantic-settings 原生支持 YAML 且提供更好的错误提示。

**实现**: 定义嵌套 model（`ServerConfig`、`CliConfig`、`WorkspaceConfig` 等），顶层 `Settings` 聚合所有子配置。

### 2. YAML 文件路径解析

**选择**: 使用 `Path(__file__)` 相对于 `config.py` 定位 `../../config.yaml`（即 `agentend/config.yaml`）。

**理由**: 不依赖 CWD，无论从哪个目录启动都能正确找到配置文件。

### 3. manifest.yaml 合并进 config.yaml

**选择**: 在 `config.yaml` 中新增 `skills` 段，包含 `builtin_dir` 和 `manifest` 字段。`provisioner.py` 从 `settings` 读取 manifest 数据，不再独立加载 `manifest.yaml`。

**理由**: 消除独立配置文件，统一入口。`builtin_dir` 可配置化以适应不同部署环境。

### 4. 配置校验失败即退出

**选择**: 在 `config.py` 模块加载时（`settings = Settings()` 处）若文件不存在或校验失败，pydantic 会抛出 `ValidationError` / `FileNotFoundError`，应用直接崩溃退出。

**理由**: Fail-fast 原则，避免带错误配置运行产生不可预期的行为。

## Risks / Trade-offs

- **[配置文件缺失导致无法启动]** → 文档中说明 `config.yaml` 为必填文件，并在仓库中提供模板（`config.example.yaml`）供参考。
- **[pydantic-settings 版本兼容]** → 确认 `yaml_file` 特性在当前使用的 pydantic-settings 版本中可用（2.0+）。若不支持，回退到手动 `yaml.safe_load` + `BaseModel` 校验。
