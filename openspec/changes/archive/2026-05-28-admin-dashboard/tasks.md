## 1. 后端 — Admin 认证与 API 基础设施

- [x] 1.1 在 `configs/config.yaml` 新增 `admin.password_hash` 字段，创建密码设置辅助脚本
- [x] 1.2 创建 `backend/internal/middleware/admin_auth.go`，实现 admin JWT 中间件（bcrypt 验证 + JWT 签发/校验）
- [x] 1.3 创建 `backend/internal/handler/admin.go`，注册 `/api/admin/*` 路由组，`POST /api/admin/auth` 免认证，其余路由挂载 admin JWT 中间件
- [x] 1.4 创建 `backend/internal/handler/admin_resource.go`，实现 `GET /api/admin/resources`（磁盘/内存/Redis）
- [x] 1.5 创建 `backend/internal/handler/admin_session.go`，实现 `DELETE /api/admin/sessions`（批量删除会话及关联数据）
- [x] 1.6 创建 `backend/internal/handler/admin_workspace.go`，实现 `GET /api/admin/workspaces` + `DELETE /api/admin/workspaces/:id`（列表 + 清理 worktree）
- [x] 1.7 创建 `backend/internal/handler/admin_agent.go`，实现 `GET /api/admin/agents`（读取 CLI 配置文件 + 脱敏）
- [x] 1.8 创建 `backend/internal/handler/admin_health.go`，实现 `GET /api/admin/services`（三端健康检测）+ Backend `/health` 端点
- [x] 1.9 创建 `backend/internal/handler/admin_stats.go`，实现 `GET /api/admin/statistics`（会话趋势、消息统计、存储趋势）
- [x] 1.10 创建 `backend/internal/handler/admin_avatar.go`，实现 `GET /api/admin/avatar` + `PUT /api/admin/avatar`（获取/保存头像 URL，复用七牛云上传）

## 2. 前端 — 管理面板布局与认证

- [x] 2.1 修改 `IconSidebar`，启用 admin tab 的点击事件，不再 disabled
- [x] 2.2 创建 `AdminMenu` 组件（180px 宽，6 个菜单项，常驻显示），顶部显示用户头像，点击头像弹出修改弹窗
- [x] 2.3 修改 `ImPage` 的三栏布局逻辑：admin tab 激活时中间栏渲染 AdminMenu 而非 ConversationList
- [x] 2.4 创建 `stores/admin.ts`（Zustand），管理 activeMenuKey、adminToken、isAuthenticated 状态
- [x] 2.5 创建 `AdminPasswordDialog` 组件（密码输入弹窗），首次进入 admin tab 及查看敏感信息时触发

## 3. 前端 — 6 个管理模块页面

- [x] 3.1 创建 `pages/admin/DashboardPage.tsx`（系统资源进度条 + 关键指标卡片）
- [x] 3.2 创建 `pages/admin/SessionCleanupPage.tsx`（筛选 + 表格 + 批量清理 + toast）
- [x] 3.3 创建 `pages/admin/WorkspacePage.tsx`（工作区表格 + 磁盘进度条 + 清理操作）
- [x] 3.4 创建 `pages/admin/AgentOverviewPage.tsx`（Agent 卡片 + 配置文件展开/收起 + 二次密码验证）
- [x] 3.5 创建 `pages/admin/ServiceHealthPage.tsx`（三端服务状态卡片 + 刷新按钮）
- [x] 3.6 创建 `pages/admin/StatisticsPage.tsx`（柱状图日/周切换 + 消息总量 + 存储趋势）

## 4. 前端 — API 集成与数据层

- [x] 4.1 在 `lib/api.ts` 中添加 admin API 请求函数（adminAuth, getResources, deleteSessions, getWorkspaces, deleteWorkspace, getAgents, getServices, getStatistics, getAvatar, updateAvatar），请求自动携带 admin JWT
- [x] 4.2 创建 `hooks/use-admin-data.ts`（TanStack React Query hooks，封装各 admin API 调用，401 时自动清除 token 触发重新验证）
- [x] 4.3 各模块页面接入 hooks，替换 mock 数据为真实 API 调用

## 5. 验证与清理

- [ ] 5.1 `make run-frontend` + `make run-backend` 启动三端服务，验证 admin tab 切换 + 密码验证流程
- [ ] 5.2 验证 6 个模块页面数据展示正确（资源、会话、工作区、Agent、健康、统计）
- [ ] 5.3 验证会话清理和工作区清理操作正常（选中 → 清理 → toast → 数据更新）
- [ ] 5.4 验证 Agent 配置查看的二次密码验证流程
- [ ] 5.5 验证 admin JWT 过期后自动弹出密码框
- [ ] 5.5 验证 admin JWT 过期后自动弹出密码框
- [ ] 5.6 验证头像修改功能（上传图片、切换 DiceBear 风格、IconSidebar 同步更新）
- [ ] 5.7 删除 `demo/admin.html`（参考原型已完成使命）
