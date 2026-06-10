## ADDED Requirements

### Requirement: VSG 图标推荐 SHALL 与开发策略统一
`visual-style-guide.md` 中图标选择 MUST 与 `development-strategy.md` 保持一致，指定 Lucide React，删除 Phosphor Icons 推荐。

#### Scenario: VSG 指定 Lucide
- **WHEN** 阅读 visual-style-guide.md 图标部分
- **THEN** 图标选择为 Lucide React，strokeWidth 建议为 1.25（模拟细线视觉），无 Phosphor 推荐

#### Scenario: 不推荐列表更新
- **WHEN** 阅读 VSG "不选" 列表
- **THEN** Phosphor 不在 "不选" 列表中（因为 VSG 已统一为 Lucide）

### Requirement: api.ts 中的手写类型 SHALL 标注 TODO 或迁移
`lib/api.ts` 中手动定义的请求/响应类型如与 `generated/` 中契约类型重叠，MUST 添加 TODO 注释标记后续迁移，或在此次变更中直接迁移。

#### Scenario: 类型来源明确
- **WHEN** 查看 `lib/api.ts` 中的类型定义
- **THEN** 与 `generated/` 契约类型重叠的部分有 TODO 注释或已完成迁移
