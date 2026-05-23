# AGENTS.md

## 项目简介

Monorepo 项目，包含前端、后端、Agent 端三个子项目。

## 目录结构

```
bytedanceai/
├── frontend/      # React 前端 → 参见 frontend/AGENTS.md
├── backend/       # Go 后端   → 参见 backend/AGENTS.md
├── agentend/      # Python Agent 端 → 参见 agentend/AGENTS.md
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


## Git 规范

详见 [docs/common/git-conventions.md](docs/common/git-conventions.md)。
