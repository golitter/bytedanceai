.PHONY: run-frontend run-backend run-agentend \
       stop stop-frontend stop-backend stop-agentend \
       restart-frontend restart-backend restart-agentend \
       status tidy

SCRIPT := ./scripts/run.sh

# 启动前端（热重载）— Vite dev server，localhost:5173
run-frontend:
	$(SCRIPT) start frontend

# 启动后端（热重载）— Air，localhost:8080
run-backend:
	$(SCRIPT) start backend

# 启动 Agent 端（热重载）— uvicorn --reload，localhost:8001
run-agentend:
	$(SCRIPT) start agentend

# 停止全部服务
stop:
	$(SCRIPT) stop

# 停止前端
stop-frontend:
	$(SCRIPT) stop frontend

# 停止后端
stop-backend:
	$(SCRIPT) stop backend

# 停止 Agent 端
stop-agentend:
	$(SCRIPT) stop agentend

# 重启前端（热重载）
restart-frontend:
	$(SCRIPT) restart frontend

# 重启后端（热重载）
restart-backend:
	$(SCRIPT) restart backend

# 重启 Agent 端（热重载）
restart-agentend:
	$(SCRIPT) restart agentend

# 查看三端运行状态（端口 + PID）
status:
	$(SCRIPT) status

# 整理 Go 依赖（go mod tidy）
tidy:
	cd backend && go mod tidy
