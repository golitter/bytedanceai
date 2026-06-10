## Requirements

### Requirement: discover_node scans skill frontmatter (L1)
discover_node SHALL scan the `skills.builtin_dir` directory, parse YAML frontmatter from each subdirectory's SKILL.md file, and produce a list of `{name, description}` dicts stored as `l1_skills` in GraphState. Only the frontmatter (between `---` markers) SHALL be parsed; SKILL.md body SHALL NOT be read.

#### Scenario: Discover two skills
- **WHEN** `src/skills/builtin/` contains `taskctl/SKILL.md` and `render/SKILL.md`
- **THEN** `l1_skills` SHALL equal `[{"name": "taskctl", "description": "多Agent协作上下文管理..."}, {"name": "render", "description": "富媒体卡片输出..."}]`

#### Scenario: Skip directories without SKILL.md
- **WHEN** a subdirectory `foo/` exists but has no `SKILL.md`
- **THEN** `foo` SHALL NOT appear in `l1_skills`

#### Scenario: Skip SKILL.md with no name in frontmatter
- **WHEN** a SKILL.md has frontmatter without a `name` field
- **THEN** that skill SHALL NOT appear in `l1_skills`

### Requirement: select_node uses LLM semantic matching
select_node SHALL use a single LLM call to select 1~N relevant skills from `l1_skills` based on the user's message. The LLM SHALL receive the L1 skill list (name + description) and the user message, and return comma-separated skill names. Returned names not present in `l1_skills` SHALL be filtered out.

#### Scenario: Select one relevant skill
- **WHEN** user message is "帮我检查任务状态" and `l1_skills` contains taskctl and render
- **THEN** select_node SHALL return `selected_skill_names` containing `["taskctl"]` (or a subset containing taskctl)

#### Scenario: No skills selected
- **WHEN** user message is irrelevant to all skills
- **THEN** select_node SHALL return `selected_skill_names = []`

#### Scenario: Invalid skill names filtered
- **WHEN** LLM returns "taskctl, unknown-skill"
- **THEN** select_node SHALL filter to only `["taskctl"]`

### Requirement: load_l2_node loads selected skill instructions
load_l2_node SHALL read the full SKILL.md body (excluding frontmatter) for each skill in `selected_skill_names` and store it in `l2_content` keyed by skill name. Skills not in `selected_skill_names` SHALL NOT be loaded.

#### Scenario: Load one skill's L2
- **WHEN** `selected_skill_names = ["taskctl"]`
- **THEN** `l2_content` SHALL contain `{"taskctl": "<SKILL.md body text>"}` and SHALL NOT contain `render`

#### Scenario: No skills selected
- **WHEN** `selected_skill_names = []`
- **THEN** `l2_content` SHALL be `{}`

#### Scenario: Selected skill has no SKILL.md
- **WHEN** `selected_skill_names = ["missing-skill"]` and no `missing-skill/SKILL.md` exists
- **THEN** `l2_content` SHALL be `{}` (skill silently skipped)

### Requirement: L3 resources loaded on-demand via tool
L3 resource files SHALL NOT be loaded by graph nodes. They SHALL be loaded only when the LLM calls the `load_skill_resource` tool during the plan_node agent loop. The tool SHALL read files from `references/` or `assets/` subdirectories of the skill directory.

#### Scenario: LLM loads L3 during planning
- **WHEN** L2 instruction for taskctl says "使用 load_skill_resource 读取 references/best-practices.md" and LLM calls the tool
- **THEN** tool SHALL return the file content and it SHALL be available in the LLM's next turn

#### Scenario: L3 never loaded without explicit tool call
- **WHEN** no tool call to `load_skill_resource` is made
- **THEN** no L3 content SHALL be loaded or injected into the LLM context

### Requirement: GraphState includes progressive disclosure fields
OrchestratorGraphState SHALL include the following fields in addition to existing fields: `l1_skills` (list of dicts), `selected_skill_names` (list of strings), `l2_content` (dict of skill_name → SKILL.md body), `l3_content` (dict of "skill:resource_path" → file content).

#### Scenario: State fields populated after discover
- **WHEN** discover_node completes
- **THEN** `l1_skills` SHALL be a non-empty list, `selected_skill_names`, `l2_content`, `l3_content` SHALL be empty/defaults

#### Scenario: State fields populated after full flow
- **WHEN** the full discover → select → load_l2 → plan flow completes with `selected_skill_names = ["taskctl"]` and LLM loads one L3 resource
- **THEN** `l1_skills` SHALL have all skills, `selected_skill_names` SHALL have taskctl, `l2_content` SHALL have taskctl's body, `l3_content` SHALL have the loaded resource

### Requirement: System prompt dynamically injects L2 instructions
The plan_node system prompt SHALL dynamically include L2 instructions from `l2_content` for selected skills. Skills not in `selected_skill_names` SHALL NOT appear in the system prompt.

#### Scenario: Prompt includes selected skill instructions
- **WHEN** `l2_content = {"taskctl": "## 概述\n..."}` and plan_node builds the system prompt
- **THEN** the prompt SHALL contain the taskctl L2 body under a "## 可用 Skills" section

#### Scenario: Prompt with no selected skills
- **WHEN** `l2_content = {}` and plan_node builds the system prompt
- **THEN** the prompt SHALL contain "## 可用 Skills\n(无)" or omit the section entirely
