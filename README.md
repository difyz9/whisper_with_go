# Whisper API Server

<div align="center">

**基于 Gin 框架的 Whisper 语音转文字 API 服务**

[![Go Version](https://img.shields.io/badge/Go-1.23.4+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[功能特性](#✨-功能特性) • [快速开始](#🚀-快速开始) • [API 文档](#📡-api-文档) • [项目架构](#🏗️-项目架构) • [配置说明](#⚙️-配置说明)

</div>

---

## 📖 简介

本项目是一个生产就绪的 Whisper 语音转文字 API 服务，采用标准的 Go Web 应用架构设计。将 OpenAI 的 Whisper 模型封装为 RESTful API，支持多种音频格式和输出格式。

### 从 CLI 到 Web API 的演进

本项目经过完整的架构重构，从简单的命令行工具升级为专业的 Web API 服务。旧版 CLI 工具代码保留在 `legacy/` 目录中供参考。

## ✨ 功能特性

- 🎯 **RESTful API** - 标准的 HTTP API 接口
- 📁 **多格式支持** - 支持 MP3, WAV, M4A, AAC, FLAC, OGG
- 📝 **多种输出** - JSON, SRT 字幕, TXT 文本
- 🌍 **自动语言检测** - 支持多语言自动识别
- 🔄 **异步转录任务** - 上传后立即返回任务 ID，后台处理转录
- 📊 **结构化日志** - 完整的请求/响应日志
- 🛡️ **错误恢复** - 自动错误恢复机制
- 🌐 **CORS 支持** - 跨域请求支持
- 📘 **Swagger 文档** - 内置 OpenAPI/Swagger UI
- ⚙️ **环境配置** - 灵活的配置管理
- 📚 **完整文档** - 详细的使用和架构文档

## 🏗️ 项目结构

```
whisper_with_go/
├── cmd/
│   └── server/            # 应用入口
│       └── main.go
├── internal/              # 内部私有代码
│   ├── handler/          # HTTP 处理器层
│   ├── service/          # 业务逻辑层
│   ├── model/            # 数据模型层
│   ├── middleware/       # 中间件（日志、CORS、恢复）
│   └── router/           # 路由配置
├── pkg/
│   └── utils/            # 可复用工具库
├── config/               # 配置管理
├── uploads/              # 上传文件目录
├── outputs/              # 输出文件目录
├── models/               # Whisper 模型文件
├── test_data/            # 测试数据
├── legacy/               # 旧版 CLI 工具（保留参考）
├── docs/                 # 文档目录
├── .env.example          # 环境变量示例
├── Makefile              # 构建脚本
└── README.md             # 本文件
```

### 架构分层

- **Handler 层**: HTTP 请求处理、参数验证
- **Service 层**: 核心业务逻辑、Whisper 调用
- **Model 层**: 数据结构定义
- **Utils 层**: 工具函数库

详细架构说明请查看 [ARCHITECTURE.md](docs/ARCHITECTURE.md)

## 🚀 快速开始

### 前置要求

1. **Go 1.23.4+**
   ```bash
   go version
   ```

2. **FFmpeg** - 用于音频格式转换
   ```bash
   # macOS
   brew install ffmpeg
   
   # Ubuntu/Debian
   sudo apt-get install ffmpeg
   ```

3. **Whisper 模型**
   ```bash
   # 下载 base 模型（推荐）
   bash download_model.sh
   ```

### 快速启动

#### 🐳 Docker 方式（推荐）

最简单的启动方式，无需配置环境：

```bash
# 1. 下载模型
bash download_model.sh

# 2. 一键启动当前项目（默认使用预构建基础镜像）
docker compose up -d

# 或启动 GPU 版本（需要 NVIDIA GPU）
./scripts/docker.sh build-gpu && ./scripts/docker.sh up-gpu
```

详细说明：[Docker 快速开始](DOCKER_QUICKSTART.md) | [Docker 部署指南](docker/README.md)

#### 拉取预构建基础镜像运行

如果你只想复用已经构建好的 Whisper 基础环境，而不是每次都重新编译 `whisper.cpp`，可以直接拉取 GitHub Actions 发布到 Docker Hub 的基础镜像，再把当前 Go 项目挂载进去运行。

```bash
# 拉取基础镜像
docker pull difyz9/whisper-go-base:latest

# 在基础镜像中直接运行当前项目
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

这个基础镜像已经内置：

- Go 1.23.4
- FFmpeg
- `whisper.cpp` 及其动态库
- CGO 相关环境变量
- 多平台镜像清单：`linux/amd64`、`linux/arm64`

如果你要发布这个基础镜像，仓库内已经提供 GitHub Actions 工作流 [build-base-image.yml](.github/workflows/build-base-image.yml)。

GitHub Actions 首次使用时，按下面配置：

1. 把仓库推到 GitHub。
2. 确认 Actions 已启用。
3. 在仓库 `Settings -> Secrets and variables -> Actions` 中添加 `DOCKERHUB_USERNAME` 和 `DOCKERHUB_TOKEN`。
4. 在 Docker Hub 预先创建仓库 `difyz9/whisper-go-base`，或把 workflow 中的镜像名改成你的仓库名。
5. 给仓库打 tag，例如 `v1.0.0`，再 push tag 触发构建。

这个 workflow 默认会发布到：

- `difyz9/whisper-go-base:v1.0.0`
- `difyz9/whisper-go-base:1.0.0`
- `difyz9/whisper-go-base:latest`

如果你想手动触发，在 GitHub 的 `Actions -> docker-hub-base-release -> Run workflow` 中填写：

- `tag`: 例如 `v1.0.0`
- `platforms`: 例如 `linux/amd64,linux/arm64`
- `push_latest`: 是否同时推送 `latest`

现在根目录默认的 [docker-compose.yml](docker-compose.yml) 已经支持直接运行当前项目：

```bash
docker compose up -d
```

默认会读取 [\.env](.env) 中的 `WHISPER_BASE_IMAGE`，当前预设为：

```bash
difyz9/whisper-go-base:latest
```

如果你仍然希望显式使用基础镜像编排文件，也可以使用 [docker-compose.base.yml](docker-compose.base.yml)：

```bash
WHISPER_BASE_IMAGE=difyz9/whisper-go-base:latest \
docker compose -f docker-compose.base.yml up
```

注意：`docker-compose.base.yml` 现在要求你显式传入 `WHISPER_BASE_IMAGE`，这样可以避免误拉取占位镜像地址。

如果你在 Windows PowerShell 下想一键启动，可以直接使用脚本 [scripts/compose-base.ps1](scripts/compose-base.ps1)：

```powershell
.\scripts\compose-base.ps1 -Action up -Image difyz9/whisper-go-base:latest
```

常用命令：

```powershell
.\scripts\compose-base.ps1 -Action up -Image difyz9/whisper-go-base:latest
.\scripts\compose-base.ps1 -Action logs
.\scripts\compose-base.ps1 -Action ps
.\scripts\compose-base.ps1 -Action down
```

如果你在 Linux 或 macOS 下使用，可以直接运行 [scripts/compose-base.sh](scripts/compose-base.sh)：

```bash
chmod +x ./scripts/compose-base.sh
./scripts/compose-base.sh --action up --image difyz9/whisper-go-base:latest
```

常用命令：

```bash
./scripts/compose-base.sh --action up --image difyz9/whisper-go-base:latest
./scripts/compose-base.sh --action logs
./scripts/compose-base.sh --action ps
./scripts/compose-base.sh --action down
```

如果你还需要保留“本地构建 CPU 镜像再启动”的旧流程，现在对应编排文件是 [docker-compose.build.yml](docker-compose.build.yml)，管理脚本 [scripts/docker.sh](scripts/docker.sh) 已经切换到这个文件。

如果你要明确拉某个平台，可以这样：

```bash
docker pull --platform linux/amd64 difyz9/whisper-go-base:latest
docker pull --platform linux/arm64 difyz9/whisper-go-base:latest
```

#### 💻 本地运行方式

#### 方式 1: 使用启动脚本（推荐）

```bash
./start.sh
```

#### 方式 2: 使用 Makefile

```bash
# 查看所有命令
make help

# 项目初始化
make setup

# 运行服务器
make run
```

#### 方式 3: 手动运行

```bash
# 安装依赖
go mod tidy

# 启动服务器
CGO_ENABLED=1 \
CGO_LDFLAGS="-L/usr/local/lib" \
CGO_CFLAGS="-I/usr/local/include" \
DYLD_LIBRARY_PATH=/usr/local/lib \
go run cmd/server/main.go
```

服务器启动后将运行在 `http://localhost:8080`

Swagger UI 默认地址：`http://localhost:8080/swagger/index.html`

如果你修改了接口注解，可以重新生成文档：

```bash
make swagger
```

## 📡 API 文档

### Swagger UI

启动服务后可直接访问：

```text
http://localhost:8080/swagger/index.html
```

项目会同时生成以下 OpenAPI 文档文件：

- `docs/swagger.json`
- `docs/swagger.yaml`

### 端点总览

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/` | API 信息 |
| POST | `/api/v1/transcribe` | 创建异步转录任务 |
| GET | `/api/v1/tasks/:task_id` | 查询任务状态和结果 |
| GET | `/api/v1/download/:filename` | 下载输出文件 |

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

响应:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "time": "2026-02-27T10:00:00Z"
}
```

### 2. 创建转录任务

**请求**

```bash
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test_data/test.mp3" \
  -F "language=auto" \
  -F "output_type=json"
```

**参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | File | ✅ | 音频文件 |
| language | String | ❌ | 语言代码（auto/zh/en等），默认 auto |
| output_type | String | ❌ | 输出格式（json/srt/txt），默认 json |
| translate | Boolean | ❌ | 是否翻译为英文，默认 false |

**响应示例**

```json
{
  "success": true,
  "message": "任务已创建",
  "data": {
    "task_id": "task_1709011200000000000",
    "status": "pending",
    "status_url": "/api/v1/tasks/task_1709011200000000000"
  }
}
```

### 3. 查询任务状态

```bash
curl http://localhost:8080/api/v1/tasks/task_1709011200000000000
```

任务处理中示例：

```json
{
  "success": true,
  "message": "查询成功",
  "data": {
    "task_id": "task_1709011200000000000",
    "status": "processing",
    "created_at": "2026-04-12T10:00:00Z",
    "started_at": "2026-04-12T10:00:01Z"
  }
}
```

任务完成示例：

```json
{
  "success": true,
  "message": "查询成功",
  "data": {
    "task_id": "task_1709011200000000000",
    "status": "completed",
    "created_at": "2026-04-12T10:00:00Z",
    "started_at": "2026-04-12T10:00:01Z",
    "completed_at": "2026-04-12T10:00:04Z",
    "result": {
      "task_id": "task_1709011200000000000",
      "filename": "test.mp3",
      "language": "zh",
      "duration_seconds": 10.5,
      "text": "这是完整的转录文本",
      "segments": [
        {
          "index": 1,
          "start": 0.0,
          "end": 3.5,
          "text": "这是第一段"
        }
      ],
      "output_file": "test.srt",
      "process_time_seconds": 2.5
    }
  }
}
```

说明：当前任务状态保存在内存中，服务重启后未完成和历史任务不会保留。

### 4. 下载文件

```bash
curl -O http://localhost:8080/api/v1/download/test.srt
```

### 依赖代理故障处理

如果你遇到类似下面的错误：

```text
reading https://goproxy.io/...: 502 Bad Gateway
```

说明当前环境里的 Go 模块代理不可用。可以直接切换到官方代理并保留直连回退：

```bash
go env -w GOPROXY=https://proxy.golang.org,direct
```

如果你在当前 shell 里只想临时生效：

```bash
export GOPROXY=https://proxy.golang.org,direct
```

本项目的 `Makefile`、`start.sh` 和 Docker 相关构建文件已经默认使用这个代理回退链。

### 测试 API

使用提供的测试脚本：

```bash
./test_api.sh
```

## ⚙️ 配置说明

### 环境变量

复制示例文件并编辑：

```bash
cp .env.example .env
```

**可配置项**

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `SERVER_PORT` | 服务器端口 | 8080 |
| `GIN_MODE` | 运行模式（debug/release） | debug |
| `WHISPER_MODEL_PATH` | 模型文件路径 | models/ggml-base.bin |
| `WHISPER_THREADS` | 处理线程数 | 4 |
| `MAX_FILE_SIZE` | 最大文件大小（字节） | 104857600 (100MB) |
| `UPLOAD_DIR` | 上传目录 | uploads |
| `OUTPUT_DIR` | 输出目录 | outputs |

### Whisper 模型选择

| 模型 | 大小 | 内存需求 | 速度 | 精度 |
|------|------|----------|------|------|
| tiny | 39M | ~1 GB | ⚡⚡⚡⚡⚡ | ⭐⭐ |
| base | 74M | ~1 GB | ⚡⚡⚡⚡ | ⭐⭐⭐ |
| small | 244M | ~2 GB | ⚡⚡⚡ | ⭐⭐⭐⭐ |
| medium | 769M | ~5 GB | ⚡⚡ | ⭐⭐⭐⭐⭐ |
| large | 1550M | ~10 GB | ⚡ | ⭐⭐⭐⭐⭐ |

**推荐**: `base` 或 `small` 模型，平衡速度和精度。

## 🧪 测试

```bash
# 运行所有测试
make test

# 测试覆盖率
make test-cover

# API 功能测试
./test_api.sh
```

## 📚 文档

- **[DOCKER_QUICKSTART.md](DOCKER_QUICKSTART.md)** - Docker 快速开始
- **[docker/README.md](docker/README.md)** - Docker 部署详细指南  
- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** - 架构设计详解
- **[QUICKSTART.md](docs/QUICKSTART.md)** - 快速入门指南
- **[API_REFERENCE.md](docs/API_REFERENCE.md)** - API 参考手册
- **[PROJECT_SUMMARY.md](docs/PROJECT_SUMMARY.md)** - 项目重构总结
- **[legacy/README.md](legacy/README.md)** - 旧版 CLI 工具说明

## 🛠️ 开发

### Makefile 命令

```bash
make help          # 显示所有命令
make setup         # 项目初始化
make build         # 构建应用
make run           # 运行服务器
make test          # 运行测试
make clean         # 清理构建文件
make fmt           # 格式化代码
make lint          # 代码检查

# Docker 命令
make docker-build      # 构建 CPU 版本镜像
make docker-build-gpu  # 构建 GPU 版本镜像
make docker-up         # 启动 CPU 版本服务
make docker-up-gpu     # 启动 GPU 版本服务
make docker-down       # 停止服务
make docker-logs       # 查看日志
make docker-clean      # 清理容器和镜像
```

### 代码规范

- 遵循 Go 官方代码风格
- 使用 `gofmt` 格式化代码
- 通过 `golangci-lint` 检查

## 🚀 生产部署

### Docker 部署（推荐）

#### CPU 版本

```bash
# 使用 docker-compose
docker compose up -d --build

# 或使用管理脚本
./scripts/docker.sh build && ./scripts/docker.sh up
```

#### GPU 版本

```bash
# 使用 docker-compose
docker compose -f docker-compose.gpu.yml up -d --build

# 或使用管理脚本
./scripts/docker.sh build-gpu && ./scripts/docker.sh up-gpu
```

详细部署文档：[docker/README.md](docker/README.md)

### 手动构建

```bash
# 设置为 release 模式
export GIN_MODE=release

# 构建
make build
```

### Docker 部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### Systemd 服务

创建 `/etc/systemd/system/whisper-api.service`:

```ini
[Unit]
Description=Whisper API Server
After=network.target

[Service]
Type=simple
User=whisper
WorkingDirectory=/opt/whisper-api
ExecStart=/opt/whisper-api/bin/whisper-server
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

[MIT License](LICENSE)

## 🔗 相关链接

- [Whisper.cpp](https://github.com/ggerganov/whisper.cpp)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [OpenAI Whisper](https://github.com/openai/whisper)

## 💡 常见问题

### Q: 如何切换模型？

A: 修改 `.env` 文件中的 `WHISPER_MODEL_PATH` 或设置环境变量。

### Q: 支持哪些音频格式？

A: MP3, WAV, M4A, AAC, FLAC, OGG

### Q: 如何提高转录速度？

A: 使用更小的模型（tiny/base）或增加 `WHISPER_THREADS` 值。

### Q: 旧版 CLI 工具还能用吗？

A: 可以，请查看 `legacy/` 目录中的文档。

---

<div align="center">

**Made with ❤️ using Go and Whisper.cpp**

[⬆ 回到顶部](#whisper-api-server)

</div>
