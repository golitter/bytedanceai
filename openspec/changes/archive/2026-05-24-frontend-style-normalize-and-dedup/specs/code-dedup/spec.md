## ADDED Requirements

### Requirement: 死代码文件移除
系统 SHALL 移除所有未被引用的源代码文件。

#### Scenario: App.tsx 移除
- **WHEN** 检查 `frontend/src/App.tsx`
- **THEN** 该文件已被移除（该文件返回 `null`，`main.tsx` 直接路由到 `ImPage`，不经过 App）

#### Scenario: use-sessions.ts 移除
- **WHEN** 检查 `frontend/src/hooks/use-sessions.ts`
- **THEN** 该文件已被移除（导出的 `useSessions()` 未被任何文件引用）

### Requirement: 未使用 CSS 变量移除
`index.css` 中 SHALL 不存在被定义但未被任何组件引用的 CSS 变量。

#### Scenario: 未使用变量清理
- **WHEN** 检查 `index.css` 的 `.dark` 块
- **THEN** `--bg-active`、`--color-success`、`--color-warning` 如无组件引用则被移除（`--bg-active` 保留需先确认是否需要用于交互态，若确认不需要则移除）

### Requirement: 无冗余 import
所有组件文件 SHALL 不包含未使用的 import 语句。

#### Scenario: 未使用 import 清理
- **WHEN** 运行 TypeScript 编译检查
- **THEN** 无 `declared but never used` 或 `unused` 相关的 import 警告

### Requirement: AGENTS.md 反映当前项目结构
`frontend/AGENTS.md` SHALL 准确描述当前前端项目的目录结构、组件清单、hooks 清单、stores 清单和依赖版本。

#### Scenario: 结构描述准确
- **WHEN** 对照 `frontend/AGENTS.md` 与实际文件系统
- **THEN** 文档中列出的目录、文件、组件均实际存在，且无遗漏

#### Scenario: 移除过时描述
- **WHEN** 检查 `frontend/AGENTS.md`
- **THEN** 不存在对已删除文件（如 `stores/app.ts`）的引用，不存在对已移除组件的描述
