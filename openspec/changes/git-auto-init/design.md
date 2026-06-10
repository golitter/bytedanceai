## Context

当前新建对话流程中，用户输入仓库路径后点击"校验"，后端检查路径是否存在且是否为 Git 仓库。若目录存在但未被 Git 初始化，校验直接失败，用户被阻断。现有代码已有 `GitOps.init_repo()` 方法（执行 `git init + add -A + commit + branch -M main`）和 `WorkspaceManager.ensure_git_repo()`，但仅在工作区创建阶段使用，未暴露给前端。

设计文档：`frontend/docs/design/11-git-auto-init.md`
UI Demo：`frontend/docs/payloads/git-auto-init-demo.html`

## Goals / Non-Goals

**Goals:**
- 让用户在前端 UI 内完成非 Git 目录的初始化，无需切换终端
- 通过要求输入目录名最后一段作为确认，防止误操作
- 三端联动：agentend 执行、Go backend 代理、前端交互
- 复用已有 `GitOps.init_repo()`，无新增依赖

**Non-Goals:**
- 不处理 `.gitignore` 配置（使用 git 默认行为）
- 不处理已有的 Git 仓库（幂等拒绝）
- 不在对话创建后补充初始化（仅在校验阶段触发）
- 不支持 Windows 反斜杠路径（agentend 运行在 Unix 环境）

## Decisions

### D1: 确认方式 — 输入目录名最后一段

**选择**: 用户必须输入路径最后一段（如 `/a/b/repo` → 输入 `repo`），精确匹配后才可提交。

**替代方案**:
- 简单 "确认" 按钮：太容易误操作
- 输入完整路径：过长，体验差
- 复用 GitHub "delete repo" 确认模式（输入名称）：业界验证过的安全模式

**理由**: 与 GitHub 删除仓库的确认模式一致，用户熟悉，安全性好。

### D2: 前端交互 — 内联确认框而非新弹窗

**选择**: 在 RepoPathInput 组件内部展示黄色确认框，不使用额外 Dialog。

**理由**: 保持单一焦点流，用户视线不需要转移。与现有组件结构一致（参考 AgentSelectList 的内联添加 Agent 交互）。

### D3: Agentend 直接使用 GitOps 而非 WorkspaceManager

**选择**: `validate.py` 直接 import `GitOps` 并调用 `init_repo()`。

**理由**: 这是一个独立的文件系统操作，不涉及 worktree、skill provision 等工作区生命周期。WorkspaceManager 带有锁和 store 依赖，过度引入。

### D4: 幂等性保护

**选择**: 端点再次检查 `.git` 是否存在，已初始化则返回错误。

**理由**: 防止前端状态过时或并发请求导致重复初始化。

## Risks / Trade-offs

- **[空目录 init 失败]** → `git commit` 在无文件时会失败，API 返回错误，用户可见。这是预期行为——空目录对代码 Agent 无意义。Mitigation: 错误信息明确提示。
- **[路径尾斜杠]** → `path.split('/').filter(Boolean).pop()` 正确处理 `/a/b/` 和 `/a/b`。已覆盖。
- **[并发请求]** → agentend 端点的幂等检查 + FastAPI async 处理，低风险。
