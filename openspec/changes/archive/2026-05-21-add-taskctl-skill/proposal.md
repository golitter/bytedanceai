## Why

多 agent（ClaudeCode / OpenCode）在同一个 task 下协作时，缺少共享上下文的访问机制。每个 agent 只能看到自己的 workspace，无法感知项目架构、任务计划、其他 agent 的输出等 task 级信息。需要一个轻量的、agent 原生集成的方案，让 agent 按需读取共享上下文，而不是由 runtime 强制注入。

## What Changes

- 新增内置 skill `taskctl`：一个 Go 编写的只读 CLI，通过文件路径自举 agent 身份，提供 help / ls / summary / common-memory / sub-memory 五个命令
- 新增 `src/skills/builtin/` 目录：存放内置 skill 的源码、预编译二进制和 instruction 文件
- 新增 `src/skills/provisioner.py`：workspace 创建时将 builtin skills 分发到 agent 的 skill 目录（`.claude/skills/` 或 `.opencode/skills/`）
- 新增 `shared/.agent/memory/` 目录结构：task 级共享存储，按 common/ 和 {agent_name}/ 隔离
- 修改 `WorkspaceManager.create()`：创建 workspace 时调用 provisioner 分发 skills 并初始化 shared 目录

## Capabilities

### New Capabilities
- `taskctl-cli`: Go 编写的只读 CLI 工具，从路径自举 agent 身份，提供 help / ls / summary / common-memory / sub-memory 命令，读取 shared/.agent/ 下的共享数据
- `skill-provisioning`: workspace 创建时将 builtin skills（taskctl 的 exe + skill.md）复制到对应 agent 的 skill 目录，并初始化 shared/.agent/memory/ 目录结构

### Modified Capabilities
- `workspace-management`: WorkspaceManager.create() 需要额外执行 skill 分发和 shared 目录初始化

## Impact

- **新增文件**：`src/skills/builtin/taskctl/`（main.go, go.mod, exe, skill.md）、`src/skills/provisioner.py`
- **修改文件**：`src/workspace/manager.py`（create 方法增加 provisioner 调用）
- **依赖**：Go 工具链（仅编译时，运行时无需）
- **存储**：worktrees 下新增 shared/ 目录和 .agent/ 数据结构
