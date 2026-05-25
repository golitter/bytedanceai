# Diff Viewer — 多文件 Diff 查看与编辑

## 实现了什么

可编辑的多文件 Diff 查看器，支持 tab 切换文件、内联 CodeMirror 编辑、react-diff-view 统一视图渲染。由 `DiffCard` 作为 Smart 组件编排，内部拆分为 Tab、View、Editor 三个子组件。

## 怎么实现的

### DiffFileTabs (`src/components/diff/DiffFileTabs.tsx`)

多文件 tab 切换组件，单文件时不渲染：

```tsx
export function DiffFileTabs({ files, activeIndex, onSelect }: DiffFileTabsProps) {
  if (files.length <= 1) return null

  return (
    <div className="flex overflow-x-auto border-b border-border bg-muted/30">
      {files.map((file, i) => (
        <button
          onClick={() => onSelect(i)}
          className={clsx(
            'flex shrink-0 items-center gap-1.5 border-b-2 px-3 py-1.5 text-xs transition-colors',
            activeIndex === i
              ? 'border-primary text-foreground'
              : 'border-transparent text-muted-foreground hover:text-foreground',
          )}
        >
          <span className="truncate max-w-32">{getFileName(file.newPath)}</span>
          <span className="shrink-0 text-[10px]">
            {file.additions > 0 && <span className="text-green-500">+{file.additions}</span>}
            {file.deletions > 0 && <span className="text-red-500">-{file.deletions}</span>}
          </span>
        </button>
      ))}
    </div>
  )
}
```

### DiffFileView (`src/components/diff/DiffFileView.tsx`)

基于 `react-diff-view` 的统一视图渲染，只接收 `ParsedDiffFile` 并渲染 hunks：

```tsx
export function DiffFileView({ file }: DiffFileViewProps) {
  return (
    <Diff viewType="unified" diffType={file.type} hunks={file.hunks}>
      {(hunks) => hunks.map((hunk) => <Hunk key={hunk.content} hunk={hunk} />)}
    </Diff>
  )
}
```

### DiffFileEditor (`src/components/diff/DiffFileEditor.tsx`)

懒加载编辑器外壳，使用 `React.lazy` + `Suspense` 按需加载 CodeMirror，加载中显示"加载编辑器..."占位：

```tsx
const DiffFileEditorInner = lazy(() => import('./DiffFileEditorInner'))
```

### DiffFileEditorInner (`src/components/diff/DiffFileEditorInner.tsx`)

CodeMirror 编辑器实现，根据文件扩展名自动加载语法高亮（js/ts/py/css/html/json），使用 `oneDark` 主题：

```typescript
function getLangExtension(fileName: string) {
  const ext = fileName.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'js': case 'jsx': case 'ts': case 'tsx':
      return javascript({ jsx: true, typescript: ext?.startsWith('t') })
    case 'py': return python()
    case 'css': case 'scss': return css()
    case 'html': case 'htm': return html()
    case 'json': return json()
    default: return []
  }
}
```

底部操作栏提供"取消"和"保存修改"按钮，保存时调用 `onSave` 回调（DiffCard 中为 `PUT /api/session/:sessionId/files/:path`）。
