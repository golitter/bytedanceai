## Why

当前 Orchestrator 只能做「LLM 调用 + 写文件」，是一个无状态的一次性脚本。没有任务分发、没有执行监控、没有结果收集、没有聚合汇报——任务要求文档明确要求「子Agent完成后，Orchestrator聚合产出并在聊天流中汇报结果」，但当前连「子Agent完成」这个信号都收不到。需要让 Orchestrator 从写文件的脚本升级为能闭环的编排器。

## What Changes

- 新增 Dispatcher 模块：将 Planner 的任务输出转换为 `@agent` 调度指令（群聊 @ 模式），由前端/Go Backend 消费执行
- 新增 Aggregator 模块：LLM 汇总多个 Agent 的执行结果为人类可读的报告
- 新增 Pin Memory：复用现有 `common/` 目录，通过 `_pins.yaml` 书签层管理动态 pin/unpin，AI 自动生成摘要，所有 Agent 可读，注入 Planner prompt
- 新增 Self-Evolution：每次编排后记录成败经验到 `evolution.yaml`，最近 N 条经验注入 Planner prompt
- 新增 RuntimeState：内存中追踪任务状态（PENDING / RUNNING / COMPLETED / FAILED）
- 升级 OrchestratorAdapter：从「写文件就结束」变为「plan → dispatch → collect → aggregate」闭环
- 升级 Planner Prompt：支持注入 Pin 约束 + 历史经验
- 新增 Pin API 端点：`/v1/pin/add`、`/v1/pin/remove`、`/v1/pin/list`

## Capabilities

### New Capabilities
- `task-dispatcher`: 将 Planner 输出转换为 @agent 调度指令，产出结构化调度 JSON
- `result-aggregator`: LLM 汇总多 Agent 执行结果为报告
- `pin-memory`: 基于 common 目录 + _pins.yaml 书签的动态约束管理，AI 摘要 + 全文可查阅
- `self-evolution`: 编排经验积累，每次编排后记录成败到 evolution.yaml
- `runtime-state`: 内存中的任务状态追踪（PENDING / RUNNING / COMPLETED / FAILED）

### Modified Capabilities
- `orchestrator-planning`: Planner prompt 升级支持 Pin 约束 + 历史经验注入，OrchestratorAdapter 从单向写文件升级为闭环编排

## Impact

- **新增文件**：`dispatcher.py`、`aggregator.py`、`pin_memory.py`、`evolution.py`、`state.py`
- **修改文件**：`adapters/orchestrator.py`（闭环逻辑）、`orchestrator/prompts.py`（Pin/Evolution 注入）、`models.py`（新增 DispatchResult、TaskResult）
- **新增 API**：`/v1/pin/add`、`/v1/pin/remove`、`/v1/pin/list`
- **Shared Workspace**：`memory/common/_pins.yaml`（Pin 书签）、`evolution.yaml`（经验记录）
- **依赖**：已有 httpx / langchain-openai / pyyaml，无新增依赖
