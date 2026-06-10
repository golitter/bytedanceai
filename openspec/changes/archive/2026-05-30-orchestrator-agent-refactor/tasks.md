## 1. plan_and_dispatch 工具定义

- [x] 1.1 在 `orchestrator/planning/tools.py` 中新增 `plan_and_dispatch` 工具函数（参数: overview + tasks），工具体只返回 "plan_generated" 字符串
- [x] 1.2 确认 `plan_and_dispatch` 已加入 `build_tools()` 返回的工具列表

## 2. REASON 节点重构

- [x] 2.1 将 `discover_skills`、`select_skills`、`load_l2_content` 保留为独立函数（从现有 discover_node/select_node/load_l2_node 提取）
- [x] 2.2 重写 `plan_node` 为 `reason_node`：合并 skill 加载为内部步骤，实现 tool-calling 循环，区分 plan_and_dispatch 调用 vs 纯文本输出
- [x] 2.3 更新 `prompts.py`：PLAN_PROMPT 改为支持双模式（闲聊直接回复 + 编排调用 plan_and_dispatch），移除"必须只输出 JSON"约束

## 3. GraphState 更新

- [x] 3.1 新增 GraphState 字段：output_type、text、dispatch_results、execution_waves、task_results（reducer: add）、task_status、needs_replan、replan_reason、summary、iteration（reducer: add_one）、max_iterations、memory_messages（reducer: add）
- [x] 3.2 定义 `route_by_output_type` 条件路由函数（text → save_mem, plan → dispatch, error → END）
- [x] 3.3 定义 `route_by_review` 条件路由函数（ok → evolve, replan → reason）

## 4. 生命周期节点

- [x] 4.1 实现 `dispatch_node`：调用 Dispatcher.dispatch() + topological_sort 生成 execution_waves
- [x] 4.2 实现 `review_node`：检查 task_results 失败情况，设置 needs_replan + replan_reason
- [x] 4.3 实现 `evolve_node`：调用 EvolutionStore.record() 记录编排经验
- [x] 4.4 实现 `save_mem_node`：将本轮交互追加到 memory_messages

## 5. Wave Executor 子图

- [x] 5.1 在 `orchestrator/execution/wave.py` 中新建 `build_execute_subgraph()`，定义 ExecuteState 和 wave_execute_node
- [x] 5.2 实现 wave 循环逻辑：按 execution_waves 顺序执行，每波内用 asyncio.gather 并行执行
- [x] 5.3 将 wave executor 子图集成到主图的 execute 节点

## 6. Dispatcher 扩展

- [x] 6.1 在 `orchestrator/execution/dispatcher.py` 中新增 `topological_sort` 函数，按 depends_on 分波

## 7. Graph 构建

- [x] 7.1 重写 `build_graph()`：注册 reason/dispatch/execute/review/evolve/save_mem 节点，添加 conditional edges，使用 MemorySaver 作为 checkpointer
- [x] 7.2 清理旧的 discover_node/select_node/load_l2_node/plan_node/write_shared_node 和相关 edge 定义

## 8. OrchestratorAdapter 重构

- [x] 8.1 重写 `OrchestratorAdapter.stream_chat()`：使用新 build_graph()，传入 config={"configurable": {"thread_id": session_id}}
- [x] 8.2 处理各节点的 stream 事件：reason(text) → TEXT, reason(plan) → PLANNING, dispatch → PLANNING, execute → TEXT(agent), review(replan) → PLANNING, save_mem → DONE
- [x] 8.3 验证 chat() 方法兼容性

## 9. models 更新

- [x] 9.1 在 `orchestrator/models.py` 中确认 PlanOutput/TaskDef/TaskResult/DispatchResult 兼容新流程

## 10. 集成验证

- [ ] 10.1 手动测试闲聊场景：发送"你好"，验证 Orchestrator 直接文本回复
- [ ] 10.2 手动测试编排场景：发送"写登录页并审查"，验证完整 REASON → DISPATCH → EXECUTE → REVIEW → EVOLVE → SAVE_MEM 流程
- [ ] 10.3 手动测试跨轮次记忆：连续发三条消息，验证第三条能引用前两条的上下文
- [ ] 10.4 手动测试重规划：模拟任务失败场景，验证 REVIEW 触发重入 REASON
