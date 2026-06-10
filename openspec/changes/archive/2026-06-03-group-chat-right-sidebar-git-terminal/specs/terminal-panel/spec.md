## ADDED Requirements

### Requirement: Terminal window layout
系统 SHALL 提供 `TerminalPanel` 组件，包含标题栏（三色圆点 + 路径）、输出区域和输入行，整体风格为深色终端。

#### Scenario: Terminal renders with connected status
- **WHEN** TerminalPanel 渲染
- **THEN** 显示标题栏（红/黄/绿三色圆点 + 工作区路径），输出区域显示欢迎消息，输入行显示 `(branch-name) $` 提示符，标题栏显示绿色 "Connected" 状态和脉冲动画圆点

### Requirement: Terminal prompt shows current branch
终端提示符 SHALL 显示当前分支名，格式为 `(<branch-name>) $`，分支名为绿色。

#### Scenario: Prompt updates on branch switch
- **WHEN** 当前分支从 main 切换到 feat/group-chat-agent
- **THEN** 提示符从 `(main) $` 更新为 `(feat/group-chat-agent) $`

### Requirement: Git commands with branch-aware output
Terminal SHALL 支持 `git status`、`git log`、`git branch` 命令，输出内容跟随当前分支变化。

#### Scenario: git status on feature branch
- **WHEN** 当前分支为 feat/group-chat-agent 且执行 `git status`
- **THEN** 输出 "On branch feat/group-chat-agent" + modified files 列表

#### Scenario: git status on main branch
- **WHEN** 当前分支为 main 且执行 `git status`
- **THEN** 输出 "On branch main" + "nothing to commit, working tree clean"

#### Scenario: git branch lists all branches
- **WHEN** 执行 `git branch`
- **THEN** 列出所有分支，当前分支前标 `*`（绿色），其他分支前标空格

#### Scenario: git log shows reachable commits only
- **WHEN** 当前分支为 main 且执行 `git log`
- **THEN** 只显示 main 分支可达的 commits（3 条）
- **WHEN** 当前分支为 feat/group-chat-agent 且执行 `git log`
- **THEN** 显示该分支可达的所有 commits（8 条，含 main 分支的共享 commits）

### Requirement: Git checkout command
Terminal SHALL 支持 `git checkout <branch>` 和 `git switch <branch>` 命令，切换分支时 SHALL 触发 `onBranchChange` 回调。

#### Scenario: Checkout to existing branch
- **WHEN** 用户输入 `git checkout main`（当前在 feat/group-chat-agent）
- **THEN** 输出 "Switched to branch 'main'"，触发 onBranchChange("main")，提示符更新为 `(main) $`

#### Scenario: Checkout to current branch
- **WHEN** 用户输入 `git checkout main`（当前已在 main）
- **THEN** 输出 "Already on 'main'"，不触发 onBranchChange

#### Scenario: Checkout to non-existent branch
- **WHEN** 用户输入 `git checkout non-existent`
- **THEN** 输出错误信息 "error: pathspec 'non-existent' did not match any branch known to git"

### Requirement: Common terminal commands
Terminal SHALL 支持以下命令：`help`、`clear`、`pwd`、`ls`、`whoami`、`npm run build`、`npm test`、`echo <text>`、`cat <file>`。

#### Scenario: help command
- **WHEN** 用户输入 `help`
- **THEN** 输出所有可用命令列表

#### Scenario: clear command
- **WHEN** 用户输入 `clear`
- **THEN** 终端输出区域清空

#### Scenario: Unknown command
- **WHEN** 用户输入未识别的命令
- **THEN** 输出 "command not found: <command>"

### Requirement: Terminal section collapsible
Terminal 区块 SHALL 支持折叠/展开，复用 RightSidebar 的 section-header + chevron 模式，Terminal 区块 SHALL 占满右侧栏剩余高度（`flex: 1`）。

#### Scenario: Toggle Terminal section
- **WHEN** 用户点击 Terminal 标题栏
- **THEN** 区块折叠/展开，chevron 图标旋转，折叠状态持久化到 localStorage
