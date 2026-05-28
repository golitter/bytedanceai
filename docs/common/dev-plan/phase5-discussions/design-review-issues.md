# 设计文档审查报告

> 审查范围: discussions/ 下 4 篇核心文档 vs 实际代码库
> 审查日期: 2026-05-28
> 涉及文档: conversation-layer-design.md, agentend-group-chat-impl.md, current-status-and-next-steps.md, orchestrator-features.md

---

## P0 — 阻塞问题（实施前必须解决）

### P0-1: 事件类型命名不一致（dot vs underscore）

**问题描述：** 4 篇文档对同一组事件类型的命名方式不统一。

| 事件 | current-status-and-next-steps.md | agentend-group-chat-impl.md |
|------|----------------------------------|------------------------------|
| 分支创建 | `workspace.branch.created`（dot） | `workspace.branch_created`（underscore） |
| 合并开始 | `workspace.merge.started`（dot） | `workspace.merge_started`（underscore） |
| 合并完成 | `workspace.merge.completed`（dot） | `workspace.merge_completed`（underscore） |
| 合并冲突 | `workspace.merge.conflict`（dot） | `workspace.merge_conflict`（underscore） |

`conversation-layer-design.md` 和 `orchestrator-features.md` 也混用了两种风格。此外 `agent.spawned` vs `runtime.agent.spawned` 在不同文档中也有分歧。

**代码证据：** 现有 `contracts/schemas/event-types.yaml` 定义了 8 种事件，全部使用无层次的单层命名（`init`, `text`, `tool_call` 等），新事件的层次命名风格尚未在契约层定义。`backend/internal/generated/events.go` 定义了对应的 `EventType` 常量。

**影响：** 前端 SSE 事件监听、Go Backend StreamWriter 的 switch case、AgentEnd 的事件发射——如果命名不统一，三端无法对接。

**建议：** 统一为 dot 风格（`workspace.branch.created`），与现有的 `EventType.TEXT`、`EventType.DONE` 等保持一致的层次语义。在 `contracts/schemas/event-types.yaml` 中定义所有新事件类型作为 Single Source of Truth。

---

### P0-2: Go Backend 变更范围自相矛盾

**问题描述：**

- `current-status-and-next-steps.md`（Phase 5 进度表）明确写着：`Go Backend ✅ 无需改动`
- `agentend-group-chat-impl.md`（第六节）列出了 4 个 Go Backend 变更：
  - `internal/generated/events.go` — 新增 20 个 RuntimeEvent 常量
  - `internal/generated/request.go` — 新增 5 个群聊字段
  - `internal/handler/task.go` — RunTask 构建 orchestrator config.agents
  - `internal/stream/writer.go` — switch case 新增事件类型处理

**实际代码验证：**

当前 `backend/internal/handler/task.go` 的 `RunTask` 确实没有构建 `config.agents` 的逻辑——构建 `AgentRequest` 时只设置了 `TaskId`、`SessionId`、`Message`、`AgentType`、`Stream` 五个字段，`Config` 字段始终为 nil。当前 `backend/internal/stream/writer.go` 的 switch case 只处理 `Text`、`Done`、`Error` 三种事件。`backend/internal/generated/events.go` 仅定义了 8 个 `EventType` 常量（init, text, tool_call, tool_result, artifact, planning, done, error），无任何 runtime.* 或 coordination.* 类型。

**影响：** 如果认为"Go 无需改动"，那么 AgentEnd 发出的新事件类型（runtime.*、coordination.*）虽然会被 `publishToRedis(line)` 透传，但 StreamWriter 的状态检测（判断是否 done/error）无法识别新事件类型。同时 `RunTask` 不会为 Orchestrator 构建 `config.agents`，前端需要自行传参。

**建议：**
- 修改 `current-status-and-next-steps.md`，将 Go Backend 状态从 "✅ 无需改动" 改为 "🔧 需少量改动（SSE 透传 + config 构建）"
- 明确 Go Backend 在 Phase 5 的改动范围，以 `agentend-group-chat-impl.md` 第六节为准

---

### P0-3: 缺少 contracts/schemas/ 更新计划

**问题描述：** 所有文档都遵循 AGENTS.md 中的"契约优先原则"，但 4 篇设计文档没有一篇提到更新 `contracts/schemas/` 的 YAML 文件。

当前契约层定义了 6 个 YAML schema：
- `event-types.yaml` — 8 种事件（EventType 枚举 + StreamEvent 结构体）→ 需扩展到 20+ 种
- `agent-request.yaml` — 10 个字段（task_id, session_id, message, agent_type, stream, system_prompt, rules, workspace_path, repo_path, config）→ 需新增 conversation_id、target_profile 等 5 个字段
- `agent-response.yaml`、`session-state.yaml`、`message.yaml`、`validate-repo-path.yaml` — 无需变更

