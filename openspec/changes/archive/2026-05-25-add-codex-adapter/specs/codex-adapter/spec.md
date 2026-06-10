## ADDED Requirements

### Requirement: CodexAdapter CLI subprocess management
`CodexAdapter` SHALL 通过 `asyncio.create_subprocess_exec` 启动 Codex CLI 进程，MUST 支持构建以下 CLI 参数：`exec --json`（非交互模式）、`-a never`（跳过审批）、`-s workspace-write`（沙箱策略）、`-C`（工作目录）、`-m`（模型覆盖）。进程 SHALL 支持 `cwd` 参数绑定到指定 workspace 目录。

#### Scenario: Build command with all parameters
- **WHEN** 调用 `_build_command` 并传入 message、cwd、model
- **THEN** SHALL 返回 `[codex_path, "exec", "--json", "-a", "never", "-s", "workspace-write", "-C", cwd, "-m", model, message]` 格式的命令列表

#### Scenario: Build command with minimal parameters
- **WHEN** 调用 `_build_command` 只传入 message
- **THEN** SHALL 返回 `[codex_path, "exec", "--json", "-a", "never", "-s", "workspace-write", message]`，不包含 `-C` 和 `-m` 参数

#### Scenario: Build resume command
- **WHEN** 调用 `_build_command` 并传入 cli_session_id 且 is_resume=True
- **THEN** SHALL 返回 `[codex_path, "exec", "resume", cli_session_id, "--json", ...]` 格式的命令列表

### Requirement: CodexAdapter NDJSON event parsing
`CodexAdapter` SHALL 将 Codex CLI 的 NDJSON stdout 逐行解析为 `StreamEvent` 对象。事件类型映射 SHALL 遵循以下规则：

| Codex 事件 | 映射目标 |
|---|---|
| `thread.started` | `INIT`（提取 `thread_id` 作为 `cli_session_id`）|
| `item.started` (type=`command_execution`) | `TOOL_CALL`（提取 `command`）|
| `item.completed` (type=`reasoning`) | `TEXT`（加 `[thinking]` 前缀）|
| `item.completed` (type=`agent_message`) | `TEXT`（正文）|
| `item.completed` (type=`command_execution`) | `TOOL_RESULT`（提取 `aggregated_output`、`exit_code`）|
| `turn.completed` | `DONE`（提取 `usage`）|

#### Scenario: Parse thread.started event
- **WHEN** CLI 输出 `{"type":"thread.started","thread_id":"abc-123"}`
- **THEN** SHALL 生成 `StreamEvent(type=INIT, content={"cli_session_id": "abc-123"})`

#### Scenario: Parse item.started command_execution event
- **WHEN** CLI 输出 `{"type":"item.started","item":{"id":"item_1","type":"command_execution","command":"ls -la","status":"in_progress"}}`
- **THEN** SHALL 生成 `StreamEvent(type=TOOL_CALL, content={"tool": "command_execution", "args": {"command": "ls -la"}})`

#### Scenario: Parse item.completed agent_message event
- **WHEN** CLI 输出 `{"type":"item.completed","item":{"id":"item_2","type":"agent_message","text":"Hello!"}}`
- **THEN** SHALL 生成 `StreamEvent(type=TEXT, content={"text": "Hello!"})`

#### Scenario: Parse item.completed reasoning event
- **WHEN** CLI 输出 `{"type":"item.completed","item":{"id":"item_0","type":"reasoning","text":"Thinking..."}}`
- **THEN** SHALL 生成 `StreamEvent(type=TEXT, content={"text": "[thinking] Thinking..."})`

#### Scenario: Parse item.completed command_execution event
- **WHEN** CLI 输出 `{"type":"item.completed","item":{"id":"item_1","type":"command_execution","command":"ls","aggregated_output":"file1.txt","exit_code":0,"status":"completed"}}`
- **THEN** SHALL 生成 `StreamEvent(type=TOOL_RESULT, content={"tool": "command_execution", "result": "file1.txt", "exit_code": 0})`

#### Scenario: Parse turn.completed event
- **WHEN** CLI 输出 `{"type":"turn.completed","usage":{"input_tokens":100,"output_tokens":50}}`
- **THEN** SHALL 生成 `StreamEvent(type=DONE, content={"usage": {"input_tokens": 100, "output_tokens": 50}})`

#### Scenario: Ignore unknown event types
- **WHEN** CLI 输出无法识别的事件类型（如 `turn.started`）或无法解析为 JSON 的行
- **THEN** SHALL 返回 None，不 SHALL 抛出异常

### Requirement: CodexAdapter session resume via exec resume
`CodexAdapter` SHALL 支持通过 `codex exec resume <thread_id>` 恢复之前的会话。首次执行时从 `thread.started` 事件捕获 `thread_id`，后续通过 `is_resume=True` 触发恢复命令。

#### Scenario: First execution captures thread_id
- **WHEN** 首次执行 `stream_chat` 且无 `cli_session_id`
- **THEN** SHALL 使用 `codex exec` 命令，从 `thread.started` 事件中捕获 `thread_id`

#### Scenario: Resume existing session
- **WHEN** 执行 `stream_chat` 且 kwargs 包含 `cli_session_id` 和 `is_resume=True`
- **THEN** SHALL 使用 `codex exec resume <cli_session_id>` 命令

### Requirement: CodexAdapter process interrupt
`CodexAdapter.interrupt` SHALL 终止指定 session 的正在运行的 CLI 进程。MUST 先发送 SIGTERM，等待 `settings.execution.process_terminate_timeout` 秒，若进程未结束则发送 SIGKILL。

#### Scenario: Interrupt running process
- **WHEN** session 有正在运行的进程，调用 `interrupt(session_id)`
- **THEN** SHALL 发送 SIGTERM，等待超时后 SIGKILL，返回 `True`

#### Scenario: Interrupt idle session
- **WHEN** session 没有正在运行的进程，调用 `interrupt(session_id)`
- **THEN** SHALL 返回 `False`

### Requirement: CodexAdapter process cleanup on error
当 CLI 进程异常退出（returncode != 0）时，系统 SHALL 捕获 stderr 输出，生成 `ERROR` 类型的 StreamEvent，并释放进程资源。

#### Scenario: Process exits with error
- **WHEN** CLI 进程以非零退出码退出
- **THEN** SHALL 从 `_processes` 字典中移除该进程，并生成 `StreamEvent(type=ERROR)` 包含 stderr 内容

### Requirement: CodexAdapter config integration
`CodexAdapter` SHALL 从 `settings.cli.codex_path` 读取 Codex CLI 可执行文件路径。该路径 SHALL 在 `config.yaml` 中配置，默认值为 `"codex"`。

#### Scenario: Use configured codex path
- **WHEN** `config.yaml` 中 `cli.codex_path` 设置为 `"/usr/local/bin/codex"`
- **THEN** `_build_command` SHALL 使用 `"/usr/local/bin/codex"` 作为命令的第一个元素

#### Scenario: Use default codex path
- **WHEN** `config.yaml` 中未设置 `cli.codex_path`
- **THEN** SHALL 使用默认值 `"codex"`

### Requirement: CodexAdapter agent_type field
所有 `CodexAdapter` 生成的 StreamEvent SHALL 在 `content` 中包含 `agent_type: "codex"` 字段，与 OpenCode adapter 行为一致。

#### Scenario: Events include agent_type
- **WHEN** CodexAdapter 生成任何 StreamEvent
- **THEN** 该 event 的 content SHALL 包含 `"agent_type": "codex"`
