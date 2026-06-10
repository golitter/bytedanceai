## Why

开发测试过程中，MySQL/Redis 数据、agentend 日志、worktree 目录和 pytest 缓存会持续积累，导致测试环境脏数据、磁盘占用增长、调试困难。需要一键清理脚本快速重置测试环境。

## What Changes

- 新增 `scripts/cleanup.sh`，提供交互式 + 无交互式两种清理模式
- 清理 Docker 中 MySQL（database: agenthub）的所有表数据
- 清理 Docker 中 Redis 的所有 key（FLUSHDB）
- 清理 agentend/logs/ 目录（session_mappings.json、workspaces.json）
- 清理 agentend/worktrees/ 目录（git worktree 隔离目录）
- 清理 agentend/.pytest_cache/ 及 __pycache__
- 支持 `--dry-run` 模式预览将被清理的内容
- 支持 `--yes` 跳过确认提示（CI 场景）

## Capabilities

### New Capabilities
- `test-data-cleanup`: 一键清理测试数据的 shell 脚本，覆盖 MySQL/Redis/日志/worktree/缓存

### Modified Capabilities

（无）

## Impact

- `scripts/` 目录新增 `cleanup.sh`
- 运行时会影响 Docker 容器中的 MySQL 和 Redis 数据（不可逆）
- 运行时会删除 agentend/logs/ 和 agentend/worktrees/ 下的文件（不可逆）
