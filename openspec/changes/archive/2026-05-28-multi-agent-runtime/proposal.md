## Why

AgentEnd 的 OrchestratorAdapter 已有完整的 LangGraph Planner（规划阶段），但 dispatch 阶段使用 mock 硬编码返回 `(mock) Task dispatched to {agent}`，未调用真实 Agent。多 Agent 协作的核心链路——Planner 生成计划 → 真实 Agent 执行 → 事件流回前端——尚未打通。设计文档已就绪（phase5-discussions 5 篇），代码 gap 集中且明确，现在需要把 mock 替换为真实运行时，让多 Agent 协作真正跑起来。

## What Changes

### 轨道 A：AgentEnd 运行时（让多 Agent 真的能跑）

1. **替换 mock dispatch → 真实 Agent 执行**（A1）
   - OrchestratorAdapter.dispatch 阶段调用 `AdapterRegistry.get(agent_type)` 获取 adapter 实例
   - 调用 `adapter.stream_chat()` 执行真实任务，收集 StreamEvent
   - 将 agent 的 StreamEvent 归一化后 yield 回 OrchestratorAdapter 的 SSE 流
   - 串行调度：按 plan 任务顺序依次执行，前一个完成再启动下一个

2. **扩展 RuntimeEvent 事件类型**（A2）
   - 在现有 8 种事件基础上新增协调相关事件：`runtime.executing`、`runtime.completed`、`coordination.start`、`coordination.message`、`coordination.done`
   - 事件命名统一使用 dot notation（`runtime.executing`）
   - 更新 contracts/schemas YAML 并重新生成三端类型

3. **简化版协调通道**（A3）
   - 项目经理向 agent 发起单轮 Q&A（不做多轮协商）
   - 使用 `adapter.chat()` 非流式调用获取回答
   - 协调结果作为 Planner 的补充输入，影响任务分配决策

### 轨道 B：Backend + Frontend（让用户能看到）

4. **Backend SSE 透传新事件**（B1）
   - SSE 端点不对事件类型做过滤，新增事件类型直接透传到 Redis Stream → 前端
   - 最小改动，可能仅需确认现有逻辑已支持

5. **前端 ChatArea 渲染扩展**（B2）
   - 根据 event.type 渲染不同 UI 组件：
     - `planning` → Plan Card（任务分配卡片，含 agent 标签 + 任务状态）
     - `tool_call` / `tool_result` → Tool Card（工具调用展示）
     - `runtime.*` → Status Badge（执行状态徽章）
     - `coordination.*` → Coordination Channel（可折叠协调通道）
   - 设计参考：`docs/common/dev-plan/phase5-notes/demo.html`

## Capabilities

### New Capabilities
- `runtime-execution-engine`: AgentEnd 执行引擎 — 替换 mock dispatch，调用真实 adapter，串行调度任务执行，归一化事件输出
- `runtime-events`: RuntimeEvent 事件体系 — 扩展事件类型，覆盖规划/协调/执行/合并全生命周期
- `coordination-channel`: 简化版协调通道 — 项目经理向 agent 单轮 Q&A，收集反馈后调整计划
- `frontend-runtime-timeline`: 前端 Runtime 渲染 — Plan Card、Tool Card、Status Badge、Coordination Channel 等 UI 组件

### Modified Capabilities
- `orchestrator-adapter`: OrchestratorAdapter — dispatch 阶段从 mock 改为真实执行，stream_chat() 输出扩展事件
- `backend-sse`: Backend SSE 转发 — 确保新事件类型透传无过滤

## Impact

- **Agent 端**: 核心改动区。OrchestratorAdapter 重构 dispatch/collect 阶段，新增 RuntimeEvent 类型，新增 CoordinationChannel 类
- **后端**: 最小改动。确认 SSE 透传逻辑，可能需要更新 generated 事件类型常量
- **前端**: ChatArea 组件扩展，新增多种事件渲染组件（PlanCard、ToolCard、StatusBadge、CoordChannel）
- **契约层**: 新增事件类型 YAML 定义，重新生成三端类型代码
- **依赖**: 无新外部依赖，复用现有 adapter/registry/workspace 组件

## Non-goals

- 多轮协商协调（Phase A3 只做单轮 Q&A）
- 并行任务调度（先串行执行）
- MergeManager 多分支合并（使用 WorkspaceManager 现有 merge）
- 引用/回复/重新生成功能（属于 Conversation Layer，后续迭代）
- 前端 @mention 交互（属于 Routing Layer，后续迭代）
- agent profile / soul.yaml 系统（使用现有 Backend 动态管理）
