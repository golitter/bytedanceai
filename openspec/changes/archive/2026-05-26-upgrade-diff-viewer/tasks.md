## 1. DiffFileView 支持 split/unified 切换

- [x] 1.1 DiffFileView 接收 `viewType: 'split' | 'unified'` prop，替换硬编码 `viewType="unified"`
- [x] 1.2 DiffCard 新增 `viewType` state（默认 `'split'`），通过 props 传给 DiffFileView

## 2. Header 切换按钮

- [x] 2.1 DiffCard Header 区域新增 Split / Unified 切换按钮，active 状态高亮
- [x] 2.2 点击按钮更新 `viewType` state，触发 diff 重新渲染

## 3. DiffFileTabs 路径 + 变更类型

- [x] 3.1 Tab 显示 `newPath` 相对路径（`max-w-40 truncate`），hover `title` 显示完整路径
- [x] 3.2 新增变更类型标签映射函数：add→A/绿, delete→D/红, modify→M/蓝, rename→R/紫, copy→C/灰
- [x] 3.3 每个 Tab 渲染类型标签 + 路径 + 增删统计

## 4. 文件信息栏

- [x] 4.1 DiffCard 内容区（Tab 下方、diff 上方）新增文件信息栏：完整路径 + 变更类型标签 + 增删统计
- [x] 4.2 信息栏样式：小字号、浅色背景、与 diff 内容区分明显
