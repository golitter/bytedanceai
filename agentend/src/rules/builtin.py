from src.rules.base import BaseRule

# Tools considered dangerous; blocked by SafetyRule
_DANGEROUS_TOOLS = {"dangerouslyDisableSandbox"}


class SafetyRule(BaseRule):
    name = "safety"
    description = "Injects safety constraints and blocks dangerous tools"
    phase = "pre"
    priority = 10

    def check(self, context: dict) -> bool:
        return True

    def enforce(self, context: dict) -> dict:
        existing_tools = set(context.get("allowed_tools", []))
        safe_tools = list(existing_tools - _DANGEROUS_TOOLS) if existing_tools else []

        result: dict = {
            "system_prompt_append": (
                "You are operating in a managed environment. "
                "Do not execute destructive commands (rm -rf, format, etc.). "
                "Do not access or modify system files outside the workspace."
            ),
        }
        if safe_tools:
            result["allowed_tools"] = safe_tools
        return result


class ScopeRule(BaseRule):
    name = "scope"
    description = "Validates workspace_path and injects scope constraints"
    phase = "pre"
    priority = 5

    def check(self, context: dict) -> bool:
        workspace_path = context.get("workspace_path")
        if not workspace_path:
            return True  # No workspace constraint is OK
        # Basic validation: workspace_path must be an absolute path
        return workspace_path.startswith("/")

    def enforce(self, context: dict) -> dict:
        workspace_path = context.get("workspace_path")
        if workspace_path:
            return {
                "system_prompt_append": f"Only modify files under: {workspace_path}",
            }
        return {}
