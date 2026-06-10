## Why

当前新建对话时，用户输入的仓库路径如果未被 Git 初始化，校验直接失败并阻止创建对话。用户必须手动切到终端执行 `git init`，体验割裂。对于新项目或从压缩包解压的代码目录，这个限制尤其常见。

## What Changes

- **Agentend 新增 `POST /v1/init-git-repo` 端点**：接收路径，校验其存在且非 Git 仓库后执行 `git init && git add -A && git commit -m "init" && git branch -M main`
- **Go Backend 新增 `/init-git-repo` 代理路由**：转发前端请求到 agentend
- **前端 RepoPathInput 组件增强**：当校验返回"不是 git 仓库"时，展示黄色确认框，要求用户输入路径最后一段目录名确认后自动初始化，成功后直接进入校验通过状态
- **前端新增 `initGitRepo` API 函数及对应 UI 文本常量**

## Capabilities

### New Capabilities
- `git-auto-init`: 非 Git 目录自动初始化能力——当校验发现目录存在但未 Git 初始化时，提供确认后自动 `git init` 的完整流程（三端联动）

### Modified Capabilities
- `repopath-validation`: 前端 RepoPathInput 组件行为变更——"路径不是 git 仓库"不再直接报错，而是进入确认初始化流程

## Impact

- **Agentend**: `validate.py` 新增端点，直接使用已有 `GitOps.init_repo()`
- **Go Backend**: `agentend_client/client.go` 新增代理方法，`task_controller.go` 新增路由
- **Frontend**: `RepoPathInput.tsx` 组件内部状态扩展，`api.ts` 新增函数，`ui-text.ts` 新增常量
- **无数据库变更、无依赖新增、无破坏性变更**
