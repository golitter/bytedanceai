# Setup Guide — AI Runtime Platform

三层架构一键搭建指南：Frontend (React + Vite) + Backend (Go) + AgentEnd (Python)。

---

## 前置依赖

| 工具       | 版本要求     | 安装                                |
| ---------- | ---------- | ----------------------------------- | ------- |
| Go         | >= 1.22    | `brew install go`                   |
| Node.js    | >= 18      | `brew install node`                 |
| pnpm       | >= 8       | `npm i -g pnpm`                     |
| Python     | >= 3.10    | 系统自带 / `brew install python`     |
| uv         | latest     | `brew install uv`                   |
| MySQL      | >= 8.0     | `brew install mysql`                |
| Redis      | >= 7       | `brew install redis`                |
| Docker     | optional   | `brew install --cask docker`        |

验证：

```bash
go version && node --version && pnpm --version && uv --version
```

---

## 目录结构（目标）

```
bytedanceai/
├── frontend/          # React + Vite + Tailwind + shadcn/ui
├── backend/           # Go + Gin + GORM + MySQL
├── agentend/          # Python FastAPI（已有）
├── docs/
├── scripts/
├── docker-compose.yml
├── Makefile
└── .env
```

---

## 1. AgentEnd（Python）— 已有，验证启动

```bash
cd agentend

# 安装依赖
uv sync

# 启动
uv run uvicorn src.api.app:app --host 0.0.0.0 --port 8001 --reload
```

验证：`curl http://localhost:8001/docs` → FastAPI Swagger 页面。

---

## 2. Backend（Go）— 从零初始化

### 2.1 初始化模块

```bash
cd backend
go mod init agenthub/backend
```

### 2.2 安装核心依赖

```bash
go get github.com/gin-gonic/gin \
       gorm.io/gorm \
       gorm.io/driver/mysql \
       gopkg.in/yaml.v3 \
       github.com/golang-jwt/jwt/v5 \
       github.com/gin-contrib/cors \
       github.com/swaggo/swag \
       github.com/swaggo/gin-swagger \
       github.com/swaggo/files \
       golang.org/x/time
```

### 2.3 创建目录结构

```bash
cd backend
mkdir -p cmd/server
mkdir -p internal/{conf,controller/impl,dao/gorm,dao/mock,model,service/impl,middleware,vo}
mkdir -p pkg/db
mkdir -p configs docs/api
```

目标结构（参考 gormlab 分层模式）：

```
backend/
├── cmd/server/main.go            # 入口：加载配置 → 连 DB → 注册路由 → 启动
├── internal/
│   ├── conf/conf.go              # 配置加载（YAML → struct）
│   ├── controller/               # 控制器接口
│   │   └── impl/                 # 控制器实现（路由注册 + Handler）
│   ├── dao/                      # 数据访问接口
│   │   ├── gorm/                 # GORM 实现
│   │   └── mock/                 # Mock 实现（测试用）
│   ├── model/                    # GORM 模型（表映射）
│   ├── service/                  # 业务逻辑接口
│   │   └── impl/                 # 业务逻辑实现
│   ├── middleware/                # Gin 中间件（CORS, Auth, Logger, RateLimit）
│   └── vo/                       # View Object（API 响应结构）
├── pkg/
│   └── db/mysql.go               # MySQL 连接（单例）
├── configs/config.yaml
├── docs/api/                     # Swagger 生成文件
├── Makefile
├── go.mod
└── go.sum
```

### 2.4 配置

创建 `backend/configs/config.yaml`：

```yaml
mysql:
  host: localhost
  port: 3306
  user: root
  password: "123456"
  dbname: agenthub

jwt:
  secret: "dev-secret-key"
  expire: 86400
```

### 2.5 配置加载

创建 `backend/internal/conf/conf.go`：

```go
package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
	JWT   JWTConfig   `yaml:"jwt"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file failed: %w", err)
	}
	return &cfg, nil
}
```

### 2.6 MySQL 连接

创建 `backend/pkg/db/mysql.go`：

```go
package db

import (
	"fmt"
	"sync"

	"agenthub/backend/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	globalDB *gorm.DB
	once     sync.Once
)

func dsn(cfg *conf.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)
}

