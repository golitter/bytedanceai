## 1. Skill Loader 基础设施

- [x] 1.1 创建 `agentend/src/orchestrator/planning/skill_loader.py`，实现 `discover_skills(builtin_dir)` 函数：扫描 builtin_dir 子目录，解析每个 SKILL.md 的 YAML frontmatter，返回 `[{name, description}]` 列表，跳过无 SKILL.md 或无 name 字段的目录
- [x] 1.2 在 skill_loader.py 中实现 `load_skill_l2(skill_name, builtin_dir)` 函数：读取指定 skill 的 SKILL.md 完整内容，去掉 frontmatter 后返回正文
- [x] 1.3 在 skill_loader.py 中实现 `load_skill_resource(skill_name, resource_path, builtin_dir)` 函数：读取 skill 目录下 references/ 或 assets/ 中的文件，路径不允许 `..`

## 2. Tool 定义

- [x] 2.1 创建 `agentend/src/orchestrator/planning/tools.py`，实现 `read_file(path: str) -> str` 的 `@tool` 定义，处理文件不存在返回错误字符串
- [x] 2.2 实现 `write_file(path: str, content: str, shared_dir: str) -> str` 的 `@tool` 定义，路径必须 resolve 到 shared_dir 内，否则拒绝写入；自动创建父目录
- [x] 2.3 实现 `list_dir(path: str) -> str` 的 `@tool` 定义，返回目录内容列表，子目录后缀 `/`
- [x] 2.4 实现 `run_skill(skill: str, command: str, args: str, shared_dir: str) -> str` 的 `@tool` 定义，skill 必须在 manifest 注册表中，执行对应二进制，cwd 锁定 shared_dir，timeout 30s，输出截断 4096 字符
- [x] 2.5 实现 `load_skill_resource(skill_name: str, resource_path: str) -> str` 的 `@tool` 定义，skill 必须在 manifest 中，路径不允许 `..`
- [x] 2.6 实现 `build_tools(shared_dir: str) -> list[BaseTool]` 工厂函数，动态从 config.yaml 的 skills.manifest 构建 run_skill 的 Literal 参数，返回工具列表

## 3. GraphState 扩展

- [x] 3.1 在 `agentend/src/orchestrator/planning/graph.py` 中扩展 `GraphState` TypedDict，新增字段：`l1_skills: list[dict]`、`selected_skill_names: list[str]`、`l2_content: dict[str, str]`、`l3_content: dict[str, str]`

## 4. 渐进式披露节点

- [x] 4.1 实现 `discover_node(state)` — 调用 `discover_skills()` 扫描 builtin_dir，填充 `state["l1_skills"]`
- [x] 4.2 实现 `select_node(state)` — 用一次 LLM 调用做语义匹配：给定 L1 列表 + user message，返回相关 skill 名称列表，过滤无效名称
- [x] 4.3 实现 `load_l2_node(state)` — 对 `selected_skill_names` 调用 `load_skill_l2()`，填充 `state["l2_content"]`

## 5. plan_node Agent Loop 改造

- [x] 5.1 改造 `plan_node(state)` — 构建 system prompt，动态注入 L2 指令（`l2_content`）到 "## 可用 Skills" 段，注入工具说明，保留 Pin Memory 和 Evolution 上下文注入
- [x] 5.2 在 plan_node 中实现 tool-calling agent loop：`llm.bind_tools(tools)` → 循环 invoke → 有 tool_calls 则执行 tools → 结果追加为 ToolMessage → 再次 invoke → 直到无 tool_calls 或达到 max_iterations=10
- [x] 5.3 实现 plan_node 的终止逻辑：LLM 无 tool_calls 时解析 PlanOutput JSON（含 markdown code block 提取），失败返回 `{"plan": None}`
- [x] 5.4 更新 `build_planner_prompt()` 函数签名，接收 `l2_content` 参数，在 prompt 中新增 "## 可用 Skills" 和 "## 可用工具" 段

## 6. LangGraph 图重构

- [x] 6.1 重构 `build_graph()` 函数：新增 discover、select、load_l2 节点，设置线性边 `START → discover → select → load_l2 → plan → write_shared → END`
- [x] 6.2 移除 `build_graph()` 中 plan → write_shared 的旧直连边，替换为新的 5 节点线性图

## 7. Adapter 适配

- [x] 7.1 在 `agentend/src/adapters/orchestrator.py` 的 `stream_chat` 中，确保 `shared_dir` 传入 GraphState 供工具使用（当前已传入，验证即可）
- [x] 7.2 验证 OrchestratorAdapter 的 PLANNING 事件流兼容新图结构（discover/select/load_l2 阶段不需要额外事件）

## 8. 测试验证

- [x] 8.1 为 skill_loader.py 编写单元测试：discover_skills、load_skill_l2、load_skill_resource 三个函数的正常/异常场景
- [x] 8.2 为 tools.py 编写单元测试：write_file 路径逃逸拒绝、run_skill 未注册 skill 拒绝、load_skill_resource 路径遍历拒绝
- [x] 8.3 手动集成测试：发送 orchestrator 请求，验证 discover → select → load_l2 → plan(agent loop with tools) → write_shared 全流程
