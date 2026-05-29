from __future__ import annotations

import logging
from typing import Annotated, TypedDict

from langgraph.graph import StateGraph

from src.orchestrator.models import DispatchResult

logger = logging.getLogger(__name__)


def _add_results(left: list, right: list) -> list:
    return left + right


class ExecuteState(TypedDict):
    execution_waves: list[list[DispatchResult]]
    task_results: Annotated[list, _add_results]
    shared_dir: str
    task_id: str
    cwd: str
    repo_path: str


def build_execute_subgraph() -> StateGraph:
    """Build a subgraph that executes tasks wave by wave.

    Each wave's tasks run in parallel; waves run sequentially.
    """
    graph = StateGraph(ExecuteState)
    graph.add_node("wave_execute", wave_execute_node)
    graph.set_entry_point("wave_execute")
    graph.set_finish_point("wave_execute")
    return graph.compile()


def wave_execute_node(state: ExecuteState) -> dict:
    """Placeholder: wave execution is driven by OrchestratorAdapter.

    The adapter iterates execution_waves and uses ExecutionEngine directly
    for async streaming. This node records the wave structure in state.
    """
    return {"task_results": []}
