# Orchestrator 架构设计

## 概述

Orchestrator 是一个基于 LangGraph 的多 Agent 编排器。它接收用户需求，决定是直接回答还是协调多个 Agent 完成任务，通过 LLM tool-calling 循环实现规划和分发。

## 核心架构

```
用户消息
  │
  ▼
┌──────────────────────────────────────────────────────────┐
│                    LangGraph State Machine                │
│                                                          │
│  skill_prepare → reason → (human_review) → dispatch →   │
│  execute → review → evolve → save_mem                    │
│                                                          │
│  GraphState: message, plan, dispatch_results,            │
│             execution_waves, task_results, ...            │
└──────────────────────────────────────────────────────────┘
  │
  ▼
子 Agent 执行（并行/串行）
```

### 关键文件

| 文件 | 职责 |
|------|------|
| `src/orchestrator/planning/graph.py` | LangGraph 定义：节点、边、条件路由 |
| `src/orchestrator/planning/prompts.py` | 系统提示词模板 + 动态构建 |
| `src/orchestrator/planning/tools.py` | LLM 工具定义（read_file, plan_and_dispatch 等） |
| `src/orchestrator/execution/dispatcher.py` | 计划 → DispatchResult 转换 + 拓扑排序 |
| `src/orchestrator/execution/engine.py` | 子 Agent 并发执行引擎 |
| `src/orchestrator/memory/pin_memory.py` | Pin 约束记忆 |
| `src/orchestrator/memory/evolution.py` | 编排经验学习 |
| `src/adapters/orchestrator.py` | OrchestratorAdapter：graph 生命周期管理 |

---

## Graph 节点详解

### 1. skill_prepare — 技能发现与提示词构建

```
输入: user message, agents, shared_dir, task_base_path
输出: system_prompt (写入 state)
```

**三级技能发现**：
1. **L1 扫描**：扫描 `shared/.orchestrator/skills/` 目录下的技能清单
2. **L2 选择**：根据用户消息语义匹配，选出相关技能
3. **L3 加载**：加载选中技能的详细内容

**系统提示词构建**（`build_reason_prompt`）：

```
┌─────────────────────────────────────┐
│  基础身份                            │
│  "你是一个对话式任务编排器"           │
├─────────────────────────────────────┤
│  {agents_desc}    可用 Agent 列表    │
│  {soul_section}   SOUL.md 身份定义   │
│  {skills_section} L2 技能内容        │
│  {tools_section}  工具说明           │
│  {pin_context}    Pin 约束           │
│  {evolution_ctx}  近期编排经验       │
│  {replan_section} 重规划上下文       │
│  {orchestrator_ctx} 跨轮对话上下文   │
│  {workspace_section} 工作区目录说明  │
├─────────────────────────────────────┤
│  规则                                │
│  - 判断逻辑（闲聊 vs 编排）          │
│  - Agent 与 Skill 区别               │
│  - main 分支合并决策                 │
│  - 通用规则（任务数 ≤ 5 等）         │
├─────────────────────────────────────┤
│  {message}  用户需求                 │
└─────────────────────────────────────┘
```

### 2. reason — LLM 工具调用循环（核心节点）

```
输入: system_prompt, memory_messages, message
输出: output_type="text"|"plan", text|plan, memory_messages
```

**流程**：

```
构建 messages = [SystemMessage(system_prompt)] + memory_messages + [HumanMessage(message)]
     │
     ▼
┌─ 循环（最多 10 次）──────────────────────┐
│                                          │
│  llm_with_tools.ainvoke(messages)        │
│       │                                  │
│       ├── 无 tool_calls?                 │
│       │   └── 返回 text 输出             │
│       │                                  │
│       ├── 有 ask_agent?                  │
│       │   └── 执行 ask_agent（同步等待）  │
│       │       append ToolMessage         │
│       │       continue 循环              │
│       │                                  │
│       ├── 有 plan_and_dispatch?          │
│       │   └── 构造 PlanOutput            │
│       │       返回 plan 输出             │
│       │                                  │
│       └── 其他 tool_calls?               │
│           └── 执行工具（read_file 等）    │
│               append ToolMessage         │
│               continue 循环              │
│                                          │
└──────────────────────────────────────────┘
```

