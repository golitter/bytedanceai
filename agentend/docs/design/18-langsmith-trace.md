# LangSmith Trace — LLM 调用可观测性

## 实现了什么

为 Orchestrator 和 CLI Adapter 接入 LangSmith Trace，自动记录所有 LLM 调用、工具执行、Graph 节点转换的完整输入/输出，包括 system prompt、message history、tool calls、token 用量。

两层 Trace 实现：
- **Phase 5.1**：Orchestrator 的 LangGraph 自动 trace — 通过环境变量开启，`reason_node` 中传播 `get_config()` 到 LLM 调用
- **Phase 5.2**：CLI Adapter 的 `StreamEvent` 生命周期 trace — 通过 `trace_stream_events()` async generator wrapper 包裹 `adapter.stream_chat()`

## 怎么实现的

### Phase 5.1：Orchestrator 自动 Trace

LangGraph 内置 LangSmith 支持，开启环境变量后自动上报。唯一代码改动是 `reason_node` 中传播 runnable config：

**`src/orchestrator/planning/graph.py`** — `reason_node`：

```python
from langchain_core.runnables import get_config

try:
    llm_config = get_config()
except RuntimeError:
    llm_config = None
response = await llm_with_tools.ainvoke(messages, config=llm_config)
```

原因：graph 节点内部手动调用 LLM 时，会脱离 LangGraph 当前 runnable config，导致 LLM 子调用没有正确挂到 graph trace 树下。显式传播 `config` 后，graph 级 + LLM 级 + tool 级全部完整挂到同一个 trace 下。

LangSmith 展示的 Trace 树：

```
Run: Orchestrator session_id=yyy
├── skill_prepare (chain)        — 0.8s
├── reason (chain)               — 12.3s
│   ├── ChatOpenAI #1 (llm)      — input/output 完整
│   ├── read_file (tool)
│   ├── ChatOpenAI #2 (llm)
│   └── plan_and_dispatch (tool)
├── dispatch (chain)
├── execute (chain)
├── evolve (chain)
└── save_mem (chain)
```

### Phase 5.2：CLI Adapter Trace

**`src/adapters/trace.py`** — `trace_stream_events()` async generator wrapper：

在 `_execute_stream` 中用 wrapper 包裹 `adapter.stream_chat()`，不改任何 Adapter 代码。

事件 → RunTree 映射：

```
StreamEvent lifecycle          →  RunTree 操作
───────────────────────────────────────────────────────
流开始                          →  创建 root RunTree (chain)
TEXT chunk × N                  →  聚合到 text_parts[]（不创建 child）
TOOL_CALL                       →  ① flush 聚合文本 → text child run
                                   ② 创建 tool:{name} child run (pending)
TOOL_RESULT                     →  LIFO 匹配 pending tool run → end + patch
DONE                            →  flush 剩余文本
ERROR                           →  flush + 标记 root.error
流结束 (finally)                →  清理残留 pending → root.end() + root.patch()
```

关键设计点：

| 约束 | 实现 |
|------|------|
| 不承诺完整 LLM prompt | 只 trace StreamEvent，不碰 CLI 内部 prompt |
| text 做聚合 | `text_parts: list[str]` 累积，按 TOOL_CALL / DONE / ERROR 边界 flush |
| Orchestrator 排除 | Orchestrator 走 Phase 5.1 的 `get_config()` 自动 trace |

### 环境变量控制

统一由 `LANGSMITH_API_KEY` 控制：

```env
LANGSMITH_TRACING=true
LANGSMITH_API_KEY=lsv2_sk_xxxxx
LANGSMITH_PROJECT=agenthub
```

删掉或清空 `LANGSMITH_API_KEY` 即可关闭追踪。

### 关键文件

| 文件 | 操作 |
|------|------|
| `src/orchestrator/planning/graph.py` | `reason_node` 中 `get_config()` 传播到 LLM |
| `src/adapters/trace.py` | `trace_stream_events()` async generator wrapper |
| `src/api/v1/agent.py` | `_execute_stream` 用 wrapper 包裹 `adapter.stream_chat()` |
