## Why

AgentHub 是一个多 Agent 协作平台，需要 Python Agent Runtime 作为 AI 执行引擎，统一管理 Agent 会话、CLI 子进程调度、规则约束和流式事件输出。当前 `agentend/` 目录为空（仅有 ruff.toml），需要从零构建整个 Runtime。

## What Changes

- 新增 FastAPI 服务骨架：入口、配置管理、依赖注入、生命周期管理
- 新增统一 Adapter 抽象层：`BaseAgentAdapter` ABC + `ClaudeCodeAdapter` 实现（CLI subprocess 驱动）
- 新增 Session Manager：内存级会话管理，跟踪进程状态与消息历史
- 新增 Rule Engine：可插拔规则系统，支持 pre-check 和 prompt/tool 约束注入
- 新增统一消息协议：`AgentRequest` / `AgentResponse` / `StreamEvent` 数据模型
- 新增 API 端点：`POST /v1/agent/stream`（SSE）、`POST /v1/agent/execute`（同步）、Session CRUD
- 新增 SSE 流式通信：解析 Claude Code CLI 的 stream-json 输出，转换为统一事件流推送给 Go Backend

## Capabilities

### New Capabilities

- `adapter-layer`: 统一 Agent 适配器抽象，支持 CLI subprocess 方式调用 Claude Code，包含进程生命周期管理
- `session-manager`: Agent 会话管理，跟踪状态（IDLE/RUNNING/COMPLETED/ERROR）、进程句柄、消息历史
- `rule-engine`: 可插拔规则系统，在请求执行前进行安全检查、作用域约束，注入 system_prompt 和 allowedTools
- `stream-protocol`: 统一消息协议与 SSE 流式输出，定义 Request/Response/Event 数据模型和 SSE 推送格式
- `api-gateway`: FastAPI 服务骨架与 HTTP 端点，作为 Go Backend 调用的入口

### Modified Capabilities

（无已有 capability 需要修改）

## Impact

- **目录**: `agentend/` 从空目录变为完整的 Python FastAPI 项目
- **API**: 新增 ~5 个 HTTP 端点，Go Backend 通过 HTTP + SSE 调用
- **依赖管理**: 使用 uv 管理 Python 版本和依赖，`uv add` 安装包，`.venv/` 虚拟环境
- **外部依赖**: 要求机器上已安装 `claude` CLI 工具
- **通信**: Go Backend → Python Runtime（HTTP），Python → Go（SSE 流式返回）
