from src.rules.base import BaseRule


class RuleEngine:
    def __init__(self, rules: list[BaseRule] | None = None) -> None:
        self._rules: list[BaseRule] = sorted(rules or [], key=lambda r: r.priority, reverse=True)

    def evaluate(self, context: dict) -> tuple[bool, dict]:
        merged: dict = {
            "system_prompt_append": [],
            "allowed_tools": [],
            "max_turns": None,
        }

        for rule in self._rules:
            if not rule.check(context):
                return False, {"error": f"Rule '{rule.name}' check failed", "rule": rule.name}

            result = rule.enforce(context)

            if result.get("system_prompt_append"):
                merged["system_prompt_append"].append(result["system_prompt_append"])
            if result.get("allowed_tools"):
                merged["allowed_tools"].extend(result["allowed_tools"])
            if result.get("max_turns") is not None and merged["max_turns"] is None:
                merged["max_turns"] = result["max_turns"]
            if result.get("blocked"):
                return False, {"error": result.get("error", "blocked by rule"), "rule": rule.name}

        return True, merged
