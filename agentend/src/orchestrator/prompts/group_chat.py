"""Group chat context prompt template for cross-agent memory.

Formats window messages (messages from other agents since this agent last spoke)
into a system prompt section injected via GroupChatRule.
"""

GROUP_CHAT_CONTEXT = """\
## 群聊上下文

你正在参与一个多 Agent 协作群聊。以下是你上次发言后，其他成员发出的消息。
请参考这些内容来执行你的任务——了解当前进展、避免重复工作、与其他成员协作。

{messages}
"""


def build_group_chat_context(cross_round_messages: list[dict] | None = None) -> str:
    """Format cross-agent window messages into a prompt section.

    Args:
        cross_round_messages: List of message dicts with keys:
            role (str): "user" or "agent"
            agent_name (str): Name of the agent
            content (str): Message content (already truncated by backend)

    Returns:
        Formatted prompt string, or empty string if no valid messages.
    """
    if not cross_round_messages:
        return ""

    lines: list[str] = []
    for msg in cross_round_messages:
        role = msg.get("role", "")
        name = msg.get("agent_name", "")
        content = msg.get("content", "")
        if role == "user":
            lines.append(f"👤 用户:\n{content}")
        elif role == "agent" and name:
            lines.append(f"🤖 {name}:\n{content}")

    if not lines:
        return ""

    return GROUP_CHAT_CONTEXT.format(messages="\n\n".join(lines))
