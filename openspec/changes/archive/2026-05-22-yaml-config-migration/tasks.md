## 1. 创建 config.yaml

- [x] 1.1 创建 `agentend/config.yaml`，包含 server、app、cli、workspace、session、execution、skills 七个分区的完整配置
- [x] 1.2 创建 `agentend/config.example.yaml` 作为模板参考

## 2. 重构 config.py

- [x] 2.1 定义嵌套 pydantic model：`ServerConfig`、`CorsConfig`、`AppConfig`、`CliConfig`、`WorkspaceConfig`、`SessionConfig`、`ExecutionConfig`、`SkillsConfig`
- [x] 2.2 顶层 `Settings` 聚合所有子配置，`model_config` 设置 `yaml_file` 指向 `../../config.yaml`（相对于 config.py）
- [x] 2.3 移除所有字段的硬编码默认值，确保缺失时报错

## 3. 更新 app/main.py

- [x] 3.1 FastAPI 实例的 title、version 改为读取 `settings.app.title`、`settings.app.version`
- [x] 3.2 CORS 中间件参数改为读取 `settings.server.cors`
- [x] 3.3 uvicorn 启动参数 host、port、reload 改为读取 `settings.server`
- [x] 3.4 workspace TTL check_interval 改为读取 `settings.workspace.ttl_check_interval`

## 4. 更新 workspace/store.py

- [x] 4.1 移除 `_DEFAULT_STORE_PATH` 硬编码常量
- [x] 4.2 `JsonFileWorkspaceStore.__init__` 中默认路径改为读取 `settings.workspace.store_path`

## 5. 更新 session/store.py

- [x] 5.1 移除 `_DEFAULT_STORE_PATH` 硬编码常量
- [x] 5.2 `SessionMappingStore.__init__` 中默认路径改为读取 `settings.session.store_path`

## 6. 更新 adapters

- [x] 6.1 `adapters/claude.py` 的 `process.terminate()` 超时改为读取 `settings.execution.process_terminate_timeout`
- [x] 6.2 `adapters/opencode.py` 的 `process.terminate()` 超时改为读取 `settings.execution.process_terminate_timeout`

## 7. 更新 skills/provisioner.py

- [x] 7.1 移除 `_BUILTIN_SKILLS_DIR` 和 `_MANIFEST_PATH` 硬编码
- [x] 7.2 builtin 目录改为读取 `settings.skills.builtin_dir`
- [x] 7.3 manifest 数据改为读取 `settings.skills.manifest`，删除 `_load_manifest` 函数

## 8. 清理

- [x] 8.1 删除 `agentend/src/skills/builtin/manifest.yaml`
- [x] 8.2 检查所有文件，确认不再有硬编码的配置值（排除不纳入范围的常量）