func Init(cfg *conf.MySQLConfig) (*gorm.DB, error) {
	var err error
	once.Do(func() {
		globalDB, err = gorm.Open(mysql.Open(dsn(cfg)), &gorm.Config{})
	})
	return globalDB, err
}

func GetDB() *gorm.DB {
	return globalDB
}
```

### 2.7 统一响应

创建 `backend/internal/vo/response.go`：

```go
package vo

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

func OK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{Code: 0, Data: data})
}

func Created(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, Response{Code: 0, Data: data})
}

func BadRequest(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, Response{Code: 400, Msg: msg})
}

func NotFound(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusNotFound, Response{Code: 404, Msg: msg})
}

func InternalError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, Response{Code: 500, Msg: msg})
}
```

### 2.8 中间件

创建 `backend/internal/middleware/cors.go`：

```go
package middleware

import (
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
```

创建 `backend/internal/middleware/logger.go`：

```go
package middleware

import (
	"log/slog"
	"time"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		ctx.Next()
		latency := time.Since(start)
		status := ctx.Writer.Status()
		attrs := []any{
			slog.Int("status", status),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.Duration("latency", latency),
			slog.String("client_ip", ctx.ClientIP()),
		}
		switch {
		case status >= 500:
			slog.Error("request", attrs...)
		case status >= 400:
			slog.Warn("request", attrs...)
		default:
			slog.Info("request", attrs...)
		}
	}
}
```

创建 `backend/internal/middleware/auth.go`：

```go
package middleware

import (
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret []byte
	jwtExpire int
)

func SetJWTConfig(secret string, expire int) {
	jwtSecret = []byte(secret)
	jwtExpire = expire
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	jwt5.RegisteredClaims
}

func GenerateToken(userID uint, name string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt5.RegisteredClaims{
			ExpiresAt: jwt5.NewNumericDate(now.Add(time.Duration(jwtExpire) * time.Second)),
			IssuedAt:  jwt5.NewNumericDate(now),
		},
	}
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		claims := &Claims{}
		token, err := jwt5.ParseWithClaims(parts[1], claims, func(t *jwt5.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		ctx.Set("userID", claims.UserID)
		ctx.Set("userName", claims.Name)
		ctx.Next()
	}
}
```

### 2.9 main.go

创建 `backend/cmd/server/main.go`：

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"agenthub/backend/internal/conf"
	"agenthub/backend/internal/middleware"
	"agenthub/backend/pkg/db"
)

func main() {
	cfg, err := conf.Load("configs/config.yaml")
	if err != nil {
		panic("load config failed: " + err.Error())
	}

	middleware.SetJWTConfig(cfg.JWT.Secret, cfg.JWT.Expire)

	database, err := db.Init(&cfg.MySQL)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	_ = database
	fmt.Println("database connected")

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	fmt.Println("server running on :8080")
	if err := r.Run(":8080"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
```

### 2.10 Makefile

创建 `backend/Makefile`：

```makefile
.PHONY: build run swagger fmt tidy

BINARY := server
MAIN   := cmd/server/main.go

build:
	go build -o $(BINARY) $(MAIN)

run: build
	./$(BINARY)

swagger:
	swag init -g $(MAIN) -o docs/api

fmt:
	gofmt -w .

tidy:
	go mod tidy
```

### 2.11 启动 Backend

```bash
cd backend
make run
```

验证：`curl http://localhost:8080/ping` → `{"message":"pong"}`

---

## 3. Frontend（Next.js）— 从零初始化

### 3.1 创建 Next.js 项目

```bash
# 在项目根目录执行
npx create-next-app@latest frontend \
  --typescript \
  --tailwind \
  --eslint \
  --app \
  --src-dir \
  --use-pnpm \
  --turbopack
```

交互项全部选 YES（import alias 选默认 `@/*`）。

### 3.2 安装核心依赖

```bash
cd frontend

# 状态管理
pnpm add zustand

# 数据请求
pnpm add @tanstack/react-query

# UI 组件库
npx shadcn@latest init
# Style: Default, Base color: Slate, CSS variables: Yes

# 常用组件
npx shadcn@latest add button card input dialog

# 图标
pnpm add lucide-react

# Markdown 渲染
pnpm add react-markdown remark-gfm
```

### 3.3 目录结构

