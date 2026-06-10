## Context

当前 Orchestrator 使用 LangGraph 实现了一个 5 节点线性管道（discover → select → load_l2 → plan → write_shared），存在以下问题：

1. **无状态**：每次请求独立，无法跨轮次引用上下文（用户说"基于上面的结果"，Orchestrator 不知道"上面"是什么）
2. **所有请求走规划**：闲聊也触发完整 skill discovery + LLM 调用流程，浪费 token
3. **LangGraph 用在错误的地方**：线性管道不需要状态机，真正需要状态编排的 Execution 阶段反而是手写 if/for
4. **OrchestratorAdapter 是上帝对象**：planning/coordination/dispatch/execute/aggregate/evolution 全部内联在一个 stream_chat 方法里

## Goals / Non-Goals

**Goals:**

- Orchestrator 成为有记忆的对话 Agent，自主判断闲聊回复 vs 任务编排
- LangGraph 覆盖完整生命周期（REASON → DISPATCH → EXECUTE → REVIEW → EVOLVE → SAVE_MEM）
- Memory 跨轮次持久化（对话历史 + 工具调用 + 结果）
- Execute 阶段按依赖分波次执行（Wave Executor）
- 任务失败时支持带上下文重规划

**Non-Goals:**

- Profile System / SOUL（留 Phase 6）
- AgentRegistry 扩展（留 Phase 6）
- MergeManager（留 Phase 6）
- 前端完整 Runtime Timeline 重构（仅最小适配）
- 持久化 Memory（跨进程/重启，当前用内存 MemorySaver）

## Decisions

### 1. REASON 节点合并 Pipeline 的 5 个步骤

**选择**：将 discover → select → load_l2 → plan → write_shared 合并为 REASON 节点的内部函数调用。

**替代方案**：保留 Pipeline 作为 REASON 的子图。
**理由**：Pipeline 是纯线性的，LangGraph 在这里没有附加价值。合并后代码更直观，减少 LangGraph 概念噪音。渐进式 skill 加载保留为函数内部步骤。

### 2. `plan_and_dispatch` 作为工具调用信号

**选择**：定义 `plan_and_dispatch(overview, tasks)` 工具。LLM 调用此工具 = 编排信号，LLM 输出纯文本（无 tool_calls）= 直接回复。

**替代方案 A**：让 LLM 自由输出文本或 JSON，后端解析区分。
**否决理由**：输出格式不稳定，解析脆弱，LLM 可能在该输出 JSON 时输出混合文本。

**替代方案 B**：用 `with_structured_output()` 强制 Union 输出类型。
**否决理由**：灵活性差，与 tool-calling 循环不兼容。

**理由**：tool_calls 语义天然区分"我要做事"和"我要说话"，LLM 对工具调用的理解成熟可靠。

### 3. Memory 使用 LangGraph MemorySaver

**选择**：使用 `langgraph.checkpoint.memory.MemorySaver` 管理跨轮次状态。

**替代方案**：自定义 Memory 存储（文件/Redis/DB）。
**理由**：当前阶段先用最简方案。MemorySaver 是内存级、零配置的，满足 MVP 需求。未来可替换为 SqliteSaver 或自定义 checkpointer。

### 4. Execute 作为 Wave Executor 子图

**选择**：将 Execute 拆为子图，内部按 waves 顺序执行，每波内并行。

**替代方案**：Execute 作为单节点，内部管理并行。
**理由**：子图结构让 LangGraph 能追踪每个 wave 的状态，未来扩展（单个任务状态追踪、中途取消）更自然。

### 5. REVIEW 节点决定重规划

**选择**：REVIEW 检查 task_results，如有失败且 iteration < max_iterations，带失败上下文重入 REASON。

**替代方案**：每次失败都重规划。
**理由**：限制重规划次数防止无限循环。失败上下文（哪些任务失败、为什么）让 REASON 能针对性调整。

## Risks / Trade-offs

- **[LLM 判断准确性]** → LLM 可能在该编排时输出文本，或反之。缓解：prompt 中明确区分规则 + few-shot 示例。
- **[Memory 增长]** → 长对话导致 messages 列表膨胀，超过 context window。缓解：后续可加 summarization 或滑动窗口，当前 MVP 先不处理。
- **[MemorySaver 非持久化]** → 进程重启丢失记忆。缓解：明确为 MVP 选择，后续替换 checkpointer。
- **[重构范围大]** → graph.py 和 orchestrator.py 都需要重写。缓解：对外接口（stream_chat/chat）不变，重构内部实现。
