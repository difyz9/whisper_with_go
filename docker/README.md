# Docker 部署指南

本目录包含 Whisper API Server 的 Docker 部署文件和工具。

## 📋 目录结构

```
.
├── Dockerfile              # CPU 版本镜像
├── Dockerfile.gpu          # GPU 版本镜像
├── docker-compose.yml      # CPU 版本编排文件
├── docker-compose.gpu.yml  # GPU 版本编排文件
├── .dockerignore          # Docker 忽略文件
└── scripts/
    └── docker.sh          # Docker 管理脚本
```

## 🚀 快速开始

### 前置要求

1. **Docker** (版本 20.10+)
   ```bash
   docker --version
   ```

2. **Docker Compose** (版本 2.0+)
   ```bash
   docker compose version
   ```

3. **GPU 支持**（仅 GPU 版本需要）
   - NVIDIA Docker Runtime
   - CUDA 12.3+
   - 安装 [NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

### 一键运行

#### CPU 版本

```bash
# 使用管理脚本（推荐）
./scripts/docker.sh build && ./scripts/docker.sh up

# 或使用 docker-compose
docker compose up -d --build
```

#### GPU 版本

```bash
# 使用管理脚本
./scripts/docker.sh build-gpu && ./scripts/docker.sh up-gpu

# 或使用 docker-compose
docker compose -f docker-compose.gpu.yml up -d --build
```

服务启动后访问：`http://localhost:8080`

## 🛠️ 管理脚本使用

### 基本命令

```bash
# 查看帮助
./scripts/docker.sh help

# 构建镜像
./scripts/docker.sh build        # CPU 版本
./scripts/docker.sh build-gpu    # GPU 版本

# 启动服务
./scripts/docker.sh up           # CPU 版本
./scripts/docker.sh up-gpu       # GPU 版本

# 停止服务
./scripts/docker.sh down

# 重启服务
./scripts/docker.sh restart

# 查看日志
./scripts/docker.sh logs

# 查看运行状态
./scripts/docker.sh ps

# 进入容器
./scripts/docker.sh shell

# 测试 API
./scripts/docker.sh test

# 下载模型
./scripts/docker.sh download-model

# 清理所有
./scripts/docker.sh clean
```

## 📦 镜像说明

### 基础环境镜像（推荐复用）

仓库提供了一个单独的基础环境镜像定义 [Dockerfile.base](../Dockerfile.base)，用于预构建以下运行环境：

- Go 1.23.4
- FFmpeg
- `whisper.cpp`
- CGO 编译所需头文件和动态库
- 多平台支持：`linux/amd64`、`linux/arm64`

这个镜像适合下面这种场景：

- CI 先构建并发布基础镜像到 Docker Hub
- 业务机器只需要 `docker pull`
- 启动时把 Go 项目目录挂载进容器，直接 `go run` 或 `go build`

示例：

```bash
docker pull difyz9/whisper-go-base:latest

docker run --rm -it \
  -p 8080:8080 \
  -v $(pwd):/workspace \
  -v $(pwd)/models:/workspace/models \
  -v $(pwd)/uploads:/workspace/uploads \
  -v $(pwd)/outputs:/workspace/outputs \
  -w /workspace \
  difyz9/whisper-go-base:latest \
  bash -lc 'go mod download && go run cmd/server/main.go'
```

对应的 GitHub Actions 工作流是 [build-base-image.yml](../.github/workflows/build-base-image.yml)。

### 配置 GitHub Actions

首次使用时，按下面步骤配置：

1. 将仓库推送到 GitHub。
2. 确认仓库启用了 GitHub Actions。
3. 在 `Settings -> Secrets and variables -> Actions` 中添加 `DOCKERHUB_USERNAME` 和 `DOCKERHUB_TOKEN`。
4. 给仓库打 tag，例如 `v1.0.0`，再 push tag 触发镜像构建。

工作流默认发布以下标签：

- `difyz9/whisper-go-base:v1.0.0`
- `difyz9/whisper-go-base:1.0.0`
- `difyz9/whisper-go-base:latest`

手动触发时支持传入：

- `tag`: 自定义版本 tag
- `platforms`: 例如 `linux/amd64,linux/arm64`
- `push_latest`: 是否推送 `latest`

仓库根目录默认的 [docker-compose.yml](../docker-compose.yml) 已经可以直接一键运行当前项目：

```bash
docker compose up -d
```

它会读取 [../.env](../.env) 中的 `WHISPER_BASE_IMAGE`，默认值已经指向当前已发布的基础镜像。

如果你想显式指定镜像地址，仓库也提供了一个直接消费该基础镜像的编排文件 [docker-compose.base.yml](../docker-compose.base.yml)，启动前只需要指定镜像地址：

```bash
WHISPER_BASE_IMAGE=difyz9/whisper-go-base:latest \
docker compose -f docker-compose.base.yml up
```

这个 compose 文件要求显式设置 `WHISPER_BASE_IMAGE`，避免默认值写成占位符后造成拉取失败。

如果你在 Windows PowerShell 下使用，推荐直接运行 [scripts/compose-base.ps1](../scripts/compose-base.ps1)：

```powershell
.\scripts\compose-base.ps1 -Action up -Image difyz9/whisper-go-base:latest
```

它会完成这些步骤：

- 检查 Docker / Docker Compose
- 自动创建 `models`、`uploads`、`outputs` 目录
- 拉取指定基础镜像
- 调用 `docker compose -f docker-compose.base.yml up -d`

Linux / macOS 下可以使用 [scripts/compose-base.sh](../scripts/compose-base.sh)：

```bash
chmod +x ./scripts/compose-base.sh
./scripts/compose-base.sh --action up --image difyz9/whisper-go-base:latest
```

### CPU 版本镜像

- **基础镜像**: Ubuntu 22.04
- **大小**: ~500MB（不含模型）
- **特性**:
  - 多阶段构建优化
  - FFmpeg 支持
  - 非 root 用户运行
  - 健康检查
  - 自动重启

### GPU 版本镜像

- **基础镜像**: NVIDIA CUDA 12.3.0
- **大小**: ~2GB（不含模型）
- **特性**:
  - CUDA 加速支持
  - 多 GPU 架构支持 (75, 80, 86, 89, 90)
  - 自动 GPU 检测
  - 其他同 CPU 版本

## ⚙️ 配置

### 环境变量

在 `docker-compose.yml` 或 `docker-compose.gpu.yml` 中配置：

```yaml
environment:
  - SERVER_PORT=8080
  - GIN_MODE=release
  - WHISPER_MODEL_PATH=/app/models/ggml-base.bin
  - WHISPER_THREADS=4
  - MAX_FILE_SIZE=104857600
  - UPLOAD_DIR=/app/uploads
  - OUTPUT_DIR=/app/outputs
```

### 卷挂载

```yaml
volumes:
  - ./models:/app/models:ro      # 模型文件（只读）
  - ./uploads:/app/uploads       # 上传目录
  - ./outputs:/app/outputs       # 输出目录
```

### 端口映射

```yaml
ports:
  - "8080:8080"  # 可修改为其他端口，如 "8888:8080"
```

## 🔧 高级配置

### 自定义 Dockerfile

如果需要自定义镜像，可以修改 `Dockerfile` 或创建新的：

```dockerfile
FROM whisper-api:latest

# 添加自定义配置
ENV CUSTOM_VAR=value

# 安装额外依赖
RUN apt-get update && apt-get install -y package-name
```

### 使用预构建镜像

如果已有预构建镜像：

```yaml
# docker-compose.yml
services:
  whisper-api:
    image: your-registry/whisper-api:latest
    # ... 其他配置
```

### 多容器部署

```yaml
# docker-compose.yml
services:
  whisper-api-1:
    build: .
    ports:
      - "8081:8080"
    # ...
  
  whisper-api-2:
    build: .
    ports:
      - "8082:8080"
    # ...
  
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    # 配置负载均衡
```

## 🧪 测试

### 健康检查

```bash
curl http://localhost:8080/health
```

### API 测试

```bash
# 使用测试脚本
./test_api.sh

# 或手动测试
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test_data/test.mp3" \
  -F "language=auto" \
  -F "output_type=json"
```

### 性能测试

```bash
# 使用 ab (Apache Bench)
ab -n 100 -c 10 http://localhost:8080/health

# 使用 wrk
wrk -t12 -c400 -d30s http://localhost:8080/health
```

## 📊 监控

### 查看日志

```bash
# 实时日志
./scripts/docker.sh logs

# 或使用 docker
docker logs -f whisper-api-cpu
```

### 资源使用

```bash
# 查看容器资源使用
docker stats whisper-api-cpu

# 使用 docker-compose
docker compose stats
```

### 健康状态

```bash
# 查看健康状态
docker inspect --format='{{.State.Health.Status}}' whisper-api-cpu
```

## 🚨 故障排除

### 常见问题

#### 1. 构建失败

```bash
# 清理缓存后重新构建
docker builder prune
./scripts/docker.sh build
```

#### 2. GPU 不可用

```bash
# 检查 NVIDIA Docker 运行时
docker run --rm --gpus all nvidia/cuda:12.3.0-base-ubuntu22.04 nvidia-smi

# 检查 GPU 可见性
docker run --rm --gpus all whisper-api-gpu:latest nvidia-smi
```

#### 3. 端口占用

```bash
# 查看端口占用
lsof -i :8080

# 修改端口映射
# 在 docker-compose.yml 中修改为 "8888:8080"
```

#### 4. 权限问题

```bash
# 确保目录权限正确
chmod -R 755 uploads outputs
chown -R 1000:1000 uploads outputs
```

#### 5. 模型文件缺失

```bash
# 下载模型
./scripts/docker.sh download-model

# 或手动下载
bash download_model.sh
```

### 日志调试

```bash
# 查看完整日志
docker logs whisper-api-cpu --tail 100

# 调试模式运行
docker compose up  # 不使用 -d 参数
```

## 🔐 安全建议

1. **使用非 root 用户**: 已默认配置
2. **限制资源使用**:
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '2'
         memory: 4G
   ```
3. **使用secrets管理敏感信息**
4. **定期更新基础镜像**
5. **扫描镜像漏洞**:
   ```bash
   docker scan whisper-api:latest
   ```

## 📈 性能优化

### CPU 版本

1. **增加线程数**: `WHISPER_THREADS=8`
2. **使用更小的模型**: tiny/base
3. **限制并发请求数**
4. **启用 HTTP/2 和压缩**

### GPU 版本

1. **确保 CUDA 版本匹配**
2. **使用更大的模型**: medium/large
3. **批处理请求**
4. **监控 GPU 使用率**: `nvidia-smi`

## 🌍 生产部署

### Kubernetes 部署

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: whisper-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: whisper-api
  template:
    metadata:
      labels:
        app: whisper-api
    spec:
      containers:
      - name: whisper-api
        image: whisper-api:latest
        ports:
        - containerPort: 8080
```

### Docker Swarm

```bash
# 初始化 Swarm
docker swarm init

# 部署服务
docker stack deploy -c docker-compose.yml whisper-stack
```

## 📚 相关资源

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [NVIDIA Container Toolkit](https://github.com/NVIDIA/nvidia-docker)
- [Whisper.cpp](https://github.com/ggerganov/whisper.cpp)

## 💡 提示

1. 首次构建可能需要较长时间（10-20分钟）
2. GPU 版本需要下载 ~2GB 的 CUDA 镜像
3. 建议使用 SSD 存储提高性能
4. 生产环境建议使用 nginx 反向代理
5. 定期清理旧的容器和镜像

---

**有问题？** 查看主项目的 [README.md](../README.md) 或提交 Issue。
