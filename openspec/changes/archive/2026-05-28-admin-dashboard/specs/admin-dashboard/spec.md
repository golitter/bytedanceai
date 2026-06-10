## ADDED Requirements

### Requirement: Admin tab activates management menu
当用户点击 IconSidebar 的 admin tab 时，系统 SHALL 将中间栏从会话列表切换为管理菜单，并默认展示总览仪表盘。

#### Scenario: Switch to admin tab
- **WHEN** 用户点击 IconSidebar 的管理 tab
- **THEN** 中间栏显示管理菜单（180px 宽），内容区显示总览仪表盘模块，admin tab 高亮

#### Scenario: Switch back to chat tab
- **WHEN** 用户从 admin tab 切换回聊天 tab
- **THEN** 中间栏恢复为会话列表，内容区恢复为聊天区域

### Requirement: Management menu with 6 modules
管理菜单 SHALL 包含 6 个菜单项：总览仪表盘、会话清理、工作区管理、Agent 概览、服务健康、数据统计。

#### Scenario: Navigate between modules
- **WHEN** 用户点击管理菜单中的任意菜单项
- **THEN** 内容区切换到对应模块，当前菜单项高亮

### Requirement: Management menu is persistent
管理菜单在 admin tab 激活期间 SHALL 始终可见，不可收起。

#### Scenario: Menu stays visible on scroll
- **WHEN** 内容区内容超出视口高度，用户滚动内容
- **THEN** 管理菜单保持固定可见，不随内容滚动
