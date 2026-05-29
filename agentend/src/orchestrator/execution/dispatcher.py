from __future__ import annotations

from src.orchestrator.models import DispatchResult, PlanOutput


class Dispatcher:
    def __init__(self, agents: list[dict]) -> None:
        self.agents = agents
        self._agent_map = {a["id"]: a for a in agents}

    def dispatch(self, plan: PlanOutput) -> list[DispatchResult]:
        results: list[DispatchResult] = []
        valid_ids = set(self._agent_map.keys())
        for task in plan.tasks:
            if task.session_id not in valid_ids:
                # Fallback: assign to first available agent if session_id is invalid (e.g. a skill name)
                fallback = next(iter(self._agent_map), None)
                if fallback:
                    task.session_id = fallback
            agent_cfg = self._agent_map.get(task.session_id, {})
            workspace_path = agent_cfg.get("workspace_path", "")
            real_session_id = agent_cfg.get("session_id", "")
            agent_type = agent_cfg.get("type", task.session_id)

            results.append(
                DispatchResult(
                    task_id=task.task_id,
                    agent=task.session_id,
                    agent_type=agent_type,
                    real_session_id=real_session_id,
                    mention=f"@{task.session_id}",
                    content=task.content,
                    depends_on=[],
                    workspace_path=workspace_path,
                )
            )
        return results


def topological_sort(dispatch_results: list[DispatchResult]) -> list[list[DispatchResult]]:
    """Sort DispatchResults into execution waves based on depends_on.

    Tasks within the same wave can execute in parallel.
    Waves execute sequentially.
    """
    if not dispatch_results:
        return []

    # Build lookup and dependency graph
    by_id: dict[str, DispatchResult] = {dr.task_id: dr for dr in dispatch_results}
    all_ids = set(by_id.keys())

    # Compute in-degree for each task
    in_degree: dict[str, int] = {tid: 0 for tid in all_ids}
    dependents: dict[str, list[str]] = {tid: [] for tid in all_ids}

    for dr in dispatch_results:
        for dep in dr.depends_on:
            if dep in all_ids:
                in_degree[dr.task_id] += 1
                dependents[dep].append(dr.task_id)

    # Kahn's algorithm
    waves: list[list[DispatchResult]] = []
    remaining = dict(in_degree)

    while remaining:
        # Find tasks with zero in-degree
        ready = [tid for tid, deg in remaining.items() if deg == 0]
        if not ready:
            # Cycle detected — put remaining tasks in one wave
            waves.append([by_id[tid] for tid in remaining])
            break

        wave = [by_id[tid] for tid in ready]
        waves.append(wave)

        for tid in ready:
            del remaining[tid]
            for dep_tid in dependents[tid]:
                if dep_tid in remaining:
                    remaining[dep_tid] -= 1

    return waves
