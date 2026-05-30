from __future__ import annotations

import json
import logging
from pathlib import Path
from typing import Annotated, TypedDict

import yaml
from langchain_core.messages import AIMessage, HumanMessage, SystemMessage, ToolMessage
from langchain_openai import ChatOpenAI
from langgraph.checkpoint.memory import MemorySaver
from langgraph.graph import END, StateGraph

from src.app.agent_config import get_agent_config_dir
from src.app.config import settings
from src.orchestrator.memory.evolution import EvolutionStore
from src.orchestrator.models import DispatchResult, PlanOutput, TaskDef
from src.orchestrator.planning.prompts import build_reason_prompt
from src.orchestrator.planning.skill_loader import (
    discover_skills,
    load_l2_content,
    select_skills,
)
from src.orchestrator.planning.tools import build_tools

logger = logging.getLogger(__name__)


def _add(left: list, right: list) -> list:
    return left + right


def _add_one(left: int, right: int) -> int:
    return left + right


class GraphState(TypedDict):
    message: str
    agents: list[dict]
    task_id: str
    shared_dir: str
    allowed_read_dirs: list[str]
    output_type: str  # "text" | "plan" | "error"
    text: str
    plan: PlanOutput | None
    dispatch_results: list[DispatchResult]
    execution_waves: list[list[DispatchResult]]
    task_results: Annotated[list, _add]
    task_status: dict
    needs_replan: bool
    replan_reason: str
    summary: str
    iteration: Annotated[int, _add_one]
    max_iterations: int
    memory_messages: Annotated[list, _add]
    # skill_prepare intermediate state
    system_prompt: str


def _build_agents_desc(agents: list[dict]) -> str:
    lines = []
    for a in agents:
        aid = a.get("id", "unknown")
        agent_type = a.get("type", aid)
        name = a.get("name", aid)
        caps = a.get("capabilities", [])
        cap_str = ", ".join(caps) if caps else "通用"
        lines.append(f"- **{aid}**（{name}，类型: {agent_type}）: {cap_str}")
    return "\n".join(lines)


def _skills_dir(shared_dir: str) -> Path:
    config_dir = get_agent_config_dir("orchestrator")
    return Path(shared_dir) / (config_dir or ".orchestrator") / "skills"


def _find_tool(tools: list, name: str):
    for t in tools:
        if t.name == name:
            return t
    return None


def _clean_ai_message(msg: AIMessage) -> AIMessage:
    if "reasoning_content" not in msg.additional_kwargs:
        return msg
    kw = {k: v for k, v in msg.additional_kwargs.items() if k != "reasoning_content"}
    return AIMessage(content=msg.content, tool_calls=msg.tool_calls, additional_kwargs=kw, id=msg.id)


def _write_shared_plan(
    shared_dir: str,
    task_id: str,
    plan: PlanOutput,
    dispatch_results: list[DispatchResult],
) -> None:
    """Write the orchestration plan into shared/.agent for taskctl consumers."""
    shared = Path(shared_dir).resolve()
    plans_dir = shared / "plans"
    plans_dir.mkdir(parents=True, exist_ok=True)
    (shared / "memory" / "common").mkdir(parents=True, exist_ok=True)

    (plans_dir / "overview.md").write_text(plan.overview, encoding="utf-8")

    config_tasks: list[dict] = []
    by_task_id = {dr.task_id: dr for dr in dispatch_results}
    for task in plan.tasks:
        plan_file = f"plans/{task.task_id}.md"
        dr = by_task_id.get(task.task_id)
        session_id = dr.real_session_id if dr and dr.real_session_id else task.session_id
        agent_id = dr.agent if dr else task.session_id
        agent_type = dr.agent_type if dr else ""

        body = "\n".join(
            [
                f"# {task.title or task.task_id}",
                "",
                f"- task_id: {task.task_id}",
                f"- agent: {agent_id}",
                f"- agent_type: {agent_type}",
                f"- session_id: {session_id}",
                "",
                "## Task",
                "",
                task.content,
                "",
            ]
        )
        (shared / plan_file).write_text(body, encoding="utf-8")
        config_tasks.append(
            {
                "task_id": task.task_id,
                "session_id": session_id,
                "agent": agent_id,
                "agent_type": agent_type,
                "file": plan_file,
            }
        )

    config = {
        "task_id": task_id,
        "overview_file": "plans/overview.md",
        "tasks": config_tasks,
    }
    (shared / "config.yaml").write_text(
        yaml.safe_dump(config, allow_unicode=True, sort_keys=False),
        encoding="utf-8",
    )


