## 1. 依赖与配置

- [x] 1.1 安装 langgraph、langchain-anthropic、langchain-core 依赖（`uv add`）
- [x] 1.2 config.yaml 末尾新增 `llm` 配置段（provider, model, api_key）
- [x] 1.3 config.py 新增 `LlmConfig` 类并加到 `Settings`

## 2. 数据模型与 Schema

- [x] 2.1 新建 `src/orchestrator/` 目录及 `__init__.py`
- [x] 2.2 新建 `src/orchestrator/models.py`，定义 `TaskDef` 和 `PlanOutput` Pydantic model
- [x] 2.3 schemas/events.py 的 EventType 枚举新增 `PLANNING`
- [x] 2.4 schemas/request.py 的 AgentType 枚举新增 `ORCHESTRATOR`

## 3. LangGraph 核心

- [x] 3.1 新建 `src/orchestrator/prompts.py`，定义 `PLAN_PROMPT`
- [x] 3.2 新建 `src/orchestrator/graph.py`，实现 `GraphState`、`plan_node`、`write_shared_node`、`build_graph()`

## 4. Adapter 桥接

- [x] 4.1 新建 `src/adapters/orchestrator.py`，实现 `OrchestratorAdapter`（BaseAgentAdapter 子类）
- [x] 4.2 `src/adapters/__init__.py` 新增 OrchestratorAdapter 导出

## 5. 注册与 API 接入

- [x] 5.1 `src/app/dependencies.py` 注册 OrchestratorAdapter 到 AdapterRegistry
- [x] 5.2 `src/api/v1/agent.py` 中 stream 和 execute 两处 kwargs 构建处加 Orchestrator config 透传

## 6. 验证

- [x] 6.1 启动服务，curl 调用 `/v1/agent/execute`（agent_type=orchestrator），验证 shared/ 目录产出（config.yaml、overview.md、tasks/*.md）
