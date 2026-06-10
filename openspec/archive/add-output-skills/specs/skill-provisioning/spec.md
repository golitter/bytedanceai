## MODIFIED Requirements

### Requirement: SkillRule registered in rule engine
The RuleEngine SHALL include SkillRule alongside existing SafetyRule, ScopeRule, and TaskctlRule. SkillRule SHALL have a priority lower than SafetyRule (10) and ScopeRule (5) to ensure it never blocks execution.

#### Scenario: SkillRule does not block other rules
- **WHEN** RuleEngine evaluates rules in priority order
- **THEN** SkillRule SHALL execute after SafetyRule and ScopeRule but its `check()` SHALL always return True
