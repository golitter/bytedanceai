from src.rules.base import BaseRule


class RuleRegistry:
    def __init__(self) -> None:
        self._rules: list[BaseRule] = []

    def register(self, rule: BaseRule) -> None:
        self._rules.append(rule)

    def list(self) -> list[BaseRule]:
        return list(self._rules)
