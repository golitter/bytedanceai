## Why

discussions/ 目录下有 4 篇核心设计文档（conversation-layer-design.md、agentend-group-chat-impl.md、current-status-and-next-steps.md、orchestrator-features.md），定义了 Phase 5 多 Agent 编排的完整架构。在开始实现前，需要对这些文档与实际代码之间的一致性、文档间的矛盾、以及设计盲点进行全面审查，避免实施阶段踩坑。

## What Changes

- **新增**: `discussions/design-review-issues.md` — 设计文档审查报告，记录 15 类问题及建议修复方案
- **不涉及代码修改** — 纯文档审查，仅产出问题清单

## Capabilities

### New Capabilities

无新能力引入。本次变更是文档审查，不引入新功能模块。

### Modified Capabilities

无。

## Impact

- **文档**: `discussions/` 目录下新增审查报告
- **后续工作**: 审查发现的问题将影响 Phase 5 实施的优先级和方案调整
