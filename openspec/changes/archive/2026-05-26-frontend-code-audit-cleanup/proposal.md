## Why

前端代码经多轮审计后仍存在显著的冗余和硬编码问题：`API_BASE` 在 4 个文件中重复声明；Agent 颜色映射在多个组件中以不同名称重复定义；文件名提取逻辑 `path.split('/').pop()` 散布 4 处；大量魔法数字（超时、尺寸、阈值）直接写在业务逻辑中；DiffCard.tsx 达 337 行远超 200 行红线。这些问题违反 development-strategy.md 的核心原则，增加维护成本和出错风险。

## What Changes

- **统一 API 基础路径**：消除 4 处 `API_BASE` 重复声明，收敛到 `lib/constants.ts`
- **收敛 Agent 颜色映射**：将 AgentAvatar、MessageBubble、DiffCard、DiffFileTabs 中分散的颜色定义统一到 `lib/constants.ts`
- **提取文件名工具函数**：将 4 处 `path.split('/').pop()` 收敛到 `lib/utils.ts`
- **消除魔法数字**：将 MessageInput、MessageList、AgentAvatar 中的超时/尺寸/阈值常量提取到命名常量
- **拆分 DiffCard.tsx**：将 337 行的 DiffCard 按职责拆分为子组件，每块 <200 行
- **统一 CSS diff 主题**：将 index.css 中散落的硬编码 diff 色（success/error 背景色）收归 CSS 变量

## Capabilities

### New Capabilities
- `constants-consolidation`: 将散落在组件中的重复常量（API_BASE、Agent 颜色、魔法数字）收敛到统一的常量文件
- `component-split-diffcard`: 将 DiffCard.tsx 按职责拆分为独立子组件（DiffStats、DiffActions、DiffContent）
- `css-diff-theme-vars`: 将 index.css 中 diff 相关的硬编码色值收归 CSS 变量，与 visual-style-guide.md 对齐

### Modified Capabilities

（无已有 spec 需要修改）

## Impact

- **代码文件**：`lib/constants.ts`、`lib/utils.ts`、`index.css`、AgentAvatar、MessageBubble、DiffCard、DiffFileTabs、AttachmentCard、ImageCard、MessageInput、MessageList
- **无 API 变更**：纯前端重构，不影响后端接口
- **无依赖变更**：不引入新包
- **风险低**：所有改动均为内部重构，对外行为不变
