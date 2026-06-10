## Why

前一轮审计（2026-05-30）已修正 rounded-full 滥用、diff 字体、Lucide strokeWidth 等表层问题。但深度审计发现仍存在系统性偏差：组件 border-radius 值不匹配规范（按钮用 10px 而非 6px、输入框 10px 而非 8px），大量 opacity 调色（`text-foreground/50`）绕过了定义色板，`transition-colors` 违规动画属性，chat store 660 行严重超限并将 server state 塞入 Zustand，API 层缺少 res.ok 校验，全应用无 Error Boundary。需要第二轮审计，将前端代码在视觉规范和实践规范上对齐到位。

## What Changes

**视觉规范修正：**
- 修正组件 border-radius：按钮改为 6px、输入框改为 8px、面板改为 12px、Badge 改为 9999px 胶囊、用户头像改为 50% 圆形
- 消除 opacity 调色模式：`text-foreground/N` 改用语义 token（`text-secondary`、`text-tertiary`）
- 品牌色收归合规场景：streaming 光标、状态指示等非品牌场景改用语义色
- `transition-colors` 改为 `transition-[transform,opacity]`，仅动画合规属性
- 消除非规范阴影值，弹出菜单统一为 `0 4px 24px rgba(0,0,0,0.4)`

**代码质量修正：**
- chat store 拆分：将 block 转换逻辑提取为纯函数到 `lib/`，server state 迁移至 TanStack Query
- API 层补充 `res.ok` 校验，消除静默失败
- 抽取魔法字符串常量：`DEFAULT_AGENT_TYPE`、`STREAMING_STATUSES`、`AGENT_ROLE` 等
- 补充 Error Boundary 包裹 ChatArea、AdminContent、ConversationList 等关键模块
- ImPage 中嵌套三元改为 guard clause 或路由驱动
- 修正 `isAuthenticated` 派生状态反模式

## Capabilities

### New Capabilities

- `visual-radius-fix`: 统一全部组件 border-radius 到 visual-style-guide.md 规范值（按钮 6px、输入框 8px、卡片 10px、面板 12px、Badge 9999px、用户头像 50%、Agent 头像 8px 方形）
- `color-token-enforcement`: 消除 opacity 调色和硬编码色值，统一使用 CSS 变量语义 token；品牌色收归合规场景
- `transition-compliance`: 全量替换 `transition-colors`/`transition-all` 为合规属性动画（仅 transform/opacity）
- `error-boundary`: 为 ChatArea、AdminContent、ConversationList 等关键模块补充 Error Boundary 兜底
- `api-resilience`: API 层补充 res.ok 校验，消除静默失败，统一错误处理

### Modified Capabilities

_(无已有 spec 需要修改行为层需求，本次为纯实现层修正)_

## Impact

- 影响范围：`frontend/src/` 下约 25-30 个组件文件
- 无 API 变更，无 breaking change
- 纯视觉、样式和代码结构调整，不影响功能逻辑
- 依赖 `visual-style-guide.md` 色彩体系和 `development-strategy.md` 状态管理规范
