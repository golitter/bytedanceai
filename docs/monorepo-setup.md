# Monorepo 工程化配置说明

## 项目结构

```
bytedanceai/
├── frontend/      # React 前端
├── backend/       # Go 后端
├── agentend/      # Python Agent 端
├── docs/          # 文档
└── scripts/       # 脚本
```

## 包管理

使用 **pnpm** 作为包管理器。

## Git Hooks（Husky）

| 钩子 | 触发时机 | 作用 |
|------|---------|------|
| pre-commit | `git commit` 前 | 运行 lint-staged，检查暂存文件的代码风格 |
| commit-msg | `git commit` 前 | 运行 commitlint，校验 commit message 格式 |

## 代码检查与格式化（lint-staged）

| 子项目 | 匹配文件 | 执行命令 |
|--------|---------|---------|
| frontend | `**/*.{ts,tsx}` | `eslint --fix --max-warnings=0` + `prettier --write` |
| backend | `**/*.go` | `gofmt -w` + `goimports -w` |
| agentend | `**/*.py` | `ruff check --fix` + `ruff format` |

## Commit 规范（commitlint）

遵循 `@commitlint/config-conventional` 约定式提交，scope 不能为空且只能是以下值：

| scope | 说明 |
|-------|------|
| frontend | 前端 |
| backend | 后端 |
| agentend | Agent 端 |
| common | 公共 |
| docs | 文档 |
| other | 其他 |

提交格式：

```
<type>(<scope>): <描述>
```

示例：

```
feat(frontend): 添加登录页面
fix(backend): 修复数据库连接超时
docs(agentend): 更新 API 文档
```

## Python 配置（ruff）

配置文件：`agentend/ruff.toml`

- 行宽：120
- Lint 规则：E, F, I, N, W, UP
- 格式化：双引号
