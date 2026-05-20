from abc import ABC, abstractmethod


class BaseRule(ABC):
    name: str = ""
    description: str = ""
    phase: str = "pre"  # "pre" | "post"
    priority: int = 0

    @abstractmethod
    def check(self, context: dict) -> bool: ...

    @abstractmethod
    def enforce(self, context: dict) -> dict: ...
