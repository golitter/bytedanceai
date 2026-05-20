# Rule Engine — 规则引擎

## 实现了什么

可插拔规则系统，在 Agent 请求执行前进行安全检查和约束注入。规则通过修改 CLI 参数控制 Agent 行为，不侵入 Agent 内部逻辑。

## 怎么实现的

### BaseRule (`src/rules/base.py`)

抽象基类，定义规则接口：

```python
class BaseRule(ABC):
    name: str = ""           # 规则名称
    description: str = ""    # 描述
    phase: str = "pre"       # 阶段："pre" | "post"
    priority: int = 0        # 优先级，数值越大越先执行

    def check(context: dict) -> bool     # 校验是否通过
    def enforce(context: dict) -> dict   # 通过后注入约束
```

### RuleRegistry (`src/rules/registry.py`)

简单的注册表，支持动态注册和查找规则实例。

### RuleEngine (`src/rules/engine.py`)

执行核心逻辑：

```
evaluate(context) → (bool, dict)
```

1. 按 `priority` 降序排列所有 Rule
2. 对每个 Rule 依次执行：
   - `check(context)` → 如果返回 `False`，立即终止，返回 `(False, {"error": ..., "rule": ...})`
   - `enforce(context)` → 如果通过，合并约束结果
3. 全部通过返回 `(True, merged_context)`

合并规则：
- `system_prompt_append`：追加合并（所有 Rule 的提示词拼接）
- `allowed_tools`：追加合并
- `max_turns`：取第一个非 None 值
- `blocked`：任一 Rule 阻断则整体失败

### 内置规则 (`src/rules/builtin.py`)

#### SafetyRule（priority=10）

- `check`：始终通过
- `enforce`：
  - 注入安全提示词："不要执行破坏性命令，不要修改系统文件"
  - 过滤危险工具（如 `dangerouslyDisableSandbox`）

#### ScopeRule（priority=5）

- `check`：
  - `workspace_path` 为空 → 通过（不限制）
  - `workspace_path` 非绝对路径 → 失败
  - 绝对路径 → 通过
- `enforce`：
  - 注入约束："只允许修改指定目录下文件"

### 约束注入流程

```
AgentRequest
  ↓
RuleEngine.evaluate(context)
  ↓
SafetyRule.enforce → system_prompt_append: "安全约束..."
                    allowed_tools: [...]
  ↓
ScopeRule.enforce  → system_prompt_append: "只允许修改 /workspace"
  ↓
合并结果 → 传入 ClaudeCodeAdapter._build_command()
  ↓
CLI 参数: --append-system-prompt "安全约束...\n只允许修改 /workspace"
          --allowedTools Read,Write
```