**影响：** 没有契约更新就没有三端类型生成（`make generate`），前端和后端无法获得新字段和新事件的类型定义。

**建议：** 在实施计划的最前面（Step 0）增加契约更新步骤：
1. 更新 `contracts/schemas/event-types.yaml` — 新增 RuntimeEvent 类型枚举
2. 更新 `contracts/schemas/agent-request.yaml` — 新增 conversation_id 等字段
3. 运行 `make generate`
4. 在 `contracts/logs/` 写变更记录

---

### P0-4: Conversation Layer 持久化方案缺失

**问题描述：** `conversation-layer-design.md` 设计了 `ConversationStore`（内存），`orchestrator-features.md` 同样标注为内存存储。但没有讨论：

- Conversation 数据如何持久化？AgentEnd 重启后对话图丢失怎么办？
- ExecutionRecord 如何持久化？当前 Backend 的 Message 表是 flat 结构，没有 execution 概念
- Message 与 Execution 的多对一关系如何在 MySQL 中建模？

**当前代码验证：**

Backend 的 `Message` 模型（`backend/internal/model/message.go`）：
```go
type Message struct {
    ID        uint      `gorm:"primarykey"`
    MessageID string    `gorm:"uniqueIndex;size:36"`
    TaskID    string    `gorm:"index;size:36"`
    SessionID string    `gorm:"size:128"`
    Role      string    `gorm:"size:16"`              // "user" or "agent"
    Content   string    `gorm:"type:longtext"`
    Status    string    `gorm:"size:16;default:completed"` // streaming/completed/failed
    LastSeq   string    `gorm:"size:64;default:''"`
    AgentType string    `gorm:"size:64,omitempty"`
    AgentName string    `gorm:"size:128,omitempty"`
    CreatedAt time.Time
}
```
SessionAgent 模型（`backend/internal/model/session_agent.go`）有 `SessionID`、`AgentType`、`AgentName`、`AvatarURL` 字段。

没有 `conversation_id`、`execution_id`、`reply_to_message_id`、`quote_message_ids` 等字段。

**影响：** 如果不规划持久化，Regenerate（Message ≠ Execution）根本无法工作——MySQL 里没有 Execution 概念。

**建议：**
- 方案 A：Conversation 持久化放 Backend（新增 Conversation + Execution 表），AgentEnd 通过 API 读写
- 方案 B：AgentEnd 自建 SQLite/JSON 持久化，Conversation 是纯 AgentEnd 内部概念
- 无论哪个方案，都需要在设计文档中明确

---

### P0-5: RuntimeEvent 与现有 StreamEvent 的关系未定义

**问题描述：** 现有 AgentEnd 使用 `StreamEvent`（8 种 EventType：init, text, tool_call, tool_result, artifact, planning, done, error），新设计引入 `RuntimeEvent`（20+ 种类型），但两套事件系统的关系不清：

- `RuntimeEvent` 是替代 `StreamEvent` 还是包装 `StreamEvent`？
- `ExecutionEngine.normalize()` 将 `StreamEvent` 转换为 `RuntimeEvent`，那 Orchestrator 自身发出的事件（runtime.started、coordination.message）用什么类型？
- 前端当前只识别 `StreamEvent` 的 8 种类型，`RuntimeEvent` 的 20+ 种类型需要前端全部新增处理逻辑

**当前代码验证：**

`agentend/src/generated/events.py` 定义了 EventType 枚举（8 种），所有 Adapter 的 `stream_chat()` 返回 `AsyncIterator[StreamEvent]`。前端 `lib/block-reducer.ts` 只处理这 8 种事件。

**建议：**
- `RuntimeEvent` 是 `StreamEvent` 的超集，Phase 5 保留 `StreamEvent` 用于单 Agent 直接调用
- Orchestrator 场景下，`ExecutionEngine.normalize()` 将 `StreamEvent` → `RuntimeEvent`
- Orchestrator 自身事件直接构造 `RuntimeEvent`
- 前端通过 `event.type` 的前缀（`runtime.`/`task.`/`agent.`/`coordination.`/`workspace.`）区分新旧事件

---

## P1 — 重要问题（影响实施质量）

### P1-1: conversation-layer-design.md 章节编号重复

**问题描述：** 文档有两个"五"章节：
- 第五章：Message ≠ Execution
- 第五章（标注为"五.五"）：Agent 间上下文共享：Skill-based

**建议：** 将"五.五"改为"六"，后续章节顺延。

---

### P1-2: Workspace 路径设计不一致

**问题描述：**

| 文档 | Workspace 路径 | 代码实现 |
|------|---------------|----------|
| current-status-and-next-steps.md | `workspaces/{task_id}/{session_agent_id}/` | 不存在 |
| agentend-group-chat-impl.md | `worktrees/{task_id}/{session_id}/` | 当前代码 |
| conversation-layer-design.md | 未讨论路径 | — |

