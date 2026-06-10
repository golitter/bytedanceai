## ADDED Requirements

### Requirement: Phase 7 演示打磨 + 交付物清单
`docs/common/dev-plan/phase7-demo-deliver.md` SHALL 包含演示打磨和交付物清单，对照 `docs/internal/任务要求.md` 的 5 类交付物。

#### Scenario: 交付物对照完整
- **WHEN** 开发者打开 `phase7-demo-deliver.md`
- **THEN** 文档列出 5 类交付物：产品设计文档、技术文档、可运行 Demo、AI 协作开发记录、3 分钟 Demo 视频
- **AND** 每类交付物有具体的产出项和负责端

### Requirement: Phase 7 包含 UI 打磨项
Phase 7 文档 SHALL 列出 UI 打磨清单：响应式适配、主题一致性、错误处理（Agent 断连/超时/重试）、Demo 场景脚本。

#### Scenario: UI 打磨项可执行
- **WHEN** 开发者按 Phase 7 执行
- **THEN** 能完成响应式适配、错误处理优化、Demo 预设对话流

### Requirement: Phase 7 预估 2 天
Phase 7 文档 SHALL 标注预估 2 天，前置条件为 Phase 5 完成（Phase 6 可选）。

#### Scenario: 时间和依赖明确
- **WHEN** 开发者查看阶段总览
- **THEN** Phase 7 预估 2 天，依赖 Phase 5（Phase 6 标记为可选）
