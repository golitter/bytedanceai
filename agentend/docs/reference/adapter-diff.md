# 适配器差异：Claude CLI vs OpenCode CLI

## 权限模型

| 特性 | Claude CLI | OpenCode CLI |
|------|-----------|--------------|
| 权限机制 | `--allowedTools` 命令行预授权 | 通过配置文件或内部机制管理，无命令行参数 |
| 适配器支持 | `allowed_tools` → `--allowedTools` | 接收但忽略 `allowed_tools` |
| 非交互模式限制 | 必须预授权工具，否则写入等操作会被拒绝 | 无此限制 |
| 危险工具过滤 | `_BLOCKED_TOOLS`（已定义，待接入） | 无 |

## 原因

Claude CLI 的 `-p` 非交互模式没有 TTY，无法弹出权限确认提示，因此必须通过 `--allowedTools` 预授权（如 `Write`、`Edit`、`Bash`），否则工具调用会被直接拒绝。可通过 `.claude/settings.local.json` 配置持久权限。

OpenCode CLI 的权限管理不依赖命令行参数，适配器目前不做工具过滤，`allowed_tools` 参数传入后被忽略。
