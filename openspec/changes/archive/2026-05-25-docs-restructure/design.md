## Context

Monorepo 项目包含 root / frontend / backend / agentend / contracts 五端，其中四端有 `docs/` 目录。当前各端 docs 内部组织方式不同：

| 端 | 现有分类 | 问题 |
|---|---------|------|
| root | `architecture/`, `common/`, `test_guide/`, `internal/`, `prompts/`, `tmp/`, `dev-plan/` | `common/` 语义模糊，混合了指南和参考 |
| frontend | `common/`, `impl/` | 只有 7 个文件，分类尚可但命名不统一 |
| backend | `common/`, `impl/` | 只有 2 个文件，与根 dev-plan 有重复 |
| agentend | `common/`, `impl/`, `playbooks/`, `todos/` | 最复杂，`todos/` 实质是 backlog/design notes |

约束：`docs/internal/`、`docs/tmp/`、`docs/prompts/`、`docs/dev-plan/` 不动。

## Goals / Non-Goals

**Goals:**

- 所有 `docs/` 采用统一的 5 分类结构：`design/`、`reference/`、`guides/`、`testing/`、`backlog/`
- 每个分类的语义明确，AI agent 可按目录名预测内容类型
- 所有 AGENTS.md 中的文档链接指向新路径
- 清除旧的空目录

**Non-Goals:**

- 不修改文档内容本身（除非路径引用）
- 不创建新文档
- 不动 `internal/`、`tmp/`、`prompts/`、`dev-plan/`
- 不处理 contracts 目录（无 docs）

## Decisions

### Decision 1: 5 分类法

采用 `design/` + `reference/` + `guides/` + `testing/` + `backlog/` 统一分类。

| 目录 | 语义 | 典型内容 |
|------|------|---------|
| `design/` | 怎么建的 | 架构设计、实现方案、数据流、组件设计 |
| `reference/` | 是什么 | 技术栈、API 端点、适配器差异、skills 列表、视觉规范 |
| `guides/` | 怎么做 | Git 规范、Makefile 指南、环境搭建、契约层使用指南 |
| `testing/` | 测试 | 手动测试手册、playbook、测试数据清理 |
| `backlog/` | 接下来 | 待办、设计笔记、RFC-like 文档 |

**替代方案**：保持各端独立命名（rejected — 无法形成统一约定，AI agent 每次都要学新的目录语义）。

### Decision 2: 平铺而非嵌套

`testing/` 目录下直接放测试文档，不再按端或类型嵌套子目录。文件数量少（每端 <10 个），嵌套反而增加认知成本。

### Decision 3: agentend playbooks 归入 testing

agentend 的 `playbooks/` 内容均为手动测试流程（curl 步骤 + 预期结果），语义更接近 `testing/` 而非 `guides/`。

### Decision 4: phase1-go-glue 重复处理

`docs/common/dev-plan/phase1-go-glue.md` 是全局路线图视角，`backend/docs/impl/phase1-go-glue.md` 是具体实现细节。两边都保留，分别迁入各自的新路径。

### Decision 5: 单次 commit 完成

所有文件移动 + AGENTS.md 更新在单次 commit 中完成，保持原子性。

## Risks / Trade-offs

- **Git blame 受影响** → 文件移动后 `git blame` 需要额外参数追踪。Mitigation: commit message 中列出移动的文件和原路径
- **外部链接失效** → 如果有 README 或其他文件引用了旧 docs 路径会失效。Mitigation: 全局 grep 更新所有引用
- **短期内目录结构变化** → 团队成员需要重新熟悉。Mitigation: 结构语义清晰，学习成本低
