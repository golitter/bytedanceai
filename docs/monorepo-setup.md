# Monorepo 工程化配置说明

## 包管理

| 子项目 | 包管理器 | 说明 |
|--------|---------|------|
| 根目录 | pnpm | Husky / commitlint / lint-staged |
| frontend | pnpm | Next.js 生态 |
| backend | go mod | Go 模块管理 |
| agentend | uv | Python 依赖管理 |

---

## Git Hooks（Husky）

| 钩子 | 触发时机 | 作用 |
|------|---------|------|
| pre-commit | `git commit` 前 | 运行 lint-staged，检查暂存文件的代码风格 |
| commit-msg | `git commit` 前 | 运行 commitlint，校验 commit message 格式 |

---

## AgentEnd（Python FastAPI）

### 前置要求

| 工具 | 版本 | 安装 |
|------|------|------|
| Python | >= 3.10 | 系统自带 / `brew install python` |
| uv | latest | `brew install uv` |

### 核心依赖

定义在 `pyproject.toml`：

| 依赖 | 用途 |
|------|------|
| fastapi | HTTP 框架 |
| uvicorn | ASGI 服务器 |
| pydantic / pydantic-settings | 数据校验 + 配置加载 |
| langchain-core / langchain-anthropic / langchain-openai | LLM 调用 |
| langgraph | Agent DAG 编排 |
| sse-starlette | SSE 流式推送 |
| httpx | 异步 HTTP 客户端 |
| pyyaml | YAML 配置解析 |
| python-dotenv | .env 加载 |

开发依赖：`pytest` + `pytest-asyncio`。

---

## Frontend（Next.js）


## Backend（Go Gin）

