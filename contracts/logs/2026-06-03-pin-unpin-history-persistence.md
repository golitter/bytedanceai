# 2026-06-03 Pin 取消事件持久化 + save_mem_node 去重

**类型**: agentend（无跨端契约变更）

## 变更说明

1. **ConversationMemoryStore** 新增 `replace_messages()` 方法，直接覆盖写入，消除 `save_messages()` 的 load + append 导致的历史重复
2. **save_mem_node** 改用 `replace_messages()` 写入全量权威数据
3. **PinMemory.unpin()** 返回值从 `bool` 改为 `dict | None`，让调用方获取被移除的 pin 元数据
4. **pin_remove API** 取消 Pin 后追加 unpin SystemMessage 到 ConversationMemoryStore，使 LLM 可感知已取消的约束

## 影响范围

- **Agentend**: 4 个文件变更（conversation_memory / graph / pin_memory / pin API）
- **Frontend**: 无影响（pin_remove 响应新增 `removed` 字段为增量兼容，不破坏现有解析）
- **Backend**: 无影响
- **contracts/schemas/**: 无变更
