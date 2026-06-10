## Why

前端代码存在两套并行的 CSS 变量体系（shadcn/ui 语义 token 与自定义 app token），部分组件使用不同约定引用同一颜色；同时 `index.css` 中存在重复定义、未使用变量，以及与 `visual-style-guide.md` 规范不一致的样式。代码层面，旧的 impl 文档已与实际实现严重脱节，且存在未清理的冗余文件和过期引用。现在进行规范化可消除视觉不一致，降低维护负担，并为后续迭代建立可靠基线。

## What Changes

- 统一 CSS 变量体系：将自定义 app token（`--bg-canvas`、`--text-primary` 等）迁移到 shadcn/ui 语义 token（`--background`、`--foreground` 等），消除双轨制
- 清理 `index.css` 中未使用的 CSS 变量和重复定义
- 逐组件排查并修复与 `visual-style-guide.md` 不一致的样式（边框、圆角、间距、动画、颜色使用）
- 移除代码中未使用的 imports、变量、hooks 和组件
- 更新 `frontend/AGENTS.md` 使其反映当前实际项目结构
- 清理 `frontend/docs/impl/` 下与实际代码不符的过期文档内容

## Capabilities

### New Capabilities

- `style-normalization`: 统一 CSS 变量体系、修复组件样式与 visual-style-guide.md 的偏差，确保所有组件遵循 Dark Utilitarian 设计规范
- `code-dedup`: 移除前端代码中的冗余定义（未使用的 CSS 变量、重复 imports、死代码），保持代码库精简

### Modified Capabilities

（无既有 spec 需要修改）

## Impact

- `frontend/src/index.css`：变量体系合并与清理
- `frontend/src/components/`：所有组件的 className / style 审查与修正
- `frontend/src/hooks/`：未使用 hooks 的移除
- `frontend/src/stores/`：冗余状态清理
- `frontend/AGENTS.md`：文档更新
- `frontend/docs/impl/`：过期内容清理
