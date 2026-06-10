## Context

Agentend 目前有 3 个 adapter：ClaudeCodeAdapter（Claude CLI subprocess）、OpenCodeAdapter（OpenCode CLI subprocess）、OrchestratorAdapter（in-process LangGraph）。所有 subprocess adapter 遵循相同模式：`_build_command` 构建命令 → `asyncio.create_subprocess_exec` 启动进程 → 逐行解析 stdout JSON → 映射为 `StreamEvent`。

Codex CLI（OpenAI，v0.132.0）提供 `codex exec --json` 子命令用于非交互式执行，输出 NDJSON 事件流。已通过实测确认事件结构：

```
thread.started   → {thread_id}
turn.started     → (无数据)
item.started     → {type: "command_execution", command, status: "in_progress"}
item.completed   → {type: "reasoning" | "agent_message" | "command_execution"}
turn.completed   → {usage}
```

Codex 有内置沙箱（`--sandbox workspace-write`）和审批控制（`-a never`），比 Claude/OpenCode 多一层安全保障。

## Goals / Non-Goals

**Goals:**
- 接入 Codex CLI 作为第四个 Agent 后端
- 复用现有 subprocess adapter 模式（Claude/OpenCode 已验证）
- 支持会话恢复（`codex exec resume`）
- 支持 Codex 特有的 `item.started` 事件（工具调用进度）

**Non-Goals:**
- 不接入 Codex Cloud / Codex Desktop
- 不接入 Codex 的 MCP server 模式
- 不修改现有 adapter 的接口或行为
- 不处理 Codex 的 `--sandbox` 外部沙箱配置（使用 `workspace-write` 默认策略）

## Decisions

### Decision 1: 使用 `codex exec --json` 而非 `codex` 交互模式

**选择**: `codex exec --json -a never -s workspace-write`
**替代方案**: 使用 `codex app-server --listen stdio://` 的 JSONL 协议
**理由**: `exec` 模式是 Codex 官方的非交互/脚本化接口，与 Claude 的 `-p` 和 OpenCode 的 `run` 对齐。app-server 是 Experimental 功能，接口不稳定。

### Decision 2: `command_execution` 映射为 TOOL_CALL / TOOL_RESULT 对

**选择**: `item.started(command_execution)` → TOOL_CALL，`item.completed(command_execution)` → TOOL_RESULT
**替代方案**: 只映射 `item.completed` 为 TOOL_RESULT（忽略 started）
**理由**: Codex 是唯一提供工具执行开始事件的 CLI，映射为 TOOL_CALL 可以让前端实时显示工具执行进度，提升 UX。

### Decision 3: reasoning 事件映射为 `[thinking]` 前缀文本

**选择**: 与 OpenCode adapter 保持一致，`item.completed(type="reasoning")` → `TEXT` 并加 `[thinking]` 前缀
**替代方案**: 映射为独立的 PLANNING 事件
**理由**: 与现有 OpenCode adapter 行为对齐，前端无需额外适配。

### Decision 4: 会话恢复使用 `codex exec resume <thread_id>`

**选择**: 首次执行从 `thread.started` 事件捕获 `thread_id`，后续使用 `codex exec resume <thread_id> --json` 恢复
**替代方案**: 使用 `--last` 恢复最近会话
**理由**: 按 thread_id 恢复更精确，与 Claude 的 `--resume` 和 OpenCode 的 `--session` 模式一致。

### Decision 5: 沙箱策略硬编码为 `workspace-write`

**选择**: 固定使用 `-s workspace-write`，不暴露为配置项
**替代方案**: 允许通过 config.yaml 配置沙箱策略
**理由**: `workspace-write` 是安全默认值，只允许写入工作区目录。当前阶段无需暴露此选项，后续有需求再开放。

## Risks / Trade-offs

- **[Codex CLI 版本兼容性]** Codex CLI 版本更新可能改变 `--json` 输出的事件结构 → 在 adapter 中对未知事件类型静默忽略，不抛异常
- **[item.type 未知类型]** 实测只确认了 `reasoning`、`agent_message`、`command_execution` 三种 item type，Codex 可能还有 `function_call`、`file_edit` 等 → `_parse_stream_line` 对未知 type 返回 None，不中断流
- **[stderr ERROR 噪音]** 模型列表刷新失败的 ERROR 日志混在 stderr 中 → adapter 只读 stdout，stderr 在进程异常退出时才读取
- **[配置目录]** Codex 配置在 `~/.codex/`（注意小写），workspace 隔离的 `config_dir_exclusion` 需要排除此目录
