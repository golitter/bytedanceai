## Why

当前四端（root / frontend / backend / agentend）的 `docs/` 目录分类不统一（`common/`、`impl/`、`playbooks/`、`todos/`、`test_guide/` 等命名各异），导致 AI coding agent 难以按统一约定定位文档，且存在内容重复（如 `phase1-go-glue.md` 同时出现在根和 backend）。需要统一为语义一致的 5 分类结构，让所有子项目遵循相同的目录约定，并同步更新 AGENTS.md 中的文档链接。

## What Changes

- **统一 `docs/` 目录分类**：所有 `docs/` 采用 `design/` + `reference/` + `guides/` + `testing/` + `backlog/` 五分类法
- **迁移现有文档**：按语义将各端现有文档移入对应分类目录
- **重命名语义不清的目录**：`todos/` → `backlog/`，`impl/` → `design/`，`common/` → 按内容拆分到 `reference/` 或 `guides/`，`playbooks/` → `testing/`，`test_guide/` → `testing/`
- **更新 AGENTS.md 文档链接**：所有 AGENTS.md 中的 `docs/` 相对路径同步更新
- **处理重复内容**：确认 `docs/common/dev-plan/phase1-go-glue.md` 与 `backend/docs/impl/phase1-go-glue.md` 的关系，合并或去重

不涉及的目录（不动）：`docs/internal/`、`docs/tmp/`、`docs/prompts/`、`docs/dev-plan/`。

## Capabilities

### New Capabilities

- `docs-unified-structure`: 统一五端 docs 目录分类规范，定义 `design/`、`reference/`、`guides/`、`testing/`、`backlog/` 的语义和内容边界
- `docs-agents-md-sync`: AGENTS.md 文档链接与实际 docs 目录的一致性维护

### Modified Capabilities

（无既有 spec 受影响）

## Impact

- **文件系统**：约 30+ 个 markdown 文件移动位置，5+ 个旧目录删除
- **AGENTS.md**：5 个 AGENTS.md 文件（root + 4 子项目）中的文档链接路径需更新
- **Git 历史**：文件移动可能影响 `git blame` 追溯，建议在单次 commit 中完成
- **无代码影响**：纯文档变更，不影响任何运行时行为
