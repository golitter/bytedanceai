## Task Breakdown

按依赖顺序排列。每个 task 预估 1-3 小时。

---

- [x] **T1: 扩展 EventType 枚举（契约层）**
  - 预估: 1h | Capability: runtime-events | 依赖: 无
  - 在 `contracts/schemas/event-types.yaml` 新增 5 种事件类型：`runtime_executing`、`runtime_completed`、`coordination_start`、`coordination_message`、`coordination_done`，每种定义 content 字段结构
  - 在 `contracts/logs/` 写入变更记录
  - 运行 `make generate` 生成三端类型（agentend/src/generated/events.py、backend/internal/generated/events.go、frontend/src/generated/events.ts）
  - 验证：三端 generated 文件均包含新事件类型，编译不报错

- [x] **T2: 新增 ExecutionEngine（AgentEnd）**
  - 预估: 3h | Capability: runtime-execution-engine | 依赖: T1
  - 新建 `agentend/src/orchestrator/execution.py`，实现 ExecutionEngine 类
  - 核心方法 `execute(dispatches) -> AsyncIterator[tuple[StreamEvent, TaskResult|None]]`
  - 串行遍历 dispatches，对每个 task：从 AdapterRegistry 获取 adapter → create_session → stream_chat → 收集事件，yield runtime_executing/runtime_completed 控制事件 + 透传 agent 流式事件
  - 错误处理：单个 task 失败不阻塞后续，yield error event + failed TaskResult
  - 新建 `agentend/tests/test_execution.py` 单元测试

- [x] **T3: 重构 OrchestratorAdapter dispatch+collect（AgentEnd）**
  - 预估: 2h | Capability: runtime-execution-engine | 依赖: T2
  - 修改 `agentend/src/adapters/orchestrator.py` 的 `stream_chat()` 方法
  - 构造函数注入 AdapterRegistry 实例
  - Phase 2 (Dispatch) 保持不变，Phase 3 (Collect) 替换为 ExecutionEngine 调用
  - Phase 4 (Aggregate) 和 Phase 5 (Record) 保持不变
  - 可能需要修改 `agentend/src/api/v1/agent.py` 传入 registry

- [x] **T4: 新增 CoordinationChannel（AgentEnd）**
  - 预估: 3h | Capability: coordination-channel | 依赖: T1, T2
  - 新建 `agentend/src/orchestrator/coordination.py`，实现 CoordinationChannel 类
  - 核心方法 `coordinate(plan, agents) -> AsyncIterator[StreamEvent]`，实现单轮 Q&A：LLM 生成问题 → adapter.chat() 获取回答 → yield coordination 系列事件
  - 提供 `summary() -> str` 返回协调结论文本
  - 在 OrchestratorAdapter 中接入为 Phase 1.5（Planning 后、Execute 前）
  - 新建 `agentend/tests/test_coordination.py` 单元测试

- [x] **T5: Backend SSE 透传确认 + 事件类型同步**
  - 预估: 1h | Capability: backend-sse | 依赖: T1
  - 检查 `backend/internal/stream/writer.go` 的事件处理 switch-case，确保新事件类型走默认透传路径
  - 如有类型枚举校验限制，改为 string pass-through 或补充 default 分支
  - 确认 `backend/internal/handler/stream.go` 的 serveStreaming() 不过滤未知事件类型
  - `make generate` 自动更新 `backend/internal/generated/events.go`

- [x] **T6: Frontend 事件处理扩展**
  - 预估: 2h | Capability: frontend-runtime-timeline | 依赖: T1
  - 扩展 `frontend/src/hooks/use-chat-stream.ts`：新增 runtime_executing/completed、coordination_start/message/done、planning(node=dispatch) 事件处理
  - 扩展 `frontend/src/stores/chat.ts`：新增 runtimeStatus 状态和协调相关 blocks 追加逻辑
  - 扩展 `frontend/src/lib/block-reducer.ts`：新增 plan、coordination、runtime_status block 类型
  - 默认分支：未知事件类型静默忽略

- [x] **T7: Frontend 新组件实现**
  - 预估: 4h | Capability: frontend-runtime-timeline | 依赖: T6
  - 新建 `frontend/src/components/cards/PlanCard.tsx`：任务分配卡片，显示 overview + tasks 列表，agent 标签颜色区分，状态图标实时更新
  - 新建 `frontend/src/components/cards/ToolCard.tsx`：工具调用卡片，显示工具名 + 命令/文件路径 + 执行结果
  - 新建 `frontend/src/components/cards/CoordChannel.tsx`：可折叠协调通道，显示协调轮次、消息流、结论摘要
  - 新建 `frontend/src/components/cards/RuntimeStatus.tsx`：状态徽章，planning(黄)/executing(橙)/completed(绿)/failed(红)，带脉动动画
  - 修改 `frontend/src/components/chat/MessageBubble.tsx` 注册新 block 类型到对应组件

- [x] **T8: 端到端集成测试**
  - 预估: 2h | Capability: runtime-execution-engine, frontend-runtime-timeline | 依赖: T3, T4, T5, T7
  - 启动三端服务，创建 orchestrator 类型任务选择 2+ 个 agent
  - 验证：Planner 生成计划 + PlanCard 渲染、CoordinationChannel 运行 + CoordChannel 渲染、真实 Agent 调用 + 事件流、RuntimeStatus 徽章状态切换、Aggregate 汇总、streaming→completed 状态转换、错误场景不阻塞

---

## Task Dependency Graph

```
T1 (契约层)
├── T2 (ExecutionEngine)
│   ├── T3 (OrchestratorAdapter 重构)
│   │   └── T8 (集成测试)
│   └── T4 (CoordinationChannel)
│       └── T8
├── T5 (Backend SSE)
│   └── T8
└── T6 (Frontend 事件处理)
    └── T7 (Frontend 组件)
        └── T8
```

## Execution Order

```
T1 → T2 → T3 → T4 ─┐
T1 → T5 ─────────────┤
T1 → T6 → T7 ───────┤
                      ▼
                     T8
```

可并行：T2 + T5 + T6 在 T1 完成后可同时进行。
串行：T3 依赖 T2，T4 依赖 T2，T7 依赖 T6，T8 依赖全部。
