## ADDED Requirements

### Requirement: 按钮 SHALL 使用 text-primary-foreground 替代 text-white
所有使用 `bg-primary` 背景的品牌色按钮 SHALL 使用 `text-primary-foreground` 作为文字颜色，MUST NOT 使用 `text-white`。这确保亮色主题下的自动适配。

#### Scenario: SkillsHubPage 上传按钮文字色
- **WHEN** SkillsHubPage 渲染「上传 Skill」按钮
- **THEN** 按钮文字颜色为 `text-primary-foreground`，不是 `text-white`

#### Scenario: SkillsHubPage 确认入库按钮文字色
- **WHEN** SkillsHubPage UploadDialog 渲染「确认入库」按钮
- **THEN** 按钮文字颜色为 `text-primary-foreground`，不是 `text-white`

#### Scenario: AgentProfilePage 按钮文字色
- **WHEN** AgentProfilePage 渲染品牌色背景按钮
- **THEN** 按钮文字颜色为 `text-primary-foreground`，不是 `text-white`

### Requirement: 错误/危险状态 SHALL 使用 text-destructive 语义类
组件中用于表示错误、删除、危险操作的红色文字和边框 SHALL 使用 `text-destructive`、`border-destructive`、`bg-destructive` 语义类，MUST NOT 使用 `text-red-500`、`border-red-500`、`bg-red-500` 等 Tailwind 直接色值。

#### Scenario: SkillsHubPage 删除按钮
- **WHEN** HubSkillCard 渲染删除按钮
- **THEN** 按钮使用 `text-destructive`、`border-destructive/20`、`bg-destructive/10`，不使用 `text-red-500`

#### Scenario: DeleteConfirmDialog 删除确认按钮
- **WHEN** DeleteConfirmDialog 渲染确认删除按钮
- **THEN** 按钮使用 `text-destructive`、`border-destructive/20`、`bg-destructive/10`

#### Scenario: SkillsHubPage 校验失败提示
- **WHEN** UploadDialog 渲染校验失败提示区域
- **THEN** 使用 `text-destructive`、`border-destructive/20`、`bg-destructive/5`

#### Scenario: DeleteConfirmDialog 危险名称高亮
- **WHEN** DeleteConfirmDialog 渲染被删除 Skill 名称
- **THEN** 名称使用 `text-destructive`，不使用 `text-red-500`

### Requirement: inline style 文字色 SHALL 替换为 Tailwind 类
组件中用于设置静态颜色的 inline `style={{ color: 'var(--text-secondary)' }}` SHALL 替换为 `text-text-secondary` Tailwind 类。

#### Scenario: SessionCleanupPage 状态文字色
- **WHEN** SessionCleanupPage 渲染 session 状态文字
- **THEN** 使用 `text-text-secondary` Tailwind 类，不使用 `style={{ color: 'var(--text-secondary)' }}`

### Requirement: 修改后零违规残留
修改完成后，在 `frontend/src/` 下的所有 TSX 文件中，SHALL 不存在以下违规模式：`text-white`（在 bg-primary 按钮中）、`text-red-500`、`border-red-500`、`bg-red-500`、inline style 中的 `var(--text-secondary)`。

#### Scenario: grep 验证零 text-white 违规
- **WHEN** 对 `frontend/src/**/*.tsx` 执行 grep 搜索 `text-white`
- **THEN** 搜索结果为零匹配（排除注释和字符串字面量）

#### Scenario: grep 验证零直接色值违规
- **WHEN** 对 `frontend/src/**/*.tsx` 执行 grep 搜索 `text-red-500`、`border-red-500`、`bg-red-500`
- **THEN** 搜索结果为零匹配
