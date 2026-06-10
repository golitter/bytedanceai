## Why

当前 Diff 查看器使用 unified（上下）模式渲染，信息密度低——文件只显示文件名而非完整路径，缺少变更类型标签（A/D/M），难以快速定位和理解变更。用户需要左右分栏的 split view 来直观对比新旧代码，同时需要更完整的文件元信息。

## What Changes

- Diff 渲染模式从固定 `unified` 改为默认 `split`（左右分栏），并支持 unified/split 切换按钮
- 文件 Tab 显示相对路径（而非仅 `getFileName()`）+ 变更类型标签（M/A/D/R）
- 内容区顶部新增文件信息栏：完整路径 + 变更类型 + 增删行统计
- Header 区域新增 unified/split 切换按钮，状态持久到当前会话

## Capabilities

### New Capabilities
- `split-view-toggle`: 左右分栏 diff 渲染模式 + unified/split 切换控制
- `diff-file-info`: 文件完整路径、变更类型标签（M/A/D/R）、增删统计的展示

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **前端组件**：`DiffFileView.tsx`、`DiffFileTabs.tsx`、`DiffCard.tsx` 共 3 个文件
- **依赖**：`react-diff-view` v3.3.3 已原生支持 `viewType="split"`，无需升级
- **不改后端**、不改 `diff-parser.ts`（数据结构已有 `oldPath`/`newPath`/`type`）
