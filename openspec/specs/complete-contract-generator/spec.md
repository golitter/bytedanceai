## ADDED Requirements

### Requirement: 契约生成器包含 message schema
`scripts/generate_contracts.py` 的 `OUTPUT_MAP` SHALL 包含 `message` 条目，将 `contracts/schemas/message.yaml` 映射到三端生成目录。

#### Scenario: make generate 生成 message 类型
- **WHEN** 执行 `make generate`
- **THEN** 三端 generated/ 目录中出现 Message、MessageRole、MessageStatus 相关类型定义

### Requirement: 契约生成器包含 validate-repo-path schema
`scripts/generate_contracts.py` 的 `OUTPUT_MAP` SHALL 包含 `validate-repo-path` 条目，将 `contracts/schemas/validate-repo-path.yaml` 映射到三端生成目录。

#### Scenario: make generate 生成 validate-repo-path 类型
- **WHEN** 执行 `make generate`
- **THEN** 三端 generated/ 目录中出现 ValidateRepoPathRequest、ValidateRepoPathResponse 相关类型定义

### Requirement: Session 状态机包含 inactive 转换
`agentend/src/session/models.py` 的 `_VALID_TRANSITIONS` SHALL 包含 `SessionState.INACTIVE: set()` 条目，允许会话进入 inactive 终态。

#### Scenario: 会话可转换为 inactive 状态
- **WHEN** 调用 `transition_to(SessionState.INACTIVE)`
- **THEN** 不抛出 ValueError，转换成功
