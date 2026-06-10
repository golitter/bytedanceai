## ADDED Requirements

### Requirement: Lucide 图标 strokeWidth SHALL 统一为 1.25
所有使用 Lucide React 图标的组件 SHALL 为图标设置 `strokeWidth={1.25}` 或省略 strokeWidth（当全局默认为 1.25 时）。组件 MUST NOT 使用 `strokeWidth={1.5}`、`strokeWidth="3"` 或其他非 1.25 的值。

#### Scenario: RightSidebar 图标 strokeWidth
- **WHEN** RightSidebar 渲染 PanelRightOpen 或其他 Lucide 图标
- **THEN** 图标的 strokeWidth 为 1.25

#### Scenario: PlanReviewCard 编辑图标 strokeWidth
- **WHEN** PlanReviewCard 渲染 PencilLine 图标
- **THEN** 图标的 strokeWidth 为 1.25

#### Scenario: ConversationItem Pin 图标 strokeWidth
- **WHEN** ConversationItem 渲染 Pin 图标
- **THEN** 图标的 strokeWidth 为 1.25

#### Scenario: ContactsPage 图标 strokeWidth
- **WHEN** ContactsPage 渲染任意 Lucide 图标
- **THEN** 图标的 strokeWidth 为 1.25

#### Scenario: AgentProfilePage Line 元素 strokeWidth
- **WHEN** AgentProfilePage 渲染 Line 图形元素
- **THEN** strokeWidth 不使用字符串 "3"，应为合理的数值

### Requirement: 新模块 SHALL 使用 Lucide 图标替代 emoji
SkillsHubPage 及其他新增功能模块 SHALL 使用 Lucide React 图标替代 emoji 字符作为图标元素。具体映射：🛡️ → Shield、📦 → Package、⚙️ → Wrench、⚠️ → AlertTriangle、✅ → CheckCircle2、❌ → XCircle。

#### Scenario: SkillsHubPage SectionLabel 图标
- **WHEN** SkillsHubPage 渲染 SectionLabel（内置/外部分区标签）
- **THEN** 使用 Lucide Shield 和 Package 图标替代 emoji，strokeWidth 为 1.25

#### Scenario: SkillsHubPage HubSkillCard 头像图标
- **WHEN** HubSkillCard 渲染 Skill 类型图标（内置/外部）
- **THEN** 使用 Lucide Wrench 或 Package 图标替代 emoji

#### Scenario: SkillsHubPage 校验结果图标
- **WHEN** UploadDialog 渲染校验成功/失败状态
- **THEN** 使用 Lucide CheckCircle2 或 XCircle 图标替代 emoji

#### Scenario: DeleteConfirmDialog 警告图标
- **WHEN** DeleteConfirmDialog 渲染删除确认警告
- **THEN** 使用 Lucide AlertTriangle 图标替代 emoji

### Requirement: 全站零 emoji 图标残留
修改完成后，在 `frontend/src/` 下的所有 TSX 文件中，SHALL 不存在用于装饰/图标目的的 emoji 字符（文字内容中的 emoji 不受限）。

#### Scenario: grep 验证零残留
- **WHEN** 对 `frontend/src/**/*.tsx` 执行 grep 搜索装饰性 emoji（🛡️📦⚙️⚠️✅❌）
- **THEN** 搜索结果为零匹配
