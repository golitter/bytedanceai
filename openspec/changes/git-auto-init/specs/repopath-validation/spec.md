## MODIFIED Requirements

### Requirement: 发消息前校验 repoPath
前端 SHALL 在每次发送消息前先调用 repoPath 校验接口，校验通过后才发送消息，校验失败则展示错误提示。当校验返回"不是 git 仓库"错误时，SHALL 展示 Git 初始化确认流程（见 `git-auto-init` spec），而非直接报错。

#### Scenario: 校验通过后发送
- **WHEN** 用户点击发送，repoPath 校验返回 valid=true
- **THEN** 前端正常发送消息到 `/api/tasks/:taskId/run`

#### Scenario: 校验失败 — 路径不存在
- **WHEN** 用户点击发送，repoPath 校验返回错误"路径不存在"
- **THEN** 前端展示红色错误信息，不发送消息

#### Scenario: 校验失败 — 不是 git 仓库
- **WHEN** 用户校验路径，返回错误"路径不是 git 仓库"
- **THEN** 前端展示黄色确认框，引导用户确认后初始化 Git，不直接报错阻断

#### Scenario: task 无 repoPath
- **WHEN** task 未设置 repoPath
- **THEN** 跳过校验，直接发送消息
