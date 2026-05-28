from __future__ import annotations

from src.orchestrator.memory.evolution import EvolutionStore
from src.orchestrator.memory.pin_memory import PinMemory

PLAN_PROMPT = """\
你是一个 AI 项目经理（Orchestrator）。你的任务是根据用户需求，将其拆解为可由不同 Agent 并行或顺序执行的具体任务。

## 可用 Agents

{agents_desc}

{skills_section}

{tools_section}

{pin_context}

{evolution_context}

## 规则

1. 每个任务的 session_id 字段必须使用上面列表中的 agent id（加粗的名称，如 Alice、Bob），不要用类型
2. 任务数量不超过 5 个
3. 每个任务的 content 必须具体、可执行，包含明确的输入/输出期望
4. 任务按执行顺序排列，如果某些任务可以并行，在 overview 中说明
5. task_id 格式为 task-NNN（如 task-001, task-002）
6. 你可以使用提供的工具来收集信息（如读取文件、执行 skill 命令），但最终必须输出 JSON 格式的计划

## 输出格式

你必须只输出一个 JSON 对象，不要包含其他文字。格式如下：

```json
{{
  "overview": "整体规划概述",
  "tasks": [
    {{
      "task_id": "task-001",
      "session_id": "Alice",
      "title": "任务标题",
      "content": "任务详细描述"
    }}
  ]
}}
```

## 用户需求

{message}
"""


def build_planner_prompt(
    agents_desc: str,
    message: str,
    shared_dir: str,
    l2_content: dict[str, str] | None = None,
) -> str:
    pin_context = ""
    evolution_context = ""

    try:
        pm = PinMemory(common_dir=f"{shared_dir}/memory/common")
        pin_context = pm.get_context()
    except Exception:
        pass

    try:
        evo = EvolutionStore(shared_dir)
        evolution_context = evo.get_recent_experience(5)
    except Exception:
        pass

    # Build skills section from L2 content
    if l2_content:
        skill_parts = []
        for name, body in l2_content.items():
            skill_parts.append(f"### {name}\n\n{body}")
        skills_section = "## 可用 Skills\n\n" + "\n\n".join(skill_parts)
    else:
        skills_section = "## 可用 Skills\n\n(无)"

    tools_section = (
        "## 可用工具\n\n"
        "你可以使用以下工具来收集信息：\n"
        "- `read_file(path)`: 读取文件内容\n"
        "- `write_file(path, content)`: 写入文件到共享目录\n"
        "- `list_dir(path)`: 列出目录内容\n"
        "- `run_skill(skill, command, args)`: 执行已注册的 skill 命令\n"
        "- `load_resource(skill_name, resource_path)`: 加载 skill 的参考资源文件\n"
    )

    return PLAN_PROMPT.format(
        agents_desc=agents_desc,
        message=message,
        pin_context=pin_context,
        evolution_context=evolution_context,
        skills_section=skills_section,
        tools_section=tools_section,
    )
