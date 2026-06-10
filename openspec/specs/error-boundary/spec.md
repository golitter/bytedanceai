## ADDED Requirements

### Requirement: 模块级 Error Boundary 包裹关键组件
ImPage 中的 ChatArea、AdminContent、ConversationList 三个关键模块 SHALL 各自被独立 Error Boundary 包裹。当任一模块抛出异常时，仅该模块显示降级 UI，不影响其他模块正常运行。

#### Scenario: ChatArea 崩溃
- **WHEN** ChatArea 组件内部抛出未捕获异常（如 SSE 连接错误、渲染异常）
- **THEN** ChatArea 区域显示错误提示 + 重试按钮，ConversationList 和侧边栏保持正常

#### Scenario: AdminContent 崩溃
- **WHEN** AdminContent 组件内部抛出未捕获异常（如认证失败、网络错误）
- **THEN** Admin 区域显示错误提示 + 重试按钮，其他 Tab 页不受影响

#### Scenario: ConversationList 崩溃
- **WHEN** ConversationList 组件内部抛出未捕获异常（如 API 错误）
- **THEN** 会话列表区域显示错误提示 + 重试按钮，聊天区域保持正常

### Requirement: 降级 UI 提供重试入口
Error Boundary 的降级 UI SHALL 包含错误描述文字和重试按钮。点击重试按钮 SHALL 重新挂载被包裹的组件。

#### Scenario: 用户点击重试
- **WHEN** Error Boundary 显示降级 UI 后用户点击重试按钮
- **THEN** 被包裹的组件重新挂载并重新初始化

### Requirement: Error Boundary 组件实现
Error Boundary SHALL 作为通用 React 组件实现在 `components/ui/error-boundary.tsx`，使用 class component（React 错误边界要求），支持 `fallback` prop 自定义降级 UI。

#### Scenario: 自定义降级 UI
- **WHEN** 使用 ErrorBoundary 组件并传入 `fallback` prop
- **THEN** 错误发生时显示传入的 fallback 内容
