## 1. 常量收敛（constants-consolidation）

- [x] 1.1 在 `lib/constants.ts` 中新增 `export const API_BASE = '/api'` 和 `export const AGENT_COLORS: Record<AgentType, string>`
- [x] 1.2 `lib/api.ts`：删除本地 `const API_BASE`，改为 `import { API_BASE } from '@/lib/constants'`
- [x] 1.3 `components/cards/DiffCard.tsx`：删除本地 `const API_BASE`，改为 import
- [x] 1.4 `components/cards/AttachmentCard.tsx`：删除本地 `const API_BASE`，改为 import；将 `path.split('/').pop()` 替换为 `getFileName(path)`（从 `@/lib/utils` import）
- [x] 1.5 `components/cards/ImageCard.tsx`：同 1.4，删除 `const API_BASE` 并 import，替换文件名提取逻辑
- [x] 1.6 `components/chat/AgentAvatar.tsx`：删除 `AGENT_COLORS` 和 `AGENT_SHADOW_COLORS`，改为 `import { AGENT_COLORS } from '@/lib/constants'`；将 `AGENT_SHADOW_COLORS[agentType]` 替换为 `AGENT_COLORS[agentType]`
- [x] 1.7 `components/chat/MessageBubble.tsx`：删除 `AGENT_STRIP_COLORS`，改为 `import { AGENT_COLORS } from '@/lib/constants'`；更新引用

## 2. 工具函数统一

- [x] 2.1 在 `lib/utils.ts` 中新增 `export function getFileName(path: string): string { return path.split('/').pop() || path }`
- [x] 2.2 `lib/diff-parser.ts`：删除本地 `getFileName` 函数定义，改为 `import { getFileName } from '@/lib/utils'`
- [x] 2.3 验证 `lib/block-reducer.ts` 中的 `path.split('/').pop()` 用法，如有则替换为 `getFileName`

## 3. 魔法数字常量化

- [x] 3.1 `components/chat/MessageInput.tsx`：提取 `MAX_INPUT_HEIGHT = 200`、`HINT_DISPLAY_DURATION = 3000`、`MIN_INPUT_HEIGHT = 48`，替换所有魔法数字
- [x] 3.2 `components/chat/MessageList.tsx`：提取 `SCROLL_BOTTOM_THRESHOLD = 60`（`VIRTUALIZE_THRESHOLD` 已有常量），替换魔法数字
- [x] 3.3 `components/chat/AgentAvatar.tsx`：提取 `STATUS_READY_DURATION = '2s'`、`STATUS_RUNNING_DURATION = '1.5s'`，替换魔法字符串

## 4. DiffCard 拆分（component-split-diffcard）

- [x] 4.1 创建 `components/diff/DiffHeader.tsx`，提取顶栏区块（文件统计 + 视图切换 + 操作按钮 + 状态 badge），DiffCard 通过 props 传递数据
- [x] 4.2 创建 `components/diff/DiffFileInfo.tsx`，提取文件信息条（路径 + 变更类型 badge + 增删统计），复用 DiffFileTabs 的 `ChangeTypeBadge`
- [x] 4.3 在 `components/diff/DiffFileTabs.tsx` 中 export `ChangeTypeBadge` 组件
- [x] 4.4 重构 `DiffCard.tsx`，引入 DiffHeader 和 DiffFileInfo，删除已提取的 JSX 块，确保行数 < 200
- [x] 4.5 验证 DiffCard 功能不变：加载快照、视图切换、接受/拒绝操作、编辑文件

## 5. CSS Diff 主题变量化（css-diff-theme-vars）

- [x] 5.1 在 `index.css` 的 `.dark` 块中新增 `--diff-insert-color`、`--diff-insert-bg`、`--diff-delete-color`、`--diff-delete-bg`、`--diff-insert-bg-strong`、`--diff-delete-bg-strong` 六个变量
- [x] 5.2 在 `@theme inline` 块中映射 `--color-diff-*` 形式的 Tailwind 变量
- [x] 5.3 替换 `.diff-gutter-insert`、`.diff-gutter-delete`、`.diff-code-insert`、`.diff-code-delete` 中的硬编码色值为 `var(--diff-*)`
- [x] 5.4 验证暗色模式下 diff 卡片的视觉表现不变

## 6. 验证

- [x] 6.1 运行 `pnpm build` 确保无编译错误
- [ ] 6.2 启动前端 `make run-frontend`，手动验证聊天界面、Diff 卡片功能正常
