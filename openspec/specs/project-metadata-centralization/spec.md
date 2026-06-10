## ADDED Requirements

### Requirement: 项目元信息 SHALL 集中定义在 constants.ts
`lib/constants.ts` SHALL 导出 `PROJECT_META` 常量对象，包含以下字段：
- `GITHUB_URL`：项目 GitHub 仓库地址
- `NAME`：项目展示名称（如 'AgentHub'）
- `DESCRIPTION`：项目简短描述（中英双语可选）

所有消费组件 MUST 从 `lib/constants.ts` 导入使用，不得内联定义。

#### Scenario: ContactsPage 使用集中常量显示 GitHub 链接
- **WHEN** ContactsPage 渲染项目 GitHub 链接
- **THEN** 从 `lib/constants.ts` 导入 `PROJECT_META.GITHUB_URL`，不内联 `https://github.com/...` 字符串

#### Scenario: ContactsPage 使用集中常量显示项目描述
- **WHEN** ContactsPage 渲染项目描述文案
- **THEN** 从 `lib/constants.ts` 导入 `PROJECT_META.DESCRIPTION`，不内联描述文本

#### Scenario: 新增项目元信息字段只需修改一处
- **WHEN** 需要在新组件中显示项目名称或描述
- **THEN** 从 `lib/constants.ts` 导入 `PROJECT_META` 即可，无需重复定义
