## ADDED Requirements

### Requirement: All docs directories use the 5-category structure

所有包含 `docs/` 目录的子项目（root、frontend、backend、agentend）SHALL 仅包含以下语义目录作为一级子目录：`design/`、`reference/`、`guides/`、`testing/`、`backlog/`。不受本变更影响的目录（`internal/`、`tmp/`、`prompts/`、`dev-plan/`）除外。

#### Scenario: Root docs structure
- **WHEN** 检查根 `docs/` 目录
- **THEN** 一级子目录仅包含 `design/`、`reference/`、`guides/`、`testing/`，以及保留的 `dev-plan/`、`internal/`、`prompts/`、`tmp/`，不存在 `architecture/`、`common/`、`test_guide/`

#### Scenario: Subproject docs structure
- **WHEN** 检查 `frontend/docs/`、`backend/docs/`、`agentend/docs/` 目录
- **THEN** 一级子目录仅包含 `design/`、`reference/`、`guides/`、`testing/`、`backlog/` 的子集，不存在 `common/`、`impl/`、`playbooks/`、`todos/`

### Requirement: design directory contains architecture and implementation docs

`design/` 目录 SHALL 仅包含描述系统如何构建的文档，包括架构设计、实现方案、数据流、组件设计、开发策略等。

#### Scenario: Root design directory
- **WHEN** 检查根 `docs/design/`
- **THEN** 包含 `three-tier-design.md`

#### Scenario: Frontend design directory
- **WHEN** 检查 `frontend/docs/design/`
- **THEN** 包含 `01-architecture.md`、`02-components.md`、`03-data-flow.md`、`04-theme.md`、`development-strategy.md`

#### Scenario: Backend design directory
- **WHEN** 检查 `backend/docs/design/`
- **THEN** 包含 `phase1-go-glue.md`

#### Scenario: Agentend design directory
- **WHEN** 检查 `agentend/docs/design/`
- **THEN** 包含 `architecture.md`、`01-schemas.md` 至 `11-orchestrator-planning.md`、`skills/taskctl.md`

### Requirement: reference directory contains factual lookup docs

`reference/` 目录 SHALL 仅包含事实性参考文档，包括技术栈说明、API 端点列表、适配器差异对照、视觉规范等。

#### Scenario: Root reference directory
- **WHEN** 检查根 `docs/reference/`
- **THEN** 包含 `skills.md`

#### Scenario: Subproject reference directories
- **WHEN** 检查 `frontend/docs/reference/`
- **THEN** 包含 `tech-stack.md`、`visual-style-guide.md`

- **WHEN** 检查 `backend/docs/reference/`
- **THEN** 包含 `tech-stack.md`

- **WHEN** 检查 `agentend/docs/reference/`
- **THEN** 包含 `details.md`、`adapter-diff.md`

### Requirement: guides directory contains how-to and convention docs

`guides/` 目录 SHALL 仅包含操作指南和规范文档，包括 Git 规范、Makefile 使用、环境搭建、契约层使用指南等。

#### Scenario: Root guides directory
- **WHEN** 检查根 `docs/guides/`
- **THEN** 包含 `git-conventions.md`、`makefile-guide.md`、`setup.md`、`monorepo-setup.md`、`contract-layer.md`

### Requirement: testing directory contains test manuals and procedures

`testing/` 目录 SHALL 仅包含测试手册、手动测试流程（playbook）、测试数据清理指南等。

#### Scenario: Root testing directory
- **WHEN** 检查根 `docs/testing/`
- **THEN** 包含 `inactive-cleanup.md`

#### Scenario: Agentend testing directory
- **WHEN** 检查 `agentend/docs/testing/`
- **THEN** 包含 `01-session-id-writeback.md`、`02-workspace-git-ops.md`、`03-taskctl-merge.md`、`04-orchestrator-planning.md`

### Requirement: backlog directory contains pending design notes

`backlog/` 目录 SHALL 仅包含待办、设计笔记、RFC-like 文档等尚未落地的规划内容。

#### Scenario: Agentend backlog directory
- **WHEN** 检查 `agentend/docs/backlog/`
- **THEN** 包含 `orchestrator-drawbacks.md`、`session-persistence.md`、`system-architecture-and-frontend-cards.md`

### Requirement: No stale empty directories remain

迁移完成后 SHALL 不存在空的旧目录（`common/`、`impl/`、`playbooks/`、`todos/`、`test_guide/`、`architecture/`）。

#### Scenario: Old directories removed
- **WHEN** 迁移完成
- **THEN** `docs/common/`、`docs/architecture/`、`docs/test_guide/`、`frontend/docs/common/`、`frontend/docs/impl/`、`backend/docs/common/`、`backend/docs/impl/`、`agentend/docs/common/`、`agentend/docs/impl/`、`agentend/docs/playbooks/`、`agentend/docs/todos/` 均不存在
