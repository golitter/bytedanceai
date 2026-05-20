# AgentEnd Runtime 实现文档

## 概述

AgentEnd Runtime 是一个 Python FastAPI 服务，作为 AgentHub 多 Agent 协作平台的 AI 执行引擎。

架构定位：

```
React 前端 → Go Backend → AgentEnd Runtime (本服务) → Claude Code CLI
```

Go Backend 通过 HTTP 调用 Runtime，Runtime 启动 Claude Code CLI 子进程执行编码任务，结果通过 SSE（Server-Sent Events）流式返回。

## 项目结构

```
agentend/
├── src/
│   ├── adapters/       # Adapter 适配器层
│   ├── api/            # FastAPI HTTP 端点
│   │   └── v1/         # v1 版本 API
│   ├── app/            # 应用入口、配置、DI
│   ├── rules/          # Rule Engine 规则引擎
│   ├── schemas/        # 数据模型
│   └── session/        # Session 会话管理
├── tests/              # 测试
├── pyproject.toml      # 项目配置与依赖
└── ruff.toml           # 代码风格
```
