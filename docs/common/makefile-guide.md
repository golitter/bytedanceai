# Makefile 使用指南

根目录 `Makefile` 通过 `scripts/run.sh` 脚本统一管理前端、后端、Agent 三端服务的启停与状态查看。

## 运行机制

- 脚本在 `.pids/` 目录下记录每个服务的 PID 文件
- 启动前检测 PID 文件 + 进程是否存活，已运行则跳过
- **启动/重启必须指定单个服务**，不允许同时启动全部
- **停止可以一次性全停**

## 命令一览

### 启动（单个）

| 命令 | 服务 | 热重载工具 | 端口 |
|------|------|-----------|------|
| `make run-frontend` | Vite dev server | Vite 内置 HMR | localhost:5173 |
| `make run-backend` | Go server | Air | localhost:8080 |
| `make run-agentend` | FastAPI server | uvicorn --reload | 见 agentend config |

### 停止

| 命令 | 说明 |
|------|------|
| `make stop` | 停止全部服务 |
| `make stop-frontend` | 停止前端 |
| `make stop-backend` | 停止后端 |
| `make stop-agentend` | 停止 Agent 端 |

### 重启（单个）

| 命令 | 说明 |
|------|------|
| `make restart-frontend` | 重启前端 |
| `make restart-backend` | 重启后端 |
| `make restart-agentend` | 重启 Agent 端 |

### 其他

| 命令 | 说明 |
|------|------|
| `make status` | 查看三端运行状态与 PID |
| `make tidy` | 执行 `go mod tidy` |

## 直接使用脚本

```bash
./scripts/run.sh start <frontend|backend|agentend>
./scripts/run.sh stop [<frontend|backend|agentend>]
./scripts/run.sh restart <frontend|backend|agentend>
./scripts/run.sh status
```
