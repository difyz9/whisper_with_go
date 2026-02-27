#!/bin/bash
# 快速开始脚本 - 下载模型并运行示例

set -e

echo "====================================="
echo "   Whisper.cpp Go 绑定快速开始"
echo "====================================="
echo ""

# 1. 检查Go是否安装
echo "[1/5] 检查Go环境..."
if ! command -v go &> /dev/null; then
    echo "❌ 错误: Go未安装"
    echo "请访问 https://golang.org/dl/ 下载安装"
    exit 1
fi
echo "✓ Go版本: $(go version)"
echo ""

# 2. 检查FFmpeg是否安装
echo "[2/5] 检查FFmpeg..."
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ 错误: FFmpeg未安装"
    echo "macOS: brew install ffmpeg"
    echo "Linux: sudo apt-get install ffmpeg"
    exit 1
fi
echo "✓ FFmpeg已安装"
echo ""

# 3. 下载依赖
echo "[3/5] 下载Go依赖..."
go mod tidy
echo "✓ 依赖下载完成"
echo ""

# 4. 下载模型
echo "[4/5] 下载Whisper模型..."
if [ ! -f "models/ggml-base.bin" ]; then
    ./download_model.sh base
else
    echo "✓ 模型已存在"
fi
echo ""

# 5. 创建示例音频说明
echo "[5/5] 准备测试..."
echo ""
echo "====================================="
echo "设置完成! 现在可以开始使用了。"
echo "====================================="
echo ""
echo "使用方法:"
echo "  go run main.go -audio your_audio.mp3"
echo ""
echo "示例："
echo "  # 如果你有MP3文件，例如 test.mp3"
echo "  go run main.go -audio test.mp3"
echo ""
echo "  # 指定语言"
echo "  go run main.go -audio test.mp3 -lang zh"
echo ""
echo "  # 使用更大的模型（需先下载）"
echo "  ./download_model.sh small"
echo "  go run main.go -audio test.mp3 -model models/ggml-small.bin"
echo ""
echo "测试音频生成（可选）："
echo "  # 使用macOS的say命令生成测试音频"
echo "  say 'Hello, this is a test of the Whisper speech recognition system' -o test.aiff"
echo "  ffmpeg -i test.aiff -ar 16000 -ac 1 test.wav"
echo "  # 然后运行: go run main.go -audio test.wav"
echo ""
