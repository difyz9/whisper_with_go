# 快速参考指南

## 快速开始

### 1. 一键设置
```bash
./setup.sh
```

### 2. 下载模型
```bash
# 下载base模型（推荐入门）
./download_model.sh base

# 其他可用模型
./download_model.sh tiny      # 最快，75MB
./download_model.sh small     # 准确度较好，466MB  
./download_model.sh medium    # 高准确度，1.5GB
```

### 3. 运行示例
```bash
# 基本用法
go run main.go -audio your_audio.mp3

# 指定中文
go run main.go -audio chinese.mp3 -lang zh

# 指定英文
go run main.go -audio english.mp3 -lang en
```

## 常用命令

### 生成测试音频（macOS）
```bash
# 生成英文测试音频
say "Hello, this is a speech recognition test" -o test.aiff
ffmpeg -i test.aiff -ar 16000 -ac 1 test.mp3

# 生成中文测试音频（需要中文语音）
say -v "Ting-Ting" "你好，这是一个语音识别测试" -o test_zh.aiff
ffmpeg -i test_zh.aiff -ar 16000 -ac 1 test_zh.mp3
```

### 编译可执行文件
```bash
# 编译
go build -o whisper-transcribe main.go

# 运行
./whisper-transcribe -audio test.mp3
```

### 批量处理
```bash
# 处理目录下所有MP3文件
for file in *.mp3; do
    echo "处理: $file"
    go run main.go -audio "$file" -lang zh > "${file%.mp3}.txt"
done
```

## 参数说明

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| `-audio` | 音频文件路径（必需） | - | `-audio test.mp3` |
| `-model` | 模型文件路径 | `models/ggml-base.bin` | `-model models/ggml-small.bin` |
| `-lang` | 语言代码 | `auto` | `-lang zh` 或 `-lang en` |
| `-threads` | CPU线程数 | `4` | `-threads 8` |

## 支持的语言代码

- `auto` - 自动检测
- `zh` - 中文
- `en` - 英语
- `ja` - 日语
- `ko` - 韩语
- `es` - 西班牙语
- `fr` - 法语
- `de` - 德语
- `ru` - 俄语
- `ar` - 阿拉伯语
- 更多语言请查看 Whisper 官方文档

## 性能对比

在Apple M1芯片上的大致性能（1分钟音频）：

| 模型 | 处理时间 | 内存使用 | 准确度 |
|------|---------|---------|--------|
| tiny | ~5秒 | 273MB | ⭐⭐⭐ |
| base | ~10秒 | 388MB | ⭐⭐⭐⭐ |
| small | ~30秒 | 852MB | ⭐⭐⭐⭐⭐ |
| medium | ~2分钟 | 2.1GB | ⭐⭐⭐⭐⭐⭐ |

## 故障排查

### 问题: 找不到ffmpeg
```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt-get install ffmpeg

# 验证安装
ffmpeg -version
```

### 问题: 模型下载慢
```bash
# 使用代理
export http_proxy=http://127.0.0.1:7890
export https_proxy=http://127.0.0.1:7890
./download_model.sh base

# 或手动从镜像下载
# https://hf-mirror.com/ggerganov/whisper.cpp
```

### 问题: 编译错误
```bash
# 确保安装了C编译器
xcode-select --install  # macOS
sudo apt-get install build-essential  # Linux

# 清理并重新编译
go clean -cache
go mod tidy
go build
```

### 问题: 识别结果为空或乱码
- 确保音频文件完整且格式正确
- 尝试指定正确的语言代码（不要使用auto）
- 使用更大的模型（如small或medium）
- 检查音频质量（清晰度、噪音等）

## 进阶用法

### 1. 自定义输出格式
修改 `main.go` 中的输出部分：
```go
// JSON格式输出
fmt.Printf(`{"start": %d, "end": %d, "text": "%s"}\n`,
    segment.Start, segment.End, segment.Text)

// SRT字幕格式
fmt.Printf("%d\n%s --> %s\n%s\n\n",
    i+1, formatSRT(segment.Start), formatSRT(segment.End), segment.Text)
```

### 2. 集成到Web服务
```go
import "net/http"

func transcribeHandler(w http.ResponseWriter, r *http.Request) {
    // 接收上传的音频文件
    file, _ := r.MultipartForm.File["audio"][0]
    // 处理转录
    // 返回结果
}

http.HandleFunc("/transcribe", transcribeHandler)
http.ListenAndServe(":8080", nil)
```

### 3. 实时音频流处理
参考 `example.go` 中的高级示例

## 更多资源

- [Whisper.cpp GitHub](https://github.com/ggerganov/whisper.cpp)
- [Go绑定文档](https://github.com/ggerganov/whisper.cpp/tree/master/bindings/go)
- [OpenAI Whisper论文](https://arxiv.org/abs/2212.04356)

## 示例项目结构
```
whisper_with_go/
├── main.go              # 主程序（命令行工具）
├── example.go           # 示例代码（库用法）
├── setup.sh             # 一键设置脚本
├── download_model.sh    # 模型下载脚本
├── README.md            # 详细文档
├── QUICKSTART.md        # 本文件
├── go.mod               # Go模块
├── go.sum               # 依赖锁定
└── models/              # 模型文件目录
    └── ggml-base.bin
```

## 获取帮助

```bash
# 查看帮助信息
go run main.go -h

# 查看模型列表
./download_model.sh

# 测试安装
./setup.sh
```
