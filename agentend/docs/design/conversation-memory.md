# Conversation Memory 设计

## 背景

Orchestrator 的 `memory_messages` 在 GraphState 中定义为 `Annotated[list, _add]`（列表拼接 reducer），但每次 `stream_chat` 调用都从 `[]` 开始，导致跨轮对话时 LLM 丢失自身的历史推理上下文（工具调用、规划决策等）。

当前的跨轮上下文 `orchestrator_context` 仅包含群聊窗口的 user/assistant 文本消息，不包含内部推理链。同时 pin/evolution/replan/orchestrator_ctx 4 个动态上下文块硬编码在系统提示词中，混在一起不利于 LLM 区分指令和上下文。

## 目标

1. 持久化 `memory_messages`，让 Orchestrator 记住跨轮推理历史
2. 将动态上下文从系统提示词拆分到消息列表，按语义使用不同消息类型
3. 动态上下文每轮重新注入，不持久化（避免重复）

## 架构变更

### 消息列表结构变更

**Before**：
```
[SystemMessage(身份+规则+pin+evolution+replan+orchestrator_ctx+工具+技能)]
+ memory_messages([])
+ [HumanMessage(用户消息)]
```

**After**：
```
[SystemMessage(身份+规则+工具+技能)]          ← 只保留不变部分
+ [SystemMessage(pin约束)]                    ← 全局约束
+ [SystemMessage(编排经验)]                   ← 全局约束
+ [HumanMessage(重规划反馈)]                  ← 仅重规划时有
+ [HumanMessage(群聊上下文)]                  ← 对话上下文
+ memory_messages(持久化的历史)               ← 含之前轮次的上下文+对话
+ [HumanMessage(用户消息)]                    ← 当前输入
+ [HumanMessage(审查反馈)]                    ← 仅审查时有
```

### 数据流

```
Turn N:
  stream_chat → load_messages() → initial_state["memory_messages"]
    → skill_prepare_node → 计算 pin/evolution → 保存到 state
    → reason_node → 注入 context_msgs + memory_messages + HumanMessage
    → LLM 调用 → 返回 [context_msgs, HumanMsg, AIMsg/ToolMsgs]
    → save_mem_node → save_messages() → 写入 JSON

Turn N+1:
  stream_chat → load_messages() → 拿到 Turn N 的完整链
    → skill_prepare_node → 计算最新 pin/evolution → state
    → reason_node → 注入最新 context + 历史记忆 + 新 HumanMessage
    → ... (LLM 看到 Turn N 的推理历史 + 最新约束)
```

重规划场景：`memory_messages` 通过 `_add` 累积所有 reason_node 返回值，`save_mem_node` 在最终完成时一次性持久化。

### 存储格式

- 文件路径：`{shared_dir}/memory/conversation_memory.json`
- 使用 LangChain 内置的 `messages_to_dict` / `messages_from_dict` 序列化
- JSON 格式，保留完整的 HumanMessage / AIMessage（含 tool_calls） / ToolMessage / SystemMessage 结构

### 保留策略

保留最近 **10 轮**对话（每轮以 HumanMessage 为起点），按轮次裁剪，不截断半轮。

### 消息类型语义

| 上下文 | 消息类型 | 语义 |
|--------|---------|------|
| Pin 约束 | `SystemMessage` | 全局约束，必须遵守 |
| 编排经验 | `SystemMessage` | 历史经验参考 |
| 重规划反馈 | `HumanMessage` | 对话上下文，上一轮失败原因 |
| 群聊上下文 | `HumanMessage` | 对话上下文，其他 Agent 消息 |

### 持久化范围

- **持久化**：对话链（HumanMessage + AIMessage + ToolMessage）— 存入 `conversation_memory.json`
- **不持久化**：4 个动态上下文消息 — 每轮从最新数据源重新构建注入，不写入 memory

这避免跨轮重复，保持职责清晰：memory 只存对话，上下文每次重新计算。

## 修改文件

| 文件 | 动作 | 说明 |
|------|------|------|
| `src/orchestrator/memory/conversation_memory.py` | 新建 | `ConversationMemoryStore` 类 |
| `src/orchestrator/planning/graph.py` | 修改 | GraphState 新增字段 + 节点逻辑调整 |
| `src/orchestrator/planning/prompts.py` | 修改 | 移除 4 个动态上下文插槽 |
| `src/adapters/orchestrator.py` | 修改 | 加载持久化记忆到 initial_state |
