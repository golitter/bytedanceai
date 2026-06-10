## Why

当前前端仅支持 Dark 主题，所有自定义语义 CSS 变量（`--bg-canvas`、`--text-primary`、`--agent-*`、`--diff-*`、`--prose-*` 等）仅在 `.dark` 块中定义，`:root`（Light 模式）只有 shadcn 基础 token 且均为无彩色 oklch 灰阶。用户无法选择亮色界面，在不同光线环境下的可读性和舒适度受限。Light 主题的配色建议值已在 `frontend/docs/reference/color-palette.md` L-八 中就绪，并有完整 demo HTML（`frontend/payloads/light-theme-demo.html`）可供视觉验证。

## What Changes

- 在 `:root` 中补全 Light 主题的所有自定义语义变量（背景色阶、文字色阶、品牌色、Agent 标识色、Diff 对比色、Prose 增强色等），覆盖 `--primary` 为 Indigo `#6366F1`
- 新增 `useTheme` hook，管理 `<html>` 元素的 `dark` class 切换，偏好持久化到 `localStorage`，首次访问跟随 `prefers-color-scheme`
- 在 `IconSidebar` 左下角新增 Settings 弹出面板（Popover），包含 Dark / Light 主题切换开关
- 调整少数硬编码暗色内联样式（`boxShadow` rgba 黑色等），使其在 Light 模式下适配

## Capabilities

### New Capabilities

- `theme-switching`: 主题切换系统——包含 Light 主题 CSS 变量补全、`useTheme` hook、Settings 弹出面板 UI、localStorage 持久化 + `prefers-color-scheme` 系统偏好跟随

### Modified Capabilities

<!-- 无现有 spec 需要修改 -->

## Impact

- **CSS**：`src/index.css` 的 `:root` 块新增约 40 个自定义变量
- **新增文件**：`src/hooks/use-theme.ts`（hook）、`src/components/layout/SettingsPanel.tsx`（Settings 弹出面板）
- **修改文件**：`src/components/layout/IconSidebar.tsx`（集成 Settings 入口）、`src/main.tsx`（初始化主题）、少数组件的内联暗色样式
- **无后端影响**：主题切换纯前端，不涉及 API 或契约层变更
