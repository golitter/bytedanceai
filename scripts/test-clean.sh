#!/usr/bin/env bash
set -euo pipefail

# ── 项目根目录（基于脚本所在位置）────────────────────
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# ── Docker 容器名（可通过环境变量覆盖）─────────────────
MYSQL_CONTAINER="${MYSQL_CONTAINER:-mysql}"
REDIS_CONTAINER="${REDIS_CONTAINER:-redis}"

# ── MySQL 连接信息 ──────────────────────────────────
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASSWORD="123456"
MYSQL_DB="agenthub"

# ── 颜色 ────────────────────────────────────────────
GREEN=$'\033[32m'
RED=$'\033[31m'
YELLOW=$'\033[33m'
CYAN=$'\033[36m'
BOLD=$'\033[1m'
RESET=$'\033[0m'

# ── 全局选项 ─────────────────────────────────────────
YES=false
DRY_RUN=false
COMMAND=""
REPO_ARG=""

# ── 工具函数 ─────────────────────────────────────────

info()  { echo -e "${CYAN}[INFO]${RESET}  $*"; }
ok()    { echo -e "${GREEN}[OK]${RESET}    $*"; }
warn()  { echo -e "${YELLOW}[WARN]${RESET}  $*"; }
err()   { echo -e "${RED}[ERROR]${RESET} $*" >&2; }

# ── Docker 检测 ──────────────────────────────────────

check_docker() {
  if ! docker info &>/dev/null; then
    err "Docker daemon 未运行，请先启动 Docker"
    return 1
  fi
}

check_mysql_container() {
  if ! docker ps --format '{{.Names}}' | grep -q "^${MYSQL_CONTAINER}$"; then
    err "MySQL 容器 '${MYSQL_CONTAINER}' 未运行"
    err "提示：通过 MYSQL_CONTAINER 环境变量指定容器名"
    return 1
  fi
}

check_redis_container() {
  if ! docker ps --format '{{.Names}}' | grep -q "^${REDIS_CONTAINER}$"; then
    err "Redis 容器 '${REDIS_CONTAINER}' 未运行"
    err "提示：通过 REDIS_CONTAINER 环境变量指定容器名"
    return 1
  fi
}

# ── 交互式确认 ────────────────────────────────────────

confirm() {
  local msg="$1"
  if [ "$DRY_RUN" = true ]; then
    echo -e "  ${YELLOW}[DRY-RUN]${RESET} $msg"
    return 0
  fi
  if [ "$YES" = true ]; then
    return 0
  fi
  echo -en "  ${BOLD}${msg} [y/N]${RESET} "
  local answer
  read -r answer
  [[ "$answer" =~ ^[yY]$ ]]
}

# ── 清理模块 ─────────────────────────────────────────

cleanup_mysql() {
  info "清理 MySQL (${MYSQL_DB})..."

  if [ "$DRY_RUN" = true ]; then
    echo -e "  ${YELLOW}[DRY-RUN]${RESET} docker exec ${MYSQL_CONTAINER} mysql -u${MYSQL_USER} -p*** -e 'DROP DATABASE ${MYSQL_DB}; CREATE DATABASE ${MYSQL_DB};'"
    return 0
  fi

  if ! confirm "即将清空 MySQL ${MYSQL_DB} 数据库（DROP + CREATE），确认？"; then
    warn "跳过 MySQL 清理"
    return 0
  fi

  check_docker && check_mysql_container || return 1

  docker exec "$MYSQL_CONTAINER" \
    mysql -u"${MYSQL_USER}" -p"${MYSQL_PASSWORD}" \
    -e "DROP DATABASE IF EXISTS \`${MYSQL_DB}\`; CREATE DATABASE \`${MYSQL_DB}\` CHARACTER SET ${MYSQL_CHARSET:-utf8mb4} COLLATE ${MYSQL_COLLATE:-utf8mb4_unicode_ci};" \
    2>/dev/null

  ok "MySQL ${MYSQL_DB} 已重建为空库"
}

cleanup_redis() {
  info "清理 Redis..."

  if [ "$DRY_RUN" = true ]; then
    echo -e "  ${YELLOW}[DRY-RUN]${RESET} docker exec ${REDIS_CONTAINER} redis-cli FLUSHDB"
    return 0
  fi

  if ! confirm "即将清空 Redis 当前数据库所有 key，确认？"; then
    warn "跳过 Redis 清理"
    return 0
  fi

  check_docker && check_redis_container || return 1

  docker exec "$REDIS_CONTAINER" redis-cli FLUSHDB

  ok "Redis 已清空"
}

