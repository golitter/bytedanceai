## Why

项目文档（`AGENTS.md` × 6 + 根 `docs/` + 三端 `docs/`）在过去多个迭代中持续累积变化，但缺乏一次系统性的对照实际代码的"对齐回合"。最近一次大规模文档结构重整发生在 `2026-05-25-docs-restructure`（统一 5 分类法），此后又新增了 `agent-memory`、`skill-import-export`、`theme-switching`、`git-graph-panel`、`contacts-page` 等大量能力，对应的 AGENTS.md 目录树、details.md 索引、`docs/design/` 实施文档是否仍然准确同步需要验证。

项目已经定义了两个 lint skill：
- `.agents/skills/agentsmd-linter/SKILL.md` — 全项目 AGENTS.md 同步精炼
- `.agents/skills/doc-linter/SKILL.md` — 全项目文档同步精炼

但目前缺少一次"按这两个 skill 的 workflow 完整跑一遍"的变更记录与验收点。本次 change 用一次结构化的 lint pass 收口：以这两个 skill 的规则为单一事实来源，对照实际代码（Go / Python / TypeScript / Makefile / package manifests）逐文件验证并修复，确保文档与代码同步至 `2026-06-09` 当前状态。

## What Changes

- **执行 `agentsmd-linter` Phase 1–4 完整流程**
  - 验证根 / `frontend` / `backend` / `agentend` / `docs` / `contracts` 共 6 个 AGENTS.md 的目录结构、Makefile 命令、文件引用
  - 验证三端 `docs/reference/details.md` 索引完整性
  - 行数硬约束：`AGENTS.md ≤ 70 行`，`details.md ≤ 150 行`
- **执行 `doc-linter` Phase 1–4 完整流程**
  - 先三端后根：`frontend/docs/` → `backend/docs/` → `agentend/docs/` → `docs/`
  - 交叉验证文档中的代码引用（文件路径、模块名、API 端点、技术栈、目录结构）
  - 三端 `docs/design/` 缺失实施文档按 OpenSpec 格式补充（`## 实现了什么` + `## 怎么实现的` + 实际代码片段）
- **修复类型**：路径过期 / 技术栈版本漂移 / 目录树过时 / 文档间描述矛盾 / 过时 TODO 残留 / 缺失的 design/ 实施文档
- **不涉及的目录**：`docs/internal/`、`docs/tmp/`、`AGENTS.md`（被两个 skill 列为只读或仅按 skill 规则可改）

## Capabilities

### New Capabilities

- `docs-agents-sync`: 定义 AGENTS.md × 6 + 根 docs/ + 三端 docs/ 与实际代码保持一致的同步约束、行数硬上限、文档分类语义与验收点

### Modified Capabilities

（无既有 spec 的需求级变更；`doc-consistency-fix` 是一次性历史变更，不在本次需求修改范围）

## Impact

- **文档文件**：预计修改 6 个 AGENTS.md + 3 个 details.md + 根 `docs/` 与三端 `docs/` 下若干 `.md` 文件；可能新增三端 `docs/design/` 下编号实施文档
- **代码引用**：使用 markdown 链接格式 `[file](path)`，便于 IDE 跳转
- **无代码影响**：纯文档变更，不修改 `.go` / `.py` / `.ts` / 依赖清单 / 配置文件
- **无运行时影响**：不影响三端服务运行、契约生成或测试
