## ADDED Requirements

### Requirement: TTL configuration
系统 SHALL 支持通过 `WORKSPACE_TTL_SECONDS` 环境变量配置 workspace 超时时长，默认值 3600（1 小时）。 SHALL 支持通过 `WORKSPACE_TTL_CHECK_INTERVAL` 配置检查间隔，默认值 300（5 分钟）。

#### Scenario: Use default TTL
- **WHEN** 未设置 `WORKSPACE_TTL_SECONDS`
- **THEN** SHALL 使用 3600 秒作为超时阈值

#### Scenario: Use custom TTL
- **WHEN** 设置 `WORKSPACE_TTL_SECONDS=7200`
- **THEN** SHALL 使用 7200 秒作为超时阈值

### Requirement: Background TTL cleanup task
系统 SHALL 在 FastAPI lifespan startup 时启动后台 asyncio task，定期扫描 ACTIVE workspace 并清理超时的。

#### Scenario: Cleanup expired workspace
- **WHEN** workspace 的 `created_at` 距当前时间超过 `WORKSPACE_TTL_SECONDS`
- **THEN** SHALL 自动执行 `cleanup(workspace_id)`，标记为 CLEANED

#### Scenario: Keep active workspace
- **WHEN** workspace 的 `created_at` 距当前时间未超过 `WORKSPACE_TTL_SECONDS`
- **THEN** SHALL 不执行任何操作

#### Scenario: Stop cleanup on shutdown
- **WHEN** FastAPI app shutdown
- **THEN** SHALL 取消 TTL cleanup 后台 task

### Requirement: TTL cleanup logging
每次 TTL 清理 SHALL 记录被清理的 workspace 数量和详情。

#### Scenario: Cleanup cycle logged
- **WHEN** TTL 检查周期执行
- **THEN** SHALL 输出日志："TTL cleanup: checked X workspaces, cleaned Y expired"