```
frontend/src/
├── app/                # Next.js App Router pages
│   ├── layout.tsx
│   ├── page.tsx
│   └── globals.css
├── components/         # 通用组件
│   └── ui/             # shadcn/ui 组件（自动生成）
├── features/           # 业务模块
│   ├── chat/
│   ├── runtime/
│   └── agent/
├── hooks/              # 自定义 hooks
├── stores/             # Zustand stores
├── services/           # API 调用封装
├── types/              # TypeScript 类型
└── lib/                # 工具函数
```

### 3.4 最小 API 服务层

创建 `frontend/src/services/api.ts`：

```typescript
const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function fetchAPI<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) throw new Error(`API error: ${res.status}`);
  return res.json();
}

export const api = {
  ping: () => fetchAPI<{ message: string }>("/ping"),
};
```

### 3.5 环境变量

创建 `frontend/.env.local`：

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 3.6 启动 Frontend

```bash
cd frontend
pnpm dev
```

验证：浏览器打开 `http://localhost:3000` → Next.js 默认页面。

---

## 4. 基础设施 — MySQL + Redis

### 4.1 使用 Docker Compose（推荐）

创建项目根目录 `docker-compose.yml`：

```yaml
services:
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "123456"
      MYSQL_DATABASE: agenthub
    volumes:
      - mysqldata:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  mysqldata:
```

启动：

```bash
docker compose up -d
```

### 4.2 手动安装（不用 Docker）

```bash
# MySQL
brew install mysql
brew services start mysql
mysql -u root -e "CREATE DATABASE agenthub;"

# Redis
brew install redis
brew services start redis
```

---

## 5. 统一环境变量

创建项目根目录 `.env`：

```env
# Backend
DATABASE_URL=root:123456@tcp(localhost:3306)/agenthub?charset=utf8mb4&parseTime=true&loc=Local
REDIS_URL=redis://localhost:6379
AGENTEND_URL=http://localhost:8001

# Frontend
VITE_API_URL=http://localhost:8080
```

---

## 6. Makefile

创建项目根目录 `Makefile`：

```makefile
.PHONY: dev frontend backend agentend db migrate

# 一键启动所有服务
dev:
	@make -j3 _frontend _backend _agentend

_frontend:
	cd frontend && pnpm dev

_backend:
	cd backend && make run

_agentend:
	cd agentend && uv run uvicorn src.api.app:app --host 0.0.0.0 --port 8001 --reload

# 单独启动
frontend:
	cd frontend && pnpm dev

backend:
	cd backend && make run

agentend:
	cd agentend && uv run uvicorn src.api.app:app --host 0.0.0.0 --port 8001 --reload

# 数据库
db:
	docker compose up -d

db-down:
	docker compose down

# 安装所有依赖
install:
	cd frontend && pnpm install
	cd backend && go mod download
	cd agentend && uv sync
```

---

## 7. 热更新（Backend）

```bash
go install github.com/air-verse/air@latest
```

创建 `backend/.air.toml`（最小配置）：

```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/server"
  bin = "./tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor"]
```

之后用 `air` 代替 `go run`：

```bash
cd backend && air
```

---

## 8. 完整启动流程

```bash
# 1. 启动基础设施
make db

# 2. 安装依赖（首次）
make install

# 3. 启动所有服务
make dev

# 或者分别启动（推荐开发时用，每个终端一个）
make backend     # → :8080
make agentend    # → :8001
make frontend    # → :3000
```

---

## 端口分配

| 服务         | 端口   |
| ------------ | ---- |
| Frontend     | 5173 |
| Backend      | 8080 |
| AgentEnd     | 8001 |
| MySQL        | 3306 |
| Redis        | 6379 |

---

## 验证清单

- [ ] `curl localhost:8080/ping` → `{"message":"pong"}`
- [ ] `curl localhost:8001/docs` → FastAPI Swagger
- [ ] 浏览器 `localhost:5173` → React 页面
- [ ] `docker compose ps` → mysql + redis running

---

## 常见问题

### MySQL 连接失败

```bash
# 检查是否运行
docker compose ps
# 或
brew services list | grep mysql

# 手动测试连接
mysql -u root -p agenthub
```

### Frontend 无法连接 Backend

检查 `frontend/.env.local` 中 `VITE_API_URL` 是否正确，重启 dev server。

### Go 依赖下载慢

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```
