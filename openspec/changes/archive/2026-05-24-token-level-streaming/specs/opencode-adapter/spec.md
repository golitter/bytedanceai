## MODIFIED Requirements

### Requirement: OpenCodeAdapter CLI command building
`OpenCodeAdapter` SHALL 通过 `_build_command` 方法组装 OpenCode CLI 命令，支持参数：`run`（子命令）、`--format json`（JSON 输出）、`--dir`（工作目录）、`--session`（会话 ID）、`--model`（模型选择）、`--agent`（agent 选择）。MUST 支持从 config 读取 CLI 路径。

#### Scenario: Build command with prompt and cwd
- **WHEN** 调用 `_build_command("修复 bug", cwd="/workspaces/task-1")`
- **THEN** SHALL 返回 `[OPENCODE_CLI_PATH, "run", "修复 bug", "--format", "json", "--dir", "/workspaces/task-1"]`

#### Scenario: Build command with rule constraint in prompt
- **WHEN** 调用 `_build_command("修复 bug", system_prompt_append="只改 tsx 文件")`
- **THEN** prompt SHALL 拼接为 `"[系统约束: 只改 tsx 文件]\n\n修复 bug"`

### Requirement: OpenCodeAdapter stream_chat
`OpenCodeAdapter.stream_chat` SHALL 执行 `opencode run "..." --format json`，通过 `async for line in process.stdout` 逐行读取 NDJSON 输出，解析每行为 StreamEvent 并立即 yield。

#### Scenario: Stream token-level text events
- **WHEN** OpenCode CLI 逐行输出 `{"type": "text", "part": {"text": "创建..."}}`
- **THEN** SHALL 逐行 yield `StreamEvent(type=TEXT, content={"text": "创建..."})`

#### Scenario: Stream with tool use events
- **WHEN** OpenCode CLI 输出 `{"type": "tool_use", "part": {"tool": "Bash", "state": {"input": {"command": "ls"}}}}`
- **THEN** SHALL yield `StreamEvent(type=TOOL_CALL, content={"tool": "Bash", "args": {"command": "ls"}})`

#### Scenario: Stream completion
- **WHEN** CLI 进程正常退出
- **THEN** SHALL yield `StreamEvent(type=DONE)`

#### Scenario: Stream error
- **WHEN** CLI 进程以非零退出码退出
- **THEN** SHALL yield `StreamEvent(type=ERROR, content={"error": stderr})`