# --- Skill Prepare Node (fast, yields SSE event within seconds) ---


def skill_prepare_node(state: GraphState) -> dict:
    """L1 → L2 skill discovery + prompt construction. Runs in seconds."""
    skills_dir_path = _skills_dir(state["shared_dir"])
    l1_skills = discover_skills(skills_dir_path)

    selection_message = state["replan_reason"] if state.get("replan_reason") else state["message"]
    selected_names = select_skills(l1_skills, selection_message)
    l2 = load_l2_content(selected_names, skills_dir_path)

    agents_desc = _build_agents_desc(state["agents"])
    system_prompt = build_reason_prompt(
        agents_desc=agents_desc,
        message=state["message"],
        shared_dir=state["shared_dir"],
        l2_content=l2,
        replan_reason=state.get("replan_reason"),
    )

    return {"system_prompt": system_prompt}


# --- REASON Node (LLM tool-calling loop) ---


def reason_node(state: GraphState) -> dict:
    """REASON node: LLM tool-calling loop.

    Determines output_type: "text" (chitchat) or "plan" (orchestration).
    """
    try:
        llm = ChatOpenAI(
            model=settings.llm.model,
            base_url=settings.llm.base_url,
            api_key=settings.llm.api_key,
        )
        tools = build_tools(state["shared_dir"], state.get("allowed_read_dirs"))
        llm_with_tools = llm.bind_tools(tools)

        # Use pre-built system prompt from skill_prepare_node
        system_prompt = state.get("system_prompt", "")
        messages: list = [SystemMessage(content=system_prompt)]

        for msg in state.get("memory_messages", []):
            messages.append(msg)

        messages.append(HumanMessage(content=state["message"]))

        if state.get("replan_reason"):
            messages.append(
                HumanMessage(content=f"[重规划请求] 以下任务执行失败，请重新规划：\n{state['replan_reason']}")
            )

        max_iterations = 10
        for i in range(max_iterations):
            response = llm_with_tools.invoke(messages)

            if not response.tool_calls:
                return {
                    "output_type": "text",
                    "text": response.content,
                    "plan": None,
                    "memory_messages": [HumanMessage(content=state["message"]), response],
                }

            plan_call = None
            other_calls = []
            for tc in response.tool_calls:
                if tc["name"] == "plan_and_dispatch":
                    plan_call = tc
                else:
                    other_calls.append(tc)

            if plan_call is not None:
                args = plan_call["args"]
                overview = args.get("overview", "")
                raw_tasks = args.get("tasks", [])

                tasks = []
                for t in raw_tasks:
                    if isinstance(t, dict):
                        tasks.append(
                            TaskDef(
                                task_id=t.get("task_id", f"task-{len(tasks) + 1:03d}"),
                                session_id=t.get("session_id", ""),
                                title=t.get("title", ""),
                                content=t.get("content", ""),
                            )
                        )

                plan = PlanOutput(overview=overview, tasks=tasks)
                tool_messages = [
                    ToolMessage(content="plan_generated", tool_call_id=plan_call["id"]),
                ]
                for tc in other_calls:
                    tool_fn = _find_tool(tools, tc["name"])
                    if tool_fn:
                        try:
                            result = tool_fn.invoke(tc["args"])
                        except Exception as e:
                            result = f"Error: {e}"
                    else:
                        result = f"Error: unknown tool '{tc['name']}'"
                    wrapped = json.dumps(
                        {"tool": tc["name"], "args": tc["args"], "output": result},
                        ensure_ascii=False,
                    )
                    tool_messages.append(ToolMessage(content=wrapped, tool_call_id=tc["id"]))

                return {
                    "output_type": "plan",
                    "text": "",
                    "plan": plan,
                    "memory_messages": [
                        HumanMessage(content=state["message"]),
                        _clean_ai_message(response),
                        *tool_messages,
                    ],
                }

            messages.append(_clean_ai_message(response))
            for tc in response.tool_calls:
                tool_fn = _find_tool(tools, tc["name"])
                if tool_fn is None:
                    result = f"Error: unknown tool '{tc['name']}'"
                else:
                    try:
                        result = tool_fn.invoke(tc["args"])
                    except Exception as e:
                        result = f"Error: {e}"
                wrapped = json.dumps(
                    {"tool": tc["name"], "args": tc["args"], "output": result},
                    ensure_ascii=False,
                )
                messages.append(ToolMessage(content=wrapped, tool_call_id=tc["id"]))

        logger.warning("Reason node reached max_iterations=%d", max_iterations)
        return {
            "output_type": "text",
            "text": "规划超时，请重新描述需求",
            "plan": None,
            "memory_messages": [HumanMessage(content=state["message"])],
        }
    except Exception:
        logger.exception("Reason node failed unexpectedly")
        return {
            "output_type": "error",
            "text": "Reasoning failed",
            "plan": None,
        }


