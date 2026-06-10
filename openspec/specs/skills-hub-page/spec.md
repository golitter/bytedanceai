## ADDED Requirements

### Requirement: SkillsHub management page
系统 SHALL 在侧边栏新增「技能」NavTab（位于通讯录和管理之间），点击展示 SkillsHub 管理页面。

#### Scenario: Navigate to skills page
- **WHEN** 用户点击侧边栏「技能」Tab
- **THEN** 导航到 SkillsHub 管理页面，展示所有 builtin + external skill 列表

#### Scenario: NavTab type extension
- **WHEN** 系统初始化
- **THEN** `NavTab` 类型包含 `'skills'`，位于 `'contacts'` 和 `'admin'` 之间

### Requirement: Skill list display
SkillsHub 页面 SHALL 分区展示 builtin skills 和 external skills，每个 skill card 显示名称、描述、来源标签（builtin/external）、已被导入的 Agent 数量。Builtin 条目只读无删除按钮，External 条目有删除按钮。

#### Scenario: Display builtin skill card
- **WHEN** 渲染一个 builtin skill（如 render）
- **THEN** 显示绿色 `builtin` 标签、名称、描述、导入数量，无删除按钮

#### Scenario: Display external skill card
- **WHEN** 渲染一个 external skill（如 my-custom-skill）
- **THEN** 显示 indigo `external` 标签、名称、描述、导入数量、删除按钮

### Requirement: Skill search filter
SkillsHub 页面 SHALL 提供搜索框，输入时实时过滤显示匹配名称的 skill cards。

#### Scenario: Filter by name
- **WHEN** 用户在搜索框输入 "api"
- **THEN** 仅显示名称含 "api" 的 skill cards

#### Scenario: Empty search shows all
- **WHEN** 搜索框为空
- **THEN** 显示所有 skill cards

### Requirement: Upload skill dialog
点击「上传 Skill」按钮 SHALL 弹出上传对话框，包含：虚线上传区域（支持点击和拖拽）、校验结果面板、名称确认输入框、确认/取消按钮。

#### Scenario: Show upload dialog
- **WHEN** 用户点击「上传 Skill」按钮
- **THEN** 弹出 modal 对话框，显示上传区域

#### Scenario: Validation success display
- **WHEN** 上传的 zip 校验通过
- **THEN** 显示绿色校验结果面板（SKILL.md 存在、frontmatter 正确、大小/文件数合规、安全检查通过）和名称确认输入框

#### Scenario: Validation failure display
- **WHEN** 上传的 zip 校验失败
- **THEN** 显示红色校验失败面板，列出具体错误

#### Scenario: Confirm upload with name conflict
- **WHEN** 用户确认的名称与 builtin skill 同名
- **THEN** 前端拦截并提示「名称与内置 Skill 冲突」，不发请求

#### Scenario: Successful confirm
- **WHEN** 用户确认入库成功
- **THEN** 关闭对话框，新 skill card 以动画出现在列表中

### Requirement: Delete skill confirmation
点击 external skill 的「删除」按钮 SHALL 弹出确认对话框，说明仅删除 hub 源文件不影响已导入副本。

#### Scenario: Show delete confirmation
- **WHEN** 用户点击 external skill 的「删除」按钮
- **THEN** 弹出确认对话框，显示 skill 名称和影响说明

#### Scenario: Confirm delete
- **WHEN** 用户确认删除
- **THEN** 发送 DELETE 请求，成功后该 card 以动画移除，显示成功 Toast

### Requirement: SkillCard external tag style
SkillCard 组件 SHALL 增加 external 标签样式：`background: rgba(99,102,241,0.08); color: #6366F1`（与 builtin 的绿色标签对应）。

#### Scenario: Builtin tag style
- **WHEN** 渲染 builtin skill 标签
- **THEN** 使用绿色样式（`background: rgba(34,197,94,0.08); color: #22C55E`）

#### Scenario: External tag style
- **WHEN** 渲染 external skill 标签
- **THEN** 使用 indigo 样式（`background: rgba(99,102,241,0.08); color: #6366F1`）
