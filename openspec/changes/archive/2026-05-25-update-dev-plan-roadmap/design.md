## Context

项目 dev-plan 文档创建于开发初期，描述 Backend ~10%、Frontend ~5%。经过持续开发，Phase 1（Go 胶水层 SSE 透传）、Phase 2（最小聊天界面）、Phase 3（IM 体验补全）已全部完成，三端实现进度远超文档描述。

同时，`docs/internal/任务要求.md` 定义了 6 大核心功能（IM 交互、Orchestrator、多 Agent、产物预览、部署发布、多端支持），当前 Phase 1-4 计划仅覆盖前 3 项的部分功能，缺少 Orchestrator 群聊协作、产物预览/部署、演示交付等关键阶段的设计文档。

约束：
- 只做 Web 端，不做桌面端/移动端
- 先完成 Phase 4 再集中做 Orchestrator（串行）
- Phase 6 只列 draft，后续再定
- Phase 5 核心内容已有设计文档 `docs/internal/orchestrator-plan-phase.md`

## Goals / Non-Goals

**Goals:**
- 同步 dev-plan README.md 状态为实际进度
- 为 Phase 5（Orchestrator）创建基于已有设计文档的详细实现计划
- 为 Phase 6（预览+部署）创建 draft 框架，明确功能边界
- 为 Phase 7（演示+交付）创建交付物清单，对照任务要求
- 建立 Phase 1-7 的完整依赖关系图

**Non-Goals:**
- 不改动 Phase 1-3 的已有文档内容
- 不实现任何代码功能
- 不为 Phase 6 做详细技术设计（仅 draft）
- 不涉及桌面端/移动端规划

## Decisions

### D1: README.md 结构保持不变，只更新数据

保留原有的「当前状态」代码块 +「阶段总览」表格 +「Phase 依赖关系」图的结构，只更新其中的数值和新增 Phase 行。理由：团队已熟悉此格式，无需重写。

### D2: Phase 5 基于 orchestrator-plan-phase.md 派生

`docs/internal/orchestrator-plan-phase.md` 已包含完整的 LangGraph 设计（models / prompts / graph / adapter / 注册流程）。Phase 5 文档将其转化为「面向开发者按天执行」的实施计划，补充 Go Backend Scheduler 和前端群聊 UI 两部分（原文档未覆盖）。

### D3: Phase 6 采用「功能范围 + 不做清单」格式

Phase 6 仅列出功能范围（Artifact Manager、ArtifactCard、iframe 预览、部署卡片）和明确的「不做」清单（Diff 视图、版本历史、容器化部署）。详细设计留待 Phase 5 完成后补充。

### D4: Phase 7 对照任务要求的交付物逐条列出

任务要求明确了 5 类交付物（产品设计文档、技术文档、可运行 Demo、AI 协作开发记录、3 分钟 Demo 视频）。Phase 7 文档逐条对应，补充 UI 打磨和稳定性保障项。

## Risks / Trade-offs

- [Phase 5 预估可能偏紧] → Orchestrator 涉及三端改动（AgentEnd LangGraph + Go Scheduler + 前端群聊 UI），3-4 天预估需要密集输出。缓解：已有完整设计文档，减少设计决策时间。
- [Phase 6 draft 可能需要大改] → 功能边界未定，后续可能需要重新评估范围。缓解：明确标记为 draft，不做过度设计。
- [Phase 7 时间被压缩] → 如果 Phase 5-6 超时，演示打磨时间会被压缩。缓解：每个 Phase 结束都有可演示成果，最坏情况下 Phase 7 可缩减为 1 天。
