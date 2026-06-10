## Context

Orchestrator 当前 LangGraph 图为 `plan → write_shared → END`。plan_node 做单次 LLM 调用生成 PlanOutput，无工具、无循环、无自主探索能力。Skill 知识未集成到 orchestrator 的规划流程中。

现有 skill 体系：`config.yaml` 的 `skills.manifest` 定义了 taskctl、render 两个 builtin skill，每个 skill 有 SKILL.md（含 frontmatter 元数据 + 正文指令）和可执行二进制。SkillProvisioner 将 skill 复制到子 Agent 的 workspace（`.claude/skills/`、`.opencode/skills/`），但 orchestrator 无 config_dir，所以被跳过。

LangGraph 已作为 orchestrator 的图引擎，原生支持 tool-calling agent 模式（`bind_tools` + 条件边 + tool node）。这为改造提供了现成的基础设施。

## Goals / Non-Goals

**Goals:**
- 让 orchestrator 在规划阶段成为 tool-calling agent loop，能自主调用工具探索上下文
- 按 L1→L2→L3 渐进式披露加载 skill 知识，控制 token 消耗
- 5 个精细工具（read_file、write_file、list_dir、run_skill、load_skill_resource）安全白名单设计
- 改动范围收束在 `orchestrator/planning/` 内，下游（engine、coordination、dispatcher、aggregator）不动

**Non-Goals:**
- 不改动 SkillProvisioner 的供给逻辑（子 Agent 的 skill 供给不变）
- 不给 orchestrator 添加 Bash 工具（不做黑名单拦截，结构性安全）
- 不改动前端/后端/契约层
- 不改变 write_shared_node 的确定性写入逻辑
- 不让 orchestrator 拥有自己的 workspace（它使用 shared_dir）

## Decisions

### D1: 图结构——线性披露节点 + agent loop plan_node

图从 2 节点扩展为 5 节点：`discover → select → load_l2 → plan → write_shared`。

前三个节点（discover、select、load_l2）是确定性的渐进式披露流程，不涉及 LLM 自主循环（除 select 做一次 LLM 语义匹配）。plan_node 是 tool-calling agent loop，LLM 在此阶段自主决定调哪些工具。

**为什么不用纯 agent loop 从头到尾？** 渐进式披露的 discover/select/load_l2 是确定性的准备工作，放在 agent loop 外可避免 LLM 浪费轮次在"发现 skill"上，也让 token 计量更精确。

**替代方案**：把整个流程做成一个大 agent loop，LLM 自己发现、选择、加载 skill。被否决——因为 discover 和 select 的逻辑是确定性的，不需要 LLM 的"创造力"，反而增加不确定性。

### D2: 工具集——5 个精细工具，无 Bash

| 工具 | 参数 | 安全约束 |
|------|------|----------|
| read_file | path: str | 无限制（只读） |
| write_file | path: str, content: str | path 必须在 shared_dir 内，防路径逃逸 |
| list_dir | path: str | 无限制（只读） |
| run_skill | skill: Literal[], command: str, args: str | skill 必须在 manifest 注册；cwd 锁定 shared_dir；timeout 30s |
| load_skill_resource | skill_name: str, resource_path: str | skill 必须在 manifest 注册；路径不允许 `..` |

**为什么不用 Bash + 黑名单？** orchestrator 是自动运行的 Agent，无人类审批环节。黑名单永远有绕过方式（`python -c "import shutil; shutil.rmtree('/')"`）。精细工具是白名单设计，结构性不可能执行危险操作。

### D3: L3 加载——作为 Tool 暴露给 LLM

L3 资源不由图节点提前加载，而是将 `load_skill_resource` 作为 tool 暴露在 plan_node 的 agent loop 中。LLM 在执行 L2 指令时，若看到"请使用 load_skill_resource 读取 references/xxx"，会主动调用。

这遵循 agentskills.io / ADK SkillToolset 规范，也符合渐进式披露的核心思想——L3 只在真正需要时才消耗 token。

### D4: select_node——LLM 语义匹配

select_node 用一次 LLM 调用做语义匹配：给定 L1 列表 + 用户任务，让 LLM 返回最相关的 1~N 个 skill 名称。不使用规则匹配或 embedding，因为 skill 数量少（当前 2 个，预期 <10 个），LLM 直接判断更准确。

### D5: plan_node agent loop 终止条件

plan_node 的 agent loop 有两个终止条件：
1. LLM 返回无 tool_calls 的响应（自然终止）
2. 达到最大迭代次数（max_iterations=10，防止死循环）

LLM 在终止时必须输出 PlanOutput JSON（与现有格式一致），保证下游 write_shared_node 无需改动。

### D6: 与现有 plan prompt 的关系

现有 PLAN_PROMPT 作为 plan_node agent loop 的系统提示基础，但在开头注入 L2 指令：

```
## 可用 Skills（L2 指令）
{selected skill 的 SKILL.md 正文}

## 可用工具
{工具列表和简要说明}

{现有 PLAN_PROMPT 内容，去掉已有的 agents_desc（改为工具动态获取）}
```

Pin Memory 和 Evolution Experience 仍通过 prompt 注入，与现有行为一致。

## Risks / Trade-offs

**[Token 超支风险]** agent loop 每轮都携带完整 context（L2 + 对话历史 + 工具结果），多轮迭代后 token 可能超过模型限制。→ 缓解：max_iterations=10；工具输出做截断（run_skill 输出限制 4096 字符）；L2 只加载选中 skill。

**[select_node 误选]** LLM 可能选中不相关的 skill，导致 L2 加载了无用的指令。→ 可接受：即使误选，影响只是多消耗几百 token L2，不会造成错误结果。随 skill 数量增长可优化 select prompt。

**[plan_node 不输出 JSON]** agent loop 中 LLM 可能在调完工具后输出自然语言而非 PlanOutput JSON。→ 缓解：系统提示明确要求"最终必须输出 JSON 格式的 Plan"；在 parse 失败时 fallback 尝试从自然语言中提取 JSON。

**[改动集中在 planning/ 内]** 好处是爆炸半径小，但 plan_node 从简单函数变为 agent loop，复杂度增加。→ 通过独立的 tools.py 和 skill_loader.py 模块化分离关注点。
