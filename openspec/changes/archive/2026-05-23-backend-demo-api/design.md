## Context

Backend 目录当前为空。monorepo-setup 和 setup 文档已定义好技术栈（Gin + GORM + MySQL）和分层架构（参考 gormlab：controller/dao/service/vo/model/middleware）。本 change 实现骨架代码和一个 demo 接口，验证链路可通。

## Goals / Non-Goals

**Goals:**

- 按 gormlab 模式搭建 controller → service → dao → model 四层骨架
- 实现 conf（YAML 配置）、pkg/db（MySQL 单例连接）、middleware（CORS / Logger / Auth）、vo（统一响应）
- 提供一个 ping 接口验证整条链路
- 提供 Makefile 支持构建和运行

**Non-Goals:**

- 不实现业务实体（User / Session / Task 等），后续 change 再做
- 不做 SSE 代理、AgentEnd 客户端、EventStore 等 runtime 层能力
- 不做 Swagger 文档生成
- 不做 RateLimit 中间件

## Decisions

1. **分层架构选 gormlab 模式**：interface + impl 分离，每层可独立测试。controller/dao/service 各自定义接口，impl 包提供实现。
2. **配置用 gopkg.in/yaml.v3 直接解析**，不引入 Viper（gormlab 实践证明轻量够用，减少依赖）。
3. **MySQL 连接用 sync.Once 单例**，通过 `pkg/db.Init(cfg)` 初始化，全局 `GetDB()` 获取。
4. **统一响应格式** `{code, data, msg}`，所有 handler 通过 `vo.OK()` / `vo.BadRequest()` 等返回。
5. **JWT Auth 中间件预置但 ping 接口不挂**，仅作为骨架代码存在，后续业务接口启用。

## Risks / Trade-offs

- **风险**: 配置路径硬编码 `configs/config.yaml` → 缓解：当前为 demo 阶段，后续可改为 flag 或环境变量
- **权衡**: 不引入 Viper 虽然减少依赖，但失去了环境变量覆盖、热加载等能力 → 后续按需加回