**LangSmith trace 接入点**：通过 `get_config()` 获取 LangGraph runnable config，传入 `ainvoke(messages, config=llm_config)`，使每次 LLM 调用自动上报 LangSmith。

### 3. human_review — 人工审查（可选）

```
输入: plan
输出: review_decision ("approve"|"discuss"|"reject"), review_message
```

- 阻塞等待外部审查结果（通过 `asyncio.Event`）
- 超时自动批准（默认 120s）
- 审查结果决定路由：
  - `approve` → dispatch
  - `discuss`/`reject` → 回到 reason（带反馈）

### 4. dispatch — 计划分发

```
输入: plan
输出: dispatch_results, execution_waves
```

- `Dispatcher` 将 `PlanOutput.tasks` 转换为 `DispatchResult` 列表
- 每个 task 映射到对应 agent（通过 session_id）
- `topological_sort` 按 `depends_on` 关系排序为执行波次：
  - 同一波次内的任务可并行执行
  - 波次之间串行执行
- 将计划写入 `shared/.agent/plans/` 目录

### 5. execute — 子 Agent 执行

```
输入: dispatch_results, execution_waves
输出: task_results
```

> 注意：graph 中的 `execute` 节点是占位符（`_execute_placeholder`），实际执行由 `OrchestratorAdapter.stream_chat` 中的 `_handle_execute` 完成。

**执行引擎**（`ExecutionEngine`）：
- 按波次依次执行
- 同波次内并发（`asyncio.create_task`）
- 每个子任务通过 `BackendClient` 调用对应 Agent 的 API
- 支持超时控制（默认 300s/任务）
- 自动为每个子 Agent 创建独立的 git worktree

### 6. review — 执行结果审查

```
输入: task_results, iteration, max_iterations
输出: needs_replan, replan_reason
```

- 分析每个子任务的成功/失败
- 如果有失败且未超过最大迭代次数 → 设置 `needs_replan=true`
- 构造重规划原因描述
- 路由：`needs_replan` → 回到 `skill_prepare`；否则 → `evolve`

### 7. evolve — 经验学习

```
输入: message, plan, task_results
输出: (无 state 变更，写文件)
```

- 将本轮编排经验记录到 `EvolutionStore`（`evolution.yaml`）
- 记录：用户需求摘要、规划摘要、执行结果、agent 表现
- 保留最近 20 条经验
- 下一轮 `skill_prepare` 时加载到提示词中

### 8. save_mem — 记忆保存

当前为空操作，预留扩展。

---

## 条件路由

```
reason ──→ route_by_output_type ──┬── "text"  → save_mem → END
                                  ├── "plan"  → human_review
                                  └── "error" → END

human_review ──→ route_by_review_decision ──┬── "approve" → dispatch
                                             └── "discuss"/"reject" → reason

execute ──→ review ──→ route_by_review ──┬── needs_replan → skill_prepare
                                          └── done → evolve → save_mem → END
```

---

## GraphState 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `message` | `str` | 用户输入 |
| `agents` | `list[dict]` | 可用 Agent 列表 |
| `task_id` | `str` | 当前任务 ID |
| `shared_dir` | `str` | 共享目录路径 |
| `allowed_read_dirs` | `list[str]` | 允许读取的目录 |
| `output_type` | `str` | 输出类型："text" / "plan" / "error" |
| `text` | `str` | 文本响应 |
| `plan` | `PlanOutput \| None` | 生成的计划 |
| `dispatch_results` | `list[DispatchResult]` | 分发结果 |
| `execution_waves` | `list[list[DispatchResult]]` | 执行波次 |
| `task_results` | `Annotated[list, _add]` | 任务执行结果（累加） |
| `iteration` | `Annotated[int, _add_one]` | 当前迭代次数 |
| `max_iterations` | `int` | 最大重规划次数（默认 3） |
| `memory_messages` | `Annotated[list, _add]` | 对话历史（累加） |
| `system_prompt` | `str` | skill_prepare 构建的系统提示词 |
| `review_decision` | `str` | 审查决定 |
| `review_message` | `str` | 审查反馈 |
| `needs_replan` | `bool` | 是否需要重规划 |
| `replan_reason` | `str` | 重规划原因 |

