在当前目录下实现，golang如何实现基于Go 绑定 https://github.com/ggml-org/whisper.cpp 的语音识别功能
写一个使用案例，让我实现对本地mp3文件的语音识别功能。要在Go中实现基于Go绑定的whisper.cpp进行语音识别功能

say "Hello, this is a test of the Whisper speech recognition system using Go" -o test.aiff && ffmpeg -i test.aiff -ar 16000 -ac 1 test.mp3 -y 2>&1 | tail -3




## ✅ 项目已完成！

### 已创建的文件

1. **main.go** - 原始Go绑定版本（需要编译whisper.cpp C库）
   - 直接使用Go绑定API
   - 需要CGO和C编译环境
   
2. **main_cli.go** - CLI封装版本（推荐使用）
   - 无需CGO，调用whisper-cli命令
   - 支持MP3文件输入（自动转换为WAV）
   - 命令行参数支持（音频路径、模型、语言、线程数）
   - 更简单易用
   
2. **example.go** - 示例代码，展示如何以库的方式使用
   - 编程方式调用API
   - 高级用法示例（批量处理、实时流、翻译等）
   
3. **README.md** - 完整的中文文档
   - 详细的安装步骤
   - 使用方法和示例
   - 故障排查指南
   
4. **QUICKSTART.md** - 快速开始指南
   - 常用命令速查
   - 快速上手教程
   
5. **setup.sh** - 一键设置脚本
   - 自动检查环境依赖
   - 下载模型和Go依赖
   
6. **download_model.sh** - 模型下载脚本
   - 支持多种模型下载（tiny、base、small、medium）
   
7. **go.mod / go.sum** - Go模块配置
   - 已配置whisper.cpp Go绑定依赖

8. **Makefile** - Make构建文件
   - 提供便捷的命令（setup、build、run、test等）
   
9. **.gitignore** - Git忽略配置
   - 忽略模型文件、临时文件等

### 快速开始

#### 方法1：使用Makefile（最简单）
```bash
# 查看所有可用命令
make help

# 一键设置环境
make setup

# 下载模型
make download-model

# 生成测试音频并运行（macOS）
make test           # 英文测试
make test-zh        # 中文测试

# 运行自己的音频文件
make run AUDIO=your_audio.mp3
make run AUDIO=your_audio.mp3 LANG=zh
```

#### 方法2：使用setup脚本
```bash
# 一键设置
./setup.sh

# 运行示例（需要有MP3文件）
go run main.go -audio your_audio.mp3
```

#### 方法3：手动步骤
```bash
# 1. 安装FFmpeg（如果未安装）
brew install ffmpeg  # macOS

# 2. 下载模型
./download_model.sh base

# 3. 运行
go run main.go -audio your_audio.mp3
```

### 使用示例

#### 使用Makefile（推荐）
```bash
# 查看帮助
make help

# 设置环境
make setup

# 下载模型
make download-model

# 自动生成测试音频并运行
make test        # 英文
make test-zh     # 中文

# 运行自己的音频
make run AUDIO=test.mp3
make run AUDIO=test.mp3 LANG=zh

# 编译
make build

# 安装到系统
make install

# 清理
make clean
```

#### 使用Go命令
```bash
# 基本用法
go run main.go -audio test.mp3

# 指定中文
go run main.go -audio chinese_audio.mp3 -lang zh

# 使用不同的模型
go run main.go -audio test.mp3 -model models/ggml-small.bin

# 使用更多线程
go run main.go -audio test.mp3 -threads 8

# 编译为可执行文件
go build -o whisper-transcribe main.go
./whisper-transcribe -audio test.mp3
```

### 生成测试音频（macOS）

```bash
# 英文测试音频
say "Hello, this is a test of speech recognition" -o test.aiff
ffmpeg -i test.aiff -ar 16000 -ac 1 test.mp3
go run main.go -audio test.mp3

# 中文测试音频
say -v "Ting-Ting" "你好，这是语音识别测试" -o test_zh.aiff  
ffmpeg -i test_zh.aiff -ar 16000 -ac 1 test_zh.mp3
go run main.go -audio test_zh.mp3 -lang zh
```

### 功能特性

✅ 完整的MP3语音识别功能
✅ 自动音频格式转换（MP3→WAV）
✅ 支持多语言识别（中文、英文等）
✅ 时间戳输出
✅ 命令行参数支持
✅ 跨平台（macOS、Linux、Windows）
✅ 详细的中文文档
✅ 一键设置脚本
✅ 示例代码

### 项目结构

