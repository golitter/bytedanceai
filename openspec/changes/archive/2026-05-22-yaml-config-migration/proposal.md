## Why

`agentend/src` 中的配置参数散落在多个文件里（`app/config.py`、`workspace/store.py`、`session/store.py`、`app/main.py`、`adapters/claude.py` 等），默认值硬编码在 Python 代码中。部署和运维时需要修改参数必须改源码，缺乏统一的配置入口。此外 `skills/builtin/manifest.yaml` 作为独立文件存在，增加了配置碎片化。

## What Changes

- 创建 `agentend/config.yaml` 作为唯一的配置文件，包含所有可配置参数
- 重构 `app/config.py`，基于 `pydantic-settings` 的 `yaml_file` 从 YAML 读取配置，移除硬编码默认值
- 将 `workspace/store.py`、`session/store.py` 中的默认存储路径改为从配置读取
- 将 `app/main.py` 中的 CORS、app title/version、reload 等参数改为从配置读取
- 将 `adapters/claude.py` 和 `adapters/opencode.py` 中的 `process_terminate_timeout` 改为从配置读取
- 将 `skills/builtin/manifest.yaml` 的内容合并进 `config.yaml`，`provisioner.py` 从配置读取 manifest 数据
- **BREAKING**: 删除 `skills/builtin/manifest.yaml` 文件
- 启动时若 `config.yaml` 不存在或配置项缺失，直接报错退出

## Capabilities

### New Capabilities
- `yaml-config`: 统一的 YAML 配置加载机制，涵盖 server、app、cli、workspace、session、execution、skills 七个配置分区

### Modified Capabilities

## Impact

- **代码**: `app/config.py`（重写）、`app/main.py`、`workspace/store.py`、`session/store.py`、`adapters/claude.py`、`adapters/opencode.py`、`skills/provisioner.py`
- **文件**: 新增 `agentend/config.yaml`，删除 `agentend/src/skills/builtin/manifest.yaml`
- **依赖**: 无新增依赖（`pydantic-settings` 已在项目中）
- **部署**: 需要确保 `config.yaml` 文件存在于正确路径
