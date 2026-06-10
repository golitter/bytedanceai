## Why

AgentHub 多 Agent 协作平台已实现 MVP（Claude Code Adapter + Session + Rule Engine + SSE 流式），但存在两个关键缺失：
1. **无工作空间隔离** — 多个 Agent 同时操作同一仓库会导致文件冲突、上下文污染、代码回滚困难
2. **仅支持 Claude Code** — 无法对接 OpenCode 等其他 Coding Agent，不符合"统一适配"目标

文档 2 明确指出 Workspace 隔离是"最容易炸的地方"，也是 AI Coding Agent 平台的核心基础设施。OpenCode 作为开源终端 AI Coding Agent，支持多模型提供商（Claude/GPT/Gemini），其 CLI 非交互模式（`opencode -p "prompt" -f json`）与 Claude Code 类似，可作为第二个 Adapter 接入。

## What Changes

- 新增 `WorkspaceManager`：基于 git worktree 为每个 Agent 任务创建独立工作空间，管理 branch 生命周期（创建 → 使用 → 提交 → 合并/清理）
- 新增 `GitOps` 工具类：封装 git worktree add/remove、branch create、commit、merge 等原子操作
- 新增 `OpenCodeAdapter`：基于 CLI subprocess 调用 OpenCode 非交互模式（`opencode -p ... -f json`），继承 `BaseAgentAdapter` 接口
- 修改 `Session` 模型：关联 `workspace_path`，Agent 执行时自动绑定到对应 worktree 目录
- 修改 `ClaudeCodeAdapter`：支持 `cwd` 参数，执行时切换到指定 workspace 目录
- 修改 API 端点：`AgentRequest` 中 `workspace_path` 字段启用，创建 workspace 并传入 Adapter
- 修改 `app/dependencies.py`：注册 `WorkspaceManager` 到 DI 容器

## Capabilities

### New Capabilities

- `workspace-isolation`: Git worktree 工作空间隔离，为每个 Agent 任务创建独立 branch + 目录，支持创建、查询、提交、合并、清理全生命周期
- `opencode-adapter`: OpenCode CLI Adapter，通过非交互模式（`opencode -p ... -f json`）调用 OpenCode Agent，支持同步/流式（基于 JSON 输出分块模拟）执行

### Modified Capabilities

- `adapter-layer`: `BaseAgentAdapter` 接口增加 `cwd` 参数支持（通过 kwargs），`ClaudeCodeAdapter` 执行时绑定 workspace 目录
- `session-manager`: `Session` 模型新增 `workspace_path` 必填关联，`SessionManager.create` 接受 `workspace_path` 参数
- `stream-protocol`: `AgentRequest` 的 `workspace_path` 字段从可选变为必填（配合 workspace 创建流程），新增 `WorkspaceInfo` 响应模型
- `api-gateway`: Agent 执行端点增加 workspace 自动创建/绑定逻辑，新增 `POST /v1/workspace/create` 和 `POST /v1/workspace/{id}/merge` 端点

## Impact

- **新增模块**: `src/workspace/` 目录（manager.py, git_ops.py, models.py）
- **新增文件**: `src/adapters/opencode.py`
- **修改文件**: `src/adapters/claude.py`（cwd 支持）、`src/session/models.py`（workspace 关联）、`src/session/manager.py`（create 参数）、`src/schemas/request.py`（字段调整）、`src/api/v1/agent.py`（workspace 绑定）、`src/app/main.py` + `src/app/dependencies.py`（注册新组件）
- **新增依赖**: 无（git worktree 通过标准库 subprocess 调用）
- **外部依赖**: 要求机器上已安装 `opencode` CLI 工具（`brew install opencode-ai/tap/opencode`）
- **目录变更**: 运行时在项目 repo 旁生成 `worktrees/` 目录（git worktree 管理）
