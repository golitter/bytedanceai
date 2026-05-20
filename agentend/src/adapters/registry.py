from src.adapters.base import BaseAgentAdapter


class AdapterRegistry:
    def __init__(self) -> None:
        self._adapters: dict[str, type[BaseAgentAdapter]] = {}

    def register(self, agent_type: str, adapter_cls: type[BaseAgentAdapter]) -> None:
        self._adapters[agent_type] = adapter_cls

    def get(self, agent_type: str) -> type[BaseAgentAdapter]:
        if agent_type not in self._adapters:
            raise ValueError(f"Unknown agent type: {agent_type}")
        return self._adapters[agent_type]

    def list(self) -> list[str]:
        return list(self._adapters.keys())
