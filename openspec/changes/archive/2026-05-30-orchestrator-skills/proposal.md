## Why

Orchestrator 当前是线性流程引擎：单次 LLM 调用生成 Plan → 写文件 → 分派。它无法自主探索上下文（如读取公共记忆、查看任务状态），也无法在规划过程中使用 skill 工具（taskctl、render）。与此同时，skill 知识以全量 SKILL.md 塞入 prompt 的方式不可扩展——随 skill 数量增长，token 消耗线性膨胀。

需要让 orchestrator 成为带工具的 Agent Loop，具备自主调用 read_file / write_file / run_skill 等工具的能力；同时引入 agentskills.io 规范的渐进式披露（L1 元数据 → L2 指令 → L3 资源），按需加载 skill 知识，控制 token 成本。

## What Changes

- LangGraph 图从 `plan → write_shared` 扩展为 `discover → select → load_l2 → plan(agent loop) → write_shared`，新增三个渐进式披露节点
- plan_node 从单次 LLM 调用改为 LangGraph tool-calling agent loop，LLM 可自主调用工具多轮迭代后输出 Plan JSON
- 新增 5 个精细工具（read_file、write_file、list_dir、run_skill、load_skill_resource），无 Bash 工具，安全白名单设计
- run_skill 只执行 `config.yaml` 的 `skills.manifest` 中注册的已知 skill 二进制，cwd 锁定 shared_dir
- OrchestratorGraphState 新增 l1_skills、selected_skill_names、l2_content、l3_content 字段
- 系统提示从硬编码全量 skill 内容改为动态注入 L1 目录 + 按需 L2/L3

## Capabilities

### New Capabilities
- `orchestrator-tool-calling`: Orchestrator Agent Loop 工具系统——定义 read_file、write_file、list_dir、run_skill、load_skill_resource 五个工具的接口、安全约束、执行语义
- `skill-progressive-disclosure`: Skill 渐进式披露加载机制——L1 元数据发现、LLM 语义选择、L2 指令按需加载、L3 资源工具调用

### Modified Capabilities
- `orchestrator-planning`: plan_node 从单次 LLM 调用改为 tool-calling agent loop；LangGraph 图新增 discover → select → load_l2 前置节点；GraphState 扩展渐进式披露字段

## Impact

- `agentend/src/orchestrator/planning/graph.py` — 图结构重构（新增节点 + 条件边）
- `agentend/src/orchestrator/planning/prompts.py` — prompt 改为接收 L2 指令动态注入
- `agentend/src/orchestrator/planning/tools.py` — 新增（5 个 @tool 定义）
- `agentend/src/orchestrator/planning/skill_loader.py` — 新增（L1/L2/L3 加载逻辑）
- `agentend/src/orchestrator/models.py` — GraphState 扩展字段
- `agentend/src/adapters/orchestrator.py` — 传递 shared_dir 给图节点
- 无前端改动，无后端改动，无契约层改动
