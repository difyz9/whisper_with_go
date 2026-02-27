# 快速启动指南

## 项目重构完成

本项目已���功重构为标准的 Gin Web API 项目！

## 新项目结构

```
whisper_with_go/
├── cmd/server/main.go          # 服务器入口
├── internal/                   # 内部代码
│   ├── handler/               # HTTP 处理器
│   ├── service/               # 业务逻辑
│   ├── model/                 # 数据模型
│   ├── middleware/            # 中间件
│   └── router/                # 路由配置
├── pkg/utils/                 # 工具函数
├── config/                    # 配置管理
├── uploads/                   # 上传目录
└── outputs/                   # 输出目录
```

## 启动步骤

### 方式 1: 使用新的 Makefile

```bash
# 查看所有命令
make -f Makefile.new help

# 直接运行（开发模式）
make -f Makefile.new run

# 或者先构建再运行
make -f Makefile.new build
DYLD_LIBRARY_PATH=/usr/local/lib ./bin/whisper-server
```

### 方式 2: 直接使用 go run

```bash
# 设置 CGO 环境变量并运行
CGO_ENABLED=1 \
CGO_LDFLAGS="-L/usr/local/lib" \
CGO_CFLAGS="-I/usr/local/include" \
DYLD_LIBRARY_PATH=/usr/local/lib \
go run cmd/server/main.go
```

### 方式 3: 创建启动脚本

创建 `start.sh` 文件:

```bash
#!/bin/bash
export CGO_ENABLED=1
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"
export DYLD_LIBRARY_PATH=/usr/local/lib

go run cmd/server/main.go
```

然后运行:
```bash
chmod +x start.sh
./start.sh
```

## 服务器启动后

服务器将在 `http://localhost:8080` 启动，你会看到类似的输出:

```
╦ ╦┬ ┬┬┌─┐┌─┐┌─┐┬─┐  ╔═╗╔═╗╦  
║║║├─┤│└─┐├─┘├┤ ├┬┘  ╠═╣╠═╝║  
╚╩╝┴ ┴┴└─┘┴  └─┘┴└─  ╩ ╩╩  ╩  
语音转文字 API 服务
Version: 1.0.0

服务器启动在 http://localhost:8080
模式: debug
模型: models/ggml-base.bin
线程数: 4
--------------------------------------
API 端点:
  健康检查: GET  http://localhost:8080/health
  转录音频: POST http://localhost:8080/api/v1/transcribe
  下载文件: GET  http://localhost:8080/api/v1/download/:filename
--------------------------------------
```

## 测试 API

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

### 2. 转录音频文件

```bash
# JSON 格式输出
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test.mp3" \
  -F "language=auto" \
  -F "output_type=json"

# SRT 字幕格式
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test.mp3" \
  -F "language=zh" \
  -F "output_type=srt"
```

### 3. 使用测试脚本

```bash
./test_api.sh
```

## 环境配置

复制环境变量示例文件并根据需要修改：

```bash
cp .env.example .env
```

可配置项:
- `SERVER_PORT`: 服务器端口（默认 8080）
- `GIN_MODE`: 运行模式（debug/release）
- `WHISPER_MODEL_PATH`: 模型文件路径
- `WHISPER_THREADS`: 处理线程数
- `MAX_FILE_SIZE`: 最大上传文件大小

## 文档

- **README_API.md** - 完整的 API 文档
- **ARCHITECTURE.md** - 项目架构说明
- **旧的 README.md** - 原 CLI 工具文档（保留参考）

## 旧版 CLI 工具

如果需要运行旧版 CLI 工具，仍然可以使用:

```bash
# 使用旧的 Makefile
make run AUDIO=test.mp3

# 或直接运行
go run main.go -audio test.mp3
```

## 常见问题

### CGO 编译错误

如果遇到 CGO 相关的编译错误，确保:

1. 已安装 whisper.cpp 库
2. 设置了正确的 CGO 环境变量
3. macOS 上需要设置 `DYLD_LIBRARY_PATH`

### 模型文件不存在

```bash
# 下载模型
bash download_model.sh
```

### FFmpeg 未安装

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt-get install ffmpeg
```

## 下一步

1. 阅读 **README_API.md** 了解详细的 API 使用方法
2. 阅读 **ARCHITECTURE.md** 了解项目架构
3. 根据需要修改 `.env` 配置文件
4. 开始使用 API！

## 项目优势

相比旧的 CLI 工具，新的 API 项目提供:

✅ **Web API 接口** - 可以通过 HTTP 调用
✅ **RESTful 设计** - 标准的 API 设计
✅ **分层架构** - 清晰的代码结构
✅ **易于扩展** - 便于添加新功能
✅ **生产就绪** - 包含日志、错误处理、CORS 等
✅ **文档完善** - 详细的使用和架构文档

开始使用吧！🚀
