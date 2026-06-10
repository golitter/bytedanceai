## Context

项目开发过程中使用 Docker 运行 MySQL（port 3306, db: agenthub）和 Redis（port 6379）。agentend 端在运行过程中会生成 logs/（session_mappings.json、workspaces.json）和 worktrees/（git worktree 隔离目录），以及 pytest 缓存。目前没有统一的清理手段，需手动逐个清理，效率低且容易遗漏。

现有 `scripts/run.sh` 负责服务启停，清理脚本将作为其补充，放在同一 `scripts/` 目录下。

## Goals / Non-Goals

**Goals:**
- 一键清理所有测试相关数据（MySQL、Redis、日志、worktree、缓存）
- 支持 dry-run 预览模式
- 支持跳过确认的 CI 模式
- 支持选择性清理（只清理某一项）
- 清理前有明确的确认提示，防止误操作

**Non-Goals:**
- 不清理 backend 的 Go 编译缓存
- 不清理 frontend 的 node_modules 或 pnpm 缓存
- 不做数据备份（用户需自行备份重要数据）

## Decisions

1. **纯 shell 脚本**：与现有 `scripts/run.sh` 风格一致，依赖仅 bash + docker + mysql/redis CLI
   - 备选方案：Make target — 但清理脚本逻辑较复杂，shell 更灵活

2. **通过 `docker exec` 操作 MySQL/Redis**：直接在容器内执行 SQL / Redis 命令，无需本机安装客户端
   - 容器名通过脚本内变量配置，用户可覆盖

3. **交互式确认 + `--yes` 跳过**：默认逐项确认，`--yes` 全部跳过
   - `--dry-run` 仅打印将要执行的操作，不实际执行

4. **子命令模式**：`cleanup.sh all` / `cleanup.sh mysql` / `cleanup.sh redis` / `cleanup.sh logs` / `cleanup.sh worktrees` / `cleanup.sh cache`，支持选择性清理

## Risks / Trade-offs

- **[数据丢失]** MySQL FLUSH 和 Redis FLUSHDB 不可逆 → 确认提示 + dry-run 预览
- **[Docker 未运行]** docker exec 会失败 → 启动前检测 Docker daemon 状态，给出友好提示
- **[容器名不匹配]** 不同开发者的容器名可能不同 → 通过环境变量或参数覆盖
