## Why

前端代码经过多轮迭代后积累了若干风格不一致和硬编码问题：部分组件使用 `text-secondary`（映射到 shadcn 的 secondary 语义色，暗色背景下近乎不可见）而非项目自定义的 `text-text-secondary`；项目元信息（GitHub 链接、项目描述）散落在组件代码中；少量组件违反视觉风格规范（过大的圆角、不必要的阴影、硬编码色值）。需要一次系统性审计修正，确保前端代码与 `visual-style-guide.md` 及 `development-strategy.md` 保持一致。

## What Changes

- 修正所有 `text-secondary` 用法为 `text-text-secondary`，确保暗色背景下次要文字可见
- 将 ContactsPage 中硬编码的 GitHub 链接、项目描述等元信息统一提取到 `lib/constants.ts`（或新建 `lib/settings.ts`）
- 将 TerminalPanel 中硬编码的语义色值（`#EF4444`、`#F59E0B`、`#22C55E`）替换为 CSS 变量或 Tailwind 语义类
- 修正 ContactsPage 中 favicon 图片的 `rounded-2xl`（16px）为合规圆角值（最大 12px = `rounded-xl`）
- 审查并修正不必要的 `shadow-lg`/`shadow-md` 用法（shadcn/ui 内置组件除外），深度应靠背景色阶表达
- 统一 `var(--text-secondary)` inline style 用法为 Tailwind 类 `text-text-secondary`

## Capabilities

### New Capabilities
- `project-metadata-centralization`: 将散落在组件中的项目元信息（GitHub URL、项目描述、平台名称）集中到配置文件，建立单一来源
- `tailwind-class-normalization`: 修正 Tailwind CSS 4 中的 `text-secondary` → `text-text-secondary` 类名映射错误，统一 inline style 与 Tailwind 类用法

### Modified Capabilities
- `style-audit`: 在现有 style-audit 基础上扩展，覆盖圆角合规性、阴影合规性、硬编码色值审计维度

## Impact

- **前端源码**：涉及 `components/im/ContactsPage.tsx`、`components/chat/TerminalPanel.tsx`、`components/layout/IconSidebar.tsx`、`pages/ImPage.tsx`、`pages/AgentProfilePage.tsx`、`pages/admin/*.tsx` 等多个文件
- **配置层**：`lib/constants.ts` 新增项目元信息常量，或新建 `lib/settings.ts`
- **无 API 变更**：纯前端样式和代码组织层面的修改，不涉及后端接口
- **无破坏性变更**：所有修改为等价替换，视觉表现不变或改善
