## Context

AgentHub 是一个多 Agent 协作平台，采用三层架构：React 前端 → Go Backend → Python Agent Runtime。当前 `agentend/` 目录仅有 `ruff.toml`，需要从零构建 Python Runtime 服务。

Runtime 定位为 **AI 执行引擎**：接收 Go Backend 的 HTTP 请求，通过 CLI subprocess 调用 Claude Code 完整 coding agent 能力，以 SSE 流式返回事件。

关键约束：
- Python 采用类型提示（type hints）
- 模块导出风格：`from src.tools.module import se_tools_list`
- 使用 uv 管理 Python 版本和依赖（`uv add`、`uv run`），虚拟环境通过 `source .venv/bin/activate` 激活
- MVP 只做 Claude Code 一个 Adapter
- Tool 执行先交给 Agent 自身处理
- 后期需扩展 Rule 强规则系统、Codex/OpenCode Adapter、Workspace 隔离

## Goals / Non-Goals

**Goals:**

- 可运行的 FastAPI 服务，Go Backend 通过 HTTP + SSE 调用
- 统一 Adapter 抽象，CLI subprocess 驱动 Claude Code
- 流式事件输出：解析 Claude Code CLI 的 stream-json → 统一 StreamEvent → SSE
- 内存级 Session 管理，跟踪进程状态
- 可插拔 Rule Engine，支持安全检查和 prompt/tool 约束注入
- 清晰的模块边界和类型提示，便于后期扩展

**Non-Goals:**

- 不做 Orchestrator（任务拆解/路由），第三阶段实现
- 不做 Workspace/Git Worktree 隔离，第二阶段实现
- 不做 Tool Runtime（统一工具执行层），Agent 自行管理 tool
- 不做 Codex/OpenCode Adapter，MVP 只做 Claude Code
- 不做持久化存储，Session 全部内存管理
- 不做容器化/Docker Sandbox
- 不做认证鉴权，内部服务间调用

## Decisions

### D1: CLI Subprocess vs SDK 直调

**选择**: CLI subprocess（`asyncio.create_subprocess_exec`）

**理由**: Claude Code 的完整 coding agent 能力（文件编辑、终端执行、多轮 tool use）只能通过 CLI 获得。SDK 只提供 LLM 对话 + tool calling，不具备 coding agent 的自主执行能力。

**替代方案**: 使用 Anthropic SDK 直接调 API — 更可控但只能拿到 LLM 能力，不是完整的 coding agent。

### D2: 进程管理策略

**选择**: 每个 Session 对应一个独立的 CLI subprocess，用 `dict[str, asyncio.subprocess.Process]` 管理。

**理由**: 简单直接，每个会话独立互不干扰。Claude Code 自带 session 机制（`--session` 参数），可跨多次调用复用会话。

**中断策略**: SIGTERM → 等待 5s → SIGKILL。

### D3: 流式输出方案

**选择**: Claude Code `--output-format stream-json` + 逐行解析 stdout + SSE 返回。

**流程**:
```
Go POST → FastAPI → ClaudeCodeAdapter → claude CLI (stream-json)
                                              ↓ stdout
                                         逐行 JSON 解析
                                              ↓
                                        StreamEvent 转换
                                              ↓
                                         SSE 推送给 Go
```

**后备**: 如果 stream-json 格式不稳定，先用 `--output-format json`（非流式）跑通，后续再切流式。

### D4: Rule Engine 注入方式

**选择**: Rule 在请求执行前评估，通过修改 CLI 参数注入约束：
- `--append-system-prompt`: 追加安全/作用域约束文本
- `--allowedTools`: 限制可用工具集
- `--max-turns`: 限制执行轮数

**理由**: 不侵入 Agent 内部逻辑，完全通过 CLI 参数控制。Rule 评估失败直接返回错误，不启动进程。

### D5: 项目结构 — src layout

**选择**: 采用 `src/` 布局，模块通过 `__init__.py` 导出公共 API。

```
agentend/
├── src/
│   ├── adapters/       # Adapter 层
│   ├── session/        # Session 管理
│   ├── rules/          # Rule 引擎
│   ├── schemas/        # 数据模型
│   ├── api/            # FastAPI 端点
│   └── app.py          # FastAPI 应用
├── tests/
├── .venv/              # uv 管理的虚拟环境
├── pyproject.toml      # uv 管理的项目配置和依赖
├── uv.lock             # uv 锁文件
└── ruff.toml
```

**理由**: src layout 防止意外 import 未安装的包，是 Python 社区推荐标准。模块导出统一通过 `__init__.py`，如 `from src.adapters import BaseAgentAdapter`。

### D7: 包管理工具 — uv

**选择**: 使用 uv 管理 Python 版本、虚拟环境和依赖。

**工作流**:
```bash
# 在 agentend/ 目录下
uv init                        # 初始化项目（生成 pyproject.toml）
uv add fastapi uvicorn ...     # 添加依赖
source .venv/bin/activate      # 激活虚拟环境
uv run python -m src.app.main  # 运行服务
```

**理由**: uv 比 pip + venv 快 10-100 倍，自带锁文件（uv.lock）保证可复现构建，一个工具覆盖 Python 版本管理、虚拟环境、依赖安装。

### D6: 通信协议

**选择**: HTTP + SSE（Server-Sent Events）

**请求**: Go → Python，`POST /v1/agent/stream`，JSON body
**响应**: Python → Go，SSE 流，`event: <type>\ndata: <json>\n\n`

**理由**: SSE 比 WebSocket 简单，单向流（Python → Go）足够。Go 侧已有 WebSocket 到前端的通道，只需要把 SSE 内容转发即可。

## Risks / Trade-offs

- **[Claude Code CLI 输出格式不稳定]** → 先用 `--output-format json` 非流式跑通，确认 stream-json 格式后再切流式；Adapter 内部做好容错解析
- **[CLI 进程 hang]** → 设置 `asyncio.wait_for` 超时，超时后 SIGTERM → SIGKILL；Session 状态标记为 ERROR
- **[并发 Session 资源消耗]** → MVP 阶段不限制并发数，后期加 Session 上限和队列
- **[Rule Engine 修改 CLI 参数可能有遗漏]** → `_build_command()` 方法集中组装所有参数，Rule 只返回约束 dict，不直接操作命令
- **[内存级 Session 不持久化]** → Runtime 重启丢失所有会话状态；可接受的 trade-off，MVP 不需要持久化
