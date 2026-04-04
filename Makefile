.PHONY: help build run test clean fmt lint swagger

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
APP_NAME=whisper-server
BUILD_DIR=bin
CMD_DIR=cmd/server
MAIN_FILE=$(CMD_DIR)/main.go

# Go 相关变量
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

help: ## 显示帮助信息
	@echo "可用的 make 命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 构建应用程序
	@echo "正在构建..."
	@mkdir -p $(BUILD_DIR)
	@echo "设置 CGO 环境变量..."
	@CGO_ENABLED=1 CGO_LDFLAGS="-L/usr/local/lib" CGO_CFLAGS="-I/usr/local/include" \
		$(GOBUILD) -o $(GOBIN)/$(APP_NAME) $(MAIN_FILE)
	@echo "构建完成: $(GOBIN)/$(APP_NAME)"
	@echo ""
	@echo "运行服务器:"
	@echo "  DYLD_LIBRARY_PATH=/usr/local/lib ./$(GOBIN)/$(APP_NAME)"

run: ## 运行应用程序
	@echo "正在启动服务器..."
	@CGO_ENABLED=1 CGO_LDFLAGS="-L/usr/local/lib" CGO_CFLAGS="-I/usr/local/include" \
		DYLD_LIBRARY_PATH=/usr/local/lib $(GORUN) $(MAIN_FILE)

dev: ## 开发模式运行（使用 air 热重载，需要先安装 air）
	@which air > /dev/null || (echo "请先安装 air: go install github.com/cosmtrek/air@latest" && exit 1)
	@air

test: ## 运行测试
	@echo "正在运行测试..."
	@$(GOTEST) -v ./...

test-cover: ## 运行测试并生成覆盖率报告
	@echo "正在运行测试并生成覆盖率..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

clean: ## 清理构建文件
	@echo "正在清理..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf uploads/*
	@rm -rf outputs/*
	@rm -f coverage.out coverage.html
	@echo "清理完成"

fmt: ## 格式化代码
	@echo "正在格式化代码..."
	@$(GOFMT) ./...
	@echo "代码格式化完成"

lint: ## 代码检查（需要安装 golangci-lint）
	@which golangci-lint > /dev/null || (echo "请先安装 golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	@echo "正在进行代码检查..."
	@golangci-lint run ./...

tidy: ## 整理依赖
	@echo "正在整理依赖..."
	@$(GOMOD) tidy
	@echo "依赖整理完成"

download: ## 下载依赖
	@echo "正在下载依赖..."
	@$(GOMOD) download
	@echo "依赖下载完成"

install-tools: ## 安装开发工具
	@echo "正在安装开发工具..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/swaggo/swag/cmd/swag@v1.16.4
	@echo "开发工具安装完成"

swagger: ## 生成 Swagger 文档
	@echo "正在生成 Swagger 文档..."
	@go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g cmd/server/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger 文档已生成: docs/swagger.json, docs/swagger.yaml"

setup: ## 项目初始设置
	@echo "正在进行项目设置..."
	@mkdir -p uploads outputs models
	@cp -n .env.example .env 2>/dev/null || true
	@$(GOMOD) tidy
	@echo "项目设置完成"
	@echo "请确保:"
	@echo "  1. 已安装 FFmpeg"
	@echo "  2. 已下载 Whisper 模型（运行 bash download_model.sh）"
	@echo "  3. 已配置 .env 文件"

model: ## 下载 Whisper 模型
	@bash download_model.sh

# Docker 相关命令
docker-build: ## 构建 Docker 镜像（CPU 版本）
	@echo "正在构建 Docker 镜像（CPU 版本）..."
	@docker build -t $(APP_NAME):latest .
	@echo "Docker 镜像构建完成"

docker-build-gpu: ## 构建 Docker 镜像（GPU 版本）
	@echo "正在构建 Docker 镜像（GPU 版本）..."
	@docker build -f Dockerfile.gpu -t $(APP_NAME)-gpu:latest .
	@echo "Docker GPU 镜像构建完成"

docker-up: ## 启动 Docker 服务（CPU 版本）
	@./scripts/docker.sh up

docker-up-gpu: ## 启动 Docker 服务（GPU 版本）
	@./scripts/docker.sh up-gpu

docker-down: ## 停止 Docker 服务
	@./scripts/docker.sh down

docker-logs: ## 查看 Docker 日志
	@./scripts/docker.sh logs

docker-ps: ## 查看 Docker 容器状态
	@./scripts/docker.sh ps

docker-shell: ## 进入 Docker 容器
	@./scripts/docker.sh shell

docker-test: ## 测试 Docker API
	@./scripts/docker.sh test

docker-clean: ## 清理 Docker 容器和镜像
	@./scripts/docker.sh clean

docker-help: ## Docker 帮助信息
	@./scripts/docker.sh help
