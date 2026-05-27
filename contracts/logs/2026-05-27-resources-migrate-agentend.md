# 2026-05-27 — 磁盘/内存监控迁移至 AgentEnd

## 变更原因

将系统资源监控（磁盘、内存）的计算逻辑从 Backend（Go）迁移至 AgentEnd（Python），Backend 改为代理转发，降低 Backend 的平台依赖（syscall/vm_stat 等 macOS 专用调用）。

## 变更文件

**无 schema 文件变更。**

Resources API 属于 Admin 单端通信，不走 `contracts/schemas/` 生成流程（参见 2026-05-27-admin-dashboard-api.md）。

## 架构变更

| 项目 | 变更前 | 变更后 |
|------|--------|--------|
| 磁盘检测 | Backend: `syscall.Statfs` | AgentEnd: `shutil.disk_usage` |
| 内存检测 | Backend: `sysctl` + `vm_stat` | AgentEnd: `sysctl`/`vm_stat`（macOS）或 `/proc/meminfo`（Linux） |
| Redis 检测 | Backend 本地查询 | **不变**，仍在 Backend 本地查询 |
| `GET /api/admin/resources` | Backend 直接计算并返回 | Backend 从 AgentEnd `GET /v1/resources` 获取 disk/memory，合并本地 redis 后返回 |

## 新增端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/v1/resources` | AgentEnd 新增，返回磁盘和内存使用情况 |

## 前端响应格式

不变，仍返回 `{ disk, memory, redis }`，前端无需修改。

## 跨端影响

| 端 | 影响 |
|------|------|
| Backend | `AdminHandler` 注入 `agentClient`；`admin_resource.go` 移除本地 disk/memory 计算，改为代理调用 |
| AgentEnd | 新增 `resources.py` 路由 + `resources_router` 注册 |
| Frontend | 无影响 |
