## Requirements

### Requirement: read_file tool
Orchestrator plan_node agent loop SHALL expose a `read_file` tool that reads the content of a file at the given path and returns it as a string. No restrictions on readable paths.

#### Scenario: Read an existing file
- **WHEN** LLM calls `read_file(path="src/skills/builtin/taskctl/SKILL.md")`
- **THEN** tool SHALL return the file content as a string

#### Scenario: Read a non-existent file
- **WHEN** LLM calls `read_file(path="/nonexistent/file.txt")`
- **THEN** tool SHALL return an error string like "Error: file not found"

### Requirement: write_file tool with path restriction
Orchestrator plan_node agent loop SHALL expose a `write_file` tool that writes content to a file. The target path MUST be within `shared_dir` (resolved from GraphState). Paths containing `..` or resolving outside shared_dir SHALL be rejected.

#### Scenario: Write to shared_dir successfully
- **WHEN** LLM calls `write_file(path="notes/planning.md", content="# Plan\n...")` and `shared_dir` is `/worktrees/task-1/shared/.agent`
- **THEN** tool SHALL write content to `/worktrees/task-1/shared/.agent/notes/planning.md`, creating parent directories if needed, and return "OK"

#### Scenario: Reject path traversal
- **WHEN** LLM calls `write_file(path="../../../etc/passwd", content="hacked")`
- **THEN** tool SHALL return "Error: path outside shared_dir" without writing any file

#### Scenario: Reject absolute path outside shared_dir
- **WHEN** LLM calls `write_file(path="/tmp/evil.txt", content="...")`
- **THEN** tool SHALL return "Error: path outside shared_dir" without writing any file

### Requirement: list_dir tool
Orchestrator plan_node agent loop SHALL expose a `list_dir` tool that lists the contents of a directory. Returns file and subdirectory names, one per line. Subdirectories SHALL be suffixed with `/`.

#### Scenario: List a non-empty directory
- **WHEN** LLM calls `list_dir(path="src/skills/builtin")`
- **THEN** tool SHALL return "render/\ntaskctl/\n" (or similar listing)

#### Scenario: List a non-existent directory
- **WHEN** LLM calls `list_dir(path="/nonexistent")`
- **THEN** tool SHALL return "Error: directory not found"

### Requirement: run_skill tool with manifest validation
Orchestrator plan_node agent loop SHALL expose a `run_skill` tool that executes a registered skill binary. The `skill` parameter SHALL be a `Literal` type restricted to skill names present in `config.yaml` `skills.manifest`. The tool SHALL execute the skill binary located at `settings.skills.builtin_dir / {skill_name} / {skill_name}` with `cwd` set to `shared_dir`. Execution timeout SHALL be 30 seconds. Output SHALL be truncated to 4096 characters.

#### Scenario: Execute a registered skill command
- **WHEN** LLM calls `run_skill(skill="taskctl", command="summary")`
- **THEN** tool SHALL execute `src/skills/builtin/taskctl/taskctl summary` with `cwd=shared_dir`, capture stdout, and return it (truncated if > 4096 chars)

#### Scenario: Reject unregistered skill
- **WHEN** LLM calls `run_skill(skill="evil-script", command="...")`
- **THEN** tool SHALL return "Error: unknown skill 'evil-script'" without executing anything

#### Scenario: Skill execution times out
- **WHEN** a skill binary runs for more than 30 seconds
- **THEN** tool SHALL terminate the process and return "Error: skill execution timed out (30s)"

#### Scenario: Skill binary not found
- **WHEN** the manifest lists a skill but the binary does not exist on disk
- **THEN** tool SHALL return "Error: skill binary not found for '{skill_name}'"

### Requirement: load_skill_resource tool with path restriction
Orchestrator plan_node agent loop SHALL expose a `load_skill_resource` tool that loads L3 resource files from a skill's `references/` or `assets/` directory. The `skill_name` MUST be in `skills.manifest`. The `resource_path` SHALL NOT contain `..`. The tool SHALL read from `settings.skills.builtin_dir / {skill_name} / {resource_path}`.

#### Scenario: Load an existing resource
- **WHEN** LLM calls `load_skill_resource(skill_name="taskctl", resource_path="references/best-practices.md")`
- **THEN** tool SHALL return the content of `src/skills/builtin/taskctl/references/best-practices.md`

#### Scenario: Reject path traversal in resource_path
- **WHEN** LLM calls `load_skill_resource(skill_name="taskctl", resource_path="../../../etc/passwd")`
- **THEN** tool SHALL return "Error: invalid resource path"

#### Scenario: Resource file not found
- **WHEN** LLM calls `load_skill_resource(skill_name="taskctl", resource_path="references/nonexistent.md")`
- **THEN** tool SHALL return "Error: resource file not found"

### Requirement: Tool definitions are auto-generated from manifest
The `run_skill` tool's `skill` Literal parameter SHALL be dynamically populated from `config.yaml` `skills.manifest` keys at graph build time. Adding a new skill to the manifest SHALL automatically make it available to the orchestrator without code changes.

#### Scenario: New skill appears in tool
- **WHEN** a new skill "analyze" is added to `skills.manifest` in config.yaml
- **THEN** the orchestrator SHALL be able to call `run_skill(skill="analyze", ...)` after restart

### Requirement: plan_node agent loop termination
The plan_node agent loop SHALL terminate when either: (1) the LLM returns a response with no `tool_calls`, or (2) the iteration count reaches `max_iterations=10`. Upon termination, the LLM's final response SHALL be parsed as a PlanOutput JSON.

#### Scenario: LLM outputs plan after using tools
- **WHEN** LLM calls `run_skill("taskctl", "summary")`, receives result, then outputs PlanOutput JSON with no tool_calls
- **THEN** plan_node SHALL parse the JSON and return `{"plan": PlanOutput(...)}`

#### Scenario: Max iterations reached
- **WHEN** LLM calls tools for 10 consecutive iterations without outputting PlanOutput
- **THEN** plan_node SHALL return `{"plan": None}`

#### Scenario: LLM outputs invalid JSON
- **WHEN** LLM terminates with a response that is not valid PlanOutput JSON
- **THEN** plan_node SHALL attempt to extract JSON from the response (handling markdown code blocks), and return `{"plan": None}` if extraction fails
