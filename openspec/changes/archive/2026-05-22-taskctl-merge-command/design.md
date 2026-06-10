## Context

taskctl 是 agent worktree 内的 CLI 工具，通过 `parsePath()` 从自身可执行文件路径推导 session 信息。当前只提供上下文读写命令（ls、summary、common-memory、sub-memory、write-sub-memory），不涉及 git 操作。

分支层级为 `main → task/{task_id} → agent/{session_id}/{task_id}`。agent 完成编码后需要合并到 task 分支，但 agent 不知道这个分支结构，默认尝试合并到 main（无权限），且合并后不会切回自己的分支。

## Goals / Non-Goals

**Goals:**
- agent 通过一条命令完成合并，无需理解分支结构
- 合并后自动回到 agent 自己的分支
- 冲突时安全回退，不留 worktree 在错误分支上

**Non-Goals:**
- 不支持合并到 main（这是 `merge_task_to_main` API 的职责）
- 不处理冲突自动解决（冲突时 abort 并报错）
- 不走 HTTP API（taskctl 是纯本地工具，直接调用 git）

## Decisions

### 1. 直接在 worktree 内执行 git 命令，不走 API

**选择**：taskctl 直接在 worktree 目录内执行 `git checkout / git merge`。

**备选**：调用 `POST /v1/workspace/{id}/merge` API。

**理由**：taskctl 是纯本地 CLI，不依赖网络。API 需要 workspace_id 和 endpoint 配置，增加了复杂度。worktree 内可以直接 checkout 不同分支执行 merge。

### 2. merge 前自动 commit

**选择**：检测到未提交改动时自动 `git add -A && git commit`。

**备选**：报错让 agent 先手动 commit。

**理由**：agent 容易忘记 commit，报错会导致它需要多轮对话才能完成 merge。自动 commit 减少 agent 卡点。

### 3. parsePath 增加 taskID 返回值

**选择**：修改 `parsePath` 返回 `(taskID, sessionID, sharedDir string, err error)`。

**备选**：新增独立函数 `parseTaskID()`。

**理由**：taskID 在路径解析时已经可以获取，不需要额外函数。现有 parsePath 内部已计算 taskID 只是没返回。改动最小。

### 4. 分支名从路径推导，不运行 git 命令

**选择**：agent 分支名通过 `fmt.Sprintf("agent/%s/%s", sessionID, taskID)` 构建。

**备选**：运行 `git rev-parse --abbrev-ref HEAD` 获取当前分支名。

**理由**：路径推导更可靠（不依赖 worktree 内 git 状态），且与 Workspace 模型中 `_generate_branch_name` 的逻辑一致。merge 前仍会用 `git branch --show-current` 做校验。

## Risks / Trade-offs

- **[并发 merge 竞态]** → 同一 task 下多个 agent 同时 merge 到 task 分支可能冲突。后续可通过 per-task 文件锁解决，当前依赖 git 自身的 merge 检测。
- **[自动 commit 信息不精确]** → 自动 commit 使用固定 message "auto: merge前自动提交"，不含具体改动描述。这是合理的 trade-off，因为 merge 命令的 commit 只是中间步骤。
- **[parsePath 签名变更]** → 修改 parsePath 返回值会影响现有调用方。需同步更新所有调用处。
