## 1. RuntimeHub Core

- [x] 1.1 创建 `backend/internal/stream/hub.go`，定义 `HubEvent` 结构体（Seq uint64, Data string, Done bool）
- [x] 1.2 实现 `RuntimeStream` 结构体：seq（atomic.Uint64）、subscribers map、mu sync.Mutex、closed bool
- [x] 1.3 实现 `RuntimeHub` 结构体：streams map[string]*RuntimeStream、mu sync.RWMutex
- [x] 1.4 实现 `RuntimeHub.Publish(key string, data string)` 方法：获取或创建 stream，递增 seq，非阻塞发送给所有 subscriber channels
- [x] 1.5 实现 `RuntimeHub.Subscribe(key string) (<-chan HubEvent, uint64)` 方法：创建 buffered channel（256），注册 subscriber，返回 channel 和 currentSeq
- [x] 1.6 实现 `RuntimeHub.Close(key string)` 方法：发送 Done event，关闭所有 subscriber channels，从 map 移除 stream
- [x] 1.7 添加全局 `var Hub = &RuntimeHub{}` 单例

## 2. StreamWriter 集成 Hub

- [x] 2.1 在 `StreamWriter` 构造函数中保存 stream key 用于 hub.Publish
- [x] 2.2 修改 `publishToRedis` 方法：在 XADD 之前调用 `stream.Hub.Publish(sw.streamKey, line)`（仅对 data: 开头的行）
- [x] 2.3 修改 `flushTextBuffer`：合并后的 SSE line 也通过 hub.Publish 发布
- [x] 2.4 修改 `finish()` 方法：在 registry.Delete 之前调用 `stream.Hub.Close(sw.streamKey)`
- [x] 2.5 修改 `PublishErrorAndFail` 函数：通过 hub.Publish 发送 error event

## 3. StreamHandler 三阶段握手

- [x] 3.1 修改 `serveStreaming` Phase 2：替换 Redis XRead 循环为 hub.Subscribe + Redis gap replay + hub channel 消费
- [x] 3.2 实现 gap replay 逻辑：subscribe 获取 currentSeq，XREAD 从 lastSeq 到最新，补齐 MySQL 历史和 hub 实时之间的 gap
- [x] 3.3 实现 hub channel 消费循环：select channel event，写入 SSE response，flush；收到 Done event 时发送 done SSE 并退出
- [x] 3.4 处理 hub.Subscribe 返回空 channel 的情况（stream 尚未创建）：定期检查 MySQL status 或等待 channel 事件

## 4. 验证与清理

- [x] 4.1 编译通过，无 race condition（`go build` + `go test -race ./internal/stream/...`）
- [ ] 4.2 启动三端服务，发送消息验证 SSE token 延迟降至 <50ms
- [ ] 4.3 测试断线重连：前端断开后重连，内容完整无丢失
- [ ] 4.4 测试 agent switch：orchestrator 切换 agent 时 SSE 流正常分段
- [ ] 4.5 测试 error 场景：agentend 返回错误时前端正确显示错误
