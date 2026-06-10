## MODIFIED Requirements

### Requirement: Diff block type carries snapshotId
diff block 类型 SHALL 从 `{ type: 'diff' }` 扩展为 `{ type: 'diff'; snapshotId: string }`。`block-types.ts` 中 MessageBlock 的 diff variant MUST 包含 `snapshotId` 字段。`block-reducer.ts` 解析 diff block 时 SHALL 从 block 内容中提取 `snapshotId` 字段。

#### Scenario: Parse diff block with snapshotId
- **WHEN** block-reducer 解析 ````aka_yhy\ntype: diff\nsnapshotId: abc-123\n````
- **THEN** SHALL 返回 `{ type: 'diff', snapshotId: 'abc-123' }`

#### Scenario: Diff block without snapshotId (backward compat)
- **WHEN** block-reducer 解析 ````aka_yhy\ntype: diff\n````（无 snapshotId 字段）
- **THEN** SHALL 返回 null（忽略此 block，降级为不渲染 diff 卡片）

## ADDED Requirements

### Requirement: BlockRenderer passes snapshotId to DiffCard
MessageBubble 的 BlockRenderer SHALL 将 diff block 的 `snapshotId` 传递给 DiffCard 组件作为 prop。

#### Scenario: Diff block rendering with snapshotId
- **WHEN** BlockRenderer 收到 `{ type: 'diff', snapshotId: 'abc-123' }` block 且 sessionId 存在
- **THEN** SHALL 渲染 `<DiffCard snapshotId="abc-123" sessionId={sessionId} />`
