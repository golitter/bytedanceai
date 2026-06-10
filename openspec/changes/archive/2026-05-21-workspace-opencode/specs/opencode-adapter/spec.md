## ADDED Requirements

### Requirement: OpenCodeAdapter inherits BaseAgentAdapter
`OpenCodeAdapter` MUST 继承 `BaseAgentAdapter` 并实现全部抽象方法：`create_session`、`chat`、`stream_chat`、`interrupt`、`destroy_session`。

#### Scenario: OpenCodeAdapter is valid subclass
- **WHEN** 实例化 `OpenCodeAdapter()`
- **THEN** SHALL 成功创建实例，具备 `create_session`、`chat`、`stream_chat`、`interrupt`、`destroy_session` 方法

### Requirement: OpenCodeAdapter CLI command building
`OpenCodeAdapter` SHALL 通过 `_build_command` 方法组装 OpenCode CLI 命令，支持参数：`-p`（prompt）、`-f json`（JSON 输出）、`-q`（静默）、`-c`（工作目录）。MUST 支持从 config 读取 `OPENCODE_CLI_PATH`。

#### Scenario: Build command with prompt and cwd
- **WHEN** 调用 `_build_command("修复 bug", cwd="/workspaces/task-1")`
- **THEN** SHALL 返回 `[OPENCODE_CLI_PATH, "-p", "修复 bug", "-f", "json", "-q", "-c", "/workspaces/task-1"]`

#### Scenario: Build command with rule constraint in prompt
- **WHEN** 调用 `_build_command("修复 bug", system_prompt_append="只改 tsx 文件")`
- **THEN** prompt SHALL 拼接为 `"[系统约束: 只改 tsx 文件]\n\n修复 bug"`

### Requirement: OpenCodeAdapter non-streaming chat
`OpenCodeAdapter.chat` SHALL 执行 `opencode -p "..." -f json -q`，等待进程完成，解析 JSON 输出为 `AgentResponse`。

#### Scenario: Successful chat execution
- **WHEN** 调用 `chat(session_id, "写一个按钮")`
- **THEN** SHALL 启动 OpenCode CLI 进程，等待完成，解析 JSON 输出为 `AgentResponse(session_id=..., content=..., artifacts=[], usage={})`

#### Scenario: Chat with process error
- **WHEN** OpenCode CLI 进程以非零退出码退出
- **THEN** SHALL 返回 `AgentResponse`，`content` 包含 stderr 错误信息

### Requirement: OpenCodeAdapter simulated stream_chat
`OpenCodeAdapter.stream_chat` SHALL 执行非交互模式获取完整 JSON 结果，然后将内容拆分为多个 `StreamEvent` 依次 yield：先 yield text 事件，最后 yield done 事件。

#### Scenario: Stream from JSON output
- **WHEN** OpenCode CLI 返回 JSON `{"content": "创建按钮...", "toolUses": [...], "usage": {...}}`
- **THEN** SHALL 依次 yield：`StreamEvent(type="text", content={"text": "创建按钮..."})`，`StreamEvent(type="done", content={"usage": {...}})`

#### Scenario: Stream fallback for non-JSON output
- **WHEN** OpenCode CLI 返回纯文本而非有效 JSON
- **THEN** SHALL yield 单个 `StreamEvent(type="text", content={"text": raw_output})`，然后 yield `StreamEvent(type="done")`

### Requirement: OpenCodeAdapter process management
`OpenCodeAdapter` SHALL 维护 `_processes: dict[str, asyncio.subprocess.Process]`，与 `ClaudeCodeAdapter` 相同的进程管理模式。

#### Scenario: Interrupt OpenCode process
- **WHEN** session 有运行中的 OpenCode 进程，调用 `interrupt(session_id)`
- **THEN** SHALL 发送 SIGTERM，等待 5s 后 SIGKILL，返回 `True`

#### Scenario: Destroy session
- **WHEN** 调用 `destroy_session(session_id)`
- **THEN** SHALL 中断进程（如有）并释放资源

### Requirement: OpenCodeAdapter config integration
`OpenCodeAdapter` SHALL 从 `src.app.config.settings` 读取 `OPENCODE_CLI_PATH`（默认 `"opencode"`）。

#### Scenario: Use default CLI path
- **WHEN** 未配置 `OPENCODE_CLI_PATH`
- **THEN** SHALL 使用 `"opencode"` 作为命令路径

#### Scenario: Use custom CLI path
- **WHEN** 设置 `OPENCODE_CLI_PATH=/usr/local/bin/opencode`
- **THEN** SHALL 使用该路径构建命令

### Requirement: Register OpenCodeAdapter in AdapterRegistry
启动时 `AdapterRegistry` SHALL 自动注册 `"opencode"` → `OpenCodeAdapter`。

#### Scenario: Lookup opencode adapter
- **WHEN** 调用 `registry.get("opencode")`
- **THEN** SHALL 返回 `OpenCodeAdapter` 类
