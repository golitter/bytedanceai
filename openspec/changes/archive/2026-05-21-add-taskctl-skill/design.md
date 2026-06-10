## Context

agentend 是一个多 AI Coding Agent Runtime，支持 ClaudeCode 和 OpenCode 两种适配器。每个 task 下，不同 agent 各自拥有独立的 git worktree workspace（`worktrees/{task_id}/{agent_name}/`），彼此隔离。

当前 workspace 创建流程：创建 task branch → 创建 agent worktree → 保存 workspace 记录。agent 启动后只能看到自己的 workspace，无法访问 task 级别的共享信息。

约束：
- agentend 是 Python 项目，CLI 工具用 Go 编写（预编译提交）
- ClaudeCode 自动扫描 `.claude/skills/` 下的文件作为 instruction
- OpenCode 自动扫描 `.opencode/skills/` 下的文件作为 instruction
- agent 通过自身的 bash 工具调用 CLI

## Goals / Non-Goals

**Goals:**
- 让 agent 能按需读取 task 级共享上下文（架构、记忆、其他 agent 输出）
- workspace 创建时自动注入 taskctl skill，agent 无需手动配置
- 通过文件路径隔离实现 sub-memory 的访问控制
- 保持 taskctl 实现极简：只读、自举、无外部依赖

**Non-Goals:**
- 不实现写入能力（现阶段全部只读）
- 不实现 Skill Runtime / Prompt 注入等重型机制
- 不实现跨 task 的共享
- 不实现 RAG / 语义搜索

## Decisions

### 1. CLI Gateway 而非 Runtime 注入

**选择**：agent 放一个 skill.md + exe，agent 按需调用 CLI 读取共享数据
**不选**：agentend 在 adapter 层把 skill 内容注入 system_prompt

理由：Runtime 注入会塞满 context，不管 agent 需不需要都在里面。CLI Gateway 让 agent 按需查询，context 更干净。且未来新增数据类型只需改 CLI，不用改 agentend Python 代码。

### 2. Go CLI 预编译提交

**选择**：Go 编写，预编译为 `exe` 提交到仓库
**不选**：Python 脚本 / Go 运行时编译

理由：Go 编译为单二进制，无运行时依赖。预编译提交避免 agentend 运行环境需要 Go 工具链。

### 3. 路径自举身份识别

**选择**：exe 从自身路径解析 task_id 和 agent_name
**不选**：环境变量注入 / 配置文件

理由：零配置。exe 放在 `.{agent}/skills/taskctl/exe`，路径本身就包含了身份信息。不依赖 adapter 启动参数。

```
exe 路径: worktrees/task-123/claude_code/.claude/skills/taskctl/exe
解析:     task_id = "task-123", agent_name = "claude_code"
shared:   worktrees/task-123/shared/.agent/
```

### 4. sub-memory 按 agent_name 隔离

**选择**：`shared/.agent/memory/{agent_name}/` 按 exe 识别的身份隔离
**不选**：按 session_id 隔离

理由：agent 不感知 session_id，但必然知道自己是什么 agent。路径隔离天然安全 — exe 硬编码只读自己名字的目录。

### 5. shared 目录放在 task 级别，不是 git worktree

**选择**：`worktrees/{task_id}/shared/` 作为普通目录
**不选**：放在 repo root / 放在 agent worktree 内

理由：task 级共享，所有 agent 都能通过相对路径 `../../shared/` 访问。不需要 git worktree，只是文件系统目录。

## Risks / Trade-offs

- **[Risk] CLI 路径解析依赖目录结构约定** → Mitigation：路径解析时做校验，找不到 worktrees 层级时报错退出
- **[Trade-off] 只读意味着写入需要其他途径** → 现阶段可接受，未来可通过 API 或 taskctl write 命令扩展
- **[Trade-off] 预编译 exe 需要为多平台编译** → 现阶段只编译 darwin/amd64 和 linux/amd64，按需扩展
