# 项目重构说明

## 概述

本项目已成功从简单的 CLI 工具重构为标准的 Gin Web API 项目，采用清晰的分层架构设计。

## 项目结构

### 核心目录说明

```
whisper_with_go/
│
├── cmd/                         # 应用程序入口目录
│   └── server/
│       └── main.go             # 主程序入口，初始化并启动服务器
│
├── internal/                    # 内部私有代码，不对外暴露
│   ├── handler/                # HTTP 处理器层（Controller）
│   │   └── whisper.go         # Whisper API 的 HTTP 处理逻辑
│   │
│   ├── service/                # 业务逻辑层（Service）
│   │   └── whisper.go         # Whisper 核心业务逻辑
│   │
│   ├── model/                  # 数据模型层（Model）
│   │   ├── request.go         # 请求数据结构
│   │   └── response.go        # 响应数据结构
│   │
│   ├── middleware/             # 中间件
│   │   ├── logger.go          # 日志记录中间件
│   │   ├── cors.go            # CORS 跨域中间件
│   │   └── recovery.go        # 错误恢复中间件
│   │
│   └── router/                 # 路由配置
│       └── router.go          # 路由定义和组装
│
├── pkg/                        # 可复用的公共库
│   └── utils/
│       ├── audio.go           # 音频处理工具函数
│       ├── file.go            # 文件操作工具函数
│       └── whisper.go         # Whisper 相关工具函数
│
├── config/                     # 配置管理
│   └── config.go              # 配置结构和加载逻辑
│
├── uploads/                    # 上传文件存储目录
├── outputs/                    # 输出文件存储目录
├── models/                     # Whisper 模型文件目录
│
├── go.mod                      # Go 模块依赖定义
├── go.sum                      # 依赖版本锁定
├── .env.example               # 环境变量示例
├── Makefile.new               # 构建和运行脚本（新版）
└── README_API.md              # API 文档
```

## 分层架构说明

### 1. Handler 层（`internal/handler/`）

**职责**:
- 处理 HTTP 请求和响应
- 参数验证和绑定
- 调用 Service 层执行业务逻辑
- 格式化返回数据

**示例**:
```go
func (h *WhisperHandler) Transcribe(c *gin.Context) {
    // 1. 接收并验证请求
    // 2. 调用 service 层
    // 3. 返回格式化的响应
}
```

### 2. Service 层（`internal/service/`）

**职责**:
- 实现核心业务逻辑
- 不依赖于 HTTP 框架
- 可以被不同的 Handler 复用
- 处理数据转换和业务规则

**示例**:
```go
func (s *whisperService) Transcribe(audioPath, language, outputType string, translate bool) (*model.TranscribeResponse, error) {
    // 1. 转换音频格式
    // 2. 加载 Whisper 模型
    // 3. 执行转录
    // 4. 返回结果
}
```

### 3. Model 层（`internal/model/`）

**职责**:
- 定义数据结构
- 请求/响应模型
- 数据传输对象（DTO）

**示例**:
```go
type TranscribeRequest struct {
    Language   string `json:"language"`
    OutputType string `json:"output_type"`
}

type TranscribeResponse struct {
    TaskID   string    `json:"task_id"`
    Text     string    `json:"text"`
    Segments []Segment `json:"segments"`
}
```

### 4. Router 层（`internal/router/`）

**职责**:
- 定义路由规则
- 组装中间件
- 创建依赖注入

**示例**:
```go
func Setup(cfg *config.Config) *gin.Engine {
    r := gin.New()
    r.Use(middleware.Logger())
    r.Use(middleware.CORS())
    
    v1 := r.Group("/api/v1")
    {
        v1.POST("/transcribe", handler.Transcribe)
    }
    
    return r
}
```

### 5. Middleware 层（`internal/middleware/`）

**职责**:
- 请求预处理
- 日志记录
- 认证授权
- 错误恢复

**示例**:
```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 记录请求信息
        c.Next()
        // 记录响应信息
    }
}
```

### 6. Utils 层（`pkg/utils/`）

**职责**:
- 提供可复用的工具函数
- 独立于业务逻辑
- 可被多个模块使用

**示例**:
```go
func ConvertToWAV(inputPath string) (string, error) {
    // 音频格式转换逻辑
}
```

### 7. Config 层（`config/`）

**职责**:
- 管理应用配置
- 环境变量加载
- 配置验证

**示例**:
```go
type Config struct {
    Server  ServerConfig
    Whisper WhisperConfig
    Upload  UploadConfig
}

func Load() *Config {
    // 从环境变量加载配置
}
```