---

## LLM 工具清单

| 工具名 | 用途 | 参数 |
|--------|------|------|
| `read_file` | 读取文件（带行号） | `path`, `start_line`, `line_count` |
| `write_file` | 写入文件到共享目录 | `path`, `content` |
| `list_dir` | 列出目录内容 | `path` |
| `run_skill` | 执行已注册的 skill 命令 | `skill`, `command`, `args` |
| `load_resource` | 加载 skill 的参考资源 | `skill_name`, `resource_path` |
| `ask_agent` | 向指定 Agent 提问 | `agent`, `question` |
| `plan_and_dispatch` | 编排多 Agent 任务 | `overview`, `tasks`, `merge_to_main` |
| `current_time` | 获取当前时间 | (无) |

---

## 记忆系统

### Pin Memory（硬约束）

- 用户可以"钉住"某些内容作为必须遵守的约束
- 存储：`shared/memory/common/_pins.yaml` + 各 pin 文件
- LLM 自动生成 1-3 句摘要
- 注入到系统提示词的 `{pin_context}` 区块

### Evolution Store（经验学习）

- 每轮编排完成后记录经验（成功/失败、agent 表现）
- 存储：`shared/evolution.yaml`
- 保留最近 20 条
- 下一轮注入到系统提示词的 `{evolution_context}` 区块

---

## 分支模型

```
main                              ← 仓库主分支
  └── task/{task_id}              ← 任务基础分支（所有 agent 共享）
        ├── agent/{session_id_1}/{task_id}  ← Agent A 的独立分支
        └── agent/{session_id_2}/{task_id}  ← Agent B 的独立分支
```

- Orchestrator 只读 `task/{task_id}` 的 task-base worktree（用于了解代码结构）
- 子 Agent 各自创建独立的 `agent/{sid}/{tid}` 分支 worktree
- 子 Agent 完成后 merge 回 `task/{task_id}`
- Orchestrator 决定是否将 `task/{task_id}` 合入 `main`

---

## Adapter 层（生命周期管理）

`OrchestratorAdapter` 负责 graph 外围的生命周期：

1. **初始化**：创建共享目录、provision skills、写入 SOUL.md
2. **运行 graph**：`graph.astream(initial_state, config, stream_mode="updates")`
3. **处理节点输出**：将 graph updates 转换为 `StreamEvent` 发送给前端
4. **执行子 Agent**：`execute` 节点由 adapter 的 `_handle_execute` 接管
5. **重规划**：执行失败时递归调用 `stream_chat`，携带重规划上下文
6. **合并到 main**：用户批准后执行 `task/{id}` → `main` 的合并

---

## 重规划机制

当子 Agent 执行失败时：

```
execute → review (检测到失败)
  → needs_replan=true, iteration < max_iterations
  → 回到 skill_prepare（带 replan_reason）
  → 重新走 reason → dispatch → execute
  → 最多重规划 3 次
```

重规划消息会追加原始用户需求 + 失败原因，要求 Orchestrator 保留已完成工作的意图，只修复失败部分。

---

## 与其他 Adapter 的区别

| 维度 | Orchestrator | Claude Code / OpenCode / Codex |
|------|---|---|
| LLM 调用 | 内部 ChatOpenAI | 外部 CLI 子进程 |
| 工具系统 | LangChain tools | CLI 内置 |
| 工作区 | 只读 task-base worktree | 独立读写 worktree |
| 输出 | Graph 节点流转 | stdout JSON 流解析 |
| Trace | LangSmith 自动 trace | 需手动 RunTree（Phase 5.2） |
