## MODIFIED Requirements

### Requirement: AgentRequest data model
系统 SHALL 定义 `AgentRequest` Pydantic 模型，包含字段：`task_id`（str）、`conversation_id`（str）、`session_id`（可选 str）、`message`（str）、`agent_type`（str，默认 "claude-code"）、`stream`（bool，默认 True）、`system_prompt`（可选 str）、`rules`（list[str]，默认空）、`workspace_path`（可选 str）、`repo_path`（可选 str）、`config`（可选 dict）。新增 `repo_path` 字段，当 `workspace_path` 为空时 MUST 与 `repo_path` 一起传入以自动创建 workspace。

#### Scenario: Request with repo_path triggers workspace creation
- **WHEN** 请求包含 `repo_path="/repos/project"` 但 `workspace_path` 为空
- **THEN** 系统 SHALL 自动创建 workspace 并将路径绑定到 session

#### Scenario: Request with workspace_path skips creation
- **WHEN** 请求包含 `workspace_path="/workspaces/task-1/frontend"`
- **THEN** 系统 SHALL 直接使用该路径，不创建新 workspace

### Requirement: StreamEvent source agent identification
`StreamEvent` SHALL 在 `content` 中包含 `agent_type` 字段，标识事件来源 Agent（如 `"claude-code"` 或 `"opencode"`）。

#### Scenario: Event from Claude Code
- **WHEN** ClaudeCodeAdapter 生成 StreamEvent
- **THEN** `content["agent_type"]` SHALL 为 `"claude-code"`

#### Scenario: Event from OpenCode
- **WHEN** OpenCodeAdapter 生成 StreamEvent
- **THEN** `content["agent_type"]` SHALL 为 `"opencode"`
