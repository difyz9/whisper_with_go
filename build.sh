#!/bin/bash
# 编译Go Whisper程序

# 设置环境变量
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"

# 编译
echo "正在编译..."
go build -o whisper-transcribe main.go

if [ $? -eq 0 ]; then
    echo "✓ 编译成功: ./whisper-transcribe"
    echo ""
    echo "使用方法:"
    echo "  DYLD_LIBRARY_PATH=/usr/local/lib ./whisper-transcribe -audio test.mp3"
    echo ""
    echo "或者创建一个别名:"
    echo "  alias whisper='DYLD_LIBRARY_PATH=/usr/local/lib $(pwd)/whisper-transcribe'"
else
    echo "✗ 编译失败"
    exit 1
fi
