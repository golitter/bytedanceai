## 1. 脚本骨架与基础设施

- [x] 1.1 创建 `scripts/cleanup.sh`，编写 shebang、颜色定义、公共变量（Docker 容器名、MySQL 连接信息）
- [x] 1.2 实现参数解析：子命令（all/mysql/redis/logs/worktrees/cache/help）+ 选项（--yes / --dry-run / --help）
- [x] 1.3 实现 Docker daemon 检测函数 `check_docker()`
- [x] 1.4 实现交互式确认函数 `confirm()`，支持 `--yes` 跳过和 `--dry-run` 模式
- [x] 1.5 实现 `help` 子命令，输出用法说明

## 2. 各清理模块实现

- [x] 2.1 实现 `cleanup_mysql()`：通过 `docker exec` 执行 `DROP DATABASE agenthub; CREATE DATABASE agenthub;`
- [x] 2.2 实现 `cleanup_redis()`：通过 `docker exec` 执行 `redis-cli FLUSHDB`
- [x] 2.3 实现 `cleanup_logs()`：删除 `agentend/logs/*`
- [x] 2.4 实现 `cleanup_worktrees()`：读取 workspaces.json → git worktree remove → git branch -D → 删除 worktrees 目录
- [x] 2.5 实现 `cleanup_cache()`：删除 `agentend/.pytest_cache` 并 `find` 递归删除 `__pycache__`（排除 .venv）

## 3. 调度与收尾

- [x] 3.1 实现 `all` 子命令，依次调用全部清理函数
- [x] 3.2 实现子命令路由（case/esac），将子命令映射到对应清理函数
- [x] 3.3 设置 `chmod +x scripts/cleanup.sh`
- [x] 3.4 端到端测试：dry-run 模式验证所有子命令输出正确
