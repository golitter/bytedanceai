#!/usr/bin/env bash
set -euo pipefail

PID_DIR=".pids"
SERVICES=(frontend backend agentend)

# ── 颜色 ──────────────────────────────────────────
GREEN='\033[32m'
RED='\033[31m'
RESET='\033[0m'

# ── 工具函数 ──────────────────────────────────────
mkdir -p "$PID_DIR"

pid_file()  { echo "$PID_DIR/$1.pid"; }

is_running() {
  local pf
  pf=$(pid_file "$1")
  [ -f "$pf" ] && kill -0 "$(cat "$pf")" 2>/dev/null
}

save_pid()  { echo $! > "$(pid_file "$1")"; }

read_pid()  { cat "$(pid_file "$1")"; }

clear_pid() { rm -f "$(pid_file "$1")"; }

# ── 启动单个服务 ──────────────────────────────────
start_service() {
  local name=$1
  if is_running "$name"; then
    echo "$name 已在运行 (PID $(read_pid "$name"))，跳过"
    return
  fi

  echo "启动 $name ..."
  case "$name" in
    frontend)
      (cd frontend && pnpm dev &)
      ;;
    backend)
      (cd backend && ~/go/bin/air -c .air.toml &)
      ;;
    agentend)
      (cd agentend && uv run uvicorn src.app.main:app --reload &)
      ;;
  esac
  save_pid "$name"
}

# ── 停止单个服务 ──────────────────────────────────
stop_service() {
  local name=$1
  if is_running "$name"; then
    echo "停止 $name (PID $(read_pid "$name"))"
    kill "$(read_pid "$name")" 2>/dev/null || true
  fi
  clear_pid "$name"
}

# ── 状态表格 ──────────────────────────────────────
show_status() {
  echo "┌──────────┬──────────┬────────────────┐"
  echo "│ 服务     │ 状态     │ PID            │"
  echo "├──────────┼──────────┼────────────────┤"
  for name in "${SERVICES[@]}"; do
    if is_running "$name"; then
      printf "│ %-8s │ ${GREEN}%-8s${RESET} │ %-14s │\n" "$name" "运行中" "$(read_pid "$name")"
    else
      printf "│ %-8s │ ${RED}%-8s${RESET} │ %-14s │\n" "$name" "未运行" "-"
    fi
  done
  echo "└──────────┴──────────┴────────────────┘"
}

# ── 入口 ──────────────────────────────────────────
cmd=${1:-help}
target=${2:-}

case "$cmd" in
  start)
    if [ -z "$target" ]; then
      echo "错误：不允许同时启动全部服务，请指定单个服务"
      echo "用法: $0 start <frontend|backend|agentend>"
      exit 1
    fi
    start_service "$target"
    ;;
  stop)
    if [ -n "$target" ]; then
      stop_service "$target"
    else
      for s in "${SERVICES[@]}"; do stop_service "$s"; done
      echo "✓ 全部已停止"
    fi
    ;;
  restart)
    if [ -z "$target" ]; then
      echo "错误：不允许同时重启全部服务，请指定单个服务"
      echo "用法: $0 restart <frontend|backend|agentend>"
      exit 1
    fi
    stop_service "$target"
    start_service "$target"
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
