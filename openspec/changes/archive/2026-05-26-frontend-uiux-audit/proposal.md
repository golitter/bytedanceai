## Why

前端代码库已有 visual-style-guide.md（Dark Utilitarian 风格）和 development-strategy.md（开发实践规范）两份权威文档，但缺少系统性审查来验证现有组件是否完全合规。在 Chat UI 进入 Phase 2 深度开发前，需要一次全面的 UI/UX 合规审计，修复已知偏差并建立可维护的规范执行机制，避免后续开发中合规性持续漂移。

## What Changes

- 审计全部现有组件（chat、im、cards、diff、markdown、ui）的视觉风格合规性
- 修复已发现的违规项：BlockRenderer 使用 index 作为 key、Agent 标识色 CSS 变量一致性检查、硬编码色值残留
- 补齐缺失的 CSS 自定义属性（Agent 标识色 `--agent-claude`/`--agent-opencode`/`--agent-orchestrator`、背景色阶等）
- 审查组件架构合规性：三层组件模型遵循度、状态管理分类正确性、Hook 拆分合理性
- 建立视觉规范合规检查清单，嵌入开发流程

## Capabilities

### New Capabilities
- `style-audit`: 全量视觉风格审计——逐组件检查色彩、字体、圆角、边框、阴影、动效、Agent 身份系统合规性，输出违规项清单与修复方案
- `dev-practice-audit`: 开发实践合规审计——逐组件检查三层架构、状态管理分类、性能模式、TypeScript 类型使用，输出违规项清单与修复方案
- `css-variables-alignment`: CSS 自定义属性与 visual-style-guide.md 对齐——补齐缺失变量、校验色值映射、统一 token 命名

### Modified Capabilities
（无已有 spec 需要修改）

## Impact

- 影响范围：`frontend/src/` 全部组件和样式文件
- 主要文件：`index.css`（CSS 变量补齐）、`components/chat/MessageBubble.tsx`（key 修复）、所有组件的 Tailwind 类名审查
- 不涉及后端或 Agent 端变更
- 不涉及 breaking change
