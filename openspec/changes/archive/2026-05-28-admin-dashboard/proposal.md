## Why

项目已有 QQ 风格三栏布局，IconSidebar 中 admin tab 为占位状态。需要实现单人管理面板，让管理员在一个页面内掌控系统状态、管理数据生命周期、查看 Agent 配置，无需手动编辑配置文件或执行 shell 命令。

## What Changes

- 激活 IconSidebar admin tab，选中时中间栏从会话列表切换为管理菜单（180px 常驻）
- 新增 6 个管理模块页面：总览仪表盘、会话清理、工作区管理、Agent 概览、服务健康、数据统计
- 后端新增 `/api/admin/*` 系列接口：系统资源查询、会话批量清理、工作区管理、服务健康检查、统计数据聚合
- Agent 概览为只读展示，读取各 CLI 工具配置文件内容并展示
- **管理页面密码保护**：进入 admin tab 需输入密码；查看敏感信息（如 Agent 配置中的 API Key）需二次验证
- **用户头像修改**：管理面板中可修改用户头像（上传图片或更换 DiceBear 风格）
- 设计参考：`demo/admin.html`

## Capabilities

### New Capabilities
- `admin-dashboard`: 管理面板整体布局 — admin tab 切换、管理菜单组件、内容区路由
- `admin-auth`: 管理页面密码保护 — 进入 admin 需密码验证，查看敏感信息需二次验证
- `admin-avatar`: 用户头像管理 — 上传图片或更换 DiceBear 风格修改头像
- `admin-resource`: 系统资源监控 — 磁盘/内存/Redis 使用率查询
- `admin-session-cleanup`: 会话清理 — 按条件筛选、批量删除会话及关联数据
- `admin-workspace`: 工作区管理 — 工作区列表、磁盘占用查询、过期清理
- `admin-agent-overview`: Agent 概览 — Agent 列表 + CLI 配置文件只读展示
- `admin-service-health`: 服务健康 — 三端服务状态、运行时长、版本信息
- `admin-statistics`: 数据统计 — 会话趋势、消息总量、存储趋势

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **前端**: 新增管理模块页面组件、admin tab 路由逻辑、管理菜单组件；复用现有 IconSidebar
- **后端**: 新增 admin API 路由组，需要系统级操作能力（磁盘查询、进程检测、文件读取）；新增密码验证中间件
- **Agent 端**: 无变更（Agent 概览通过后端读取配置文件，不直接调用 Agent 端）
- **依赖**: 后端可能需要 `psutil` 等系统监控库的 Go 等价物
