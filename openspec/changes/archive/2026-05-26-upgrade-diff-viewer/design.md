## Context

前端 Diff 查看器当前使用 `react-diff-view` v3.3.3 的 `viewType="unified"` 模式。文件 Tab 仅显示 `getFileName()`（即 `path.split('/').pop()`），缺少完整路径和变更类型。用户需要左右分栏的 split view 直观对比新旧代码。

现有组件层次：`DiffCard` → `DiffFileTabs` + `DiffFileView` / `DiffFileEditor`。`diff-parser.ts` 已解析出 `oldPath`、`newPath`、`type`（`add/delete/modify/rename/copy`）、`additions`、`deletions` 等完整字段。

## Goals / Non-Goals

**Goals:**
- 默认左右分栏（split view），可切换回 unified
- 显示文件完整相对路径和变更类型标签
- 文件信息栏展示完整元信息

**Non-Goals:**
- 不改后端 API 或 diff 生成逻辑
- 不改 diff-parser.ts
- 不做上下文行数调整
- 不做 side-by-side 编辑（仅编辑模式保留 unified）

## Decisions

### 1. viewType 状态提升到 DiffCard

`viewType` 状态由 `DiffCard` 持有，通过 props 传递给 `DiffFileView`。切换按钮放在 DiffCard Header。

**替代方案**：在 DiffFileView 内部自管状态 → 放弃，因为 Header 上的按钮和 DiffFileView 不在同一组件层级。

### 2. Tab 显示相对路径，tooltip 显示完整路径

Tab 空间有限，使用 `max-w-40 truncate` 显示相对路径，hover 时通过 `title` 显示完整路径。

**替代方案**：只显示文件名 → 放弃，多文件同名时无法区分。

### 3. 变更类型标签使用单字母缩写

| type | 标签 | 颜色 |
|------|------|------|
| add | A | 绿色 |
| delete | D | 红色 |
| modify | M | 蓝色 |
| rename | R | 紫色 |
| copy | C | 灰色 |

### 4. 文件信息栏作为 DiffCard 内容区的一部分

在 Tab 下方、diff 内容上方插入一行，显示完整路径 + 类型标签 + 增删统计。不作为独立组件，直接写在 DiffCard 内。

## Risks / Trade-offs

- **[Split 模式窄屏]** → 移动端 split view 水平空间不足。Mitigation：CSS `min-width` 确保最小可读宽度，overflow-x auto 横向滚动。
- **[编辑模式]** → 编辑模式（DiffFileEditor/CodeMirror）保持使用 `newContent`，不受 split/unified 切换影响。
