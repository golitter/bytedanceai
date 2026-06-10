## Why

dev-plan 文档仍停留在项目初期状态（Backend ~10%、Frontend ~5%），但实际 Phase 1-3 已全部完成，三端实现进度远超计划描述。需同步文档与现实状态，并根据 `docs/internal/任务要求.md` 的 6 大核心功能补齐 Phase 5-7，为后续 Orchestrator 群聊协作、产物预览、演示交付提供明确的执行路线。

## What Changes

- **更新 `docs/common/dev-plan/README.md`**：重写「当前状态」和「阶段总览」，反映 Phase 1-3 已完成、Phase 4 待执行、Phase 5-7 新增
- **新增 `phase5-orchestrator.md`**：Orchestrator 群聊协作实现计划（基于 `docs/internal/orchestrator-plan-phase.md` 已有设计）
- **新增 `phase6-preview-deploy.md`**：产物预览 + 部署发布 draft 框架（仅列功能范围，不做详细设计）
- **新增 `phase7-demo-deliver.md`**：演示打磨 + 交付物清单（对照任务要求的交付物列表）
- **Phase 1-3 文件不动**：已完成，内容仍具参考价值

## Capabilities

### New Capabilities
- `dev-plan-phase5`: Orchestrator 群聊协作的实现计划 — LangGraph 接入、plan → write_shared graph、OrchestratorAdapter、Go Scheduler 顺序调度、前端群聊 UI
- `dev-plan-phase6`: 产物预览 + 部署发布的 draft 框架 — Artifact Manager、ArtifactCard、iframe 预览、部署状态卡片
- `dev-plan-phase7`: 演示打磨 + 交付物清单 — 产品/技术文档整理、Demo 视频、AI 协作记录

### Modified Capabilities
- `dev-plan-roadmap`: README.md 当前状态从 "Phase 1 待开始" 更新为 "Phase 1-3 完成，Phase 4 进行中"；阶段总览从 4 Phase 扩展为 7 Phase

## Impact

- 文档变更：`docs/common/dev-plan/` 下 1 个更新 + 3 个新增
- 不涉及代码改动
- 不影响现有功能
- 为后续 `/opsx:apply` 实施 Phase 4-7 提供文档基础
