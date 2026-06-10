## ADDED Requirements

### Requirement: Admin tab requires password
用户首次点击 admin tab 时 SHALL 弹出密码输入框，验证通过后才能进入管理页面。

#### Scenario: First access to admin tab
- **WHEN** 用户点击 admin tab 且未持有有效 admin JWT
- **THEN** 弹出密码输入对话框，管理菜单和内容区不可见，直到密码验证通过

#### Scenario: Correct password
- **WHEN** 用户输入正确密码并提交
- **THEN** 后端签发 admin JWT（1 小时有效期），前端存储 token，管理页面正常展示

#### Scenario: Wrong password
- **WHEN** 用户输入错误密码并提交
- **THEN** 显示错误提示"密码错误"，密码框不清空，允许重试

#### Scenario: Admin session expires
- **WHEN** admin JWT 过期（超过 1 小时）
- **THEN** 下次访问 admin tab 或调用 admin API 时返回 401，前端重新弹出密码输入框

### Requirement: Sensitive data requires re-authentication
查看 Agent 配置文件等敏感信息时 SHALL 要求二次密码验证，即使已持有有效 admin JWT。

#### Scenario: View agent config requires re-auth
- **WHEN** 用户点击"查看配置"展开 Agent 配置文件内容
- **THEN** 弹出密码确认框，输入正确密码后才展示配置内容

#### Scenario: Re-auth session is short-lived
- **WHEN** 用户通过二次验证查看配置
- **THEN** 该次验证的有效期为 10 分钟，10 分钟内再次查看其他 Agent 配置无需重新输入

### Requirement: Admin auth API
后端 SHALL 提供密码验证接口。

#### Scenario: Verify password and issue token
- **WHEN** 前端发送 `POST /api/admin/auth` 且 body 包含 `{ password: "xxx" }`
- **THEN** 后端使用 bcrypt 对比密码哈希，正确则返回 `{ token: "jwt...", expires_in: 3600 }`，错误则返回 401

### Requirement: Admin JWT middleware
所有 `/api/admin/*` 路由（除 `POST /api/admin/auth` 外）SHALL 要求有效 admin JWT。

#### Scenario: Request without token
- **WHEN** 请求 admin API 但 Authorization header 无 token
- **THEN** 返回 401 Unauthorized

#### Scenario: Request with valid token
- **WHEN** 请求 admin API 且携带有效 admin JWT
- **THEN** 正常处理请求

### Requirement: Password stored in config
管理密码 SHALL 存储在后端配置文件中（bcrypt 哈希），不存储明文。

#### Scenario: Config stores hashed password
- **WHEN** 后端启动读取 `configs/config.yaml`
- **THEN** 从 `admin.password_hash` 字段读取 bcrypt 哈希值
