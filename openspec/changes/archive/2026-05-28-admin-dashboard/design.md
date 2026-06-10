## Context

项目是 QQ 风格三栏布局的多 Agent 聊天系统（React 19 + Go 后端 + Python Agent 端）。IconSidebar 已有 admin tab 但为占位状态。选中 admin tab 后，中间栏（原 ConversationList 180px）需切换为管理菜单，右侧内容区渲染对应管理模块。

当前后端无管理相关 API，所有系统级操作需通过 shell 命令（`make status`、`scripts/test-clean.sh`）完成。Agent 配置文件（`~/.claude/settings.json`、`~/.codex/config.toml`）需手动编辑。

设计参考：`demo/admin.html` — 纯 HTML/CSS/JS 实现的完整原型。

## Goals / Non-Goals

**Goals:**
- 激活 admin tab，实现管理菜单 + 内容区四栏布局
- 6 个管理模块全部实现，连接真实后端 API
- 单人管理场景，不需要角色权限系统
- Agent 概览只读展示 CLI 配置文件内容
- 管理页面密码保护：进入需验证密码，查看敏感信息需二次验证

**Non-Goals:**
- 不做用户管理、角色权限、审计日志
- 不做 Agent 配置的在线编辑（只读）
- 不做 LLM API Key 管理、全局参数、数据备份导出
- 不做移动端适配

## Decisions

### D1: 前端管理菜单与聊天列表共用中间栏

**选择**: admin tab 激活时，中间栏从 ConversationList 切换为 AdminMenu 组件。
**原因**: 复用现有三栏布局结构，避免增加第四栏导致空间紧张。与 QQ 的联系人/设置切换模式一致。
**替代方案**: 独立路由页面（`/admin`）— 但会离开聊天上下文，体验割裂。

### D2: 后端新增 `/api/admin/*` 路由组

**选择**: 后端新增 admin 路由组，不复用现有 task/session 路由。
**原因**: 管理 API 语义清晰，权限隔离方便（未来加认证时统一拦截）。
**路由规划**:
- `POST /api/admin/auth` — 密码验证，签发 admin JWT
- `GET /api/admin/resources` — 系统资源（需 admin JWT）
- `DELETE /api/admin/sessions` — 批量清理会话（需 admin JWT）
- `GET /api/admin/workspaces` — 工作区列表
- `DELETE /api/admin/workspaces/:id` — 清理单个工作区
- `GET /api/admin/agents` — Agent 列表 + 配置文件内容
- `GET /api/admin/services` — 服务健康状态
- `GET /api/admin/statistics` — 统计数据

### D3: 系统资源查询使用 Go 标准库

**选择**: 使用 `runtime` + `os` + `syscall` 获取内存/磁盘信息，Redis 信息通过 `INFO memory` 命令。
**原因**: 无需引入第三方库（如 `gopsutil`），Go 标准库足够覆盖磁盘和内存查询。
**替代方案**: 引入 `github.com/shirou/gopsutil/v3` — 功能更强但对管理页面来说过重。

### D4: Agent 配置文件通过后端读取

**选择**: 后端读取 `~/.claude/settings.json`、`~/.codex/config.toml` 等文件并返回内容。
**原因**: 前端无法直接读取文件系统。后端与 Agent 端同机器部署，可直接访问。
**安全**: 脱敏敏感字段（API Key 等），用 `***` 替代。

### D5: 服务健康检测基于进程 + HTTP

**选择**: Backend/AgentEnd 通过 HTTP 健康端点检测（`/health`），Frontend 通过进程 PID 检测。
**原因**: AgentEnd 已有 `/health` 端点。Backend 新增 `/health` 端点。Frontend 无后端，用 PID 判断。

### D6: 管理页面密码保护方案

**选择**: 后端配置文件存储管理密码（bcrypt 哈希），前端首次进入 admin tab 弹出密码输入框，验证通过后签发短期 JWT（1 小时），后续 admin API 请求携带该 token。查看 Agent 配置文件等敏感信息时，即使已登录也需二次输入密码确认。
**原因**: 单人管理场景不需要完整用户体系，一个密码足够。bcrypt + 短期 JWT 平衡安全与便捷。敏感操作二次验证防止页面被他人误操作。
**替代方案**:
- 每次请求都验证密码 — 体验太差
- 纯前端密码验证（硬编码）— 不安全，可被绕过
**配置**: `configs/config.yaml` 新增 `admin.password_hash` 字段，首次启动时若无密码则通过 CLI 设置。

### D7: 用户头像修改

**选择**: 管理面板顶部管理菜单中增加头像修改入口。支持两种方式：上传自定义图片（存储到七牛云，复用现有 `POST /api/avatar/upload` 接口）或更换 DiceBear 风格/seed。
**原因**: 项目已有七牛云上传能力（`backend/internal/handler/avatar.go`）和 DiceBear 头像生成，直接复用。
**头像存储**: 头像 URL 存储在后端配置或数据库（单用户场景，配置文件即可）。

## Risks / Trade-offs

- **[配置文件路径]** 不同环境下 Agent 配置文件路径可能不同 → 通过配置文件指定路径，默认 `~/.claude` 等
- **[资源查询性能]** 不做自动轮询刷新，所有数据仅在用户主动点击菜单切换或点击刷新按钮时请求，避免不必要的后端压力
- **[批量清理安全]** 误删会话数据不可恢复 → 清理前返回预览（将删除多少条），前端二次确认
- **[配置文件敏感信息]** 展示配置文件可能泄露 API Key → 后端脱敏处理，正则匹配 key/token/secret 字段
