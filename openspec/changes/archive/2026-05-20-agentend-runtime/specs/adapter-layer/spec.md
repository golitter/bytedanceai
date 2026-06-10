## ADDED Requirements

### Requirement: BaseAgentAdapter abstract interface
系统 SHALL 定义 `BaseAgentAdapter` 抽象基类，包含以下抽象方法：`create_session`、`chat`、`stream_chat`、`interrupt`、`destroy_session`。所有 Adapter 实现 MUST 继承此基类并实现全部抽象方法。

#### Scenario: Adapter implements all required methods
- **WHEN** 一个新的 Adapter 类继承 `BaseAgentAdapter`
- **THEN** 该类 MUST 实现 `create_session`、`chat`、`stream_chat`、`interrupt`、`destroy_session` 全部五个抽象方法，否则实例化时 SHALL 抛出 `TypeError`

### Requirement: AdapterRegistry registration and lookup
系统 SHALL 提供 `AdapterRegistry`，支持通过 agent 类型名称注册和查找 Adapter 类。Registry MUST 在启动时自动注册内置 Adapter。

#### Scenario: Register and retrieve adapter
- **WHEN** 调用 `registry.register("claude-code", ClaudeCodeAdapter)` 注册后
- **THEN** `registry.get("claude-code")` SHALL 返回 `ClaudeCodeAdapter` 类

#### Scenario: Lookup unregistered adapter
- **WHEN** 调用 `registry.get("unknown-agent")` 查找不存在的 Adapter
- **THEN** SHALL 抛出 `ValueError` 并包含 agent 类型名称信息

### Requirement: ClaudeCodeAdapter CLI subprocess management
`ClaudeCodeAdapter` SHALL 通过 `asyncio.create_subprocess_exec` 启动 Claude Code CLI 进程，MUST 支持构建以下 CLI 参数：`-p`（prompt）、`--session`（会话复用）、`--output-format`（输出格式）、`--system-prompt`、`--append-system-prompt`、`--allowedTools`、`--max-turns`。

#### Scenario: Build command with all parameters
- **WHEN** 调用 `_build_command` 并传入 session_id、message、system_prompt_append、allowed_tools、max_turns
- **THEN** SHALL 返回包含所有对应 CLI 参数的命令列表

#### Scenario: Build command with minimal parameters
- **WHEN** 调用 `_build_command` 只传入 message，不传可选参数
- **THEN** SHALL 返回仅包含 `-p <message>` 和 `--output-format stream-json` 的最小命令

### Requirement: ClaudeCodeAdapter stream_chat yields StreamEvent
`ClaudeCodeAdapter.stream_chat` SHALL 启动 CLI 进程并以 `AsyncIterator[StreamEvent]` 形式逐行解析 stdout 的 JSON 输出，转换为统一的 `StreamEvent` 对象。

#### Scenario: Parse stream-json output lines
- **WHEN** CLI 进程输出 `{"type":"assistant","content":[...]}` 格式的 JSON 行
- **THEN** SHALL 解析为对应类型的 `StreamEvent`（text / tool_call / tool_result / done）

#### Scenario: Handle non-JSON stdout line
- **WHEN** CLI 进程输出无法解析为 JSON 的行
- **THEN** SHALL 将该行作为 `type="text"` 的 `StreamEvent` 包装输出，不 SHALL 抛出异常

### Requirement: ClaudeCodeAdapter process interrupt
`ClaudeCodeAdapter.interrupt` SHALL 终止指定 session 的正在运行的 CLI 进程。MUST 先发送 SIGTERM，等待最多 5 秒，若进程未结束则发送 SIGKILL。

#### Scenario: Interrupt running process
- **WHEN** session 有正在运行的进程，调用 `interrupt(session_id)`
- **THEN** SHALL 发送 SIGTERM，等待最多 5 秒，未结束则 SIGKILL，返回 `True`

#### Scenario: Interrupt idle session
- **WHEN** session 没有正在运行的进程，调用 `interrupt(session_id)`
- **THEN** SHALL 返回 `False`

### Requirement: ClaudeCodeAdapter process cleanup on error
当 CLI 进程异常退出（returncode != 0）时，系统 SHALL 捕获 stderr 输出，将进程状态标记为 ERROR，并释放进程资源。

#### Scenario: Process exits with error
- **WHEN** CLI 进程以非零退出码退出
- **THEN** SHALL 从 `_processes` 字典中移除该进程，并生成 `type="error"` 的 StreamEvent 包含 stderr 内容
