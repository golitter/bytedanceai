## REMOVED Requirements

### Requirement: TTL configuration
**Reason**: TTL 自动过期机制被 inactive-cleanup 轮询机制替代。workspace 不再有基于时间的自动过期，只有用户手动停用 session 后才触发清理。
**Migration**: 删除 `WORKSPACE_TTL_SECONDS` 和 `WORKSPACE_TTL_CHECK_INTERVAL` 配置项，替换为 `cleanup_interval` 配置项。

### Requirement: Background TTL cleanup task
**Reason**: TTL 后台清理任务被 inactive-cleanup 轮询任务替代。清理触发条件从"超时"变为"用户手动停用"。
**Migration**: 删除 TTL cleanup task，替换为 inactive-cleanup 轮询 task。

### Requirement: TTL cleanup logging
**Reason**: TTL 清理日志格式不再适用，被 inactive-cleanup 的日志格式替代。
**Migration**: 使用 inactive-cleanup spec 中定义的日志格式。
