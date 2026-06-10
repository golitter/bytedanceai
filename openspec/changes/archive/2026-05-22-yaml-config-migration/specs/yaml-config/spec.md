## ADDED Requirements

### Requirement: YAML 配置文件为应用唯一配置源
系统 SHALL 从 `agentend/config.yaml` 读取所有可配置参数。系统 SHALL NOT 从环境变量或硬编码默认值读取配置参数。

#### Scenario: 配置文件存在且合法
- **WHEN** `config.yaml` 文件存在且内容符合 schema 定义
- **THEN** 系统正常加载所有配置并启动

#### Scenario: 配置文件不存在
- **WHEN** `config.yaml` 文件不存在
- **THEN** 系统 SHALL 在启动时报错退出，输出明确的错误信息提示文件缺失

#### Scenario: 配置字段缺失或类型错误
- **WHEN** `config.yaml` 文件存在但缺少必填字段或字段类型不匹配
- **THEN** 系统 SHALL 在启动时报错退出，输出 pydantic 校验错误详情

### Requirement: 配置结构包含七个分区
配置文件 SHALL 包含以下七个顶层分区：`server`、`app`、`cli`、`workspace`、`session`、`execution`、`skills`。

#### Scenario: 配置分区完整
- **WHEN** `config.yaml` 包含所有七个分区及其必填字段
- **THEN** 系统正确解析各分区配置

### Requirement: server 分区配置
`server` 分区 SHALL 包含 `host`（str）、`port`（int）、`reload`（bool）和 `cors` 子对象。`cors` 子对象 SHALL 包含 `origins`（list[str]）、`credentials`（bool）、`methods`（list[str]）、`headers`（list[str]）。

#### Scenario: server 配置用于 uvicorn 启动
- **WHEN** 应用启动
- **THEN** uvicorn 使用 `server.host` 和 `server.port` 绑定地址，使用 `server.reload` 控制热加载

#### Scenario: CORS 配置用于中间件
- **WHEN** FastAPI 应用初始化 CORS 中间件
- **THEN** 中间件参数来源于 `server.cors` 配置

### Requirement: app 分区配置
`app` 分区 SHALL 包含 `title`（str）和 `version`（str）。

#### Scenario: app 元信息配置
- **WHEN** FastAPI 实例创建
- **THEN** 使用 `app.title` 和 `app.version` 作为应用元信息

### Requirement: cli 分区配置
`cli` 分区 SHALL 包含 `claude_path`（str）和 `opencode_path`（str）。

#### Scenario: CLI 路径用于构建命令
- **WHEN** adapter 构建 CLI 命令
- **THEN** 使用 `cli.claude_path` 或 `cli.opencode_path` 作为可执行文件路径

### Requirement: workspace 分区配置
`workspace` 分区 SHALL 包含 `base_dir`（str）、`ttl_seconds`（int）、`ttl_check_interval`（int）、`store_path`（str）、`git_default_branch`（str）。

#### Scenario: workspace 存储路径配置
- **WHEN** workspace store 初始化
- **THEN** 使用 `workspace.store_path` 作为 JSON 持久化文件路径

#### Scenario: workspace TTL 配置
- **WHEN** TTL 清理任务启动
- **THEN** 使用 `workspace.ttl_seconds` 和 `workspace.ttl_check_interval` 控制清理策略

### Requirement: session 分区配置
`session` 分区 SHALL 包含 `store_path`（str）。

#### Scenario: session 存储路径配置
- **WHEN** session store 初始化
- **THEN** 使用 `session.store_path` 作为 JSON 持久化文件路径

### Requirement: execution 分区配置
`execution` 分区 SHALL 包含 `max_turns`（int）、`timeout`（int）、`process_terminate_timeout`（float）。

#### Scenario: 执行超时配置
- **WHEN** agent 执行任务
- **THEN** 使用 `execution.timeout` 作为最大执行时间，`execution.max_turns` 作为最大轮次

#### Scenario: 进程终止超时配置
- **WHEN** adapter 终止子进程
- **THEN** 等待 `execution.process_terminate_timeout` 秒后强制 kill

### Requirement: skills 分区配置
`skills` 分区 SHALL 包含 `builtin_dir`（str）和 `manifest`（dict）。

#### Scenario: manifest 数据从配置读取
- **WHEN** provisioner 加载 skill manifest
- **THEN** 从 `skills.manifest` 读取 skill 清单数据，不再读取独立 manifest.yaml 文件

#### Scenario: builtin 目录可配置
- **WHEN** provisioner 查找 builtin skills 目录
- **THEN** 使用 `skills.builtin_dir` 定位目录

### Requirement: 删除独立 manifest.yaml
系统 SHALL 移除 `skills/builtin/manifest.yaml` 文件。所有 manifest 数据 SHALL 从 `config.yaml` 的 `skills.manifest` 字段读取。

#### Scenario: manifest.yaml 不再被引用
- **WHEN** provisioner 初始化
- **THEN** 不访问 `skills/builtin/manifest.yaml`，所有 manifest 数据来源于 config
