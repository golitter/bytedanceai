## ADDED Requirements

### Requirement: Phase 6 Draft 框架文档
`docs/common/dev-plan/phase6-preview-deploy.md` SHALL 包含产物预览 + 部署发布的 draft 框架，标记为 DRAFT 状态。

#### Scenario: Draft 文档存在
- **WHEN** 开发者打开 `phase6-preview-deploy.md`
- **THEN** 文档顶部标记 "DRAFT — 待 Phase 5 完成后细化"
- **AND** 文档包含功能范围列表和「不做」清单

### Requirement: Phase 6 功能范围明确
Phase 6 文档 SHALL 列出以下功能范围（不做详细设计）：
- Artifact Manager（AgentEnd 内存存储）
- ArtifactCard 前端组件（图片/代码/文件）
- iframe 网页预览（sandbox）
- 部署指令 + 状态卡片
- Go Backend Artifact 代理层

#### Scenario: 功能边界可追溯
- **WHEN** 开发者查看 Phase 6 文档
- **THEN** 能看到功能范围与任务要求 #4 #5 的对应关系
- **AND** 能看到明确的「不做」清单（Diff 视图、版本历史、容器化部署）

### Requirement: Phase 6 预估时间标记为 TBD
Phase 6 文档 SHALL 将预估时间标记为 TBD，不做精确预估。

#### Scenario: 时间不锁定
- **WHEN** 开发者查看阶段总览
- **THEN** Phase 6 的预估列显示 "TBD"
