## ADDED Requirements

### Requirement: WorkspaceStoreProtocol abstract interface
系统 SHALL 定义 `WorkspaceStoreProtocol` 抽象接口，包含方法：`load_all`、`save`、`delete`、`query_by_task`、`query_by_status`。所有 workspace 持久化操作 MUST 通过此接口。

#### Scenario: Load all workspaces from store
- **WHEN** 调用 `store.load_all()`
- **THEN** SHALL 返回 `dict[str, Workspace]`，key 为 workspace_id

#### Scenario: Save workspace to store
- **WHEN** 调用 `store.save(workspace)`
- **THEN** SHALL 将 workspace 数据持久化，已存在则更新

#### Scenario: Delete workspace from store
- **WHEN** 调用 `store.delete(workspace_id)`
- **THEN** SHALL 从存储中移除对应记录

#### Scenario: Query by task_id
- **WHEN** 调用 `store.query_by_task("task-123")`
- **THEN** SHALL 返回该 task 下所有 Workspace 列表

#### Scenario: Query by status
- **WHEN** 调用 `store.query_by_status(WorkspaceStatus.ACTIVE)`
- **THEN** SHALL 返回所有指定状态的 Workspace 列表

### Requirement: JsonFileWorkspaceStore implementation
系统 SHALL 提供 `JsonFileWorkspaceStore` 实现 `WorkspaceStoreProtocol`，使用 JSON 文件存储。存储路径默认为 `logs/workspaces.json`。

#### Scenario: Persist workspace to JSON file
- **WHEN** 调用 `store.save(workspace)` 且 store 路径为 `logs/workspaces.json`
- **THEN** SHALL 将 workspace 序列化为 JSON 写入文件，目录不存在时自动创建

#### Scenario: Load workspaces from existing file
- **WHEN** store 初始化且 `logs/workspaces.json` 已存在
- **THEN** SHALL 从文件加载所有 workspace 记录到内存

#### Scenario: Handle corrupted store file
- **WHEN** JSON 文件内容损坏无法解析
- **THEN** SHALL 记录警告日志，返回空 dict，不抛异常

#### Scenario: Empty store file on first run
- **WHEN** store 初始化且文件不存在
- **THEN** SHALL 创建空 dict，不创建文件

### Requirement: WorkspaceStore concurrent write safety
`JsonFileWorkspaceStore` MUST 内部使用 `asyncio.Lock` 保护文件写操作，避免多协程并发写同一文件。

#### Scenario: Concurrent save operations
- **WHEN** 两个协程同时调用 `store.save()`
- **THEN** SHALL 串行化写操作，后写的不丢失数据
