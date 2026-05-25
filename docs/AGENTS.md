# AGENTS.md — docs/

本目录是项目级文档的统一入口。所有 `docs/` 目录（根 + 子项目）遵循相同的 **5 分类约定**。

## 分类约定

```
docs/
├── design/        # 怎么建的 — 架构设计、实现方案、数据流、组件设计
├── reference/     # 是什么  — 技术栈、API 端点、适配器差异、视觉规范
├── guides/        # 怎么做  — Git 规范、Makefile 指南、环境搭建、契约层使用
├── testing/       # 测试    — 手动测试手册、测试流程、测试数据清理
└── backlog/       # 接下来  — 待办、设计笔记、RFC-like 文档（按需创建）
```

| 分类 | 语义 | AI Agent 何时查阅 |
|------|------|-------------------|
| `design/` | 系统如何构建 | 需要理解架构、修改实现方案、排查跨模块问题 |
| `reference/` | 事实性信息 | 需要查技术栈、API 端点、适配器差异、视觉规范 |
| `guides/` | 操作步骤 | 需要执行某项操作（发版、搭建环境、使用契约层） |
| `testing/` | 测试相关 | 需要手动验证功能、清理测试数据 |
| `backlog/` | 未落地的规划 | 需要了解待做事项、设计演进方向（按需创建） |

## 本目录文档索引

### design/

- [three-tier-design.md](design/three-tier-design.md) — 三层架构设计（React + Go + Python），冻结边界与状态迁移 authority

### reference/

- [skills.md](reference/skills.md) — Claude Code Skills 配置说明（openspec、design-taste、minimalist-ui 等）

### guides/

- [git-conventions.md](guides/git-conventions.md) — Git 提交规范（Conventional Commits）
- [makefile-guide.md](guides/makefile-guide.md) — Makefile 命令说明（三端服务管理）
- [contract-layer.md](guides/contract-layer.md) — 契约层使用指南（YAML → 三端代码生成）
- [setup.md](guides/setup.md) — 环境搭建
- [monorepo-setup.md](guides/monorepo-setup.md) — Monorepo 配置

### testing/

- [inactive-cleanup.md](testing/inactive-cleanup.md) — 会话停用功能的手动测试手册
- [code-audit-report.md](testing/code-audit-report.md) — 代码审计报告

### dev-plan/（独立保留）

- [dev-plan/](common/dev-plan/) — 开发路线图（Phase 1-4，串行叠代）

## 子项目 docs/

各子项目 `docs/` 同样遵循 5 分类，具体文档见各自的 AGENTS.md：

- [frontend/docs/](../frontend/docs/) — 前端设计、组件、数据流、主题、技术栈
- [backend/docs/](../backend/docs/) — 后端实现方案、技术栈
- [agentend/docs/](../agentend/docs/) — Agent 端架构、API 端点、适配器差异、测试流程、待办

## 新增文档规则

1. 按内容语义放入对应分类，不要新建子目录（平铺）
2. 文件名使用 kebab-case（如 `session-persistence.md`）
3. 子项目专属文档放子项目 `docs/`，跨端文档放根 `docs/`
4. 添加后在本文件索引中补充条目
