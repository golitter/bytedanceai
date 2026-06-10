## 1. 契约层更新

- [x] 1.1 在 `contracts/schemas/agent-request.yaml` 的 AgentType 枚举中添加 `codex` 值（description: "Codex CLI"）
- [x] 1.2 在 `contracts/logs/` 创建变更记录文件 `2026-05-25-add-codex-agent-type.md`
- [x] 1.3 运行 `make generate` 生成三端类型文件，确认 Python/TypeScript/Go 的 AgentType 枚举包含 `codex`

## 2. 配置层

- [x] 2.1 在 `agentend/src/app/config.py` 的 CLISettings 中添加 `codex_path: str = "codex"` 字段
- [x] 2.2 在 `agentend/config.yaml` 中添加 `cli.codex_path: "codex"` 配置项

## 3. CodexAdapter 核心实现

- [x] 3.1 创建 `agentend/src/adapters/codex.py`，实现 `BaseAgentAdapter` 接口骨架（`__init__`、`create_session`、`destroy_session`）
- [x] 3.2 实现 `_build_command` 方法：构建 `codex exec --json -a never -s workspace-write` 命令，支持 `-C`（cwd）、`-m`（model）、`exec resume` 参数
- [x] 3.3 实现 `_parse_stream_line` 方法：NDJSON 事件解析，映射 thread.started→INIT、item.started(command_execution)→TOOL_CALL、item.completed(reasoning/agent_message)→TEXT、item.completed(command_execution)→TOOL_RESULT、turn.completed→DONE
- [x] 3.4 实现 `stream_chat` 方法：启动子进程，逐行解析 stdout，yield StreamEvent，处理进程异常退出
- [x] 3.5 实现 `chat` 方法：聚合 stream_chat 事件，返回 AgentResponse
- [x] 3.6 实现 `interrupt` 方法：SIGTERM + 超时 SIGKILL，复用现有 terminate_timeout 配置

## 4. 注册与集成

- [x] 4.1 在 `agentend/src/app/dependencies.py` 的 `create_adapter_registry` 中注册 `CodexAdapter`（key: `AgentType.CODEX`）
- [x] 4.2 确认 workspace 隔离排除 `~/.codex/` 配置目录（检查 `config_dir_exclusion` 相关代码）

## 5. 验证

- [ ] 5.1 启动 agentend 服务，通过 API 发送 `agent_type: "codex"` 的请求，验证流式响应正常
- [ ] 5.2 测试会话恢复：发送第二条消息到同一 session，确认使用 `exec resume` 命令
- [ ] 5.3 测试中断：发送长时间任务后调用 interrupt，确认进程被终止
- [ ] 5.4 测试错误处理：使用无效 workspace 路径，确认 ERROR 事件正确返回
