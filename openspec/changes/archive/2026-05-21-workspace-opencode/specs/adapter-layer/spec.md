## MODIFIED Requirements

### Requirement: ClaudeCodeAdapter CLI subprocess management
`ClaudeCodeAdapter` SHALL 通过 `asyncio.create_subprocess_exec` 启动 Claude Code CLI 进程，MUST 支持构建以下 CLI 参数：`-p`（prompt）、`--session`（会话复用）、`--output-format`（输出格式）、`--system-prompt`、`--append-system-prompt`、`--allowedTools`、`--max-turns`。进程 SHALL 支持 `cwd` 参数绑定到指定 workspace 目录。

#### Scenario: Build command with all parameters
- **WHEN** 调用 `_build_command` 并传入 session_id、message、system_prompt_append、allowed_tools、max_turns
- **THEN** SHALL 返回包含所有对应 CLI 参数的命令列表

#### Scenario: Execute with workspace cwd
- **WHEN** `stream_chat` 的 kwargs 包含 `cwd="/workspaces/task-1/frontend"`
- **THEN** SHALL 将 `cwd` 传入 `asyncio.create_subprocess_exec`，进程在 workspace 目录内执行

#### Scenario: Execute without workspace cwd
- **WHEN** kwargs 不包含 `cwd`
- **THEN** SHALL 使用默认行为（当前工作目录），不传入 cwd 参数
