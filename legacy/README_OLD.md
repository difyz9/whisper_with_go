# Go Whisper.cpp 语音识别示例

这个项目演示了如何在Go中使用whisper.cpp的Go绑定来实现本地MP3文件的语音识别功能。

## 功能特性

- ✅ 支持MP3格式音频文件（自动转换为WAV）
- ✅ 基于whisper.cpp的高性能语音识别
- ✅ 支持多种语言识别（中文、英文等）
- ✅ 时间戳输出
- ✅ 跨平台支持（macOS、Linux、Windows）

## 系统要求

- Go 1.19或更高版本
- FFmpeg（用于音频格式转换）
- C/C++编译器（用于编译whisper.cpp）

## 安装步骤

### 1. 安装依赖

#### macOS
```bash
# 安装FFmpeg
brew install ffmpeg

# 安装编译工具（通常已安装）
xcode-select --install
```

#### Linux (Ubuntu/Debian)
```bash
# 安装FFmpeg和编译工具
sudo apt-get update
sudo apt-get install ffmpeg build-essential
```

#### Windows
```bash
# 安装FFmpeg
# 下载: https://ffmpeg.org/download.html
# 安装MinGW-w64用于编译
```

### 2. 克隆whisper.cpp并下载模型

```bash
# 创建models目录
mkdir -p models

# 下载whisper.cpp仓库中的模型下载脚本
curl -o models/download-ggml-model.sh https://raw.githubusercontent.com/ggerganov/whisper.cpp/master/models/download-ggml-model.sh

# 赋予执行权限（macOS/Linux）
chmod +x models/download-ggml-model.sh

# 下载模型（可选: tiny, base, small, medium, large）
bash models/download-ggml-model.sh base
```

或者手动从 [Hugging Face](https://huggingface.co/ggerganov/whisper.cpp) 下载模型文件到`models/`目录。

### 3. 安装Go依赖

```bash
go mod tidy
```

## 使用方法

### 基本用法

```bash
# 识别MP3文件
go run main.go -audio your_audio.mp3

# 指定模型文件
go run main.go -audio your_audio.mp3 -model models/ggml-base.bin

# 指定语言
go run main.go -audio your_audio.mp3 -lang zh

# 指定线程数
go run main.go -audio your_audio.mp3 -threads 8
```

### 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-audio` | (必需) | MP3音频文件路径 |
| `-model` | `models/ggml-base.bin` | Whisper模型文件路径 |
| `-lang` | `auto` | 语言代码 (auto, zh, en, ja, etc.) |
| `-threads` | `4` | 使用的CPU线程数 |

### 示例

```bash
# 识别中文音频
go run main.go -audio chinese_speech.mp3 -lang zh

# 识别英文音频，使用更大的模型以提高准确度
go run main.go -audio english_speech.mp3 -model models/ggml-small.bin -lang en

# 使用所有可用CPU核心
go run main.go -audio my_audio.mp3 -threads 8
```

## 编译为可执行文件

```bash
# 编译
go build -o whisper-transcribe main.go

# 运行
./whisper-transcribe -audio your_audio.mp3
```

## 项目结构

```
whisper_with_go/
├── main.go              # 主程序
├── go.mod               # Go模块配置
├── go.sum               # Go依赖锁定文件
├── README.md            # 说明文档
├── models/              # 模型文件目录
│   ├── download-ggml-model.sh
│   └── ggml-base.bin
└── dev.md               # 开发文档
```

## 支持的模型

| 模型 | 大小 | 内存使用 | 相对速度 | 准确度 |
|------|------|----------|----------|--------|
| tiny | 75 MB | ~273 MB | 最快 | 一般 |
| base | 142 MB | ~388 MB | 快 | 较好 |
| small | 466 MB | ~852 MB | 中等 | 好 |
| medium | 1.5 GB | ~2.1 GB | 慢 | 很好 |
| large | 2.9 GB | ~3.9 GB | 最慢 | 最好 |

推荐：
- 日常使用：`base`或`small`
- 高准确度需求：`medium`或`large`
- 快速测试：`tiny`

## 输出格式

程序会输出带时间戳的转录文本：

```
音频文件: my_audio.mp3
临时WAV文件: my_audio_temp.wav
模型文件: models/ggml-base.bin
语言: auto
开始转录...
------------------------------------------------------------
[00:00:00.000 --> 00:00:03.000]  这是一段测试音频
[00:00:03.000 --> 00:00:06.500]  用于演示语音识别功能
[00:00:06.500 --> 00:00:10.000]  识别结果将包含时间戳信息
------------------------------------------------------------
转录完成!
```

## 常见问题

### 1. 找不到ffmpeg
```
错误: ffmpeg未安装，请先安装ffmpeg
```
**解决方案**: 按照上述安装步骤安装FFmpeg

### 2. 模型文件不存在
```
错误: 模型文件不存在: models/ggml-base.bin
```
**解决方案**: 运行模型下载脚本或手动下载模型文件

### 3. 编译错误
```
错误: C编译器未找到
```
**解决方案**: 
- macOS: `xcode-select --install`
- Linux: `sudo apt-get install build-essential`
- Windows: 安装MinGW-w64

### 4. 音频识别不准确
- 尝试使用更大的模型（small、medium、large）
- 确保音频质量良好
- 指定正确的语言代码

## 性能优化建议

1. **使用适当的模型**: 根据准确度和速度需求选择合适的模型
2. **调整线程数**: 使用`-threads`参数设置为CPU核心数
3. **音频质量**: 使用清晰、低噪音的音频文件可提高准确度

## 技术实现

本项目使用以下技术：

- **Go语言**: 主程序语言
- **whisper.cpp**: OpenAI Whisper模型的C++实现
- **FFmpeg**: 音频格式转换
- **cgo**: Go与C/C++互操作

## 参考资源

- [whisper.cpp GitHub](https://github.com/ggerganov/whisper.cpp)
- [OpenAI Whisper](https://github.com/openai/whisper)
- [whisper.cpp Go绑定](https://github.com/ggerganov/whisper.cpp/tree/master/bindings/go)

## 许可证

本项目遵循MIT许可证。whisper.cpp也使用MIT许可证。

## 贡献

欢迎提交Issue和Pull Request！
