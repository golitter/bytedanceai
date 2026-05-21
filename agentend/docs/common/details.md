# 详细文档

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/agent/stream` | 流式 Agent 响应（SSE） |
| POST | `/v1/agent/execute` | 同步 Agent 执行 |
| GET | `/v1/session` | 列出所有会话 |
| GET | `/v1/session/{id}` | 获取会话详情 |
| POST | `/v1/session/{id}/interrupt` | 中断运行中会话 |
| DELETE | `/v1/session/{id}` | 删除会话 |
| POST | `/v1/workspace/create` | 创建工作区（Git Worktree） |
| POST | `/v1/workspace/{id}/commit` | 提交工作区变更 |
| POST | `/v1/workspace/{id}/merge` | 合并工作区到目标分支 |
| DELETE | `/v1/workspace/{id}` | 清理工作区 |
| GET | `/v1/workspace` | 列出所有工作区 |
| GET | `/health` | 健康检查 |

## 核心架构

- **执行流程**：请求到达 → 规则引擎评估 → 适配器注册表解析 → 会话管理器跟踪状态 → 适配器执行 → 结果流式/同步返回
- **会话状态机**：`IDLE → RUNNING → COMPLETED / INTERRUPTED / ERROR`
- **适配器模式**：通过抽象基类支持不同 Agent 类型，当前实现 Claude CLI 与 OpenCode CLI 适配器
- **规则引擎**：执行前评估 Safety（阻止危险工具）、Scope（校验工作区路径）等规则，可修改 system prompt 和工具白名单
- **会话持久化**：API session_id 与 CLI session_id 映射持久化至 `logs/session_mappings.json`
- **工作区管理**：基于 Git Worktree 的任务级隔离，支持自动创建任务分支（`task/{task_id}`）、提交、合并与清理，含 TTL 自动回收与启动恢复

## 配置

通过环境变量配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `CLAUDE_CLI_PATH` | `claude` | Claude CLI 可执行文件路径 |
| `OPENCODE_CLI_PATH` | `opencode` | OpenCode CLI 可执行文件路径 |
| `WORKSPACE_BASE_DIR` | `./worktrees` | Worktree 存放根目录 |
| `DEFAULT_MAX_TURNS` | `20` | 默认最大对话轮次 |
| `EXECUTION_TIMEOUT` | `300` | 执行超时（秒） |
| `WORKSPACE_TTL_SECONDS` | `3600` | 工作区 TTL（秒） |
| `WORKSPACE_TTL_CHECK_INTERVAL` | `300` | 工作区清理检查间隔（秒） |
| `HOST` | `0.0.0.0` | 监听地址 |
| `PORT` | `8001` | 监听端口 |
