## ADDED Requirements

### Requirement: Memory persists across conversation turns

The Orchestrator SHALL maintain a conversation history that persists across turns within the same session. Each turn's user message, LLM response, tool calls, and tool results SHALL be accumulated.

#### Scenario: Turn 2 references Turn 1 results
- **WHEN** user says "基于上面的结果细化" in Turn 2, and Turn 1 produced plan results
- **THEN** REASON node loads Turn 1 history from memory, and LLM can reference previous results

#### Scenario: Chitchat turn is also remembered
- **WHEN** user says "你好" in Turn 1 (text reply) and "刚才聊了什么" in Turn 2
- **THEN** REASON node loads Turn 1 chitchat from memory, and LLM can reference it

### Requirement: Memory uses LangGraph MemorySaver

The system SHALL use `langgraph.checkpoint.memory.MemorySaver` as the checkpointer. The session_id SHALL be used as the thread_id for checkpoint scoping.

#### Scenario: Memory scoped to session
- **WHEN** two different sessions exist with session_id "session-a" and "session-b"
- **THEN** each session has independent memory, and cross-session access is not possible

### Requirement: Memory messages accumulate with reducer

The `memory_messages` field in GraphState SHALL use an additive reducer (`Annotated[list, add]`). Each turn's messages (Human + AI + ToolMessages) are appended, not replaced.

#### Scenario: Three turns accumulate
- **WHEN** Turn 1 has 2 messages, Turn 2 has 4 messages, Turn 3 has 2 messages
- **THEN** memory_messages contains all 8 messages in chronological order
