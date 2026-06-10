## ADDED Requirements

### Requirement: AGENTS.md doc links match actual file paths

所有 AGENTS.md 文件中引用的 `docs/` 路径 SHALL 与迁移后的实际文件位置一致。不存在指向已移除路径的死链接。

#### Scenario: Root AGENTS.md links
- **WHEN** 读取根 `AGENTS.md`
- **THEN** 所有 `docs/` 相对路径指向根 `docs/` 下存在的文件（如 `docs/common/git-conventions.md` 更新为 `docs/guides/git-conventions.md`）

#### Scenario: Frontend AGENTS.md links
- **WHEN** 读取 `frontend/AGENTS.md`
- **THEN** 所有 `docs/` 路径指向 `frontend/docs/` 下存在的文件（如 `docs/common/tech-stack.md` 更新为 `docs/reference/tech-stack.md`）

#### Scenario: Backend AGENTS.md links
- **WHEN** 读取 `backend/AGENTS.md`
- **THEN** 所有 `docs/` 路径指向 `backend/docs/` 下存在的文件

#### Scenario: Agentend AGENTS.md links
- **WHEN** 读取 `agentend/AGENTS.md`
- **THEN** 所有 `docs/` 路径指向 `agentend/docs/` 下存在的文件（如 `docs/common/details.md` 更新为 `docs/reference/details.md`）

#### Scenario: Contracts AGENTS.md links
- **WHEN** 读取 `contracts/AGENTS.md`
- **THEN** 该文件无 `docs/` 路径引用（contracts 无 docs 目录），内容不变

### Requirement: No stale references elsewhere

项目中其他文件（如 README、注释、配置）引用旧 docs 路径的 SHALL 一并更新。

#### Scenario: Cross-references updated
- **WHEN** 全局搜索 `docs/common/`、`docs/impl/`、`docs/playbooks/`、`docs/todos/`、`docs/test_guide/`、`docs/architecture/`
- **THEN** 除保留目录（`docs/internal/`、`docs/tmp/`、`docs/prompts/`、`docs/dev-plan/`）外，无匹配结果