实际代码（`agentend/src/workspace/models.py`）使用 `worktrees/` 前缀：
```
{repo_parent}/worktrees/{task_id}/{session_id}/
```
分支名模板为 `agent/{session_id}/{task_id}`，任务分支为 `task/{task_id}`。

新设计将其改为 `workspaces/` 并引入子目录结构（`orchestrator/`、`shared/`、`{session_agent_id}/`）。

**影响：** 路径变更影响 workspace 创建、文件代理、前端预览、Git worktree 管理等全链路。

**建议：** 统一为 `workspaces/`（新命名），明确迁移路径。在实施时先保留 `worktrees/` 兼容，通过 WorkspaceManager 抽象路径生成。

---

### P1-3: Profile System 存储与 Backend 数据冲突

**问题描述：** 设计文档提出 Profile 基于 `soul.yaml` 静态文件（`agentend/src/profiles/`），但 Backend 已有动态 Agent 管理能力：

- `backend/internal/model/session_agent.go` — SessionAgent 表存储 agent_type、agent_name、avatar_url
- `backend/internal/handler/agent_profile.go` — Agent Profile API 支持动态创建/更新 agent profile
- 前端 `AgentProfilePage.tsx` 和 `AgentEditDialog.tsx` — 用户可自定义 Agent 名称、头像

**影响：** 如果 Profile 是静态文件，那前端已有的 Agent 自定义功能（改名、换头像）就无效了。如果 Profile 是动态的，那 `soul.yaml` 的身份/权限/约束如何与 Backend 的数据库记录同步？

**建议：**
- `soul.yaml` 定义 Profile 模板（身份、职责、权限、约束）— 静态，AgentEnd 内部
- Backend 数据库存储 Profile 实例（名称、头像、adapter 映射）— 动态，用户可编辑
- 启动时从 `soul.yaml` 加载 Profile 定义，通过 API 同步给 Backend

---

### P1-4: Coordinator 多轮规划的 Token 成本和超时未讨论

**问题描述：** 协调通道（Coordination Channel）支持最多 10 轮规划，每轮涉及：
- 1 次 Planner LLM 调用（LangGraph）
- N 次 Agent CLI 调用（通过 `adapter.chat()` 非流式）

最坏情况：10 轮 × (1 次 Planner + N 次 Agent 咨询) = 大量 LLM 调用。

**影响：**
- Token 成本可能非常高（10 轮 × 每轮 3 个 Agent 咨询 = 30+ 次 LLM 调用）
- 总耗时可轻松超过 Backend 的 30 分钟 streaming 超时（`backend/internal/stream/writer.go` 的 `StreamTimeout`）
- 用户等待时间过长，体验差

**建议：**
- 将默认最大轮次从 10 降为 3-5
- 增加总规划时间上限（如 5 分钟）
- 增加 token 预算控制
- 超时后强制生成最终 ExecutionPlan

---

### P1-5: adapter.chat() 可用性未验证

**问题描述：** Coordinator 在协调阶段使用 `adapter.chat()`（非流式，获取完整回答），但：

- `BaseAgentAdapter` 的 `chat()` 方法在大多数适配器中是简单封装，可能没有独立的 session 管理
- Claude Code CLI 的 `chat()` 调用是否支持快速问答（不创建/恢复 session）？
- 协调阶段的 chat 调用是否会与正式执行阶段的 session 冲突？

**当前代码验证：** `agentend/src/adapters/base.py` 定义了 `chat()` 方法，但 Claude 适配器实际运行 CLI 子进程，每次调用会启动新的 CLI 进程。

**建议：**
- 验证 Claude CLI 是否支持无 session 模式的快速问答
- 如果不支持，协调阶段使用独立的临时 session（用完即弃）
- 考虑协调阶段使用轻量级 LLM 直接调用（不经过 CLI），降低延迟

---

### P1-6: 前端 Runtime Timeline 实施细节缺失

**问题描述：** 4 篇文档对前端改动都只有高层描述（"渲染 Runtime Coordination Timeline"），缺少具体实施方案：

- `stores/chat.ts` 如何处理 20+ 种新事件类型？
- `MessageBubble.tsx` 如何区分渲染单 Agent 消息 vs 多 Agent 编排 Timeline？
- Coordination Channel 的折叠面板 UI 如何实现？
- Agent Presence（Profile 名称 + 颜色标签）的组件结构？

**当前前端代码验证：** `frontend/src/lib/block-reducer.ts` 处理当前 8 种事件，`frontend/src/stores/chat.ts` 管理 streaming 状态。新增 20+ 种事件需要大幅重构这两个核心模块。

