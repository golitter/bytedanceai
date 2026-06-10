## 1. parsePath 签名变更

- [x] 1.1 修改 `parsePath` 返回值为 `(taskID, sessionID, sharedDir string, err error)`，增加 taskID 返回值
- [x] 1.2 更新所有 parsePath 调用处以适配新签名

## 2. merge 子命令实现

- [x] 2.1 实现 `runMerge` 函数：通过 parsePath 获取 taskID 和 sessionID，推导分支名
- [x] 2.2 实现自动提交逻辑：`git status --porcelain` 检测未提交改动，有则 `git add -A && git commit`
- [x] 2.3 实现合并流程：`git checkout task/{taskID}` → `git merge agent/{sessionID}/{taskID}`
- [x] 2.4 实现冲突处理：merge 失败时 `git merge --abort`，确保切回 agent 分支
- [x] 2.5 实现成功路径：merge 成功后 `git checkout agent/{sessionID}/{taskID}`，输出结果
- [x] 2.6 在 main.go 的 dispatch 中注册 `merge` 命令

## 3. 文档更新

- [x] 3.1 在 SKILL.md 中添加 `merge` 命令说明，包含用法和示例
