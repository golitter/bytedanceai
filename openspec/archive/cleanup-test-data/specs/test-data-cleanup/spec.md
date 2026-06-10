## ADDED Requirements

### Requirement: 全量清理命令
脚本 SHALL 提供 `cleanup.sh all` 子命令，依次执行 MySQL、Redis、日志、worktree、缓存的全部清理。

#### Scenario: 执行全量清理
- **WHEN** 用户运行 `scripts/cleanup.sh all --yes`
- **THEN** 脚本依次清空 MySQL agenthub 数据库所有表数据、Redis FLUSHDB、删除 agentend/logs/*、删除 agentend/worktrees/*、删除 agentend/.pytest_cache 和 __pycache__

#### Scenario: dry-run 全量预览
- **WHEN** 用户运行 `scripts/cleanup.sh all --dry-run`
- **THEN** 脚本仅打印每一步将要执行的操作，不实际执行任何清理

### Requirement: MySQL 清理
脚本 SHALL 提供 `cleanup.sh mysql` 子命令，通过 `docker exec` 在 MySQL 容器中执行 `DROP DATABASE agenthub; CREATE DATABASE agenthub;` 重建空数据库。

#### Scenario: MySQL 清理成功
- **WHEN** 用户运行 `scripts/cleanup.sh mysql --yes`
- **THEN** MySQL 容器内的 agenthub 数据库被重建为空库

#### Scenario: MySQL 容器未运行
- **WHEN** 用户运行 `scripts/cleanup.sh mysql`，但 Docker daemon 未运行或 MySQL 容器不存在
- **THEN** 脚本输出错误提示并退出，exit code 1

### Requirement: Redis 清理
脚本 SHALL 提供 `cleanup.sh redis` 子命令，通过 `docker exec` 在 Redis 容器中执行 `FLUSHDB`。

#### Scenario: Redis 清理成功
- **WHEN** 用户运行 `scripts/cleanup.sh redis --yes`
- **THEN** Redis 容器内当前数据库的所有 key 被清除

### Requirement: 日志清理
脚本 SHALL 提供 `cleanup.sh logs` 子命令，删除 `agentend/logs/` 目录下的所有文件。

#### Scenario: 日志清理成功
- **WHEN** 用户运行 `scripts/cleanup.sh logs --yes`
- **THEN** `agentend/logs/session_mappings.json` 和 `agentend/logs/workspaces.json` 及其他日志文件被删除

### Requirement: Worktree 清理
脚本 SHALL 提供 `cleanup.sh worktrees` 子命令，删除 `agentend/worktrees/` 目录下的所有内容。

#### Scenario: Worktree 清理成功
- **WHEN** 用户运行 `scripts/cleanup.sh worktrees --yes`
- **THEN** `agentend/worktrees/` 目录下的所有子目录和文件被删除

### Requirement: 缓存清理
脚本 SHALL 提供 `cleanup.sh cache` 子命令，删除 `agentend/.pytest_cache/` 及递归删除 `agentend/` 下所有 `__pycache__` 目录。

#### Scenario: 缓存清理成功
- **WHEN** 用户运行 `scripts/cleanup.sh cache --yes`
- **THEN** `agentend/.pytest_cache/` 被删除，所有 `__pycache__` 目录被递归删除

### Requirement: 交互式确认
脚本在默认模式下（无 `--yes`）SHALL 在执行每项清理前提示用户确认。

#### Scenario: 用户确认清理
- **WHEN** 用户运行 `scripts/cleanup.sh mysql`（无 --yes）
- **THEN** 脚本显示 "即将清空 MySQL agenthub 数据库，确认？[y/N]" 并等待用户输入

#### Scenario: 用户拒绝清理
- **WHEN** 用户输入非 y/Y
- **THEN** 脚本跳过该步骤并继续

### Requirement: 帮助信息
脚本 SHALL 提供 `cleanup.sh help` 子命令和 `cleanup.sh --help`，显示用法说明。

#### Scenario: 查看帮助
- **WHEN** 用户运行 `scripts/cleanup.sh help`
- **THEN** 脚本显示所有子命令、选项说明和使用示例
