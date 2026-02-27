#!/bin/bash

# Whisper API Server 启动脚本

echo "================================"
echo "Whisper API Server"
echo "================================"
echo ""

# 设置 CGO 环境变量
export CGO_ENABLED=1
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"
export DYLD_LIBRARY_PATH=/usr/local/lib

# 检查模型文件
if [ ! -f "models/ggml-base.bin" ]; then
    echo "⚠️  警告: 模型文件不存在"
    echo "请运行: bash download_model.sh"
    echo ""
fi

# 检查 FFmpeg
if ! command -v ffmpeg &> /dev/null; then
    echo "⚠️  警告: FFmpeg 未安装"
    echo "请安装 FFmpeg: brew install ffmpeg (macOS)"
    echo ""
fi

# 确保目录存在
mkdir -p uploads outputs

# 启动服务器
echo "正在启动服务器..."
echo ""
go run cmd/server/main.go
