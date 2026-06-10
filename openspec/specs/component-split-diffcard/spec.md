## ADDED Requirements

### Requirement: DiffCard 拆分为 DiffHeader 子组件
系统 SHALL 将 DiffCard.tsx 中的顶栏区块（文件统计摘要、视图模式切换按钮、编辑/接受/拒绝操作按钮、状态 badge）提取为独立的 `DiffHeader` 组件，位于 `components/diff/DiffHeader.tsx`。

DiffCard 保留快照 CRUD 逻辑和状态管理，通过 props 向 DiffHeader 传递数据和回调。

#### Scenario: DiffCard 行数低于 200
- **WHEN** 统计 `DiffCard.tsx` 的总行数
- **THEN** 行数 MUST 低于 200

#### Scenario: DiffHeader 接收正确的 props
- **WHEN** DiffHeader 渲染
- **THEN** 它接收 `summary`、`viewType`、`onViewTypeChange`、`snapshotStatus`、`isSettled`、`hasSession`、`editingFile`、`onEdit`、`onAccept`、`onReject`、`actionStatus` 等 props，并根据这些 props 正确渲染顶栏 UI

#### Scenario: 视图切换功能不变
- **WHEN** 用户点击 Split/Unified 切换按钮
- **THEN** diff 内容区在 split view 和 unified view 之间切换，行为不变

### Requirement: DiffCard 文件信息条提取为 DiffFileInfo
系统 SHALL 将 DiffCard.tsx 中的文件信息条（文件路径、变更类型 badge、增删统计）提取为独立的 `DiffFileInfo` 组件，位于 `components/diff/DiffFileInfo.tsx`。

#### Scenario: 文件信息条显示正确
- **WHEN** 用户查看 diff 中某个文件的详情
- **THEN** 文件路径、变更类型字母（A/D/M/R/C）、增删行数均正确显示，与拆分前完全一致

### Requirement: 拆分后组件行数约束
拆分后各组件文件 MUST 均低于 200 行：`DiffCard.tsx`、`DiffHeader.tsx`、`DiffFileInfo.tsx`。

#### Scenario: 所有拆分组件行数合规
- **WHEN** 统计 `DiffCard.tsx`、`DiffHeader.tsx`、`DiffFileInfo.tsx` 的行数
- **THEN** 每个文件行数均低于 200 行
