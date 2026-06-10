## ADDED Requirements

### Requirement: 空 store 文件 SHALL 被清理
`stores/app.ts` 如无实际状态定义，MUST 删除或在其中补充注释说明保留理由。

#### Scenario: 空 Zustand store 处理
- **WHEN** `stores/app.ts` 中无已定义的 state 或 action
- **THEN** 该文件被删除，且无其他文件导入它

### Requirement: 未使用组件 SHALL 被清理
导出但从未被其他组件导入的组件 MUST 删除，避免误导开发者。

#### Scenario: SessionList 未使用清理
- **WHEN** `SessionList.tsx` 无任何导入引用
- **THEN** 该文件被删除

### Requirement: 未使用的 API mutations SHALL 被清理
`use-sessions.ts` 中未被任何组件调用的 mutation（如 delete、update）MUST 删除，减少认知负担。

#### Scenario: 未使用 mutation 清理
- **WHEN** 检查 `use-sessions.ts` 导出的 mutation
- **THEN** 仅保留被组件实际使用的 mutation

### Requirement: 未使用的 store 字段 SHALL 被清理
`stores/chat.ts` 中 `currentTaskId` 如无消费方，MUST 删除。

#### Scenario: currentTaskId 无消费方
- **WHEN** `currentTaskId` 在 `stores/chat.ts` 中定义但无组件通过 selector 订阅它
- **THEN** 从 store 中删除该字段

### Requirement: 冗余常量定义 SHALL 合并
`AgentAvatar.tsx` 中的 `AGENT_LABELS` 与 `lib/constants.ts` 中的 `AGENT_NAMES` 表达相同语义，MUST 合并到 `lib/constants.ts`。

#### Scenario: AgentAvatar 使用统一常量
- **WHEN** AgentAvatar 需要 Agent 名称标签
- **THEN** 从 `lib/constants.ts` 导入 `AGENT_NAMES`，不内联定义 `AGENT_LABELS`
