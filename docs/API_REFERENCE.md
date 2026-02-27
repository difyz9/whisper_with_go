# Whisper API Server

基于 Gin 框架的 Whisper 语音转文字 API 服务，采用标准的分层架构设计。

## 项目结构

```
whisper_with_go/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── handler/                 # HTTP 处理器层
│   │   └── whisper.go
│   ├── service/                 # 业务逻辑层
│   │   └── whisper.go
│   ├── model/                   # 数据模型层
│   │   ├── request.go
│   │   └── response.go
│   ├── middleware/              # 中间件
│   │   ├── logger.go
│   │   ├── cors.go
│   │   └── recovery.go
│   └── router/                  # 路由配置
│       └── router.go
├── pkg/                         # 可复用的公共库
│   └── utils/
│       ├── audio.go             # 音频处理工具
│       ├── file.go              # 文件处理工具
│       └── whisper.go           # Whisper 工具
├── config/                      # 配置管理
│   └── config.go
├── uploads/                     # 上传文件目录
├── outputs/                     # 输出文件目录
├── models/                      # Whisper 模型文件
├── go.mod
├── go.sum
├── .env.example                 # 环境变量示例
├── Makefile                     # 构建和运行脚本
└── README.md
```

## 特性

- ✅ 标准的 Gin 框架项目结构
- ✅ 清晰的分层架构（handler → service → model）
- ✅ RESTful API 设计
- ✅ 支持多种音频格式（MP3, WAV, M4A, AAC, FLAC, OGG）
- ✅ 支持多种输出格式（JSON, SRT, TXT）
- ✅ 自动语言检测
- ✅ 完整的日志记录
- ✅ CORS 支持
- ✅ 错误恢复机制
- ✅ 环境变量配置

## 前置要求

1. **Go 1.23.4+**
2. **FFmpeg** - 用于音频格式转换
   ```bash
   # macOS
   brew install ffmpeg
   
   # Ubuntu/Debian
   sudo apt-get install ffmpeg
   
   # CentOS/RHEL
   sudo yum install ffmpeg
   ```

3. **Whisper 模型文件**
   ```bash
   # 下载模型（例如 base 模型）
   bash download_model.sh
   ```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

```bash
# 复制环境变量示例文件
cp .env.example .env

# 编辑 .env 文件根据需要调整配置
```

### 3. 下载 Whisper 模型

```bash
bash download_model.sh
```

### 4. 启动服务器

```bash
# 方式1: 使用 go run
go run cmd/server/main.go

# 方式2: 使用 Makefile
make run

# 方式3: 构建后运行
make build
./bin/whisper-server
```

服务器将在 `http://localhost:8080` 启动

## API 文档

### 1. 健康检查

**GET** `/health`

响应:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "time": "2026-02-27T10:00:00Z"
}
```

### 2. 转录音频

**POST** `/api/v1/transcribe`

请求参数:
- `file` (必需): 音频文件（multipart/form-data）
- `language` (可选): 语言代码，默认 "auto"
  - `auto`: 自动检测
  - `zh`: 中文
  - `en`: 英文
  - 其他语言代码...
- `output_type` (可选): 输出格式，默认 "json"
  - `json`: JSON 格式（包含详细的片段信息）
  - `srt`: SRT 字幕格式
  - `txt`: 纯文本格式
- `translate` (可选): 是否翻译为英文，默认 false

示例 (使用 curl):
```bash
# JSON 输出（默认）
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test.mp3" \
  -F "language=auto" \
  -F "output_type=json"

# SRT 字幕输出
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test.mp3" \
  -F "language=zh" \
  -F "output_type=srt"

# 翻译为英文
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test.mp3" \
  -F "translate=true"
```

响应 (JSON):
```json
{
  "success": true,
  "message": "转录成功",
  "data": {
    "task_id": "task_1709011200000000000",
    "filename": "uploads/test_20260227_abc123.mp3",
    "language": "zh",
    "duration_seconds": 10.5,
    "text": "这是转录的完整文本内容",
    "segments": [
      {
        "index": 1,
        "start": 0.0,
        "end": 3.5,
        "text": "这是第一段文本"
      },
      {
        "index": 2,
        "start": 3.5,
        "end": 7.0,
        "text": "这是第二段文本"
      }
    ],
    "output_file": "outputs/test.srt",
    "process_time_seconds": 2.5
  }
}
```

### 3. 下载输出文件

**GET** `/api/v1/download/:filename`

示例:
```bash
curl -O http://localhost:8080/api/v1/download/test.srt
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `SERVER_PORT` | 服务器端口 | `8080` |
| `GIN_MODE` | Gin 运行模式 | `debug` |
| `WHISPER_MODEL_PATH` | Whisper 模型文件路径 | `models/ggml-base.bin` |
| `WHISPER_THREADS` | 处理线程数 | `4` |
| `MAX_FILE_SIZE` | 最大上传文件大小（字节） | `104857600` (100MB) |
| `UPLOAD_DIR` | 上传文件目录 | `uploads` |
| `OUTPUT_DIR` | 输出文件目录 | `outputs` |

### 模型说明

Whisper 提供多种模型，按大小和精度分为:

| 模型 | 参数 | 所需显存 | 速度 | 精度 |
|------|------|----------|------|------|
| tiny | 39M | ~1 GB | 最快 | 较低 |
| base | 74M | ~1 GB | 快 | 中等 |
| small | 244M | ~2 GB | 中等 | 好 |
| medium | 769M | ~5 GB | 慢 | 很好 |
| large | 1550M | ~10 GB | 最慢 | 最好 |

推荐使用 `base` 或 `small` 模型以获得速度和精度的平衡。

## 开发

### 项目架构

本项目采用标准的三层架构:

1. **Handler 层** (`internal/handler/`)
   - 处理 HTTP 请求和响应
   - 参数验证
   - 调用 Service 层

2. **Service 层** (`internal/service/`)
   - 核心业务逻辑
   - 调用 Whisper API
   - 数据处理

3. **Model 层** (`internal/model/`)
   - 数据结构定义
   - 请求/响应模型

4. **Utils 层** (`pkg/utils/`)
   - 可复用的工具函数
   - 音频处理
   - 文件操作

### Makefile 命令

```bash
make help       # 显示帮助信息
make build      # 构建应用
make run        # 运行应用
make test       # 运行测试
make clean      # 清理构建文件
make fmt        # 格式化代码
make lint       # 代码检查
```

## 测试

```bash
# 运行所有测试
make test

# 运行特定测试
go test -v ./internal/service/...

# 测试覆盖率
go test -cover ./...
```

## 生产部署

### 1. 构建生产版本

```bash
# 设置为 release 模式
export GIN_MODE=release

# 构建
make build
```

### 2. 使用 Docker (可选)

```dockerfile
FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o whisper-server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ffmpeg
WORKDIR /root/
COPY --from=builder /app/whisper-server .
COPY models/ models/
CMD ["./whisper-server"]
```

### 3. Systemd 服务 (Linux)

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

## 常见问题

### 1. FFmpeg 未安装

**错误**: `ffmpeg 未安装，请先安装 ffmpeg`

**解决**: 安装 FFmpeg（见前置要求部分）

### 2. 模型文件不存在

**错误**: `模型文件不存在`

**解决**: 运行 `bash download_model.sh` 下载模型

### 3. 文件上传失败

**错误**: `文件过大`

**解决**: 调整 `MAX_FILE_SIZE` 环境变量

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 相关链接

- [Whisper.cpp](https://github.com/ggerganov/whisper.cpp)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [OpenAI Whisper](https://github.com/openai/whisper)