cleanup_logs() {
  local logs_dir="${PROJECT_ROOT}/agentend/logs"
  info "清理 agentend/logs/..."

  if [ "$DRY_RUN" = true ]; then
    if [ -d "$logs_dir" ]; then
      echo -e "  ${YELLOW}[DRY-RUN]${RESET} 将删除以下文件："
      ls -1 "$logs_dir" 2>/dev/null | sed 's/^/    /'
    else
      echo -e "  ${YELLOW}[DRY-RUN]${RESET} ${logs_dir} 不存在，跳过"
    fi
    return 0
  fi

  if ! confirm "即将删除 agentend/logs/ 下所有文件，确认？"; then
    warn "跳过日志清理"
    return 0
  fi

  if [ -d "$logs_dir" ]; then
    rm -f "$logs_dir"/*
    ok "agentend/logs/ 已清空"
  else
    warn "agentend/logs/ 目录不存在，跳过"
  fi
}

cleanup_repo() {
  local ws_json="${PROJECT_ROOT}/agentend/logs/workspaces.json"
  info "清理 repo（git worktree + branch + 目录）..."

  # ── 确定要清理的 repo 列表 ──
  local repo_paths=()

  # 优先使用命令行传入的 repo 路径
  if [ -n "$REPO_ARG" ]; then
    repo_paths+=("$REPO_ARG")
  fi

  # 否则从 workspaces.json 提取 repo_path
  if [ ${#repo_paths[@]} -eq 0 ] && [ -f "$ws_json" ] && command -v python3 &>/dev/null; then
    while IFS= read -r rp; do
      [ -n "$rp" ] && repo_paths+=("$rp")
    done < <(python3 -c "
import json, sys
try:
  data = json.load(open('$ws_json'))
  if isinstance(data, dict):
    items = data.values()
  elif isinstance(data, list):
    items = data
  else:
    items = []
  seen = set()
  for ws in items:
    rp = ws.get('repo_path', '')
    if rp and rp not in seen:
      seen.add(rp)
      print(rp)
except Exception:
  pass
" 2>/dev/null)
  fi

  if [ ${#repo_paths[@]} -eq 0 ]; then
    warn "未指定 repo 路径且 workspaces.json 为空，无法清理"
    warn "用法: $(basename "$0") repo --repo /path/to/repo"
    return 0
  fi

  # ── dry-run 模式 ──
  if [ "$DRY_RUN" = true ]; then
    for rp in "${repo_paths[@]}"; do
      echo -e "  ${YELLOW}[DRY-RUN]${RESET} 仓库: ${rp}"
      echo -e "    git worktree remove --force (所有附加 worktree)"
      echo -e "    git branch -D (除 main 外所有分支)"
      echo -e "    git reset --hard origin/main"
      local wt_base
      wt_base="$(dirname "$rp")/worktrees"
      if [ -d "$wt_base" ]; then
        echo -e "    删除 ${wt_base}/"
      fi
    done
    return 0
  fi

  if ! confirm "即将清理所有 git worktree、分支和 worktrees 目录，确认？"; then
    warn "跳过 worktree 清理"
    return 0
  fi

  # ── 清理每个 repo 的 worktree 和分支 ──
  for rp in "${repo_paths[@]}"; do
    if [ ! -d "$rp/.git" ] && [ ! -f "$rp/.git" ]; then
      warn "仓库 ${rp} 不存在，跳过"
      continue
    fi

    info "清理仓库: ${rp}"

    # 1. 移除所有附加 worktree
    local wt_count
    wt_count=$(cd "$rp" && git worktree list | tail -n +2 | wc -l | tr -d ' ')
    if [ "$wt_count" -gt 0 ]; then
      (cd "$rp" && git worktree list | tail -n +2 | awk '{print $1}' | while read -r wt; do
        git worktree remove "$wt" --force 2>/dev/null || true
      done)
      ok "  已移除 ${wt_count} 个 worktree"
    else
      info "  无附加 worktree"
    fi

    # 2. 清理残留 worktree 元数据 + 删除除 main 外的所有分支
    (cd "$rp" && git worktree prune 2>/dev/null || true)
    local branches
    branches=$(cd "$rp" && git branch --format='%(refname:short)' | grep -v '^main$' || true)
    if [ -n "$branches" ]; then
      (cd "$rp" && git branch --format='%(refname:short)' | grep -v '^main$' | xargs git branch -D 2>/dev/null || true)
      ok "  已清理分支: $(echo "$branches" | tr '\n' ' ')"
    else
      info "  无需清理分支"
    fi

    # 3. main 回到 origin/main 状态
    if (cd "$rp" && git remote get-url origin &>/dev/null); then
      (cd "$rp" && git reset --hard origin/main 2>/dev/null || true)
      ok "  main 已重置到 origin/main"
    fi

    # 4. 清理 repo 旁的 worktrees 目录
    local wt_base
    wt_base="$(dirname "$rp")/worktrees"
    if [ -d "$wt_base" ]; then
      rm -rf "$wt_base"
      ok "  已删除 ${wt_base}"
    fi
  done

  ok "worktree 清理完成"
}

cleanup_cache() {
  local agentend_dir="${PROJECT_ROOT}/agentend"
  info "清理测试缓存与测试目录..."

  if [ "$DRY_RUN" = true ]; then
    local items=()
    [ -d "${agentend_dir}/.pytest_cache" ] && items+=("${agentend_dir}/.pytest_cache/")
    [ -d "${agentend_dir}/tests" ] && items+=("${agentend_dir}/tests/")
    if [ ${#items[@]} -gt 0 ]; then
      echo -e "  ${YELLOW}[DRY-RUN]${RESET} 将删除："
      printf '    %s\n' "${items[@]}"
    else
      echo -e "  ${YELLOW}[DRY-RUN]${RESET} 无需清理"
    fi
    return 0
  fi

  if ! confirm "即将删除 .pytest_cache 和 agentend/tests/，确认？"; then
    warn "跳过缓存清理"
    return 0
  fi

  rm -rf "${agentend_dir}/.pytest_cache"
  rm -rf "${agentend_dir}/tests"

  ok ".pytest_cache 和 tests/ 已清理"
}

# ── all 子命令 ────────────────────────────────────────

cleanup_all() {
  info "========== 全量清理 =========="
  echo ""
  cleanup_mysql     || true
  cleanup_redis     || true
  cleanup_repo || true
  cleanup_logs      || true
  cleanup_cache   || true
  echo ""
  ok "========== 清理完成 =========="
}

# ── 帮助信息 ─────────────────────────────────────────

show_help() {
  cat <<EOF
${BOLD}用法:${RESET}  $(basename "$0") <子命令> [选项]

${BOLD}子命令:${RESET}
  all         执行全部清理（mysql + redis + logs + repo + cache）
  mysql       清理 Docker MySQL 数据库（DROP + CREATE agenthub）
  redis       清理 Docker Redis（FLUSHDB）
  logs        清理 agentend/logs/ 目录
  repo   清理 git worktree + 分支 + repo 目录（自动读取 workspaces.json）
  cache       清理 .pytest_cache 和 agentend/tests/
  help        显示此帮助信息

${BOLD}选项:${RESET}
  --yes       跳过所有确认提示（CI 模式）
  --dry-run   仅预览将被清理的内容，不实际执行
  --repo PATH 指定 repo 路径（repo 子命令专用，覆盖 workspaces.json）
  --help      显示此帮助信息

${BOLD}环境变量:${RESET}
  MYSQL_CONTAINER  MySQL Docker 容器名（默认: mysql）
  REDIS_CONTAINER  Redis Docker 容器名（默认: redis）

${BOLD}示例:${RESET}
  $(basename "$0") all --dry-run     # 预览全量清理
  $(basename "$0") all --yes         # 全量清理（跳过确认）
  $(basename "$0") mysql             # 仅清理 MySQL（交互确认）
  $(basename "$0") repo --repo /path/to/repo --yes  # 清理指定 repo 的 worktree
  MYSQL_CONTAINER=db $(basename "$0") mysql --yes   # 指定容器名
EOF
}

# ── 参数解析 ─────────────────────────────────────────

parse_args() {
  if [ $# -eq 0 ]; then
    show_help
    exit 0
  fi

  while [ $# -gt 0 ]; do
    case "$1" in
      --yes|-y)    YES=true; shift ;;
      --dry-run)   DRY_RUN=true; shift ;;
      --repo)      REPO_ARG="${2:-}"; [ -z "$REPO_ARG" ] && { err "--repo 需要路径参数"; exit 1; }; shift 2 ;;
      --help|-h)   show_help; exit 0 ;;
      all|mysql|redis|logs|repo|cache|help)
        COMMAND="$1"; shift ;;
      *)
        err "未知参数: $1"
        show_help
        exit 1
        ;;
    esac
  done

  if [ -z "$COMMAND" ]; then
    err "请指定子命令"
    show_help
    exit 1
  fi
}

# ── 入口 ─────────────────────────────────────────────

main() {
  parse_args "$@"

  case "$COMMAND" in
    all)       cleanup_all ;;
    mysql)     cleanup_mysql ;;
    redis)     cleanup_redis ;;
    logs)      cleanup_logs ;;
    repo) cleanup_repo ;;
    cache)     cleanup_cache ;;
    help)      show_help ;;
    *)
      err "未知子命令: $COMMAND"
      show_help
      exit 1
      ;;
  esac
}

main "$@"
