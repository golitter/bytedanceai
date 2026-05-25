# Adapter Layer — 适配器层

## 实现了什么

统一 Agent 适配器抽象，当前 MVP 实现 `ClaudeCodeAdapter`，通过 CLI subprocess 驱动 Claude Code。

## 怎么实现的

### BaseAgentAdapter (`src/adapters/base.py`)

抽象基类，定义所有 Adapter 必须实现的接口：

```python
class BaseAgentAdapter(ABC):
    async def create_session(session_id) -> None
    async def chat(session_id, message, **kwargs) -> AgentResponse      # 同步调用
    async def stream_chat(session_id, message, **kwargs) -> AsyncIterator[StreamEvent]  # 流式调用
    async def interrupt(session_id) -> bool                              # 中断执行
    async def destroy_session(session_id) -> None                        # 销毁会话
```

### AdapterRegistry (`src/adapters/registry.py`)

通过 agent 类型名称注册和查找 Adapter 类：

```python
registry.register("claude-code", ClaudeCodeAdapter)
adapter_cls = registry.get("claude-code")  # 返回类，由调用方实例化
```

查找不存在的类型时抛出 `ValueError`。

### ClaudeCodeAdapter (`src/adapters/claude.py`)

核心实现，通过 `asyncio.create_subprocess_exec` 启动 Claude CLI 进程。

#### 命令构建 (`_build_command`)

将请求参数组装为 CLI 命令：

```python
claude -p "<message>" --output-format stream-json \
    [--session <id>] \
    [--append-system-prompt "<text>"] \
    [--allowedTools Read,Write] \
    [--max-turns 20]
```

参数来源：
- `-p`：用户消息
- `--session`：复用会话（对应 request.session_id）
- `--append-system-prompt`：Rule Engine 注入的约束文本
- `--allowedTools`：Rule Engine 限制的工具集
- `--max-turns`：执行轮数限制

#### 流式输出解析 (`_parse_stream_line`)

逐行读取 CLI 的 stdout，解析 JSON 并转换为 StreamEvent：

1. 空行 → 忽略
2. 合法 JSON → 按 `type` 字段映射为对应 StreamEvent
3. 非法 JSON → 包装为 `text` 类型事件，不抛异常

#### 流式调用 (`stream_chat`)

```python
async def stream_chat(self, session_id, message, **kwargs) -> AsyncIterator[StreamEvent]:
    cmd = self._build_command(message, ...)
    process = await asyncio.create_subprocess_exec(*cmd, stdout=PIPE, stderr=PIPE)
    self._processes[session_id] = process  # 记录进程句柄

    async for line in process.stdout:
        yield self._parse_stream_line(line.decode())

    # 进程退出后清理
    self._processes.pop(session_id, None)
```

进程通过 `dict[str, Process]` 管理，key 为 session_id。

#### 同步调用 (`chat`)

内部调用 `stream_chat`，收集所有 text 事件拼接后返回 `AgentResponse`。

#### 中断机制 (`interrupt`)

```
SIGTERM → 等待 5 秒 → SIGKILL
```

1. 查找 session 对应的进程
2. 发送 SIGTERM
3. `asyncio.wait_for(process.wait(), timeout=5)`
4. 超时则发送 SIGKILL
5. 从 `_processes` 中移除
