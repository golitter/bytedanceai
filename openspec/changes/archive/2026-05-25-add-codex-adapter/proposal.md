## Why

项目需要接入 OpenAI Codex CLI 作为第四个 Agent 后端，以满足用户对多 Agent 平台（Claude Code、OpenCode、Codex）的接入需求。Codex CLI 提供了 `codex exec --json` 非交互模式和内置沙箱机制，与现有 subprocess adapter 模式高度对齐。

## What Changes

- 契约层 `AgentType` 枚举增加 `codex` 值，三端类型自动生成
- 新建 `CodexAdapter`，实现 `BaseAgentAdapter` 接口，通过 `codex exec --json` 子进程桥接
- 配置层增加 `cli.codex_path` 字段
- 注册 `CodexAdapter` 到 `AdapterRegistry`
- Workspace 隔离排除 `~/.codex/` 配置目录

## Capabilities

### New Capabilities
- `codex-adapter`: Codex CLI adapter 实现，包含命令构建、NDJSON 事件解析（thread.started / item.started / item.completed / turn.completed）、会话恢复（exec resume）、进程中断

### Modified Capabilities
- `adapter-registry`: 注册 `CodexAdapter` 到 `AgentType.CODEX`
- `contract-schema`: `AgentType` 枚举增加 `codex` 值

## Impact

- **契约层**: `contracts/schemas/agent-request.yaml` AgentType 枚举新增值，需运行 `make generate` 更新三端生成代码
- **Agentend**: 新增 `adapters/codex.py`，修改 `config.py`、`config.yaml`、`dependencies.py`
- **Frontend**: `AgentType` 枚举自动更新，前端 Agent 选择器需支持 `codex` 类型
- **Backend**: `AgentType` 枚举自动更新
- **依赖**: 需要服务器上安装 `codex` CLI（`codex exec` 可执行）
