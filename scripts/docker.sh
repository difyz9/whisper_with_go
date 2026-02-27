#!/bin/bash

# Whisper API Docker 管理脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 显示帮助信息
show_help() {
    cat << EOF
Whisper API Docker 管理工具

用法: $0 [命令] [选项]

命令:
  build         构建 Docker 镜像
  build-gpu     构建 GPU 版本镜像
  up            启动服务（CPU 版本）
  up-gpu        启动服务（GPU 版本）
  down          停止服务
  restart       重启服务
  logs          查看日志
  ps            查看运行状态
  shell         进入容器 shell
  test          测试 API
  clean         清理所有容器和镜像
  download-model 下载 Whisper 模型

选项:
  -h, --help    显示帮助信息

示例:
  $0 build              # 构建 CPU 版本镜像
  $0 build-gpu          # 构建 GPU 版本镜像
  $0 up                 # 启动 CPU 版本服务
  $0 up-gpu             # 启动 GPU 版本服务
  $0 logs               # 查看日志
  $0 test               # 测试 API

EOF
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        printf "${RED}错误: Docker 未安装${NC}\n"
        printf "请访问 https://docs.docker.com/get-docker/ 安装 Docker\n"
        exit 1
    fi
}

# 检查 Docker Compose 是否安装
check_docker_compose() {
    if ! docker compose version &> /dev/null; then
        printf "${RED}错误: Docker Compose 未安装${NC}\n"
        exit 1
    fi
}

# 检查模型文件
check_model() {
    if [ ! -f "models/ggml-base.bin" ]; then
        printf "${YELLOW}警告: 模型文件不存在${NC}\n"
        printf "正在下载模型...\n"
        download_model
    fi
}

# 下载模型
download_model() {
    mkdir -p models
    if [ -f "download_model.sh" ]; then
        bash download_model.sh
    else
        printf "${YELLOW}请手动下载模型文件到 models 目录${NC}\n"
        printf "或使用: bash download_model.sh\n"
    fi
}

# 构建镜像
build() {
    printf "${GREEN}正在构建 CPU 版本镜像...${NC}\n"
    docker build -t whisper-api:latest .
    printf "${GREEN}✓ 构建完成${NC}\n"
}

# 构建 GPU 镜像
build_gpu() {
    printf "${GREEN}正在构建 GPU 版本镜像...${NC}\n"
    docker build -f Dockerfile.gpu -t whisper-api-gpu:latest .
    printf "${GREEN}✓ 构建完成${NC}\n"
}

# 启动服务
up() {
    check_model
    printf "${GREEN}正在启动 CPU 版本服务...${NC}\n"
    docker compose up -d
    printf "${GREEN}✓ 服务已启动${NC}\n"
    printf "访问: http://localhost:8080\n"
    printf "健康检查: http://localhost:8080/health\n"
}

# 启动 GPU 服务
up_gpu() {
    check_model
    printf "${GREEN}正在启动 GPU 版本服务...${NC}\n"
    docker compose -f docker-compose.gpu.yml up -d
    printf "${GREEN}✓ 服务已启动${NC}\n"
    printf "访问: http://localhost:8080\n"
    printf "健康检查: http://localhost:8080/health\n"
}

# 停止服务
down() {
    printf "${YELLOW}正在停止服务...${NC}\n"
    docker compose down 2>/dev/null || true
    docker compose -f docker-compose.gpu.yml down 2>/dev/null || true
    printf "${GREEN}✓ 服务已停止${NC}\n"
}

# 重启服务
restart() {
    down
    sleep 2
    up
}

# 查看日志
logs() {
    if docker ps --format '{{.Names}}' | grep -q whisper-api-gpu; then
        docker compose -f docker-compose.gpu.yml logs -f
    elif docker ps --format '{{.Names}}' | grep -q whisper-api-cpu; then
        docker compose logs -f
    else
        printf "${RED}没有运行中的服务${NC}\n"
        exit 1
    fi
}

# 查看状态
ps() {
    printf "${GREEN}运行中的容器:${NC}\n"
    docker ps --filter "name=whisper-api"
}

# 进入容器
shell() {
    CONTAINER=$(docker ps --format '{{.Names}}' | grep whisper-api | head -n 1)
    if [ -z "$CONTAINER" ]; then
        printf "${RED}没有运行中的容器${NC}\n"
        exit 1
    fi
    printf "${GREEN}进入容器: $CONTAINER${NC}\n"
    docker exec -it "$CONTAINER" /bin/bash
}

# 测试 API
test_api() {
    printf "${GREEN}测试 API...${NC}\n"
    
    # 健康检查
    printf "\n1. 健康检查:\n"
    curl -s http://localhost:8080/health | jq . || printf "${RED}失败${NC}\n"
    
    # API 信息
    printf "\n2. API 信息:\n"
    curl -s http://localhost:8080/ | jq . || printf "${RED}失败${NC}\n"
    
    printf "\n${GREEN}测试完成${NC}\n"
}

# 清理
clean() {
    printf "${YELLOW}正在清理容器和镜像...${NC}\n"
    down
    docker rmi whisper-api:latest 2>/dev/null || true
    docker rmi whisper-api-gpu:latest 2>/dev/null || true
    printf "${GREEN}✓ 清理完成${NC}\n"
}

# 主逻辑
main() {
    check_docker
    check_docker_compose
    
    case "${1:-}" in
        build)
            build
            ;;
        build-gpu)
            build_gpu
            ;;
        up)
            up
            ;;
        up-gpu)
            up_gpu
            ;;
        down)
            down
            ;;
        restart)
            restart
            ;;
        logs)
            logs
            ;;
        ps)
            ps
            ;;
        shell)
            shell
            ;;
        test)
            test_api
            ;;
        clean)
            clean
            ;;
        download-model)
            download_model
            ;;
        -h|--help|help)
            show_help
            ;;
        *)
            printf "${RED}错误: 未知命令 '${1:-}'${NC}\n\n"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
