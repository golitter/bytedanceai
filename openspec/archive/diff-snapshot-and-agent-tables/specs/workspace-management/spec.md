## ADDED Requirements

### Requirement: Agent diff block outputs snapshotId
Agent 的 render diff skill（`card_diff.go`）输出 diff block 标记时 SHALL 生成一个 UUID v4 作为 `snapshotId`，并写入 block 内容。输出格式为 ````aka_yhy\ntype: diff\nsnapshotId: {uuid}\n````。

#### Scenario: Diff block contains snapshotId
- **WHEN** agent 调用 `render diff` 命令
- **THEN** SHALL 输出 ````aka_yhy\ntype: diff\nsnapshotId: {UUID}\n````，其中 UUID 为新生成的 v4 UUID

#### Scenario: Each diff invocation generates unique snapshotId
- **WHEN** agent 在同一 session 中多次调用 `render diff`
- **THEN** 每次 SHALL 生成不同的 snapshotId UUID
