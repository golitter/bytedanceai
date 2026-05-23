#!/usr/bin/env bash
set -euo pipefail

# ── 服务定义 ──────────────────────────────────────
# name:port 格式（均为热重载模式）
SERVICES=(frontend:5173 backend:8080 agentend:8001)

# ── 颜色 ──────────────────────────────────────────
GREEN='\033[32m'
RED='\033[31m'
RESET='\033[0m'

# ── 工具函数 ──────────────────────────────────────
name_of() { cut -d: -f1 <<< "$1"; }
port_of() { cut -d: -f2 <<< "$1"; }

pid_on_port() {
  lsof -i ":$1" -sTCP:LISTEN -t 2>/dev/null | head -1
}

is_running() {
  local port
  port=$(port_of "$1")
  [ -n "$(pid_on_port "$port")" ]
}

# ── 启动单个服务 ──────────────────────────────────
start_service() {
  local entry=$1
  local name
  name=$(name_of "$entry")
  local port
  port=$(port_of "$entry")

  if is_running "$entry"; then
    echo "$name 已在运行 (port $port, PID $(pid_on_port "$port"))，跳过"
    return
  fi

  echo "启动 $name (port $port) ..."
  case "$name" in
    frontend)
      (cd frontend && exec pnpm dev) &
      ;;
    backend)
      (cd backend && exec ~/go/bin/air -c .air.toml) &
      ;;
    agentend)
      (cd agentend && exec uv run uvicorn src.app.main:app --reload --port "$port") &
      ;;
  esac

  # 等待端口就绪（最多 10 秒）
  local waited=0
  while [ $waited -lt 100 ] && ! is_running "$entry"; do
    sleep 0.1
    waited=$((waited + 1))
  done

  if is_running "$entry"; then
    echo "$name 启动成功 (PID $(pid_on_port "$port"))"
  else
    echo "$name 启动超时，请检查日志"
  fi
}

# ── 停止单个服务 ──────────────────────────────────
stop_service() {
  local entry=$1
  local name
  name=$(name_of "$entry")
  local port
  port=$(port_of "$entry")

  if ! is_running "$entry"; then
    echo "$name 未运行"
    return
  fi

  local pid
  pid=$(pid_on_port "$port")
  echo "停止 $name (PID $pid)"
  kill "$pid" 2>/dev/null || true
}

# ── 状态表格 ──────────────────────────────────────
show_status() {
  echo "┌──────────┬──────────┬──────┬─────────┐"
  echo "│ 服务     │ 状态     │ 端口 │ PID     │"
  echo "├──────────┼──────────┼──────┼─────────┤"
  for entry in "${SERVICES[@]}"; do
    local name
    name=$(name_of "$entry")
    local port
    port=$(port_of "$entry")
    if is_running "$entry"; then
      printf "│ %-8s │ ${GREEN}%-8s${RESET} │ %-4s │ %-7s │\n" "$name" "运行中" "$port" "$(pid_on_port "$port")"
    else
      printf "│ %-8s │ ${RED}%-8s${RESET} │ %-4s │ %-7s │\n" "$name" "未运行" "$port" "-"
    fi
  done
  echo "└──────────┴──────────┴──────┴─────────┘"
}

# ── 入口 ──────────────────────────────────────────
cmd=${1:-help}
target=${2:-}

# 根据 name 查找完整 entry
find_entry() {
  for e in "${SERVICES[@]}"; do
    [ "$(name_of "$e")" = "$1" ] && echo "$e" && return
  done
  echo ""
}

case "$cmd" in
  start)
    if [ -z "$target" ]; then
      echo "错误：不允许同时启动全部服务，请指定单个服务"
      echo "用法: $0 start <frontend|backend|agentend>"
      exit 1
    fi
    entry=$(find_entry "$target")
    [ -z "$entry" ] && { echo "未知服务: $target"; exit 1; }
    start_service "$entry"
    ;;
  stop)
    if [ -n "$target" ]; then
      entry=$(find_entry "$target")
      [ -z "$entry" ] && { echo "未知服务: $target"; exit 1; }
      stop_service "$entry"
    else
      for e in "${SERVICES[@]}"; do stop_service "$e"; done
      echo "✓ 全部已停止"
    fi
    ;;
  restart)
    if [ -z "$target" ]; then
      echo "错误：不允许同时重启全部服务，请指定单个服务"
      echo "用法: $0 restart <frontend|backend|agentend>"
      exit 1
    fi
    entry=$(find_entry "$target")
    [ -z "$entry" ] && { echo "未知服务: $target"; exit 1; }
    stop_service "$entry"
    start_service "$entry"
    ;;
  status)
    show_status
    ;;
  help|*)
    echo "用法: ./scripts/run.sh <start|stop|restart|status> <frontend|backend|agentend>"
    echo ""
    echo "  start    启动单个服务（必须指定）"
    echo "  stop     停止服务（不指定则全部停止）"
    echo "  restart  重启单个服务（必须指定）"
    echo "  status   查看三端运行状态"
    ;;
esac
