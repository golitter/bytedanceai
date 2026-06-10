## MODIFIED Requirements

### Requirement: Dev-plan README 反映实际进度
`docs/common/dev-plan/README.md` SHALL 包含以下更新：
- 「当前状态」代码块 SHALL 显示 AgentEnd ~85%、Backend ~80%、Frontend ~70%
- 「阶段总览」表格 SHALL 包含 7 个 Phase（1-7），每个标记完成状态
- 「Phase 依赖关系」图 SHALL 反映 Phase 1-3 已完成、Phase 4 待执行、Phase 5-7 新增

#### Scenario: README 状态同步
- **WHEN** 开发者打开 `docs/common/dev-plan/README.md`
- **THEN** 看到当前状态反映 Phase 1-3 完成、Phase 4 待执行
- **AND** 阶段总览包含 Phase 1-7 的完整列表

#### Scenario: 依赖关系图完整性
- **WHEN** 开发者查看 Phase 依赖关系
- **THEN** 看到 Phase 4 → Phase 5 → Phase 6/7（Phase 6 和 Phase 7 可并行）
