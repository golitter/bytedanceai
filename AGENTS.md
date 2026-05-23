# AGENTS.md

## 项目简介

Monorepo 项目，包含前端、后端、Agent 端三个子项目。

## 目录结构

```
bytedanceai/
├── frontend/      # React 前端 → 参见 frontend/AGENTS.md
├── backend/       # Go 后端   → 参见 backend/AGENTS.md
├── agentend/      # Python Agent 端 → 参见 agentend/AGENTS.md
├── contracts/     # 三端共享契约（schemas + logs）→ 参见 contracts/AGENTS.md
├── docs/          # 项目文档
└── scripts/       # 工程脚本
```

子项目的框架选型、构建命令、测试方式等详情，请查阅对应目录下的 AGENTS.md。

## Makefile

通过 `make <命令>` 管理三端服务（热重载），详情参见 [docs/common/makefile-guide.md](docs/common/makefile-guide.md)。

| 命令 | 说明 |
|------|------|
| `make run-frontend` | 启动前端 |
| `make run-backend` | 启动后端 |
| `make run-agentend` | 启动 Agent 端 |
| `make stop` | 停止全部 |
| `make status` | 查看运行状态 |
| `make generate` | 从 contracts/schemas/ 生成三端类型文件 |

## 契约优先原则

跨端协议的类型定义以 `contracts/schemas/` 为单一来源。修改协议时：

1. 先更新 `contracts/schemas/*.yaml`
2. 运行 `make generate` 生成三端类型
3. 在 `contracts/logs/` 写入变更记录

详见 [contracts/AGENTS.md](contracts/AGENTS.md)。


## Git 规范

详见 [docs/common/git-conventions.md](docs/common/git-conventions.md)。
