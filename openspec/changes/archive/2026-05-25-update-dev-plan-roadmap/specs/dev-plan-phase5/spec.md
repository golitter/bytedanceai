## ADDED Requirements

### Requirement: Phase 5 Orchestrator 实现计划文档
`docs/common/dev-plan/phase5-orchestrator.md` SHALL 包含完整的 Orchestrator 群聊协作实现计划，基于 `docs/internal/orchestrator-plan-phase.md` 设计文档。

#### Scenario: 文档结构完整性
- **WHEN** 开发者打开 `phase5-orchestrator.md`
- **THEN** 文档包含：目标、交付标准、实现步骤（9 步）、文件清单、验证流程

#### Scenario: 三端实现覆盖
- **WHEN** 开发者按计划执行
- **THEN** 计划覆盖 AgentEnd（LangGraph + OrchestratorAdapter）、Go Backend（Scheduler 顺序调度）、Frontend（群聊 UI）三端

### Requirement: Phase 5 包含 AgentEnd LangGraph 实现
Phase 5 文档 SHALL 详细描述 AgentEnd 侧的 LangGraph 实现，包括 models.py、prompts.py、graph.py、OrchestratorAdapter、schemas 扩展、注册接入。

#### Scenario: LangGraph 实现步骤可执行
- **WHEN** 开发者按 Step 1-7 执行 AgentEnd 改动
- **THEN** 能完成 LangGraph 接入、plan → write_shared graph、OrchestratorAdapter 注册

### Requirement: Phase 5 包含 Go Backend Scheduler
Phase 5 文档 SHALL 描述 Go Backend 新增 Scheduler 功能：读取 `shared/.agent/tasks/*.md`，按 config.yaml 中数组顺序依次调度对应 Agent 执行任务。

#### Scenario: Scheduler 顺序调度
- **WHEN** Go Scheduler 读取 config.yaml 中 tasks 数组
- **THEN** 按 task-001 → task-002 顺序依次调用对应 Agent 的 stream 接口
- **AND** 将每个 Agent 的 SSE 流透传给前端

### Requirement: Phase 5 包含前端群聊 UI
Phase 5 文档 SHALL 描述前端群聊模式 UI：Orchestrator 消息带来源 Agent 标签、不同 Agent 用颜色区分、协作流展示。

#### Scenario: 群聊消息区分
- **WHEN** 用户选择 orchestrator 类型并发送消息
- **THEN** 聊天流中显示 Orchestrator 规划过程
- **AND** 每个 Agent 执行结果带 Agent 名称标签和颜色区分
