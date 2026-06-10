## Why

Backend 目录当前为空，需要按照 monorepo-setup 和 setup 文档中定义的 gormlab 分层架构搭建 Go 后端骨架，验证 Gin + GORM + MySQL 技术栈可运行。

## What Changes

- 初始化 backend Go 模块及全部依赖
- 按 controller/dao/service/vo/model/middleware 分层创建项目骨架
- 实现 MySQL 连接（单例）、YAML 配置加载、统一响应、中间件（CORS / Logger / Auth）
- 实现一个 demo 接口（ping），验证整条链路可通

## Capabilities

### New Capabilities

- `backend-skeleton`: Go 后端项目骨架 — 分层目录结构、配置加载、MySQL 连接、中间件、统一响应、Makefile

### Modified Capabilities

## Impact

- 新增 `backend/` 目录下全部 Go 源文件
- 依赖 MySQL 服务运行（Docker 或本地）
- 端口 8080
