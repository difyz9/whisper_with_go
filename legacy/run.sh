#!/bin/bash
# 运行Go Whisper程序的便捷脚本
# 自动设置必要的环境变量

# 设置动态库路径
export DYLD_LIBRARY_PATH=/usr/local/lib:$DYLD_LIBRARY_PATH
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_CFLAGS="-I/usr/local/include"

# 运行程序
# 支持所有参数: -audio, -model, -lang, -threads, -srt, -output
go run main.go "$@"
