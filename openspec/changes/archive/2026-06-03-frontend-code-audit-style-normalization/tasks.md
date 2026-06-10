## 1. 项目元信息集中化

- [x] 1.1 在 `lib/constants.ts` 中新增 `PROJECT_META` 常量对象（`GITHUB_URL`、`NAME`、`DESCRIPTION_ZH`、`DESCRIPTION_EN`）
- [x] 1.2 修改 `components/im/ContactsPage.tsx`：将硬编码的 GitHub URL 替换为 `PROJECT_META.GITHUB_URL`
- [x] 1.3 修改 `components/im/ContactsPage.tsx`：将硬编码的项目描述文案替换为 `PROJECT_META.DESCRIPTION_EN` / `PROJECT_META.DESCRIPTION_ZH`

## 2. text-secondary 类名修正

- [x] 2.1 修改 `components/layout/IconSidebar.tsx`：将 `text-secondary` 替换为 `text-text-secondary`（2 处）
- [x] 2.2 修改 `pages/ImPage.tsx`：将 `text-secondary` 替换为 `text-text-secondary`（1 处）
- [x] 2.3 修改 `pages/AgentProfilePage.tsx`：将 `text-secondary` 替换为 `text-text-secondary`（1 处）

## 3. inline style 统一为 Tailwind 类

- [x] 3.1 修改 `components/layout/AdminPasswordDialog.tsx`：移除 `style={{ color: 'var(--text-secondary)' }}`，改用 `text-text-secondary` 类
- [x] 3.2 修改 `pages/admin/StatisticsPage.tsx`：移除 `var(--text-secondary)` inline style，改用 `text-text-secondary` 类
- [x] 3.3 修改 `pages/admin/ServiceHealthPage.tsx`：移除 `var(--text-secondary)` inline style（4 处），改用 `text-text-secondary` 类
- [x] 3.4 修改 `pages/admin/SessionCleanupPage.tsx`：移除 `var(--text-secondary)` inline style（2 处），改用 `text-text-secondary` 类

## 4. 硬编码色值替换

- [x] 4.1 修改 `components/chat/TerminalPanel.tsx`：将 `bg-[#EF4444]` 替换为 `bg-destructive`，`bg-[#F59E0B]` 替换为 `bg-[var(--color-warning)]`，`bg-[#22C55E]` 替换为 `bg-[var(--color-success)]`

## 5. 视觉规范修正

- [x] 5.1 修改 `components/im/ContactsPage.tsx`：将 favicon 图片的 `rounded-2xl`（16px）替换为 `rounded-xl`（12px），符合面板最大圆角规范
- [x] 5.2 修改 `components/chat/GitGraphPanel.tsx`：移除 tooltip 的 `shadow-lg`，用 `bg-popover` 背景色差 + border 表达层级
- [x] 5.3 修改 `components/layout/IconSidebar.tsx`：移除 tooltip 的 `shadow-lg`，用背景色差 + border 表达层级
- [x] 5.4 修改 `components/chat/HistorySearch.tsx`：移除 `shadow-shadow-lg`（2 处），确认是否为拼写错误并修正

## 6. 验证

- [x] 6.1 运行 `pnpm build` 确认无编译错误
- [x] 6.2 启动开发服务器，检查所有修改页面的视觉效果：次要文字可见性、圆角合规、tooltip 层级感
- [x] 6.3 全局搜索确认无残留 `text-secondary`（排除 shadcn/ui 组件和 `text-text-secondary`）
