#!/usr/bin/env bash

set -euo pipefail

ACTION="up"
IMAGE="${WHISPER_BASE_IMAGE:-}"
NO_PULL="false"
FOREGROUND="false"

info() {
    printf '[INFO] %s\n' "$1"
}

success() {
    printf '[OK] %s\n' "$1"
}

warn() {
    printf '[WARN] %s\n' "$1"
}

fail() {
    printf '[ERROR] %s\n' "$1" >&2
    exit 1
}

show_help() {
    cat <<'EOF'
Usage: ./scripts/compose-base.sh [options]

Options:
  -a, --action <up|down|restart|logs|ps|pull>
  -i, --image <ghcr-image>
      --no-pull
  -f, --foreground
  -h, --help

Examples:
    ./scripts/compose-base.sh --action up --image difyz9/whisper-go-base:latest
  ./scripts/compose-base.sh --action logs
    WHISPER_BASE_IMAGE=difyz9/whisper-go-base:latest ./scripts/compose-base.sh
EOF
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -a|--action)
                ACTION="$2"
                shift 2
                ;;
            -i|--image)
                IMAGE="$2"
                shift 2
                ;;
            --no-pull)
                NO_PULL="true"
                shift
                ;;
            -f|--foreground)
                FOREGROUND="true"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                fail "Unknown argument: $1"
                ;;
        esac
    done
}

ensure_docker() {
    command -v docker >/dev/null 2>&1 || fail 'Docker 未安装或未加入 PATH。'
    docker compose version >/dev/null
}

ensure_project_directories() {
    for path in models uploads outputs; do
        if [[ ! -d "$path" ]]; then
            mkdir -p "$path"
            info "已创建目录: $path"
        fi
    done
}

warn_if_model_missing() {
    if [[ ! -f models/ggml-base.bin ]]; then
        warn '未找到 models/ggml-base.bin，服务可以启动，但转录请求会因缺少模型失败。'
    fi
}

require_image() {
    [[ -n "$IMAGE" ]] || fail '请通过 --image 参数或 WHISPER_BASE_IMAGE 环境变量提供基础镜像地址，例如 difyz9/whisper-go-base:latest。'
}

compose() {
    docker compose -f docker-compose.base.yml "$@"
}

parse_args "$@"
ensure_docker
ensure_project_directories
warn_if_model_missing

case "$ACTION" in
    pull)
        require_image
        export WHISPER_BASE_IMAGE="$IMAGE"
        info "拉取基础镜像: $IMAGE"
        docker pull "$IMAGE"
        success '镜像拉取完成。'
        ;;
    up)
        require_image
        export WHISPER_BASE_IMAGE="$IMAGE"
        if [[ "$NO_PULL" != "true" ]]; then
            info "拉取基础镜像: $IMAGE"
            docker pull "$IMAGE"
        fi
        if [[ "$FOREGROUND" == "true" ]]; then
            info '以前台模式启动基础镜像服务。'
            compose up
        else
            info '以后台模式启动基础镜像服务。'
            compose up -d
            success '服务已启动。'
            printf '访问: http://localhost:8080\n'
            printf '健康检查: http://localhost:8080/health\n'
        fi
        ;;
    down)
        info '停止基础镜像服务。'
        compose down
        success '服务已停止。'
        ;;
    restart)
        require_image
        export WHISPER_BASE_IMAGE="$IMAGE"
        info '重启基础镜像服务。'
        compose down
        if [[ "$NO_PULL" != "true" ]]; then
            info "拉取基础镜像: $IMAGE"
            docker pull "$IMAGE"
        fi
        compose up -d
        success '服务已重启。'
        ;;
    logs)
        info '查看服务日志。'
        compose logs -f
        ;;
    ps)
        info '查看服务状态。'
        compose ps
        ;;
    *)
        fail "Unsupported action: $ACTION"
        ;;
esac