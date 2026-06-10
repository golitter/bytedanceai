## 1. Workspace Core Models

- [x] 1.1 Create `agentend/src/workspace/__init__.py` with public exports (`WorkspaceManager`, `Workspace`, `WorkspaceStatus`)
- [x] 1.2 Create `agentend/src/workspace/models.py` — `Workspace` dataclass (id, task_id, agent_name, repo_path, worktree_path, branch_name, session_id?, status, created_at) + `WorkspaceStatus` enum (ACTIVE, MERGED, CLEANED) + branch name generation logic (`agent/{agent_name}/{task_id}`)

## 2. GitOps Utility

- [x] 2.1 Create `agentend/src/workspace/git_ops.py` — `GitOps` class with async methods: `worktree_add`, `worktree_remove`, `branch_create`, `add_and_commit`, `merge_branch`, `get_current_branch`，全部通过 `asyncio.create_subprocess_exec` 调用 git 命令，解析 stdout/stderr 返回 bool/str

## 3. WorkspaceManager

- [x] 3.1 Create `agentend/src/workspace/manager.py` — `WorkspaceManager` class with in-memory `_workspaces: dict[str, Workspace]` 和 `GitOps` 实例
- [x] 3.2 Implement `create(repo_path, task_id, agent_name)` — 自动生成 branch name，调用 `GitOps.worktree_add`，存储 Workspace 对象
- [x] 3.3 Implement `get(workspace_id)`, `list()`, `cleanup(workspace_id)`, `cleanup_by_task(task_id)`
- [x] 3.4 Implement `commit(workspace_id, message)` — 调用 `GitOps.add_and_commit`
- [x] 3.5 Implement `merge(workspace_id, target_branch)` — 调用 `GitOps.merge_branch`，冲突时 abort 并返回 False

## 4. OpenCode Adapter

- [x] 4.1 Create `agentend/src/adapters/opencode.py` — `OpenCodeAdapter` 继承 `BaseAgentAdapter`
- [x] 4.2 Implement `_build_command(message, cwd, system_prompt_append)` — 组装 `[OPENCODE_CLI_PATH, "-p", prompt, "-f", "json", "-q"]` + 可选 `-c cwd`，system_prompt_append 拼接到 prompt 前缀
- [x] 4.3 Implement `_parse_json_output(raw_json)` — 解析 OpenCode JSON 输出，提取 text / tool_use / usage 信息，转换为 StreamEvent 列表
- [x] 4.4 Implement `chat(session_id, message, **kwargs)` — 执行 CLI，等待完成，解析 JSON 返回 AgentResponse
- [x] 4.5 Implement `stream_chat(session_id, message, **kwargs)` — 执行 CLI，拿到完整结果后拆分为 StreamEvent 序列 yield，支持 cwd 绑定
- [x] 4.6 Implement `create_session`, `interrupt`, `destroy_session` — 与 ClaudeCodeAdapter 相同模式
- [x] 4.7 Update `agentend/src/adapters/__init__.py` — 导出 `OpenCodeAdapter`

## 5. Modify Existing Components

- [x] 5.1 Update `agentend/src/app/config.py` — 新增 `OPENCODE_CLI_PATH` 配置项（默认 `"opencode"`），新增 `WORKSPACE_BASE_DIR` 配置项（默认 `"./worktrees"`）
- [x] 5.2 Update `agentend/src/session/models.py` — `workspace_path` 从 `str | None` 改为必填 `str`
- [x] 5.3 Update `agentend/src/session/manager.py` — `create` 方法接受 `workspace_path: str` 必填参数
- [x] 5.4 Update `agentend/src/adapters/claude.py` — `stream_chat` 和 `chat` 从 kwargs 取 `cwd` 传入 `asyncio.create_subprocess_exec`
- [x] 5.5 Update `agentend/src/schemas/request.py` — 新增 `repo_path: str | None` 字段
- [x] 5.6 Update `agentend/src/schemas/events.py` — StreamEvent content 中包含 `agent_type` 字段（由 Adapter 注入）

## 6. API Endpoints

- [x] 6.1 Create `agentend/src/api/v1/workspace.py` — workspace CRUD 端点：`POST /v1/workspace/create`、`POST /v1/workspace/{id}/commit`、`POST /v1/workspace/{id}/merge`、`DELETE /v1/workspace/{id}`、`GET /v1/workspace`
- [x] 6.2 Update `agentend/src/api/v1/agent.py` — agent 执行端点增加 workspace 逻辑：检查 workspace_path/repo_path，自动创建 workspace，将 cwd 传入 Adapter kwargs
- [x] 6.3 Update `agentend/src/api/__init__.py` — 导出 workspace router

## 7. App Wiring

- [x] 7.1 Update `agentend/src/app/dependencies.py` — 新增 `create_workspace_manager` DI 工厂，`create_adapter_registry` 注册 OpenCodeAdapter
- [x] 7.2 Update `agentend/src/app/main.py` — lifespan 中初始化 WorkspaceManager 并存储到 app.state，注册 workspace router，shutdown 时 cleanup 所有 workspace

## 8. Tests

- [x] 8.1 Create `agentend/tests/test_workspace.py` — 测试 WorkspaceManager create/get/list/cleanup/commit/merge，测试 GitOps 原子操作（mock subprocess）
- [x] 8.2 Create `agentend/tests/test_opencode_adapter.py` — 测试 OpenCodeAdapter._build_command 参数组装（prompt 拼接约束、cwd 绑定）、_parse_json_output 解析
- [x] 8.3 Update `agentend/tests/test_adapter.py` — 测试 ClaudeCodeAdapter cwd 参数传入 subprocess
- [x] 8.4 Update `agentend/tests/test_api.py` — 测试 workspace CRUD 端点、agent 执行时 workspace 自动创建逻辑
