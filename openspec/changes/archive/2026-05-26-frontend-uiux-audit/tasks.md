## 1. CSS 变量对齐

- [x] 1.1 在 `index.css` 中补齐 5 级背景色阶变量（`--bg-canvas`、`--bg-sidebar`、`--bg-card`、`--bg-hover`、`--bg-active`）
- [x] 1.2 在 `index.css` 中补齐 3 级文字色阶变量（`--text-primary`、`--text-secondary`、`--text-tertiary`）
- [x] 1.3 在 `index.css` 中统一定义 Agent 标识色变量（`--agent-claude`、`--agent-opencode`、`--agent-orchestrator`）
- [x] 1.4 在 `index.css` 中补齐语义色变量（`--color-brand`、`--color-success`、`--color-warning`、`--color-error`）
- [x] 1.5 验证所有新增变量与现有 Tailwind/shadcn 语义 token 不冲突，视觉无变化

## 2. 色彩合规扫描与修复

- [x] 2.1 扫描全部 TSX/CSS 文件中的渐变（`linear-gradient`、`radial-gradient`）使用，记录违规项
- [x] 2.2 扫描全部 TSX/CSS 文件中的 `backdrop-blur` 使用，记录违规项
- [x] 2.3 扫描全部 TSX/CSS 文件中的纯白（`#FFFFFF`/`text-white`）和纯黑（`#000000`）文字，记录违规项
- [x] 2.4 扫描组件中 Agent 标识色是否仅用于身份标识（头像 + 色条），记录违规项
- [x] 2.5 修复所有 Critical 级色彩违规

## 3. 字体与排版合规扫描

- [x] 3.1 扫描全部 TSX 文件中的 `font-family` 使用，验证均为 Geist Sans / Geist Mono
- [x] 3.2 扫描字号使用，验证在定义的 6 级字号体系内（11/12/13/14/20px + 对应场景）
- [x] 3.3 修复字号违规项（如有）

## 4. 圆角与边框合规扫描

- [x] 4.1 扫描全部组件的 `border-radius` / `rounded-*` 类名，记录超出规范值（>12px，Badge 除外）
- [x] 4.2 扫描全部组件的边框样式，验证使用 `rgba(255,255,255,0.06)` 或对应 CSS 变量
- [x] 4.3 扫描全部组件的 `box-shadow` 使用，验证仅用于弹出菜单/下拉框
- [x] 4.4 修复圆角与边框违规项

## 5. 动效合规扫描

- [x] 5.1 扫描全部 CSS/TSX 中的 animation/transition，验证时长在 120-300ms 范围内
- [x] 5.2 扫描是否存在禁止动效（bounce、弹簧物理、渐变流光、粒子效果）
- [x] 5.3 验证所有动画仅使用 `transform` 和 `opacity`，不动画布局属性
- [x] 5.4 修复动效违规项（如有）

## 6. Agent 身份系统合规扫描

- [x] 6.1 验证 Agent 头像：圆角方形（8px）、尺寸（32px/24px）、标识色背景
- [x] 6.2 验证 Agent 状态灯：颜色（绿/黄/灰/红）和动效（脉冲/旋转/无/无）匹配规范
- [x] 6.3 验证消息色条：3px 宽、使用 Agent 标识色
- [x] 6.4 重构 MessageBubble 中的 Agent 标识色引用，从 `var(--agent-${type})` 动态拼接改为 `var(--agent-claude)` 等固定引用

## 7. 开发实践合规扫描

- [x] 7.1 扫描组件分层：验证 Dumb 组件不直接访问全局状态或发起请求
- [x] 7.2 扫描状态管理分类：验证 Server State 未进入 Zustand、派生状态未用 useState 存储
- [x] 7.3 扫描列表 key 使用：验证所有 `.map()` 渲染使用稳定唯一 ID，不使用 index
- [x] 7.4 扫描 Hook 合规性：验证纯计算函数未写成 Hook、Hook 体积未超标
- [x] 7.5 扫描 TypeScript 类型使用：检测 `any` 类型、魔法字符串、三元嵌套

## 8. 修复已确认的违规项

- [x] 8.1 修复 MessageBubble 中 BlockRenderer 的 index key 问题（使用 block.id 或生成唯一 ID）
- [x] 8.2 修复所有扫描中发现的 Critical 级违规
- [x] 8.3 修复所有扫描中发现的 Warning 级违规（可选）
- [x] 8.4 验证修复后所有组件视觉无回归（启动 dev server 逐组件检查）