**建议：** 在实施前补充前端设计文档，至少包含：
1. 新事件类型的处理策略（扩展 block-reducer 还是新写 runtime-event-handler）
2. Runtime Timeline 组件树
3. Coordination Panel 组件规格
4. 消息模型扩展（如何标识 multi-agent message）

---

### P1-7: AgentRequest 新增字段缺少默认值和验证规则

**问题描述：** `conversation-layer-design.md` 为 AgentRequest 新增了 5 个字段，但没有定义默认值和验证规则：

- `conversation_id: str` — 标记为必填，但现有的单 Agent 聊天不传此字段怎么办？
- `target_profile: str | None` — None 时走 Orchestrator，那与现有 `agent_type: AgentType` 的关系是什么？
- `quote_message_ids: list[str] = []` — 空列表 vs None 的语义区别？

**建议：** 所有新增字段应为 `Optional` + `None` 默认值，确保向后兼容。`conversation_id` 在单 Agent 模式下为 None（回退到当前行为）。

---

## P2 — 改进建议（可后续处理）

### P2-1: 群聊 Agent 来源三种模式的优先级未排序

`orchestrator-features.md`（第五节）定义了三种 Agent 来源模式（建群指定 / Orchestrator 自动选 / 混合动态拉入），但没有明确 Phase 5 先实现哪种。

**建议：** Phase 5 只实现模式 A（建群时指定 Agent 列表），这是最简单的且已有 `config.agents` 支持。模式 B 和 C 留 Phase 6。

---

### P2-2: RoutingPolicy._is_simple() 实现策略未定

`conversation-layer-design.md` 的 `RoutingPolicy._is_simple()` 标注为"简单启发式判断，或用 LLM 快速分类"，但没有确定用哪种方式。

**建议：** Phase 5 用硬编码规则（消息长度 < 100 字 + 包含 @mention = 简单），Phase 6 考虑 LLM 分类。

---

### P2-3: 文档中的文件路径与实际代码不完全对应

- `conversation-layer-design.md` 提到 `src/schemas/request.py`，实际应为 `src/api/v1/agent.py` 或 `src/schemas/` 下的对应文件
- `agentend-group-chat-impl.md` 提到 `src/api/dependencies.py`，实际依赖注入在 `src/app/dependencies.py`

**建议：** 实施前根据实际代码更新文档中的文件路径引用。

---

### P2-4: 测试策略缺失

4 篇文档都没有包含测试计划。当前代码库也没有测试文件。

**建议：** 至少为以下模块编写单元测试：
- `ExecutionEngine.normalize()` — 事件转换映射
- `Coordinator.plan_with_coordination()` — 多轮协调循环
- `Scheduler.run()` — 串行执行编排
- `ConversationService` — 上下文加载和 @mention 解析

---

### P2-5: Regenerate 设计与现有 Message 流不兼容

**问题描述：** `conversation-layer-design.md` 设计 Regenerate 为"基于同一触发消息创建新 Execution"（Message ≠ Execution）。但当前 Backend 的消息流是：

1. 用户发消息 → `Message(role=user)`
2. Agent 回复 → `Message(role=agent, status=streaming)`
3. 流式更新 → `Message.status=completed`

Regenerate 需要：
1. 找到原始 Agent Message
2. 创建新 Agent Message（新的 execution）
3. 前端需要知道显示哪次 execution

**影响：** 需要在 Backend Message 表新增字段或在 AgentEnd 新增 Execution 表。

**建议：** Phase 5 暂不实现 Regenerate，先跑通核心路径（Conversation → Routing → Execution）。Regenerate 留 Phase 6。

---

### P2-6: Orchestrator Agent 来源模式与 Planner Prompt 的关系未细化

`orchestrator-features.md` 第五节提到 Planner Prompt 应展示 "Agent 的名字"而非"角色描述"，但 `config.agents[].profile` 字段和 `soul.yaml` 的 `identity.role` 之间的关系需要更精确的定义。

**建议：** 明确 Planner Prompt 中 Agent 展示的优先级链：`config.agents[].name` > `soul.yaml.identity.role` > `profile_id`。

---

## 总结

| 等级 | 数量 | 概述 |
|------|------|------|
| P0（阻塞） | 5 | 事件命名不一致、Go 变更范围矛盾、缺少契约更新、持久化缺失、新旧事件关系未定 |
| P1（重要） | 7 | 章节编号重复、路径不一致、Profile 冲突、成本/超时、chat() 可用性、前端细节缺失、字段验证 |
| P2（建议） | 6 | 模式优先级、路由策略、文件路径、测试、Regenerate、Prompt 细化 |

**核心结论：** 设计文档的架构思路清晰（三层分离、Profile ≠ Adapter ≠ RuntimeAgent、Skill-based 上下文共享），但在与代码库对齐、跨文档一致性、以及实施细节（契约更新、持久化、前端改动）方面需要补充。建议先解决 P0 的 5 个问题再开始编码。
