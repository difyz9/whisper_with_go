# Whisper.cpp Go 绑定使用说明

## 重要说明

Go绑定需要先编译whisper.cpp的C/C++库。有两种方式使用：

### 方式1：使用预编译的CLI工具（推荐，最简单）

不需要Go，直接使用whisper.cpp编译好的命令行工具：

```bash
# 1. 克隆whisper.cpp仓库
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp

# 2. 编译
cmake -B build
cmake --build build -j

# 3. 下载模型
bash ./models/download-ggml-model.sh base

# 4. 转换MP3为WAV（whisper需要WAV格式）
ffmpeg -i your_audio.mp3 -ar 16000 -ac 1 -c:a pcm_s16le output.wav

# 5. 运行识别
./build/bin/whisper-cli -m models/ggml-base.bin -f output.wav -l zh
```

### 方式2：使用Go封装CLI（实用方案）

创建一个Go程序来调用whisper.cpp的CLI工具：

```go
// simple_transcribe.go
package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("用法: go run simple_transcribe.go <audio.mp3> [language]")
        os.Exit(1)
    }

    audioFile := os.Args[1]
    language := "auto"
    if len(os.Args) > 2 {
        language = os.Args[2]
    }

    // 转换为WAV
    wavFile := strings.TrimSuffix(audioFile, filepath.Ext(audioFile)) + ".wav"
    fmt.Println("转换音频格式...")
    cmd := exec.Command("ffmpeg", "-i", audioFile, "-ar", "16000", "-ac", "1", 
                        "-c:a", "pcm_s16le", wavFile, "-y")
    if err := cmd.Run(); err != nil {
        fmt.Printf("转换失败: %v\n", err)
        os.Exit(1)
    }
    defer os.Remove(wavFile)

    // 调用whisper-cli
    fmt.Println("开始识别...")
    whisperCmd := exec.Command("whisper-cli", 
        "-m", "models/ggml-base.bin",
        "-f", wavFile,
        "-l", language)
    
    whisperCmd.Stdout = os.Stdout
    whisperCmd.Stderr = os.Stderr
    
    if err := whisperCmd.Run(); err != nil {
        fmt.Printf("识别失败: %v\n", err)
        os.Exit(1)
    }
}
```

### 方式3：使用Go绑定（需要编译whisper.cpp）

如果确实需要使用Go绑定，需要先构建whisper.cpp：

```bash
# 1. 克隆并编译whisper.cpp
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp
cmake -B build
cmake --build build -j
sudo cmake --install build  # 安装到系统路径

# 2. 设置环境变量
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"

# 3. 在你的Go项目中使用
cd /path/to/your/go/project
go mod tidy
go build
```

## 推荐方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| CLI工具 | 最简单，无需Go | 功能受限于CLI | 简单脚本、一次性任务 |
| Go封装CLI | 简单，易集成 | 需要安装CLI | Go项目中的音频处理 |
| Go绑定 | 完全控制，高性能 | 配置复杂 | 需要深度集成的项目 |

## 快速开始（方式1 - 推荐）

```bash
# 完整流程
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp
cmake -B build && cmake --build build -j
bash ./models/download-ggml-model.sh base

# 使用
ffmpeg -i your_audio.mp3 -ar 16000 -ac 1 -c:a pcm_s16le output.wav
./build/bin/whisper-cli -m models/ggml-base.bin -f output.wav -l zh
```

这是最稳定、最简单的方案！
