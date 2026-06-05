# Codex Adapter DONE 事件补发 + Resume 省略空 Prompt

## 变更原因

Codex CLI 进程在 `--resume` 模式下可能正常退出（exit code 0）但未输出 `turn.completed` 事件，导致前端 SSE 流永远收不到 `done` 事件，会话卡在 "运行中" 状态。此外，resume 模式下传递空字符串 prompt 参数会导致 CLI 行为异常。

## 变更文件

- `agentend/src/adapters/codex.py`：行为修复（非 schema 变更）

## 对比结果

| 项目 | 修复前 | 修复后 |
|------|--------|--------|
| 进程正常退出但无 `turn.completed` | 无 DONE 事件，前端卡住 | 自动补发 `done` 事件 |
| Resume + 空 prompt | 追加空字符串 `""` 到 CLI 命令 | 跳过空 prompt 参数 |

## 跨端影响

- **Frontend**：无需改动。前端已正确处理 `done` 事件，修复后行为与已有逻辑完全兼容。
- **Backend**：无需改动。Backend 不直接处理 Codex 的 SSE 流。
- **AgentEnd**：`CodexAdapter` 内部行为修复，接口签名不变。

## 契约变更

无。`done` 事件类型已在 `contracts/schemas/event-types.yaml` 中定义（`EventType.done`），本次仅修正 Adapter 实现，确保在所有退出路径下都能正确发出该事件。

## 附加变更

- `scripts/run.sh`：uvicorn `--reload` 限定到手写源码目录，排除 `src/generated/`，避免契约生成触发不必要的重载。
- `agentend/tests/test_codex_adapter.py`：新增单元测试覆盖上述两项修复。