# --- DISPATCH Node ---


def dispatch_node(state: GraphState) -> dict:
    """Convert PlanOutput to DispatchResults and topologically sort into waves."""
    from src.orchestrator.execution.dispatcher import Dispatcher, topological_sort

    plan = state["plan"]
    if not plan:
        return {"dispatch_results": [], "execution_waves": []}

    dispatcher = Dispatcher(state["agents"])
    dispatch_results = dispatcher.dispatch(plan)
    waves = topological_sort(dispatch_results)
    try:
        _write_shared_plan(state["shared_dir"], state["task_id"], plan, dispatch_results)
    except Exception:
        logger.exception("Failed to write orchestrator plan into shared dir")

    return {"dispatch_results": dispatch_results, "execution_waves": waves}


# --- REVIEW Node ---


def review_node(state: GraphState) -> dict:
    """Check task results for failures. Set needs_replan if failures exist and under iteration limit."""
    task_results = state.get("task_results", [])
    failed = [tr for tr in task_results if not tr.get("success", True)]

    iteration = state.get("iteration", 0)
    max_iterations = state.get("max_iterations", 3)

    if not failed:
        return {"needs_replan": False, "replan_reason": ""}

    if iteration >= max_iterations:
        logger.warning("Review: max_iterations=%d reached, accepting partial results", max_iterations)
        return {"needs_replan": False, "replan_reason": ""}

    failure_details = []
    for tr in failed:
        failure_details.append(
            f"- 任务 {tr.get('task_id', '?')} (agent: {tr.get('agent', '?')}): {tr.get('content', '')[:200]}"
        )
    replan_reason = "以下任务执行失败，请重新规划：\n" + "\n".join(failure_details)

    return {"needs_replan": True, "replan_reason": replan_reason, "iteration": 1}


# --- EVOLVE Node ---


def evolve_node(state: GraphState) -> dict:
    """Record orchestration experience in EvolutionStore."""
    try:
        evolution = EvolutionStore(state["shared_dir"])
        plan = state.get("plan")
        overview = plan.overview if plan else ""
        task_results = state.get("task_results", [])
        all_success = all(tr.get("success", True) for tr in task_results)
        results_summary = "; ".join(
            f"{tr.get('task_id', '?')}: {'✅' if tr.get('success') else '❌'}" for tr in task_results
        )
        evolution.record(
            message=state["message"],
            plan_summary=overview[:200],
            results_summary=results_summary[:200],
            success=all_success,
            agent_performance=[
                {
                    "agent_id": tr.get("agent", ""),
                    "success": tr.get("success", False),
                    "duration": tr.get("duration", 0),
                }
                for tr in task_results
            ],
        )
    except Exception:
        logger.exception("Evolve node failed")
    return {}


# --- SAVE_MEM Node ---


def save_mem_node(state: GraphState) -> dict:
    return {}


# --- Conditional Routers ---


def route_by_output_type(state: GraphState) -> str:
    output_type = state.get("output_type", "error")
    if output_type == "text":
        return "save_mem"
    elif output_type == "plan":
        return "dispatch"
    else:
        return END


def route_by_review(state: GraphState) -> str:
    if state.get("needs_replan", False):
        return "skill_prepare"
    return "evolve"


# --- Graph Builder ---


def build_graph() -> StateGraph:
    graph = StateGraph(GraphState)

    graph.add_node("skill_prepare", skill_prepare_node)
    graph.add_node("reason", reason_node)
    graph.add_node("dispatch", dispatch_node)
    graph.add_node("execute", _execute_placeholder)
    graph.add_node("review", review_node)
    graph.add_node("evolve", evolve_node)
    graph.add_node("save_mem", save_mem_node)

    graph.set_entry_point("skill_prepare")

    graph.add_edge("skill_prepare", "reason")
    graph.add_conditional_edges("reason", route_by_output_type)
    graph.add_edge("dispatch", "execute")
    graph.add_edge("execute", "review")
    graph.add_conditional_edges("review", route_by_review)
    graph.add_edge("evolve", "save_mem")
    graph.set_finish_point("save_mem")

    memory = MemorySaver()
    return graph.compile(checkpointer=memory)


def _execute_placeholder(state: GraphState) -> dict:
    """Placeholder execute node — actual execution is handled by OrchestratorAdapter."""
    return {"task_results": []}
