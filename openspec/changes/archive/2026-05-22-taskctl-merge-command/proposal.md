## Why

Agent（Claude Code / OpenCode）完成编码后需要将变更合并到集成分支，但 agent 不知道分支层级结构（main → task → agent），默认尝试合并到 main 导致失败（权限不足），需要反复询问用户。即使合并成功，agent 也不会切回自己的分支，导致后续操作在错误的分支上执行。

## What Changes

- 为 taskctl 新增 `merge` 子命令，一键完成：自动提交 → 合并到 task 分支 → 切回 agent 分支
- `parsePath` 返回值增加 `taskID`，供 merge 命令推导分支名
- 更新 SKILL.md 文档，告知 agent 使用 `./taskctl merge` 完成合并
- 不提供 `--main` 参数，目标始终是 `task/{taskID}`，防止 agent 直接操作 main

## Capabilities

### New Capabilities

- `taskctl-merge`: taskctl merge 子命令，支持自动提交、合并到 task 分支、冲突处理、分支回退

### Modified Capabilities

- `taskctl-cli`: `parsePath` 需要返回 `taskID` 以支持 merge 命令推导分支名

## Impact

- `agentend/src/skills/builtin/taskctl/main.go` — 新增 merge 子命令，修改 parsePath 签名
- `agentend/src/skills/builtin/taskctl/SKILL.md` — 添加 merge 命令文档