```
whisper_with_go/
├── main.go              # 主程序（Go绑定版本）⭐
├── main_cli.go          # CLI封装版本
├── example.go           # API使用示例
├── run.sh               # 便捷运行脚本 ⭐
├── build.sh             # 编译脚本
├── extract_audio.sh     # MP4音频提取工具 ⭐
├── setup.sh             # 一键设置脚本
├── download_model.sh    # 模型下载脚本
├── Makefile             # Make构建文件
├── README.md            # 详细中文文档
├── QUICKSTART.md        # 快速参考指南
├── INSTALL_GUIDE.md     # 安装指南
├── dev.md               # 本文件（项目说明）
├── go.mod               # Go模块定义
├── go.sum               # 依赖锁定文件
├── .gitignore           # Git忽略配置
└── models/              # 模型文件目录
    ├── ggml-tiny.bin    # 74M
    ├── ggml-base.bin    # 141M
    ├── ggml-small.bin   # 465M
    └── ggml-medium.bin  # 1.4G
```

### 技术栈

- **语言**: Go 1.19+
- **核心库**: github.com/ggerganov/whisper.cpp/bindings/go
- **音频处理**: FFmpeg
- **模型**: OpenAI Whisper (GGML格式)

### 下一步

1. 运行 `./setup.sh` 完成环境设置
2. 准备一个MP3音频文件
3. 运行 `go run main.go -audio your_file.mp3`
4. 查看README.md了解更多功能

### ✅ 方案3已成功配置！(使用Go绑定)

**main.go已经可以正常运行了！** 通过编译whisper.cpp并安装到系统，现在可以直接使用Go绑定。

#### 快速使用（三种方式）：

**方式1：使用run.sh脚本（最方便）**
```bash
./run.sh -audio test.mp3                        # 纯文本输出
./run.sh -audio test.mp3 -lang zh               # 指定中文
./run.sh -audio test.mp3 -srt                   # 生成SRT字幕 ⭐
./run.sh -audio test.mp3 -srt -lang zh          # SRT字幕（中文）
./run.sh -audio test.mp3 -model models/ggml-small.bin
```

**方式2：直接运行（需要设置环境变量）**
```bash
DYLD_LIBRARY_PATH=/usr/local/lib go run main.go -audio test.mp3
```

**方式3：编译后使用**
```bash
# 使用build.sh编译
./build.sh

# 运行编译后的程序
DYLD_LIBRARY_PATH=/usr/local/lib ./whisper-transcribe -audio test.mp3
```

#### 测试结果示例：

```bash
$ ./run.sh -audio test.mp3
音频文件: test.mp3
临时WAV文件: test_temp.wav
模型文件: models/ggml-base.bin
语言: auto
开始转录...
------------------------------------------------------------
[    0s -->   5.6s]  Hello, this is the test of the Whisper Speech Recognition System using Go.
------------------------------------------------------------
转录完成!
```

#### 其他可用方案：

1. **main_cli.go** - CLI封装版本（不需要CGO）
   - 调用whisper-cli命令
   - 更简单，无需动态库
   
2. **直接使用whisper-cli**
   ```bash
   whisper-cli -m models/ggml-base.bin -f audio.wav
   ```

#### 新增文件：

- **run.sh** ⭐ - 便捷运行脚本（自动设置环境变量）
- **build.sh** - 编译脚本（自动设置编译环境）
- **extract_audio.sh** ⭐ - MP4音频提取工具（支持批量处理和直接转录）
- **INSTALL_GUIDE.md** - 三种方案的详细对比

### 🎬 从视频提取音频和生成字幕

使用 `extract_audio.sh` 脚本可以轻松从MP4视频中提取音频并生成字幕：

```bash
# 基本用法
./extract_audio.sh video.mp4              # 提取为video.mp3
./extract_audio.sh video.mp4 wav          # 提取为video.wav (Whisper格式)

# 自定义输出文件名
./extract_audio.sh video.mp4 mp3 audio.mp3
./extract_audio.sh video.mp4 wav audio.wav

# 提取并生成SRT字幕（一步到位）⭐
./extract_audio.sh video.mp4 -srt              # 自动检测语言
./extract_audio.sh video.mp4 -srt zh           # 指定中文
./extract_audio.sh video.mp4 -srt en           # 指定英文

# 提取并转录为纯文本
./extract_audio.sh video.mp4 -transcribe        # 自动检测语言
./extract_audio.sh video.mp4 -transcribe zh     # 指定中文

# 批量处理当前目录所有MP4文件
./extract_audio.sh -batch mp3    # 批量提取为MP3
./extract_audio.sh -batch wav    # 批量提取为WAV
```

**完整工作流示例：**
```bash
# 1. 从视频生成SRT字幕（推荐）
./extract_audio.sh meeting.mp4 -srt zh

# 2. 或纯文本转录
./extract_audio.sh meeting.mp4 -transcribe zh

# 3. 或者分步操作
./extract_audio.sh meeting.mp4 mp3              # 提取音频
./run.sh -audio meeting.mp3 -srt -lang zh       # 生成字幕
```

### 需要帮助？

- 查看 INSTALL_GUIDE.md 获取详细安装步骤
- 查看 README.md 获取完整文档
- 查看 QUICKSTART.md 快速上手
- 查看 example.go 了解API用法