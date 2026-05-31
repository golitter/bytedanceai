import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from src.adapters.orchestrator import _child_result_text
from src.orchestrator.models import TaskResult
from src.orchestrator.reporting.aggregator import build_final_summary_block


def _extract_json_block(text: str) -> dict:
    line = next(line for line in text.splitlines() if line.startswith("json: "))
    return json.loads(line[len("json: ") :])


def test_child_result_text_uses_structured_failure_block() -> None:
    result = TaskResult(
        task_id="task-004",
        agent="执行者",
        success=False,
        content="已经完成部分改造",
        duration=12.5,
        error_type="timeout",
        error_message="Task task-004 exceeded 300.0s",
    )

    text = _child_result_text(result)

    assert "[Timeout]" not in text
    assert "[Error]" not in text
    assert "type: task_failure" in text
    payload = _extract_json_block(text)
    assert payload["task_id"] == "task-004"
    assert payload["failureType"] == "timeout"


def test_final_summary_block_is_summary_first() -> None:
    block = build_final_summary_block(
        [
            TaskResult(
                task_id="task-001",
                agent="执行者",
                success=True,
                content="完成权限页面",
                duration=5,
            ),
            TaskResult(
                task_id="task-002",
                agent="检查者",
                success=False,
                content="",
                duration=300,
                error_type="timeout",
                error_message="Task task-002 exceeded 300.0s",
            ),
        ]
    )

    payload = _extract_json_block(block)
    assert payload["status"] == "partial"
    assert payload["completed"] == 1
    assert payload["failed"] == 1
    assert payload["details"][1]["status"] == "failed"
    assert "[Timeout]" not in block