## 数据流向

```
HTTP Request
    ↓
Router (路由匹配)
    ↓
Middleware (日志、CORS、恢复等)
    ↓
Handler (参数验证)
    ↓
Service (业务逻辑)
    ↓
Utils (工具函数)
    ↓
Service (返回结果)
    ↓
Handler (格式化响应)
    ↓
Middleware (记录日志)
    ↓
HTTP Response
```

## API 端点

### 1. 健康检查
```
GET /health
```

### 2. 转录音频
```
POST /api/v1/transcribe
Content-Type: multipart/form-data

参数:
- file: 音频文件
- language: 语言代码（可选）
- output_type: 输出格式（可选）
- translate: 是否翻译（可选）
```

### 3. 下载文件
```
GET /api/v1/download/:filename
```

## 依赖管理

### 核心依赖

1. **github.com/gin-gonic/gin** - Web 框架
2. **github.com/ggerganov/whisper.cpp/bindings/go** - Whisper Go 绑定
3. **github.com/sirupsen/logrus** - 日志库

### 安装依赖

```bash
go mod tidy
```

## 运行项目

### 开发环境

```bash
# 方式1: 直接运行
go run cmd/server/main.go

# 方式2: 使用 Makefile
make run
```

### 生产环境

```bash
# 构建
make build

# 运行
./bin/whisper-server
```

## 环境变量

项目支持以下环境变量配置:

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务器端口 | 8080 |
| GIN_MODE | 运行模式 | debug |
| WHISPER_MODEL_PATH | 模型路径 | models/ggml-base.bin |
| WHISPER_THREADS | 线程数 | 4 |
| MAX_FILE_SIZE | 最大文件大小 | 104857600 |
| UPLOAD_DIR | 上传目录 | uploads |
| OUTPUT_DIR | 输出目录 | outputs |

## 代码规范

### 1. 命名规范

- **文件名**: 小写，下划线分隔（如 `whisper_service.go`）
- **包名**: 小写，简短有意义（如 `handler`, `service`）
- **接口**: I 开头或 er 结尾（如 `WhisperService`）
- **结构体**: 大写驼峰（如 `TranscribeRequest`）
- **函数**: 大写驼峰（公开）或小写驼峰（私有）

### 2. 目录规范

- `cmd/`: 应用程序入口
- `internal/`: 内部私有代码
- `pkg/`: 可复用的公共库
- `config/`: 配置相关
- `api/`: API 定义（如 Swagger）
- `docs/`: 文档

### 3. 错误处理

- 使用 `error` 返回值
- 在 Handler 层统一处理错误响应
- Service 层专注于业务逻辑错误

### 4. 日志记录

- 使用 logrus 记录结构化日志
- 记录请求、响应、错误信息
- 生产环境使用 JSON 格式

## 测试

### 单元测试

```bash
# 运行所有测试
make test

# 测试覆盖率
make test-cover
```

### API 测试

使用 curl 或 Postman 测试 API:

```bash
# 健康检查
curl http://localhost:8080/health

# 转录测试
curl -X POST http://localhost:8080/api/v1/transcribe \
  -F "file=@test.mp3" \
  -F "language=auto" \
  -F "output_type=json"
```

## 后续改进建议

1. **数据库集成**: 添加数据库存储转录记录
2. **异步处理**: 使用队列处理长时间转录任务
3. **缓存机制**: 缓存常见请求结果
4. **认证授权**: 添加 JWT 或 OAuth2 认证
5. **限流**: 添加请求限流保护
6. **监控**: 集成 Prometheus/Grafana 监控
7. **文档**: 添加 Swagger API 文档
8. **Docker**: 完善 Docker 部署方案
9. **CI/CD**: 添加自动化测试和部署
10. **WebSocket**: 支持实时转录推送

## 旧文件说明

以下文件为原 CLI 版本的文件，已被新架构替代：

- `main.go` - 原 CLI 主程序（保留作为参考）
- `main_cli.go` - 原 CLI 工具（保留作为参考）
- `example.go` - 示例代码（保留作为参考）

建议：这些文件可以移动到 `legacy/` 目录或删除。

## 总结

本次重构完成了以下目标：

✅ 清晰的目录结构
✅ 标准的分层架构
✅ RESTful API 设计
✅ 完整的错误处理
✅ 日志和监控支持
✅ 配置管理
✅ 完善的文档

项目现在具有良好的可维护性、可扩展性和可测试性。
