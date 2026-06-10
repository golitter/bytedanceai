## ADDED Requirements

### Requirement: SkillRule injects output skill prompt
SkillRule SHALL append a system prompt describing the `aka_yhy` block format and 5 available card types to the agent's system prompt. The prompt SHALL be injected via the existing `system_prompt_append` mechanism.

#### Scenario: SkillRule enforce returns system_prompt_append
- **WHEN** `SkillRule().enforce()` is called
- **THEN** SHALL return a dict with key `system_prompt_append` containing the output skill prompt text

#### Scenario: SkillRule check always passes
- **WHEN** `SkillRule().check()` is called with any kwargs
- **THEN** SHALL return `True` (never blocks execution)

### Requirement: SkillRule priority is lower than safety rules
SkillRule SHALL have a priority value lower than SafetyRule (10) and ScopeRule (5), so it never blocks execution.

#### Scenario: SkillRule does not block safety rules
- **WHEN** RuleEngine evaluates rules in priority order
- **THEN** SkillRule SHALL execute after SafetyRule and ScopeRule

### Requirement: Output skill prompt defines 5 card types
The output skill prompt SHALL describe exactly 5 card types with their `aka_yhy` format: `html-render` (HTML content in block), `image` (path field), `attachment` (path field), `diff` (no extra fields), `preview` (url field).

#### Scenario: Prompt mentions all 5 types
- **WHEN** the injected system prompt is inspected
- **THEN** SHALL contain descriptions and format examples for: html-render, image, attachment, diff, preview
